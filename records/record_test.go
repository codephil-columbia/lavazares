package records

import (
	"fmt"
	"lavazares/auth"
	"lavazares/content"
	"log"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	testLesson = content.Lesson{
		LessonID:  "200",
		ChapterID: "1",
	}

	testChapter = content.Chapter{
		ChapterID: "1",
	}

	testUser1 = auth.User{
		UID:        "1",
		Password:   "password",
		Email:      "tester@gmail.com",
		Occupation: "student",
		FirstName:  "Tester",
		LastName:   "Tester",
		Username:   "TestUser1",
	}

	testUser2 = auth.User{
		UID:        "2",
		Password:   "password",
		Email:      "tester2@gmail.com",
		Occupation: "Student",
		FirstName:  "Tester",
		LastName:   "Tester",
		Username:   "TestUser2",
	}

	testReadOnlyLessonRecord = LessonRecord{
		UID:       "1",
		LessonID:  "200",
		ChapterID: "1",
		Accuracy:  "100",
		WPM:       "100",
	}

	testReadOnlyChapterRecord = ChapterRecord{
		UID:       "1",
		ChapterID: "1",
	}

	conn       = "port=5432 host=localhost sslmode=disable user=postgres dbname=postgres"
	cleanUpStr = "DELETE FROM LessonsCompleted WHERE LessonID=:lessonid AND UID=:uid"

	setupData = []testData{
		{
			`INSERT INTO Users(UID, Password, Email, Occupation, FirstName, LastName, Username)
					VALUES(:uid, :password, :email, :occupation, :firstname, :lastname, :username);`,
			testUser1,
		},
		{
			`INSERT INTO Users(UID, Password, Email, Occupation, FirstName, LastName, Username)
					VALUES(:uid, :password, :email, :occupation, :firstname, :lastname, :username);`,
			testUser2,
		},
		{
			`INSERT INTO Chapters(ChapterID) VALUES(:chapterid)`,
			testChapter,
		},
		{
			`INSERT INTO Lessons(LessonID, ChapterID) VALUES(:lessonid, :chapterid)`,
			testLesson,
		},
		{
			`INSERT INTO LessonsCompleted(UID, LessonID, ChapterID, Accuracy, WPM) 
					VALUES(:uid, :lessonid, :chapterid, :accuracy, :wpm)`,
			testReadOnlyLessonRecord,
		},
		{
			`INSERT INTO ChaptersCompleted(UID, ChapterID)
					VALUES(:uid, :chapterid)`,
			testReadOnlyChapterRecord,
		},
	}

	// Ordering here is important since this follows our DB schema dependencies
	cleanUpData = []testData{
		{
			`DELETE FROM Users WHERE uid=:uid`,
			testUser2,
		},
		{
			`DELETE FROM LessonsCompleted WHERE lessonid=:lessonid AND chapterid=:chapterid AND uid=:uid`,
			testReadOnlyLessonRecord,
		},
		{
			`DELETE FROM ChaptersCompleted WHERE uid=:uid AND chapterid=:chapterid`,
			testReadOnlyChapterRecord,
		},
		{
			`DELETE FROM Lessons WHERE lessonid=:lessonid`,
			testLesson,
		},
		{
			`DELETE FROM Chapters WHERE chapterid=:chapterid`,
			testChapter,
		},
		{
			`DELETE FROM Users WHERE uid=:uid`,
			testUser1,
		},
	}
)

type (
	testDB struct {
		*sqlx.DB
	}

	cleanUpError struct {
		FailedVal string
	}

	testData struct {
		stmt string
		val  interface{}
	}
)

func (err cleanUpError) Error() string {
	return fmt.Sprintf("Failed to clean up testing db value %v", err.FailedVal)
}

func newTestDB(connStr string) *testDB {
	db := sqlx.MustConnect("postgres", connStr)

	tx, err := db.Beginx()
	if err != nil {
		log.Fatalln(err)
	}

	for _, setup := range setupData {
		_, err := tx.NamedExec(setup.stmt, setup.val)
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}
	}
	tx.Commit()
	return &testDB{db}
}

func cleanUpDB(db *testDB, toClean []testData) {
	tx, _ := db.Beginx()
	for _, data := range toClean {
		_, err := tx.NamedExec(data.stmt, data.val)
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}
	}
	tx.Commit()
}

func cleanUpTest(db *sqlx.DB, cleanUpStr string, bindFrom interface{}) error {
	namedStmt, err := db.PrepareNamed(cleanUpStr)
	if err != nil {
		return err
	}

	_, err = namedStmt.Exec(bindFrom)
	if err != nil {
		return err
	}

	return nil
}

