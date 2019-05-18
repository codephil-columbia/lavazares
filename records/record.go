package records

import (
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"lavazares/content"
	"strconv"
)

type (
	LessonsRecordManager struct {
		lessonRecordStore
	}

	ChapterRecordManager struct {
		chapterRecordStoreI
	}
)

const (
	QUERY_LESSON_RECORD      = "SELECT * FROM lessonscompleted WHERE lessonid=:lessonID AND chapterid=:chapterID AND uid=:uid"
	QUERY_ALL_LESSON_RECORDS = "SELECT * FROM lessonscompleted WHERE uid=:uid"

	QUERY_CHAPTER_RECORD      = "SELECT * FROM chapterscompleted WHERE chapterid=:chapterid AND uid=:uid"
	QUERY_ALL_CHAPTER_RECORDS = "SELECT * FROM chapterscompleted WHERE uid=:uid"

	QUERY_LESSON_RECORD_STATS  = "SELECT WPM, Accuracy, UID, LessonID FROM lessonscompleted WHERE lessonid=:lessonid AND uid=:uid"
	QUERY_LESSONS_RECORD_STATS = "SELECT WPM, Accuracy, UID, LessonID FROM lessonscompleted WHERE uid=:uid"
)

const (
	firstLessonID  = "1"
	firstChapterID = "1"
)

var (
	ErrUserHasCompletedAllLessons  = errors.New("User has completed all lessons")
	ErrUserHasCompletedAllChapters = errors.New("User has completed all chapters")
)

var (
	chapterOrder = []string{
		"Chapter 0: The Basics",
		"Chapter 1: Home Row",
		"Chapter 2: Shift and Basic Punctuations",
		"Chapter 3: Top Row",
		"Chapter 4: Bottom Row",
		"Chapter 5: Number Row",
		"Chapter 6: Advanced Content",
	}

	// lessonid:chapterid
	lastLessonInChapter = map[string]string{
		"1": "6",
		"2": "12",
		"3": "25",
		"4": "30",
		"5": "36",
		"6": "40",
	}
)

func QueryLessonRecord(queryer recordQueryer, dest LessonRecord, params QueryLessonRecordParams) error {
	return query(queryer, QUERY_LESSON_RECORD, dest, params)
}

func QueryLessonRecords(queryer recordsQueryer, dest *[]LessonRecord, params QueryLessonsRecordParams) error {
	return queryer.QueryAll(QUERY_ALL_LESSON_RECORDS, dest, params)
}

func SaveLessonRecord(saver, chapterSaver recordSaver, record LessonRecord) error {
	if wasLastLessonInChapter(record.LessonID, record.ChapterID) {
		if err := SaveChapterRecord(chapterSaver, ChapterRecord{UID: record.UID, ChapterID: record.ChapterID}); err != nil {
			return err
		}
	}
	return saver.Save(record)
}

func SaveChapterRecord(saver recordSaver, record ChapterRecord) error {
	return saver.Save(record)
}

func UpdateLessonRecord(updater recordUpdater, record LessonRecord) error {
	return updater.Update(record)
}

func QueryChapterRecord(queryer recordQueryer, dest ChapterRecord, params QueryChapterRecordParams) error {
	return query(queryer, QUERY_CHAPTER_RECORD, dest, params)
}

func QueryChaptersRecords(queyer recordsQueryer, dest *[]ChapterRecord, params QueryChaptersRecordParams) error {
	return queryAll(queyer, QUERY_ALL_CHAPTER_RECORDS, dest, params)
}

func QueryLessonStats(queryer statQueryer, dest LessonStats, params QueryLessonStatsParams) error {
	return query(queryer, QUERY_LESSON_RECORD_STATS, dest, params)
}

func QueryLessonsStats(queryer statsQueryer, dest *[]LessonRecord, params QueryLessonsStatsParams) error {
	return queryAll(queryer, QUERY_LESSONS_RECORD_STATS, dest, params)
}

func QueryChapterProgress(queyer recordsQueryer, params QueryChaptersRecordParams) (*chapterProgress, error) {
	return queryChapterProgressPercentage(queyer, params)
}

func QueryAvgTutorialStats(queryer statsQueryer, uid string) (*avgTutorialStats, error) {
	stats := []LessonRecord{}
	err := queryer.QueryAll(QUERY_ALL_LESSON_RECORDS, &stats, QueryLessonsRecordParams{UID: uid})
	if err != nil {
		return nil, err
	}

	wpms := []string{}
	accuracies := []string{}
	for _, stat := range stats {
		wpms = append(wpms, stat.WPM)
		accuracies = append(accuracies, stat.Accuracy)
	}

	avgWPM, err := getAvg(wpms)
	if err != nil {
		return nil, err
	}

	avgAccuracy, err := getAvg(accuracies)
	if err != nil {
		return nil, err
	}

	fmt.Println(avgAccuracy, avgWPM)
	return &avgTutorialStats{
		WPM:      strconv.Itoa(avgWPM),
		Accuracy: strconv.Itoa(avgAccuracy),
		UID:      uid,
	}, nil
}

