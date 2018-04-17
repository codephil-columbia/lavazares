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
	LessonText         []string   `json:"lessonContent" db:"lessontext"`
	MinimumScoreToPass [][]int    `json:"minScore" db:"minimumscoretopass"`
	ChapterID          string     `json:"chapterID" db:"chapterid"`
	Image              [][]string `json:"imageID" db:"image"`
	LessonDescriptions [][]string `json:"lessonDescriptions" db:"lessondescriptions"`
}

// Chapter metadata. Maps directly to SQL definition in DB.
type Chapter struct {
	ChapterID          string `json:"-" db:"ChapterID"`
	ChapterName        string `json:"chapterName" db:"ChapterName"`
	ChapterDescription string `json:"ChapterDescription" db:"ChapterDescription"`
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
		lesson.LessonID, lesson.LessonName, pq.Array(lesson.LessonText),
		lesson.MinimumScoreToPass, lesson.ChapterID)

	if err != nil {
		return nil, err
	}

	return &lesson, nil
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

	uncompletedLessons, err := GetUncompletedLessons(forUser)
	if err != nil {
		return nil, err
	}

	bulkRespForUser := make(map[string]interface{})
	bulkRespForUser["lessons"] = lessons
	bulkRespForUser["chapters"] = chapters
	bulkRespForUser["uncompletedLessons"] = uncompletedLessons

	fmt.Println(bulkRespForUser)

	return &bulkRespForUser, err
}
