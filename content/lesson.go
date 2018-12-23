package content

import (
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

// Lesson struct
type Lesson struct {
	LessonID           string         `db:"lessonid" json:"lessonID"`
	CreatedAt          time.Time      `db:"createdat" json:"createdAt"`
	UpdatedAt          time.Time      `db:"updatedat" json:"updatedAt"`
	DeletedAt          *time.Time     `db:"deletedat" json:"deletedAt"`
	LessonName         string         `db:"lessonname" json:"lessonName"`
	LessonText         pq.StringArray `db:"lessontext" json:"lessonText"`
	LessonDescriptions pq.StringArray `db:"lessondescriptions" json:"lessonDescriptions"`
	MinimumScoreToPass pq.Int64Array  `db:"minimumscoretopass" json:"minimumScoreToPass"`
	ChapterID          string         `db:"chapterid" json:"chapterID"`
	Image              pq.StringArray `db:"image" json:"image"`
}

// // ContentStatsManager
// type ContentStatsManager interface {
// 	SetStatsForContent(stats *Statistic, uid, contentID string) error
// }

// DefaultLessonManager handles most of the basic operations on generic
// Lesson objects
type DefaultLessonManager struct {
	store  *lessonStore
	logger *log.Logger
}

const defaultLessonManagerLoggerName = "DefaultLessonManager"

// NewDefaultLessonManager creates a new DefaultLessonManager.
// db is a pointer to the already initialized sqlx DB object
func NewDefaultLessonManager(db *sqlx.DB) *DefaultLessonManager {
	return &DefaultLessonManager{
		store:  newLessonStore(db),
		logger: log.New(os.Stdout, defaultLessonManagerLoggerName, log.Lshortfile),
	}
}

// // MarkLessonAsComplete marks a Lesson completed by a User
// func (dlm *DefaultLessonManager) MarkLessonAsComplete(uid, lessonID string) error {
// 	l, err := dlm.store.Query(lessonID)
// }

// lessonStore is the object used to interact with underlying Lesson objects
// in the database.
type lessonStore struct {
	db *sqlx.DB
}

func newLessonStore(db *sqlx.DB) *lessonStore {
	return &lessonStore{db: db}
}

func (store *lessonStore) Query(ID string) (*Lesson, error) {
	var l Lesson
	err := store.db.QueryRowx("SELECT * FROM Lessons WHERE LessonID = $1", ID).StructScan(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (store *lessonStore) QueryAll() ([]*Lesson, error) {
	var all []*Lesson
	rows, err := store.db.Queryx("SELECT * FROM Lessons")
	for rows.Next() {
		var l *Lesson
		err = rows.StructScan(&l)
		if err != nil {
			rows.Close()
			return nil, err
		}
		all = append(all, l)
	}
	return all, nil
}