func queryChapterProgressPercentage(queryer recordsQueryer, params QueryChaptersRecordParams) (*chapterProgress, error) {
	var records []ChapterRecord
	err := queryer.QueryAll(QUERY_ALL_CHAPTER_RECORDS, &records, params)
	if err != nil {
		return nil, err
	}

	return &chapterProgress{
		PercentComplete: strconv.Itoa(len(records) / len(chapterOrder)),
	}, nil
}

func query(queryer recordQueryer, stmt string, dest, params interface{}) error {
	return queryer.Query(stmt, &dest, params)
}

func queryAll(queyer recordsQueryer, stmt string, dest, params interface{}) error {
	return queyer.QueryAll(stmt, &dest, params)
}

func CurrentLesson(queryer recordsQueryer, lessonManager *content.LessonManager, uid string) (*content.Lesson, error) {
	records := []LessonRecord{}
	err := queryer.QueryAll(QUERY_ALL_LESSON_RECORDS, &records, QueryLessonsRecordParams{UID: uid})
	if err != nil {
		return nil, err
	}

	// If the user has no lesson records, then they haven't completed any lessons -> return the first lesson
	if len(records) == 0 {
		return lessonManager.GetLesson(firstLessonID)
	}

	lastCompletedRecord := findLastCompletedLessonRecord(records)
	lastCompletedLesson, err := lessonManager.GetLesson(lastCompletedRecord.LessonID)
	if err != nil {
		return nil, err
	}

	if *lastCompletedLesson.NextLessonID == "" {
		return nil, ErrUserHasCompletedAllLessons
	} else {
		return lessonManager.GetLesson(*lastCompletedLesson.NextLessonID)
	}
}

func CurrentChapter(queryer recordsQueryer, manager *content.ChapterManager, uid string) (*content.Chapter, error) {
	var records []ChapterRecord
	err := queryer.QueryAll(QUERY_ALL_CHAPTER_RECORDS, &records, QueryChaptersRecordParams{UID: uid})
	if err != nil {
		return nil, err
	}

	// If user does not have any ChapterRecords, then they haven't completed any chapters -> give them the first one
	if len(records) == 0 {
		return manager.GetChapter(firstChapterID)
	}

	lastCompletedRecord := findLastCompletedChapterRecord(records)
	lastCompletedChapter, err := manager.GetChapter(lastCompletedRecord.ChapterID)
	if err != nil {
		return nil, err
	}

	if *lastCompletedChapter.NextChapterID == "" {
		return nil, ErrUserHasCompletedAllChapters
	} else {
		return manager.GetChapter(*lastCompletedChapter.NextChapterID)
	}
}

func wasLastLessonInChapter(lessonID, chapterID string) bool {
	lastLessonInChapterID, ok := lastLessonInChapter[chapterID]
	if !ok {
		return false
	}

	// In case the completed Lesson was the last one for the Chapter, make sure to save a Chapter record
	// since the Chapter was also completed
	if lastLessonInChapterID == lessonID {
		return true
	}
	return false
}

func getAvg(elems []string) (int, error) {
	if len(elems) == 0 {
		return 0, nil
	}

	sum := 0
	for _, val := range elems {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		sum += intVal
	}

	return sum / len(elems), nil
}

func findLastCompletedLessonRecord(records []LessonRecord) LessonRecord {
	lastCompletedID := ""
	var lastCompletedRecord LessonRecord

	for _, record := range records {
		if record.LessonID > lastCompletedID {
			lastCompletedID = record.LessonID
			lastCompletedRecord = record
		}
	}
	return lastCompletedRecord
}

func findLastCompletedChapterRecord(completedRecords []ChapterRecord) ChapterRecord {
	lastCompletedID := ""
	var lastCompletedRecord ChapterRecord

	for _, record := range completedRecords {
		if record.ChapterID > lastCompletedID {
			lastCompletedID = record.ChapterID
			lastCompletedRecord = record
		}
	}
	return lastCompletedRecord
}

func NewLessonRecordManager(db *sqlx.DB) *LessonsRecordManager {
	return newLessonRecordManagerStore(lessonRecordStore{db})
}

func newLessonRecordManagerStore(store lessonRecordStore) *LessonsRecordManager {
	return &LessonsRecordManager{
		store,
	}
}

func NewChapterRecordManager(db *sqlx.DB) *ChapterRecordManager {
	return newChapterRecordManager(chapterRecordStoreI{DB: db})
}

func newChapterRecordManager(store chapterRecordStoreI) *ChapterRecordManager {
	return &ChapterRecordManager{
		store,
	}
}

type (
	conversionError struct {
		err string
	}

	storeErr struct {
		err string
	}
)

