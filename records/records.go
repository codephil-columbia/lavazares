package records

import (
	"encoding/json"
	"lavazares/utils"

	"github.com/jmoiron/sqlx"
)

// LessonRecord represents a completed Lesson finished
// by a User in the Tutorial mode
// type LessonRecord struct {
// 	LessonID  string `json:"lessonID" db:"lessonid"`
// 	ChapterID string `json:"chapterID" db:"chapterid"`
// 	UID       string `json:"uid" db:"uid"`
// 	WPM       string `json:"wpm" db:"wpm"`
// 	Accuracy  string `json:"accuracy" db:"accuracy"`
// }

type record interface {
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

// NewLessonRecord creates a new LessonRecord
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

type TutorialRecordManager struct {
	lessonRecordStore recordStore
	db                sqlx.DB
}

type recordStore interface {
	save(r record) error
	update(r record) error
	exists(r record) (bool, error)
}

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
	default:
		return nil
	}
}
