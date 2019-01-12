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

// DefaultContentManager is in charge of stuff that depend on both
// chapter and lesson managers ie finding the next sequential lesson/chapter,
// finding the next lesson for a User
type DefaultContentManager struct {
	chapterManager *DefaultChapterManager
	lessonManager  *DefaultLessonManager
	logger         *log.Logger
}

const defaultContentManagerLogger = "DefaultContentManager"

// NewContentManager returns a new ContentManager
func NewContentManager(db *sqlx.DB) *DefaultContentManager {
	return &DefaultContentManager{
		chapterManager: NewDefaultChapterManager(db),
		lessonManager:  NewDefaultLessonManager(db),
		logger:         log.New(os.Stdout, "DefaultContentManager", log.Lshortfile),
	}
}

// GetNextLesson returns the next lesson in the sequential order the Users
// must complete.
// Unfortunately all of our Lesson and Chapter IDS are created at runtime by the DB,
// so we have to manually find the next through iteration.
// TODO: Maybe at some point we should embed pointers to the next lesson within each
// Lesson struct? At this point we don't have that many Lessons so it might not matter
// too much
func (manager *DefaultContentManager) GetNextLesson(lessonID string) (*Lesson, error) {
	lesson, err := manager.lessonManager.GetLesson(lessonID)
	if err != nil {
		return nil, err
	}

	lessonsInChapter, err := manager.getLessonsInChapter(lesson.ChapterID)
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

	nextChapter, err := manager.getNextChapter(lesson.ChapterID)
	if err != nil {
		if err == errCompletedAllChapters {
			return nil, errCompletedAllLessons
		}
		return nil, err
	}

	lessonsInNextChapter, err := manager.getLessonsInChapter(nextChapter.ChapterID)
	if err != nil {
		return nil, err
	}

	return lessonsInNextChapter[0], nil
}

func (manager *DefaultContentManager) getNextChapter(chapterID string) (*Chapter, error) {
	chapters, err := manager.chapterManager.GetChapters()
	if err != nil {
		return nil, err
	}

	for i, chapter := range chapters {
		if chapter.ChapterID == chapterID {
			if i < len(chapters)-1 {
				return chapters[i+1], nil
			}
			return nil, errCompletedAllChapters
		}
	}
	return nil, errors.New("Unexpected Err")
}

func (manager *DefaultContentManager) getLessonsInChapter(chapterID string) ([]*Lesson, error) {
	lessons, err := manager.lessonManager.GetLessons()
	if err != nil {
		return nil, err
	}

	lessonsInChapter := lessons
	for _, lesson := range lessons {
		if lesson.ChapterID == chapterID {
			lessonsInChapter = append(lessonsInChapter, lesson)
		}
	}
	return lessons, nil
}
