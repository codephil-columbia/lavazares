package content

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

// Lesson stores information about a specific lesson
// Most lessons have been hardcoded atm.
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

type lessons []*Lesson

func (l lessons) Len() int      { return len(l) }
func (l lessons) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l lessons) Less(i, j int) bool {
	return comapreLessonNames(l[i], l[j])
}

// comapreLessonNames does a custom sort of []Lesson. Note that this sort
// will only work for groups of Lessons that are all belong to the same Chapter,
// and would not work as expected otherwise.
func comapreLessonNames(l1 *Lesson, l2 *Lesson) bool {
	// At the end of all but the first Chapter, the last chronological lesson
	// in a list of Lessons for a Chapter will be of the form "Foo Test"
	if strings.Contains(l1.LessonName, "Test") {
		return false
	}
	if strings.Contains(l2.LessonName, "Test") {
		return true
	}

	// If neither lesson has the form "Foo Test", we can use the built in string
	// Compare which sorts them lexographically.
	if l1.LessonName > l2.LessonName {
		return false
	}
	return true
}

// DefaultLessonManager handles most of the basic operations on generic
// Lesson objects
// All operations on Lessons through the DefaultLessonManager for now
// are read only
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

// SortChrono sorts a list of Lessons chronologically
func SortChrono(l *lessons) {
	sort.Sort(*l)
}

// GetLesson returns a lesson by id
func (manager *DefaultLessonManager) GetLesson(id string) (*Lesson, error) {
	return manager.store.Query(id)
}

// GetLessons returns a slice to all lessons
func (manager *DefaultLessonManager) GetLessons() ([]*Lesson, error) {
	return manager.store.QueryAll()
}

type lessonManager interface {
	GetLessons() ([]*Lesson, error)
	GetLesson(id string) (*Lesson, error)
}

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
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		l := Lesson{}
		err = rows.StructScan(&l)
		if err != nil {
			return nil, err
		}
		all = append(all, &l)
	}
	return all, nil
}