func TestRecordSavers(t *testing.T) {
	testDB := newTestDB(conn)
	defer cleanUpDB(testDB, cleanUpData)

	tests := []struct {
		desc       string
		saver      recordSaver
		record     interface{}
		recordType string
	}{
		{
			"Should be able to save valid record",
			newLessonRecordStore(testDB.DB),
			LessonRecord{
				LessonID:  testLesson.LessonID,
				ChapterID: testChapter.ChapterID,
				UID:       "2",
			},
			"LessonRecord",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			testRecordSaver(t, tc.saver, tc.record)
			err := cleanUpTest(testDB.DB, cleanUpStr, tc.record)
			if err != nil {
				t.Error(cleanUpError{FailedVal: tc.recordType}.Error(), tc.record)
			}
		})
	}
}

func TestRecordUpdaters(t *testing.T) {
	testDB := newTestDB(conn)
	defer cleanUpDB(testDB, cleanUpData)

	tests := []struct {
		desc       string
		updater    recordUpdater
		record     interface{}
		recordType string
	}{
		{
			"Should be able to update valid record",
			newLessonRecordStore(testDB.DB),
			LessonRecord{
				LessonID:  testLesson.LessonID,
				ChapterID: testChapter.ChapterID,
				UID:       testUser1.UID,
				WPM:       "100",
				Accuracy:  "100",
			},
			"LessonRecord",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			testRecordUpdater(t, tc.updater, tc.record)
		})
	}
}

func TestRecordQueryers(t *testing.T) {
	db := newTestDB(conn)
	defer cleanUpDB(db, cleanUpData)

	tests := []struct {
		desc        string
		queryer     recordQueryer
		record      interface{}
		queryParams interface{}
		stmt        string
		testFunc    func(t *testing.T, queryer recordQueryer, stmt string, record, queryParams interface{})
	}{
		{
			"Should be be able to query valid lesson record",
			newLessonRecordStore(db.DB),
			testReadOnlyLessonRecord,
			QueryLessonRecordParams{UID: "1", ChapterID: "1", LessonID: "200"},
			QUERY_LESSON_RECORD,
			testLessonRecordQueryer,
		},
		{
			"Should be able to query valid chapter record",
			newChapterRecordStoreI(db.DB),
			testReadOnlyChapterRecord,
			QueryChapterRecordParams{UID: "1", ChapterID: "1"},
			QUERY_CHAPTER_RECORD,
			testChapterRecordQueryer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tc.testFunc(t, tc.queryer, tc.stmt, tc.record, tc.queryParams)
		})
	}
}

func TestRecordsQueryers(t *testing.T) {
	db := newTestDB(conn)
	defer cleanUpDB(db, cleanUpData)

	tests := []struct {
		desc       string
		queryer    recordsQueryer
		records    interface{}
		queryParam interface{}
		stmt       string
		testFunc   func(t *testing.T, queryer recordsQueryer, stmt string, records, params interface{})
	}{
		{
			"Should be able to query all valid lesson records",
			newLessonRecordStore(db.DB),
			[]LessonRecord{
				testReadOnlyLessonRecord,
			},
			QueryLessonsRecordParams{UID: testReadOnlyLessonRecord.UID},
			QUERY_ALL_LESSON_RECORDS,
			testLessonRecordsQueryer,
		},
		{
			"Should be able to query all valid chapter records",
			newChapterRecordStoreI(db.DB),
			[]ChapterRecord{
				testReadOnlyChapterRecord,
			},
			QueryChaptersRecordParams{UID: testReadOnlyChapterRecord.UID},
			QUERY_ALL_CHAPTER_RECORDS,
			testChapterRecordsQueyer,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			tc.testFunc(t, tc.queryer, tc.stmt, tc.records, tc.queryParam)
		})
	}
}

func testLessonRecordsQueryer(t *testing.T, queryer recordsQueryer, stmt string, records, params interface{}) {
	dest := []LessonRecord{}

	err := queryer.QueryAll(stmt, &dest, params)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(dest, records.([]LessonRecord)) {
		t.Errorf("Records do not match got: %v, expected %v", dest, records)
	}
}

func testChapterRecordsQueyer(t *testing.T, queryer recordsQueryer, stmt string, records, params interface{}) {
	dest := []ChapterRecord{}

	err := queryer.QueryAll(stmt, &dest, params)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(dest, records.([]ChapterRecord)) {
		t.Errorf("Records do not match, got %v, expected %v", dest, records)
	}
}

func testLessonRecordQueryer(t *testing.T, queryer recordQueryer, stmt string, record, queryParams interface{}) {
	dest := LessonRecord{}

	err := queryer.Query(stmt, &dest, queryParams)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(dest, record) {
		t.Errorf("Records do not match got: %v, expected %v", dest, record)
	}
}

