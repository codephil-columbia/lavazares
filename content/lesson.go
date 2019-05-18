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
	Image              pq.StringArray `db:"image" json:"lessonImages"`
	NextLessonID       *string        `db:"nextlessonid"`
}

// Implements sort interface for []Lesson
type byLessonName []Lesson

func (l byLessonName) Len() int      { return len(l) }
func (l byLessonName) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byLessonName) Less(i, j int) bool {
	return comapreLessonNames(l[i], l[j])
}

// comapreLessonNames does a custom sort of []Lesson. Note that this sort
// will only work for groups of Lessons that are all belong to the same Chapter,
// and would not work as expected otherwise.
func comapreLessonNames(l1 Lesson, l2 Lesson) bool {
	// At the end of all but the first Chapter, the last chronological lesson
	// in a list of Lessons for a Chapter will be of the form "<Foo> Test"
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

// LessonManager handles operations on tutorial lessons
type LessonManager struct {
	store  lessonStore
	logger *log.Logger
}

const lessonManagerLoggerName = "LessonManager"

// NewLessonManager creates a new LessonManager.
// db is a pointer to the already initialized sqlx DB object
func NewLessonManager(db *sqlx.DB) *LessonManager {
	store := newLessonStore(db)
	return newLessonManager(store)
}

func newLessonManager(store lessonStore) *LessonManager {
	return &LessonManager{
		store:  store,
		logger: log.New(os.Stdout, lessonManagerLoggerName, log.Lshortfile),
	}
}

// SortLessonsChrono sorts a list of Lessons chronologically by LessonName
func SortLessonsByName(l []Lesson) {
	sort.Sort(byLessonName(l))
}

// GetLesson returns a lesson by id
func (m *LessonManager) GetLesson(id string) (*Lesson, error) {
	return m.store.Query(id)
}

// GetLessons returns a slice to all lessons
func (m *LessonManager) GetLessons() ([]*Lesson, error) {
	return m.store.QueryAll()
}

func (m *LessonManager) GetLessonsInChapter(chapterID string) ([]Lesson, error) {
	lessons, err := m.GetLessons()
	if err != nil {
		return nil, err
	}

	lessonsInChapter := []Lesson{}
	for _, l := range lessons {
		if l.ChapterID == chapterID {
			lessonsInChapter = append(lessonsInChapter, *l)
		}
	}

	sort.Sort(byLessonName(lessonsInChapter))

	return lessonsInChapter, nil
}

type LessonManagerI interface {
	GetLessons() ([]*Lesson, error)
	GetLesson(id string) (*Lesson, error)
}

// lessonStoreImpl is the object used to interact with underlying Lesson objects
// in the database.
type lessonStoreImpl struct {
	db *sqlx.DB
}

func newLessonStore(db *sqlx.DB) *lessonStoreImpl {
	return &lessonStoreImpl{db: db}
}

type lessonStore interface {
	Query(id string) (*Lesson, error)
	QueryAll() ([]*Lesson, error)
}

func (s *lessonStoreImpl) Query(ID string) (*Lesson, error) {
	var l Lesson
	err := s.db.QueryRowx("SELECT * FROM Lessons WHERE LessonID = $1", ID).StructScan(&l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (s *lessonStoreImpl) QueryAll() ([]*Lesson, error) {
	var all []*Lesson
	rows, err := s.db.Queryx("SELECT * FROM Lessons")
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
