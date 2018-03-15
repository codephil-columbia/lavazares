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
	LessonID           string     `json:"-" db:"LessonID"`
	CreatedAt          time.Time  `json:"-" db:"CreatedAt"`
	UpdatedAt          time.Time  `json:"-" db:"UpdatedAt"`
	DeletedAt          *time.Time `json:"-" db:"DeletedAt"`
	LessonName         string     `json:"lessonName" db:"LessonName"`
	LessonContent      []string   `json:"lessonContent" db:"LessonContent"`
	MinimumScoreToPass int        `json:"minScore" db:"MinimumScoreToPass"`
	ChapterID          string     `json:"chapterID" db:"ChapterID"`
}

// Chapter metadata. Maps directly to SQL definition in DB.
// UnitID is a foreign key reference to the parent Unit that Chapter belongs to.
type Chapter struct {
	ChapterID          string `json:"-" db:"ChapterID"`
	ChapterName        string `json:"chapterName" db:"ChapterName"`
	ChapterDescription string `json:"ChapterDescription" db:"ChapterDescription"`
	UnitID             string `json:"unitID" db:"UnitID"`
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

func GetUncompletedLessons(uid string) (*[]Lesson, error) {
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

	lessons := []Lesson{}

	for rows.Next() {
		lesson := Lesson{}
		err = rows.StructScan(&lesson)
		lessons = append(lessons, lesson)
	}
	fmt.Println(lessons)
	return &lessons, nil
}