func testChapterRecordQueryer(t *testing.T, queryer recordQueryer, stmt string, record, queryParams interface{}) {
	dest := ChapterRecord{}

	err := queryer.Query(stmt, &dest, queryParams)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(dest, record) {
		t.Errorf("Records do not match, got: %v, expected %v", dest, record)
	}
}

func testRecordUpdater(t *testing.T, updater recordUpdater, record interface{}) {
	err := updater.Update(record)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func testRecordSaver(t *testing.T, saver recordSaver, record interface{}) {
	err := saver.Save(record)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestFindLastCompletedRecord(t *testing.T) {
	cases := []struct {
		name     string
		records  []LessonRecord
		expected LessonRecord
	}{
		{
			"Should find last completed LessonRecord",
			[]LessonRecord{
				{
					LessonID: "1",
				},
				{
					LessonID: "2",
				},
				{
					LessonID: "5",
				},
				{
					LessonID: "4",
				},
			},
			LessonRecord{
				LessonID: "5",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lastCompletedRecord := findLastCompletedLessonRecord(tc.records)
			if !reflect.DeepEqual(lastCompletedRecord, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, lastCompletedRecord)
			}
		})
	}
}

func TestQueryAvgTutorialStats(t *testing.T) {
	db := newTestDB(conn)
	defer cleanUpDB(db, cleanUpData)

	cases := []struct {
		name        string
		uid         string
		queyer      recordsQueryer
		expected    avgTutorialStats
		expectedErr bool
	}{
		{
			"Should be able to get correct avg stats for user",
			"1",
			newLessonRecordStore(db.DB),
			avgTutorialStats{
				WPM:      "100",
				Accuracy: "100",
				UID:      "1",
			},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			avgStats, err := QueryAvgTutorialStats(tc.queyer, tc.uid)
			if err != nil {
				if !tc.expectedErr {
					t.Error(err)
				}
			} else {
				fmt.Println(avgStats)
				if !reflect.DeepEqual(*avgStats, tc.expected) {
					t.Errorf("Expected %v, got %v", tc.expected, avgStats)
				}
			}
		})
	}
}

//func TestCurrentLesson(t *testing.T) {
//	db := newTestDB(conn)
//	defer cleanUpDB(db, cleanUpData)
//
//	addRecords := func(db *sqlx.DB, queryStmt string, records ...interface{}) {
//		for _, record := range records {
//			_, err := db.NamedExec(queryStmt, record)
//			if err != nil {
//				t.Error(err)
//			}
//		}
//	}
//
//	removeRecords := func(db *sqlx.DB, queryStmt string, records ...interface{}) {
//		for _, record := range records {
//			_, err := db.NamedExec(queryStmt, record)
//			if err != nil {
//				t.Error(err)
//			}
//		}
//	}
//
//	cases := []struct {
//		name string
//		uid string
//		lessonManager *content.LessonManager
//		queyer recordsQueryer
//		addRecordFunc func(db *sqlx.DB, queryStmt string, realRecords ...interface{})
//		removeRecordFunc func(db *sqlx.DB, queryStmt string, records ...interface{})
//
//		expected content.Lesson
//		expectedErr bool
//	} {
//		{
//			"Should be able to get current lesson for user",
//			"1",
//			content.NewLessonManager(db.DB),
//			newLessonRecordStore(db.DB),
//			addRecords,
//			removeRecords,
//			content.Lesson{
//				LessonID: "3",
//			},
//			false,
//		},
//	}
//
//	for _, tc := range cases {
//		tc.addRecordFunc(db.DB,"INSERT INTO lessonscompleted(LessonID, UID) VALUES(:lessonid, :uid)", LessonRecord{UID: "1", LessonID:"1"}, LessonRecord{UID:"1", LessonID:"2"})
//		defer tc.removeRecordFunc(db.DB,"INSERT INTO lessonscompleted(LessonID, UID) VALUES(:lessonid, :uid)", LessonRecord{UID: "1", LessonID:"1"}, LessonRecord{UID: "1", LessonID:"2"})
//
//		t.Run(tc.name, func(t *testing.T) {
//			lesson, err := CurrentLesson(tc.queyer, tc.lessonManager, "1")
//			if err != nil {
//				if !tc.expectedErr {
//					t.Error(err)
//				}
//			} else {
//				if lesson.LessonID != tc.expected.LessonID {
//					t.Errorf("Expected %v, got %v", tc.expected.LessonID, lesson.LessonID)
//				}
//			}
//		})
//	}
//}
