package records

import (
	"encoding/json"
	"errors"
	"lavazares/content"
	"lavazares/utils"

	"github.com/jmoiron/sqlx"
)

type lessonID string
type uid string

var errUserCompletedAllLessons = errors.New("User has no uncompleted lessons.")

// TutorialRecordManager manages Tutorial Records
// Amongst those include LessonCompleted records,
// ChapterCompleted records.
type TutorialRecordManager struct {
	lessonRecordStore  *lessonRecordStore
	chapterRecordStore *chapterRecordStore
	lessonManager      *content.DefaultLessonManager
	db                 *sqlx.DB
}

func NewTutorialRecordManager(db *sqlx.DB) *TutorialRecordManager {
	return &TutorialRecordManager{
		lessonRecordStore:  newLessonRecordStore(db),
		chapterRecordStore: newChapterRecordStore(db),
		lessonManager:      content.NewDefaultLessonManager(db),
		db:                 db,
	}
}

// Save saves a record
func (manager *TutorialRecordManager) Save(record interface{}) error {
	switch r := record.(type) {
	case LessonRecord:
		exists, err := manager.lessonRecordStore.exists(&r)
		if err != nil {
			return err
		}
		if exists {
			return manager.lessonRecordStore.update(&r)
		}
		return manager.lessonRecordStore.save(&r)
	case ChapterRecord:
		err := manager.chapterRecordStore.save(&r)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

// GetNextNoncompletedLesson returns the next lesson that a User hasn't completed
// If the User has completed all Lessons, returns a nil lesson, and returns
// func (manager *TutorialRecordManager) GetNextNoncompletedLesson(id uid) (*content.Lesson, error) {
// 	// completed, err := manager.lessonRecordStore
// }

type tutorialRecord interface {
	uid() string
	id() string
}

type MissingRequiredFieldErr struct {
	missingFields []string
}

type recordStore interface {
	// TODO: interface parameters should at some point not be interface
	// its interface for now because stats are currently embedded in Lesson records
	// future records should separate id/uid from stats(wpm/time/accuracy)
	save(record tutorialRecord) error
	update(record tutorialRecord) error
	exists(record tutorialRecord) (bool, error)
	query(id lessonID, uid uid) (tutorialRecord, error)
	queryAll(uid uid) (tutorialRecord, error)
}

type ChapterRecord struct {
	ChapterID string `json:"chapterID"`
	UID       string `json:"uid"`
}

func (c ChapterRecord) uid() string {
	return c.UID
}

func (c ChapterRecord) id() string {
	return c.ChapterID
}

type chapterRecordStore struct {
	db *sqlx.DB
}

func newChapterRecordStore(db *sqlx.DB) *chapterRecordStore {
	return &chapterRecordStore{db: db}
}

// NewChapterRecord creates a new ChapterRecord from an incoming
// JSON object.
func NewChapterRecord(args utils.RequestJSON) (*ChapterRecord, error) {
	var record ChapterRecord
	err := json.Unmarshal(args, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (store *chapterRecordStore) save(record *ChapterRecord) error {
	_, err := store.db.Exec(
		`INSERT INTO ChaptersCompleted(ChapterID, UID)
		VALUES($1, $2)`,
		record.ChapterID,
		record.UID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *chapterRecordStore) update(record *ChapterRecord) error {
	return nil
}

func (store *chapterRecordStore) exists(record *ChapterRecord) (bool, error) {
	var count int
	err := store.db.Get(
		&count,
		`SELECT count(*) FROM ChaptersCompleted 
		 WHERE chapterid=$1 and uid=$2`,
		record.ChapterID,
		record.UID,
	)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}

func (store *chapterRecordStore) queryAll(uid string) ([]*ChapterRecord, error) {
	var all []*ChapterRecord
	rows, err := store.db.Queryx("SELECT * FROM ChaptersCompleted WHERE uid = $1", uid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := ChapterRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			rows.Close()
			return nil, err
		}
		all = append(all, &record)
	}

	return all, nil
}

func (store *chapterRecordStore) query(id, uid string) (*ChapterRecord, error) {
	var record ChapterRecord
	err := store.db.QueryRowx(
		"SELECT * FROM ChaptersCompleted WHERE uid = $1 AND chapterid = $2",
		uid, id).StructScan(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

type LessonRecord struct {
	LessonID  string `json:"lessonID"`
	ChapterID string `json:"chapterID"`
	UID       string `json:"uid"`
	WPM       string `json:"wpm"`
	Accuracy  string `json:"accuracy"`
}

func (l LessonRecord) uid() string {
	return l.UID
}

func (l LessonRecord) id() string {
	return l.LessonID
}

type lessonRecordStore struct {
	db *sqlx.DB
}

func newLessonRecordStore(db *sqlx.DB) *lessonRecordStore {
	return &lessonRecordStore{db: db}
}

// NewLessonRecord creates a new LessonRecord from an
// incoming JSON object
func NewLessonRecord(args utils.RequestJSON) (*LessonRecord, error) {
	var record LessonRecord
	err := json.Unmarshal(args, &record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (store *lessonRecordStore) save(record *LessonRecord) error {
	_, err := store.db.Exec(
		`INSERT INTO 
		LessonsCompleted(LessonID, UID, WPM, Accuracy, ChapterID)
		VALUES($1, $2, $3, $4, $5)`,
		record.LessonID,
		record.UID,
		record.WPM,
		record.Accuracy,
		record.ChapterID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *lessonRecordStore) update(record *LessonRecord) error {
	_, err := store.db.Exec(
		`UPDATE LessonsCompleted SET wpm=$1, accuracy=$2 
		 WHERE uid=$3 AND lessonid=$4`,
		record.WPM,
		record.Accuracy,
		record.UID,
		record.LessonID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (store *lessonRecordStore) exists(record *LessonRecord) (bool, error) {
	var count int
	err := store.db.Get(
		&count,
		`SELECT count(*) FROM LessonsCompleted 
		 WHERE lessonid=$1 and uid=$2`,
		record.LessonID,
		record.UID,
	)
	if err != nil {
		return false, err
	}
	return count != 0, nil
}
