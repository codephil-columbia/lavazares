package records

import (
	"lavazares/auth"
	"lavazares/content"
	"lavazares/utils"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"
)

var (
	testDB                 *sqlx.DB
	testLessonRecordStore  *lessonRecordStore
	testChapterRecordStore *chapterRecordStore
	validLessonRecord      = LessonRecord{
		LessonID:  "1",
		ChapterID: "2",
		UID:       "123",
		WPM:       "100",
		Accuracy:  "100",
	}
	validChapterRecord = ChapterRecord{
		ChapterID: "2",
		UID:       "123",
	}
)

const (
	conn = `port=5432 
		host=localhost 
		sslmode=disable 
		user=postgres 
		dbname=postgres 
		password=postgres`
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = sqlx.Open("postgres", conn)
	if err != nil {
		log.Fatalln(err)
	}

	testLessonRecordStore = newLessonRecordStore(testDB)
	testChapterRecordStore = newChapterRecordStore(testDB)

	os.Exit(m.Run())
}

func TestNewChapterRecord(t *testing.T) {
	cases := []struct {
		name           string
		json           utils.RequestJSON
		expectedRecord ChapterRecord
		expectedErr    bool
	}{
		{
			"Should create valid record successfully",
			[]byte(`{
				"chapterID":"2",
				"uid":"123"
			}`),
			validChapterRecord,
			false,
		}, {
			"Should throw error when missing required field",
			[]byte(`{
				"chapterID":"2"
			}`),
			ChapterRecord{},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record, err := NewChapterRecord(tc.json)
			t.Log(err, "errr")
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected err: [%v]", err)
				}
			} else {
				if !reflect.DeepEqual(*record, tc.expectedRecord) {
					t.Errorf("Expected [%v], got [%v]", tc.expectedRecord, record)
				}
			}
		})
	}
}

func TestChapterRecordStoreSave(t *testing.T) {
	cases := []struct {
		name        string
		record      ChapterRecord
		expectedErr bool
	}{
		{
			"Should be able to insert valid record",
			validChapterRecord,
			false,
		},
		{
			"Should throw err when record is in valid",
			ChapterRecord{
				// ChapterID not valid
				ChapterID: "456",
				UID:       "123",
			},
			true,
		},
	}

	addTestChapter()
	addTestUser()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := testChapterRecordStore.save(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected err: [%v]", err)
				}
			} else {
				removeTestChapterRecord(tc.record.ChapterID)
			}
		})
	}
	removeTestUser()
	removeTestChapter()
}

func TestChapterRecordStoreExists(t *testing.T) {
	cases := []struct {
		name        string
		record      ChapterRecord
		expectedErr bool
		exists      bool
	}{
		{
			"Should find existing record",
			validChapterRecord,
			false,
			true,
		}, {
			"Should throw err when record is invalid",
			ChapterRecord{
				UID: "123",
			},
			true,
			false,
		}, {
			"Should return false when record does not exists",
			ChapterRecord{
				UID:       "123",
				ChapterID: "1",
			},
			false,
			false,
		},
	}

	addTestChapter()
	addTestUser()
	addTestChapterRecord()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := testChapterRecordStore.exists(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected err: [%v]", err)
				}
			} else {
				if exists != tc.exists {
					t.Errorf("Exists returned [%v], expected [%v]", exists, tc.exists)
				}
			}
		})
	}
	removeTestChapterRecord(validChapterRecord.ChapterID)
	removeTestUser()
	removeTestChapter()
}

func TestChapterRecordStoreQueryAll(t *testing.T) {
	cases := []struct {
		name        string
		uid         string
		expected    []*ChapterRecord
		expectedErr bool
	}{
		{
			"Should be able to query all records for valid user",
			"123",
			[]*ChapterRecord{
				&validChapterRecord,
			},
			false,
		},
		{
			"Should return error for nonexistent user",
			"321",
			nil,
			true,
		},
	}

	addTestChapter()
	addTestUser()
	addTestChapterRecord()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := testChapterRecordStore.queryAll(tc.uid)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected err: [%v]", err)
				}
			} else {
				if !reflect.DeepEqual(records, tc.expected) {
					t.Errorf("Expected records [%v], got [%v]", tc.expected, records)
				}
			}
		})
	}
	removeTestChapterRecord(validChapterRecord.ChapterID)
	removeTestUser()
	removeTestChapter()
}

func TestChapterRecordStoreQuery(t *testing.T) {
	cases := []struct {
		name        string
		id          string
		uid         string
		expected    ChapterRecord
		expectedErr bool
	}{
		{
			"Should be able to query existing record",
			"2",
			"123",
			validChapterRecord,
			false,
		},
		{
			"Should return error for nonexistant id",
			"9",
			"321",
			ChapterRecord{},
			true,
		},
	}

	addTestChapter()
	addTestUser()
	addTestChapterRecord()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record, err := testChapterRecordStore.query(tc.id, tc.uid)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected error [%v]", err)
				}
			} else {
				if !reflect.DeepEqual(*record, tc.expected) {
					t.Errorf("Expected record [%v], got [%v]", tc.expected, record)
				}
			}
		})
	}
	removeTestChapterRecord(validChapterRecord.ChapterID)
	removeTestUser()
	removeTestChapter()
}

