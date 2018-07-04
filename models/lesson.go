package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/rs/xid"
)

// Lesson metadata. Maps directly to SQL definition of Lesson in db.
// LessonContent for now defined as string array, where every element in array
// is an individual row that a person has to type.
// ChapterID is a foreign key reference to it's parent Chapter that it belongs to.
type Lesson struct {
	LessonID           string         `db:"lessonid"`
	CreatedAt          time.Time      `db:"createdat"`
	UpdatedAt          time.Time      `db:"updatedat"`
	DeletedAt          *time.Time     `db:"deletedat"`
	LessonName         string         `db:"lessonname"`
	LessonText         pq.StringArray `db:"lessontext"`
	LessonDescriptions pq.StringArray `db:"lessondescriptions"`
	MinimumScoreToPass pq.Int64Array  `db:"minimumscoretopass"`
	ChapterID          string         `db:"chapterid"`
	Image              string         `db:"image"`
}

// Chapter metadata. Maps directly to SQL definition in DB.
type Chapter struct {
	ChapterID          string `db:"chapterid"`
	ChapterName        string `db:"chaptername"`
	ChapterDescription string `db:"chapterdescription"`
	ChapterImage       string `db:"chapterimage"`
}

// LessonsComplete holds the LessonID and UserID for all of the Lessons that a
// User has completed.
type LessonsComplete struct {
	LessonID  string  `db:"lessonid"`
	ChapterID string  `db:"chapterid"`
	UID       *string `db:"uid"`
	WPM       float64 `db:"wpm"`
	Accuracy  float64 `db:"accuracy"`
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
// func NewLesson(lessonRequest []byte) (*Lesson, error) {
// 	lesson := Lesson{}
// 	err := json.Unmarshal(lessonRequest, &lesson)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lesson.LessonID = xid.New().String()
// 	lesson.DeletedAt = nil

// 	// Can't user NamedQuery because sqlx does not support Arrays :(, instead convert to postgres array
// 	// obj (note: does not support nested Arrays) and insert manually
// 	_, err = db.Queryx(
// 		`INSERT INTO Lessons
// 		(LessonID, LessonName, LessonContent, MinimumScoreToPass, ChapterID)
// 		VALUES($1, $2, $3, $4, $5)`,
// 		lesson.LessonID, lesson.LessonName, pq.Array(lesson.LessonContent),
// 		lesson.MinimumScoreToPass, lesson.ChapterID)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &lesson, nil
// }

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
	fmt.Println(hasCompletedLesson(lessonComplete))
	if hasCompletedLesson(lessonComplete) {
		_, err := db.Query("UPDATE LessonsCompleted SET wpm=$1, accuracy=$2 WHERE uid=$3 and lessonid=$4", lessonComplete.WPM, lessonComplete.Accuracy, lessonComplete.UID, lessonComplete.LessonID)
		return err
	}
	_, err := db.NamedQuery(
		`INSERT INTO LessonsCompleted(LessonID, UID, WPM, Accuracy, ChapterID)
		VALUES(:lessonid, :uid, :wpm, :accuracy, :chapterid)`,
		lessonComplete)
	return err
}

func hasCompletedLesson(lessonComplete LessonsComplete) bool {
	err := db.QueryRow("select count(1) from LessonsCompleted where lessonid=$1 and uid=$2", lessonComplete.LessonID, lessonComplete.UID).Scan()
	return (err != sql.ErrNoRows)
}

func AllLessons() (*map[string]interface{}, error) {
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

	bulkRespForUser := make(map[string]interface{})
	bulkRespForUser["lessons"] = lessons
	bulkRespForUser["chapters"] = chapters

	fmt.Println(bulkRespForUser)

	return &bulkRespForUser, err
}

func GetLesson(lessonID string) (*Lesson, error) {
	lesson := Lesson{}
	rows, err := db.Queryx("SELECT * FROM Lessons where lessonid=$1", lessonID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		rows.StructScan(&lesson)
	}

	fmt.Println(lesson)

	return &lesson, nil
}

func NextLessonForStudent(uid string) (map[string]interface{}, error) {
	lessonInfo := make(map[string]interface{})
	s, err := GetStudent(uid)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowx("select chapterimage, lessonname, chaptername, L.lessontext, L.lessondescriptions, L.lessonid, C.chapterid from lessons L, chapters C where C.chaptername = $1 and C.chapterid = L.chapterid and L.lessonid = $2",
		s.CurrentChapterName, s.CurrentLessonID).MapScan(lessonInfo)
	if err != nil {
		return nil, err
	}

	return lessonInfo, nil
}

func GetAllChapterNames() (*[]string, error) {
	chapters := []string{}
	err := db.Select(&chapters, "select chaptername from chapters order by chaptername asc")
	if err != nil {
		return nil, err
	}
	return &chapters, nil
}

func GetAllLessonsChapters() (*[]map[string]interface{}, error) {
	allInfo := []map[string]interface{}{}

	allChapters := []Chapter{}
	err := db.Select(&allChapters, "select chaptername, chapterid from chapters order by chaptername asc")
	if err != nil {
		return nil, err
	}

	for _, c := range allChapters {
		lessons := []Lesson{}
		err = db.Select(&lessons,
			`select L.lessonname, L.lessonid, L.chapterid from lessons L where L.chapterid = $1 order by L.lessonname asc
			`, c.ChapterID)

		if err != nil {
			return nil, err
		}

		// if chapter contains chapter test then
		// since order by asc puts the chapter test before lesson,
		// swap it so that its at the end, shift over everything
		if strings.Contains(lessons[0].LessonName, "Chapter") {
			l := lessons[0]
			lessons = append(lessons, l)
			lessons = append(lessons[:0], lessons[1:]...)
		}

		chapterLessons := make(map[string]interface{})
		chapterLessons["chapterName"] = c.ChapterName
		chapterLessons["lessons"] = lessons

		allInfo = append(allInfo, chapterLessons)
	}

	return &allInfo, nil
}

func GetCompletedLessonsForUser(uid string) (*[]LessonsComplete, error) {
	lc := []LessonsComplete{}
	err := db.Select(&lc, "select lessonid, wpm, accuracy, uid, chapterid from lessonscompleted where uid = $1", uid)
	if err != nil {
		return nil, err
	}
	fmt.Println(lc)
	return &lc, nil
}

func GetOverallWPMAndAccuracy(uid string) (*map[string]interface{}, error) {
	stats := make(map[string]interface{})

	err := db.QueryRowx(
		`select AVG(accuracy) as avgAccuracy, AVG(wpm) as avgWPM 
		from lessonscompleted 
		where uid=$1`, uid).MapScan(stats)

	if err != nil {
		return nil, err
	}

	fmt.Println(stats)
	return &stats, nil
}

func GetProgressForCurrentUserLesson(uid string) (*map[string]interface{}, error) {
	progress := make(map[string]interface{})

	err := db.QueryRowx(
		`select count(*) as totalLessons
		from chapters C, lessons L, students S 
		where S.currentchaptername = C.chaptername
		and L.chapterid = C.chapterid
		and S.uid=$1`, uid).MapScan(progress)
	if err != nil {
		return nil, err
	}

	err = db.QueryRowx(
		`select Count(*) as compCount
		from students S, lessonscompleted LC, chapters C
		where S.currentchaptername = C.chaptername
		and C.chapterid = LC.chapterid
		and S.uid=$1`, uid).MapScan(progress)
	if err != nil {
		return nil, err
	}

	return &progress, nil
}
