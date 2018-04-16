package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/rs/xid"
)

// Lesson metadata. Maps directly to SQL definition of Lesson in db.
// LessonContent for now defined as string array, where every element in array
// is an individual row that a person has to type.
// ChapterID is a foreign key reference to it's parent Chapter that it belongs to.
type Lesson struct {
	LessonID           string     `json:"-" db:"lessonid"`
	CreatedAt          time.Time  `json:"-" db:"createdat"`
	UpdatedAt          time.Time  `json:"-" db:"updatedat"`
	DeletedAt          *time.Time `json:"-" db:"deletedat"`
	LessonName         string     `json:"lessonName" db:"lessonname"`
	LessonContent      []string   `json:"lessonContent" db:"lessoncontent"`
	MinimumScoreToPass int        `json:"minScore" db:"minimumscoretopass"`
	ChapterID          string     `json:"chapterID" db:"chapterid"`
}

// Chapter metadata. Maps directly to SQL definition in DB.
// UnitID is a foreign key reference to the parent Unit that Chapter belongs to.
type Chapter struct {
	ChapterID          string `json:"-" db:"ChapterID"`
	ChapterName        string `json:"chapterName" db:"ChapterName"`
	ChapterDescription string `json:"ChapterDescription" db:"ChapterDescription"`
}

//Unit metadata. Maps directly to SQL definition in DB.
type Unit struct {
	UnitName        string     `json:"unitName" db:"UnitName"`
	UnitDescription string     `json:"unitDescription" db:"UnitDescription"`
	UnitID          string     `json:"-" db:"UnitID"`
	CreatedAt       time.Time  `json:"-" db:"CreatedAt"`
	UpdatedAt       time.Time  `json:"-" db:"UpdatedAt"`
	DeletedAt       *time.Time `json:"-" db:"DeletedAt"`
}

// LessonsComplete holds the LessonID and UserID for all of the Lessons that a
// User has completed.
type LessonsComplete struct {
	LessonID string `json:"lessonID" db:"LessonID"`
	UID      string `json:"uid" db:"UID"`
}

// ChapterComplete holds the ChapterID and User ID for all of the chapters
// that a User has completed
type ChapterComplete struct {
	ChapterID string `json:"chapterID" db:"ChapterID"`
	UID       string `json:"uid" db:"UID"`
}

// UnitComplete holds UnitID and User ID for all of the units that
// a User has completed
type UnitComplete struct {
	UnitID string `json:"unitID" db:"UnitID"`
	UID    string `json:"uid" db:"UID"`
}

type NextLessonReq struct {
	UnitID    string `json:"unitID"`
	ChapterID string `json:"chapterID"`
}