type (
	QueryLessonRecordParams struct {
		UID       string `db:"uid"`
		ChapterID string `db:"chapterID"`
		LessonID  string `db:"lessonID"`
	}

	QueryLessonsRecordParams struct {
		UID string `db:"uid"`
	}

	QueryChapterRecordParams struct {
		UID       string `db:"uid"`
		ChapterID string `db:"chapterid"`
	}

	QueryChaptersRecordParams struct {
		UID string `db:"uid"`
	}

	QueryLessonStatsParams struct {
		LessonID string `db:"lessonid"`
		UID      string `db:"uid"`
	}

	QueryLessonsStatsParams struct {
		UID string `db:"uid"`
	}
)

func (e conversionError) Error() string {
	return fmt.Sprintf("Error converting: %v", e.err)
}

func (e storeErr) Error() string {
	return fmt.Sprintf("Error executing sql: %v", e.err)
}

type recordQueryer interface {
	Query(stmt string, dest, queryParams interface{}) error
}

type recordsQueryer interface {
	QueryAll(stmt string, dest, queryParams interface{}) error
}

type recordSaver interface {
	Save(record interface{}) error
}

type recordUpdater interface {
	Update(record interface{}) error
}

type lessonRecordStore struct {
	*sqlx.DB
}

func newLessonRecordStore(db *sqlx.DB) *lessonRecordStore {
	return &lessonRecordStore{db}
}

func (store *lessonRecordStore) Save(record interface{}) error {
	r, ok := record.(LessonRecord)
	if !ok {
		return conversionError{err: "Could not convert Lesson record"}
	}

	if r.Accuracy == "" {
		r.Accuracy = "0"
	}

	if r.WPM == "" {
		r.WPM = "0"
	}

	_, err := store.NamedExec(`INSERT INTO 
		LessonsCompleted(LessonID, UID, WPM, Accuracy, ChapterID)
		VALUES(:lessonid, :uid, :wpm, :accuracy, :chapterid)`,
		r,
	)
	if err != nil {
		return storeErr{err: err.Error()}
	}
	return nil
}

func (store *lessonRecordStore) Update(record interface{}) error {
	r, ok := record.(LessonRecord)
	if !ok {
		return conversionError{err: "Could not convert to LessonRecord"}
	}

	_, err := store.NamedExec(
		`UPDATE LessonsCompleted SET WPM=:wpm, accuracy=:accuracy 
		 WHERE uid=:uid AND lessonid=:lessonid`,
		r,
	)
	if err != nil {
		return storeErr{err: err.Error()}
	}

	return nil
}

func (store *lessonRecordStore) exists(record *LessonRecord) (bool, error) {
	var count int
	err := store.Get(
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

func (store *lessonRecordStore) QueryAll(stmt string, dest, queryParam interface{}) error {
	namedStmt, err := store.PrepareNamed(stmt)
	if err != nil {
		return storeErr{err: err.Error()}
	}

	err = namedStmt.Select(dest, queryParam)
	if err != nil {
		return storeErr{err: err.Error()}
	}

	return nil
}

func (store *lessonRecordStore) Query(stmt string, dest interface{}, queryParams interface{}) error {
	namedStmt, err := store.PrepareNamed(stmt)
	if err != nil {
		return err
	}

	err = namedStmt.Get(dest, queryParams)
	if err != nil {
		return storeErr{err: err.Error()}
	}
	return nil
}

type chapterRecordStoreI struct {
	*sqlx.DB
}

func newChapterRecordStoreI(db *sqlx.DB) *chapterRecordStoreI {
	return &chapterRecordStoreI{
		db,
	}
}

func (store *chapterRecordStoreI) Query(stmt string, dest interface{}, queryParams interface{}) error {
	namedStmt, err := store.PrepareNamed(stmt)
	if err != nil {
		return err
	}

	err = namedStmt.Get(dest, queryParams)
	if err != nil {
		return storeErr{err: err.Error()}
	}

	return nil
}

func (store *chapterRecordStoreI) QueryAll(stmt string, dest, queryParams interface{}) error {
	namedStmt, err := store.PrepareNamed(stmt)
	if err != nil {
		return err
	}

	err = namedStmt.Select(dest, queryParams)
	if err != nil {
		return storeErr{err: err.Error()}
	}

	return nil
}

// ChapterRecord represents a User finishing a specific chapter
type ChapterRecord struct {
	ChapterID string `json:"chapterID" db:"chapterid"`
	UID       string `json:"uid" db:"uid"`
}

// LessonRecord represents a User finishing a specific Lesson
// Note that this struct also contains information on per Lesson
// User typing statistics
type LessonRecord struct {
	LessonID  string `json:"lessonID" db:"lessonid"`
	ChapterID string `json:"chapterID" db:"chapterid"`
	UID       string `json:"uid" db:"uid"`
	WPM       string `json:"wpm" db:"wpm"`
	Accuracy  string `json:"accuracy" db:"accuracy"`
}

func (store *chapterRecordStoreI) Save(record interface{}) error {
	_, err := store.NamedExec("INSERT INTO chapterscompleted(uid, chapterid) VALUES(:uid, :chapterid)", record)
	if err != nil {
		return err
	}
	return nil
}