func TestNewLessonRecord(t *testing.T) {
	cases := []struct {
		name        string
		json        utils.RequestJSON
		expected    LessonRecord
		expectedErr bool
	}{
		{
			"Can create record with valid fields",
			[]byte(`{
				"lessonID":"1",
				"chapterID":"2",
				"uid":"123",
				"wpm":"100",
				"accuracy":"100"
			}`),
			validLessonRecord,
			false,
		}, {
			"Should reject faulty record due to missing fields",
			[]byte(`{
				"lessonID":"1"
			}`),
			LessonRecord{},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record, err := NewLessonRecord(tc.json)
			t.Log(err, "ERRR")
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Unexpected error [%v]", err)
				}
			} else {
				if !reflect.DeepEqual(*record, tc.expected) {
					t.Errorf("Expected: [%v], got [%v]", tc.expected, record)
				}
			}
		})
	}
}

func TestLessonRecordStoreSave(t *testing.T) {
	cases := []struct {
		name        string
		record      LessonRecord
		expectedErr bool
	}{
		{
			"Should be able to insert valid record",
			validLessonRecord,
			false,
		}, {
			"Should reject record with missing fields",
			LessonRecord{LessonID: "1"},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			addTestLesson()
			err := testLessonRecordStore.save(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("expected nil error, got [%v]", err)
				}
			} else {
				removeTestLessonRecord(tc.record.LessonID)
			}
			removeTestLesson()
		})
	}
}

func TestLessonRecordStoreUpdate(t *testing.T) {
	cases := []struct {
		name        string
		record      LessonRecord
		expectedErr bool
	}{
		{
			"Should be able to update existing record",
			validLessonRecord,
			false,
		}, {
			// LessonID does not exist
			"Should not accept record that does not exist",
			LessonRecord{
				LessonID:  "2",
				ChapterID: "2",
				WPM:       "100",
				Accuracy:  "100",
				UID:       "123",
			},
			true,
		},
	}

	addTestLesson()
	addTestLessonRecord()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := testLessonRecordStore.update(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Expected no error, got [%v]", err)
				}
			}
		})
		removeTestLessonRecord(validLessonRecord.LessonID)
		removeTestLesson()
	}
}

func TestLessonRecordStoreExists(t *testing.T) {
	cases := []struct {
		name        string
		record      LessonRecord
		expectedErr bool
		exists      bool
	}{
		{
			"Should find existing record",
			validLessonRecord,
			false,
			true,
		}, {
			"Should return false if record does not exist",
			LessonRecord{
				LessonID: "1234",
				UID:      "1234",
			},
			true,
			false,
		}, {
			"Should return err due to missing fields in record",
			LessonRecord{UID: "1234"},
			true,
			false,
		},
	}

	addTestLesson()
	addTestLessonRecord()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			exists, err := testLessonRecordStore.exists(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Did not expect err, got [%v]", err)
				}
			} else {
				if tc.exists != exists {
					t.Errorf("Exists returned [%v], expected [%v]", exists, tc.exists)
				}
			}
		})
	}
	removeTestLessonRecord(validLessonRecord.LessonID)
	removeTestLesson()
}

func addTestChapter() {
	c := content.Chapter{ChapterID: "2"}
	_, err := testDB.Exec(
		`INSERT INTO Chapters(ChapterID)
		VALUES($1)`,
		c.ChapterID,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func addTestUser() {
	u := auth.User{
		UID:        "123",
		Username:   "neil wang wang",
		Password:   "password",
		FirstName:  "Wang",
		LastName:   "Chen",
		Email:      "someemail@email.com",
		Occupation: "student",
	}
	_, err := testDB.Exec(
		`INSERT INTO Users(UID, Username, Password, FirstName, LastName, Email, Occupation) 
		VALUES($1, $2, $3, $4, $5, $6, $7)`,
		u.UID,
		u.Username,
		u.Password,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Occupation,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func addTestLessonRecord() {
	err := testLessonRecordStore.save(&validLessonRecord)
	if err != nil {
		log.Fatalln(err)
	}
}

func addTestChapterRecord() {
	err := testChapterRecordStore.save(&validChapterRecord)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeTestLessonRecord(id string) {
	_, err := testDB.Exec("DELETE FROM LessonsCompleted WHERE LessonID=$1", id)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeTestChapterRecord(id string) {
	_, err := testDB.Exec("DELETE FROM ChaptersCompleted WHERE ChapterID=$1", id)
	if err != nil {
		log.Fatalln(err)
	}
}

func addTestLesson() {
	l := content.Lesson{LessonID: "1", ChapterID: "2"}

	addTestChapter()
	addTestUser()

	_, err := testDB.Exec(
		`INSERT INTO Lessons(LessonID, ChapterID)
		VALUES($1, $2)`,
		l.LessonID,
		l.ChapterID,
	)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeTestLesson() {
	_, err := testDB.Exec(`DELETE FROM Lessons WHERE lessonid=$1`, "1")
	if err != nil {
		log.Fatalln(err)
	}
	removeTestUser()
	removeTestChapter()
}

func removeTestUser() {
	_, err := testDB.Exec("DELETE FROM Users WHERE uid=$1", "123")
	if err != nil {
		log.Fatalln(err)
	}
}

func removeTestChapter() {
	_, err := testDB.Exec(`DELETE FROM Chapters WHERE chapterid=$1`, "2")
	if err != nil {
		log.Fatalln(err)
	}
}
