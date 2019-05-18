package records

import (
	"encoding/json"
	"errors"
	"lavazares/content"
	"lavazares/utils"

	"github.com/jmoiron/sqlx"
)

var (
	errUserCompletedAllLessons = errors.New("User has no uncompleted lessons")
)

// Save saves a record using a record saver
func Save(saver recordSaver, record interface{}) error {
	return saver.Save(record)
}

// TutorialRecordManager manages Tutorial Records
// Amongst those include LessonCompleted records,
// ChapterCompleted records.
type TutorialRecordManager struct {
	chapterRecordStore *chapterRecordStore
	lessonManager      *content.LessonManagerI
	chapterManager     *content.ChapterManager
	contentManager     *content.ContentManager
}

// NewTutorialRecordManager returns an initialized TutorialRecordManager
func NewTutorialRecordManager(
	lessonRecordStore *lessonRecordStore,
	chapterRecordStore *chapterRecordStore,
	lessonManager *content.LessonManagerI,
	contentManager *content.ContentManager,
	chapterManager *content.ChapterManager,
	db *sqlx.DB,
) *TutorialRecordManager {
	return &TutorialRecordManager{
		chapterRecordStore: newChapterRecordStore(db),
		lessonManager:      lessonManager,
		contentManager:     contentManager,
		chapterManager:     chapterManager,
	}
}

// Export functions (Save, Query, QueryAll)

//func (m *TutorialRecordManager) save(record interface{}) error {
//	switch rec := record.(type) {
//	case ChapterRecord:
//		return m.saveChapterRecord(&rec)
//	case LessonRecord:
//		return m.saveLessonRecord(&rec)
//	default:
//		return errors.New("passed in unreconizable tutorial record")
//	}
//}

//func (m *TutorialRecordManager) saveLessonRecord(record *LessonRecord) error {
//	return m.lessonRecordStore.Save(record)
//}

func (m *TutorialRecordManager) SaveChapterRecord(record *ChapterRecord) error {
	return m.chapterRecordStore.save(record)
}

//// LessonStats returns stats on a User's individual lesson
//func (m *TutorialRecordManager) LessonStats(lessonID, uid string) (*LessonStats, error) {
//	record, err := m.lessonRecordStore.query(lessonID, uid)
//	if err != nil {
//		return nil, err
//	}
//
//	return &LessonStats{
//		WPM:      record.WPM,
//		Accuracy: record.Accuracy,
//		ID:       record.LessonID,
//		UID:      uid,
//	}, nil
//}

// GetNextNoncompletedLesson returns the next lesson that a User hasn't completed
// If the User has completed all Lessons, returns a nil lesson, and returns
//func (m *TutorialRecordManager) GetNextNoncompletedLesson(userid string) (*content.Lesson, error) {
//	// Get most recent chapter record, find the most recent lesson record in that chapter
//	currentChapter, err := m.GetNextNoncompletedChapter(userid)
//	if err != nil {
//		return nil, err
//	}
//
//	lessonsInChapter, err := m.contentManager.GetLessonsInChapter(currentChapter.ChapterID)
//	if err != nil {
//		return nil, err
//	}
//
//	completedLessons, err := m.lessonRecordStore.queryAll(userid)
//	if err != nil {
//		return nil, err
//	}
//
//	// ERROR IS HERE, U RETURN NILL AND NILL
//	intersection := lessonIntersection(completedLessons, lessonsInChapter)
//	if len(intersection) == 0 {
//		return nil, err
//	}
//
//	content.SortLessonsChrono(intersection)
//	return intersection[0], nil
//}

// GetNextNoncompletedChapter returns the most
func (m *TutorialRecordManager) GetNextNoncompletedChapter(userid string) (*content.Chapter, error) {
	chapterRecords, err := m.chapterRecordStore.queryAll(userid)
	if err != nil {
		return nil, err
	}

	chapters, err := m.chapterManager.GetChapters()
	if err != nil {
		return nil, err
	}

	intersection := chapterIntersection(chapterRecords, chapters)
	if len(intersection) == 0 {
		return nil, errors.New("User has completed all lessons")
	}

	content.SortChaptersChrono(intersection)
	return intersection[0], nil
}

// chapterIntersection takes the intersection between a list of User records and all
// available chapters.
func chapterIntersection(records []*ChapterRecord, chapters []*content.Chapter) []*content.Chapter {
	complete := make(map[string]bool)
	missing := []*content.Chapter{}

	for _, record := range records {
		complete[record.ChapterID] = true
	}

	for _, chapter := range chapters {
		if _, ok := complete[chapter.ChapterID]; !ok {
			missing = append(missing, chapter)
		}
	}

	if len(missing) == 0 {
		return nil
	}
	return missing
}

func lessonIntersection(records []*LessonRecord, lessons []*content.Lesson) []*content.Lesson {
	complete := make(map[string]bool)
	missing := []*content.Lesson{}

	for _, record := range records {
		complete[record.LessonID] = true
	}

	for _, lesson := range lessons {
		if _, ok := complete[lesson.LessonID]; !ok {
			missing = append(missing, lesson)
		}
	}

	if len(missing) == 0 {
		return nil
	}

	return missing
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
