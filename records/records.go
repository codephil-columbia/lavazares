package records

import (
	"encoding/json"
	"lavazares/utils"

	"github.com/jmoiron/sqlx"
)

// TutorialRecordManager manages Tutorial Records
// Amongst those include LessonCompleted records,
// ChapterCompleted records.
type TutorialRecordManager struct {
	lessonRecordStore  recordStore
	chapterRecordStore recordStore
	db                 sqlx.DB
}

// Save saves a record
func (manager *TutorialRecordManager) Save(record record) error {
	switch r := record.(type) {
	case LessonRecord:
		exists, err := manager.lessonRecordStore.exists(r)
		if err != nil {
			return err
		}
		if exists {
			return manager.lessonRecordStore.update(record)
		}
		return manager.lessonRecordStore.save(record)
	case ChapterRecord:
		err := manager.chapterRecordStore.save(r)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

type record interface {
}

type recordStore interface {
	save(r record) error
	update(r record) error
	exists(r record) (bool, error)
}

type ChapterRecord struct {
	ChapterID string `db:"chapterid" json:"chapterID"`
	UID       string `db:"uid" json:"uid"`
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

type LessonRecord struct {
	LessonID  string `json:"lessonID"`
	ChapterID string `json:"chapterID"`
	UID       string `json:"uid"`
	WPM       string `json:"wpm"`
	Accuracy  string `json:"accuracy"`
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
