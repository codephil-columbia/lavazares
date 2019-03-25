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

// TutorialRecordManager manages Tutorial Records
// Amongst those include LessonCompleted records,
// ChapterCompleted records.
type TutorialRecordManager struct {
	lessonRecordStore  *lessonRecordStore
	chapterRecordStore *chapterRecordStore
	lessonManager      *content.LessonManager
	chapterManager     *content.ChapterManager
	contentManager     *content.ContentManager
	db                 *sqlx.DB
}

// NewTutorialRecordManager returns an initialized TutorialRecordManager
func NewTutorialRecordManager(db *sqlx.DB) *TutorialRecordManager {
	return &TutorialRecordManager{
		lessonRecordStore:  newLessonRecordStore(db),
		chapterRecordStore: newChapterRecordStore(db),
		lessonManager:      content.NewLessonManager(db),
		contentManager:     content.NewContentManager(db),
		chapterManager:     content.NewChapterManager(db),
		db:                 db,
	}
}

// Save saves a record
func (m *TutorialRecordManager) Save(record interface{}) error {
	switch r := record.(type) {
	case LessonRecord:
		exists, err := m.lessonRecordStore.exists(&r)
		if err != nil {
			return err
		}
		if exists {
			return m.lessonRecordStore.update(&r)
		}
		return m.lessonRecordStore.save(&r)
	case ChapterRecord:
		err := m.chapterRecordStore.save(&r)
		if err != nil {
			return err
		}
		return nil
	default:
		return nil
	}
}

func (m *TutorialRecordManager) LessonStats(lessonID, uid string) (*LessonStats, error) {
	record, err := m.lessonRecordStore.query(lessonID, uid)
	if err != nil {
		return nil, err
	}

	return &LessonStats{
		WPM:      record.WPM,
		Accuracy: record.Accuracy,
		ID:       record.LessonID,
		UID:      uid,
	}, nil
}

// TutorialHollisticStats returns hollistic stats for a User's Tutorial Records
func (m *TutorialRecordManager) LessonsStats(uid string) ([]*LessonStats, error) {
	records, err := m.lessonRecordStore.queryAll(uid)
	if err != nil {
		return nil, err
	}

	stats := []*LessonStats{}
	for _, r := range records {
		stats = append(stats, &LessonStats{
			WPM:      r.WPM,
			Accuracy: r.Accuracy,
			ID:       r.LessonID,
			UID:      uid,
		})
	}

	return stats, nil
}

// GetNextNoncompletedLesson returns the next lesson that a User hasn't completed
// If the User has completed all Lessons, returns a nil lesson, and returns
func (m *TutorialRecordManager) GetNextNoncompletedLesson(userid string) (*content.Lesson, error) {
	// Get most recent chapter record, find the most recent lesson record in that chapter
	currentChapter, err := m.GetNextNoncompletedChapter(userid)
	if err != nil {
		return nil, err
	}

	lessonsInChapter, err := m.contentManager.GetLessonsInChapter(currentChapter.ChapterID)
	if err != nil {
		return nil, err
	}

	completedLessons, err := m.lessonRecordStore.queryAll(userid)
	if err != nil {
		return nil, err
	}

	intersection := lessonIntersection(completedLessons, lessonsInChapter)
	if len(intersection) == 0 {
		return nil, err
	}

	content.SortLessonsChrono(intersection)
	return intersection[0], nil
}

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

// GetLessonRecords returns all of the LessonRecords for a User
func (m *TutorialRecordManager) GetLessonRecords(uid string) ([]*LessonRecord, error) {
	return m.lessonRecordStore.queryAll(uid)
}

type tutorialRecord interface {
	uid() string
	id() string
}

type recordStore interface {
	// TODO: interface parameters should at some point not be interface
	// its interface for now because stats are currently embedded in Lesson records
	// future records should separate id/uid from stats(wpm/time/accuracy)
	save(record tutorialRecord) error
	update(record tutorialRecord) error
	exists(record tutorialRecord) (bool, error)
	query(id string, uid string) (tutorialRecord, error)
	queryAll(uid string) (tutorialRecord, error)
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

func (store *lessonRecordStore) query(lessonID, uid string) (*LessonRecord, error) {
	var record LessonRecord
	err := store.db.QueryRowx(
		"SELECT * FROM LessonsCompleted WHERE uid = $1 AND lessonid = $2",
		uid, lessonID).StructScan(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (store *lessonRecordStore) queryAll(uid string) ([]*LessonRecord, error) {
	var all []*LessonRecord
	rows, err := store.db.Queryx("SELECT * FROM LessonsCompleted WHERE uid = $1", uid)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := LessonRecord{}
		err = rows.StructScan(&record)
		if err != nil {
			rows.Close()
			return nil, err
		}
		all = append(all, &record)
	}

	return all, nil
}
