package content

import (
	"errors"
	"log"
	"os"
	"sort"

	"github.com/jmoiron/sqlx"
)

var (
	errCompletedAllLessons  = errors.New("Completed all Lessons")
	errCompletedAllChapters = errors.New("Completed all Chapters")
)

// ContentManager is in charge of stuff that depend on both
// chapter and lesson managers (to avoid coupling the two)
// ie finding the next sequential lesson/chapter,
// finding the next lesson for a User
type ContentManager struct {
	chapterManager chapterManager
	lessonManager  lessonManager
	logger         *log.Logger
}

const contentManagerLogger = "ContentManager"

// NewContentManager returns a new ContentManager
func NewContentManager(db *sqlx.DB) *ContentManager {
	return &ContentManager{
		chapterManager: NewChapterManager(db),
		lessonManager:  NewLessonManager(db),
		logger:         log.New(os.Stdout, contentManagerLogger, log.Lshortfile),
	}
}

// GetNextLesson returns the next lesson in the sequential order the Users
// must complete.
// Unfortunately all of our Lesson and Chapter IDS are created at runtime by the DB,
// so we have to manually find the next through iteration.
// TODO: Maybe at some point we should embed pointers to the next lesson within each
// Lesson struct? At this point we don't have that many Lessons so it might not matter
// too much
func (m *ContentManager) GetNextLesson(lessonID string) (*Lesson, error) {
	lesson, err := m.lessonManager.GetLesson(lessonID)
	if err != nil {
		return nil, err
	}

	lessonsInChapter, err := m.GetLessonsInChapter(lesson.ChapterID)
	if err != nil {
		return nil, err
	}

	sort.Sort(lessons(lessonsInChapter))
	for i, lesson := range lessonsInChapter {
		if lesson.LessonID == lessonID {
			if i < len(lessonsInChapter)-1 {
				return lessonsInChapter[i+1], nil
			}
		}
	}

	nextChapter, err := m.getNextChapter(lesson.ChapterID)
	if err != nil {
		if err == errCompletedAllChapters {
			return nil, errCompletedAllLessons
		}
		return nil, err
	}

	lessonsInNextChapter, err := m.GetLessonsInChapter(nextChapter.ChapterID)
	if err != nil {
		return nil, err
	}

	return lessonsInNextChapter[0], nil
}

func (m *ContentManager) getNextChapter(chapterID string) (*Chapter, error) {
	lst, err := m.chapterManager.GetChapters()
	if err != nil {
		return nil, err
	}

	sort.Sort(chapters(lst))
	for i, chapter := range lst {
		if chapter.ChapterID == chapterID {
			if i < len(lst)-1 {
				return lst[i+1], nil
			}
			return nil, errCompletedAllChapters
		}
	}
	return nil, errors.New("Unexpected Err")
}

// GetLessonsInChapter returns a list of Lessons for a given ChapterID
func (m *ContentManager) GetLessonsInChapter(chapterID string) ([]*Lesson, error) {
	lessons, err := m.lessonManager.GetLessons()
	if err != nil {
		return nil, err
	}

	lessonsInChapter := []*Lesson{}
	for _, lesson := range lessons {
		if lesson.ChapterID == chapterID {
			lessonsInChapter = append(lessonsInChapter, lesson)
		}
	}
	return lessonsInChapter, nil
}
