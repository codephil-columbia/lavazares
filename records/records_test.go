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

func TestNewLessonRecord(t *testing.T) {
	cases := []struct {
		name        string
		args        utils.RequestJSON
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
				"lessonID:"1"
			}`),
			LessonRecord{},
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			record, err := NewLessonRecord(tc.args)
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
			addTestData()
			err := testLessonRecordStore.save(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("expected nil error, got [%v]", err)
				}
			} else {
				removeTestRecord(tc.record.LessonID)
			}
			removeTestData()
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

	addTestData()
	addTestRecord()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := testLessonRecordStore.update(&tc.record)
			if err != nil {
				if !tc.expectedErr {
					t.Errorf("Expected no error, got [%v]", err)
				}
			}
		})
		removeTestRecord(validLessonRecord.LessonID)
		removeTestData()
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

	addTestData()
	addTestRecord()
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
	removeTestRecord(validLessonRecord.LessonID)
	removeTestData()
}

func addTestRecord() {
	err := testLessonRecordStore.save(&validLessonRecord)
	if err != nil {
		log.Fatalln(err)
	}
}

func removeTestRecord(id string) {
	_, err := testDB.Exec("DELETE FROM LessonsCompleted WHERE LessonID=$1", id)
	if err != nil {
		log.Fatalln(err)
	}
}

func addTestData() {
	l := content.Lesson{LessonID: "1", ChapterID: "2"}
	c := content.Chapter{ChapterID: "2"}
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
		`INSERT INTO Chapters(ChapterID)
		VALUES($1)`,
		c.ChapterID,
	)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = testDB.Exec(
		`INSERT INTO Lessons(LessonID, ChapterID)
		VALUES($1, $2)`,
		l.LessonID,
		l.ChapterID,
	)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = testDB.Exec(
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

func removeTestData() {
	_, err := testDB.Exec(`DELETE FROM Lessons WHERE lessonid=$1`, "1")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = testDB.Exec(`DELETE FROM Chapters WHERE chapterid=$1`, "2")
	if err != nil {
		log.Fatalln(err)
	}
	_, err = testDB.Exec("DELETE FROM Users WHERE uid=$1", "123")
	if err != nil {
		log.Fatalln(err)
	}
}