// NewLesson creates a new Lesson from a lessonRequest and inserts it into DB.
// We don't have to worry about populating any time fields since Postgres will fill
// with current time if we leave it empty.
func NewLesson(lessonRequest []byte) (*Lesson, error) {
	lesson := Lesson{}
	err := json.Unmarshal(lessonRequest, &lesson)
	if err != nil {
		return nil, err
	}

	lesson.LessonID = xid.New().String()
	lesson.DeletedAt = nil

	// Can't user NamedQuery because sqlx does not support Arrays :(, instead convert to postgres array
	// obj (note: does not support nested Arrays) and insert manually
	_, err = db.Queryx(
		`INSERT INTO Lessons 
		(LessonID, LessonName, LessonContent, MinimumScoreToPass, ChapterID)
		VALUES($1, $2, $3, $4, $5)`,
		lesson.LessonID, lesson.LessonName, pq.Array(lesson.LessonContent),
		lesson.MinimumScoreToPass, lesson.ChapterID)

	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

// NewUnit creates a Unit from a unitReq and inserts it into DB.
func NewUnit(unitReq []byte) (*Unit, error) {
	unit := Unit{}
	err := json.Unmarshal(unitReq, &unit)
	if err != nil {
		return nil, err
	}

	unit.UnitID = xid.New().String()

	_, err = db.NamedQuery(
		`INSERT INTO Units(UnitName, UnitID, UnitDescription)
		VALUES(:UnitName, :UnitID, :UnitDescription)`,
		unit)
	if err != nil {
		return nil, err
	}

	return &unit, err
}

// NewChapter creates a Chapter from a chapterReq and inserts it into DB.
func NewChapter(chapterReq []byte) (*Chapter, error) {
	chapter := Chapter{}
	err := json.Unmarshal(chapterReq, &chapter)
	if err != nil {
		return nil, err
	}

	chapter.ChapterID = xid.New().String()

	_, err = db.NamedQuery(
		`INSERT INTO Chapters(ChapterID, ChapterName, ChapterDescription)
		VALUES(:ChapterID, :ChapterName, :ChapterDescription)`,
		chapter)
	if err != nil {
		return nil, err
	}

	return &chapter, nil
}

// UserCompletedUnit inserts UnitID User ID pair when a user completes a unit
func UserCompletedUnit(unitComplete UnitComplete) error {
	_, err := db.NamedQuery(
		`INSERT INTO UnitsCompleted(UnitID, UID)
		VALUES (:UnitID, :UID)`,
		unitComplete)

	return err
}

// UserCompletedChapter takes UID of User and ChapterID of the chapter that
// they completed and inserts the pair into the DB.
func UserCompletedChapter(chapterComplete ChapterComplete) error {
	_, err := db.NamedQuery(
		`INSERT INTO ChaptersCompleted(ChapterID, UID)
		VALUES(:ChapterID, :UID)`,
		chapterComplete)
	return err
}

// UserCompletedLesson takes UID of User and LessonID of the Lesson that
// they completed and inserts the pair into DB.
func UserCompletedLesson(lessonComplete LessonsComplete) error {
	_, err := db.NamedQuery(
		`INSERT INTO LessonsCompleted(LessonID, UID)
		VALUES(:LessonID, :UID)`,
		lessonComplete)
	return err
}

func GetUncompletedLessons(uid string) (*[]map[string]interface{}, error) {
	rows, err := db.Queryx(
		`select * from lessons L 
		where L.lessonid not in (
			select LC.lessonid 
			from lessonscompleted LC 
			where LC.uid = $1
		)`, uid)

	if err != nil {
		return nil, err
	}

	lessons := []map[string]interface{}{}

	for rows.Next() {
		lesson := make(map[string]interface{})
		err = rows.MapScan(lesson)
		lessons = append(lessons, lesson)
	}
	fmt.Println(lessons)
	return &lessons, nil
}

func GetBulkInfo(forUser string) (*map[string]interface{}, error) {
	lessons := []map[string]interface{}{}
	rows, err := db.Queryx("SELECT * FROM Lessons")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		lesson := make(map[string]interface{})
		rows.MapScan(lesson)
		lessons = append(lessons, lesson)
	}

	chapters := []map[string]interface{}{}
	rows, err = db.Queryx("SELECT * FROM Chapters")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		chapter := make(map[string]interface{})
		rows.MapScan(chapter)
		chapters = append(chapters, chapter)
	}

	units := []map[string]interface{}{}
	rows, err = db.Queryx("SELECT * FROM Units")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		unit := make(map[string]interface{})
		rows.MapScan(unit)
		units = append(units, unit)
	}

	uncompletedLessons, err := GetUncompletedLessons(forUser)
	if err != nil {
		return nil, err
	}

	bulkRespForUser := make(map[string]interface{})
	bulkRespForUser["lessons"] = lessons
	bulkRespForUser["chapters"] = chapters
	bulkRespForUser["units"] = units
	bulkRespForUser["uncompletedLessons"] = uncompletedLessons

	fmt.Println(bulkRespForUser)

	return &bulkRespForUser, err
}
