package models

import (
	"encoding/json"
	"time"

	"github.com/rs/xid"
)

type Lesson struct {
	LessonID           string     `json:"-" db:"LessonID"`
	CreatedAt          time.Time  `json:"-" db:"CreatedAt"`
	UpdatedAt          time.Time  `json:"-" db:"UpdatedAt"`
	DeletedAt          *time.Time `json:"-" db:"DeletedAt"`
	LessonName         string     `json:"lessonName" db:"LessonName"`
	LessonContent      string     `json:"lessonContent" db:"LessonContent"`
	MinimunScoreToPass int        `json:"minScore" db:"MinimumScoreToPass"`
	ChapterID          string     `json:"-" db:"ChapterID"`
}

type LessonComplete struct {
	LessonID string `json:"lessonID" db:"LessonID"`
	UID      string `json:"uid" db:"UID"`
}

type Chapter struct {
	ChapterID          string `json:"-" db:"ChapterID"`
	ChapterName        string `json:"chapterName" db:"ChapterName"`
	ChapterDescription string `json:"ChapterDescription" db:"ChapterDescription"`
}

func NewLesson(lessonRequest []byte) (*Lesson, error) {
	lesson := Lesson{}
	err := json.Unmarshal(lessonRequest, &lesson)
	if err != nil {
		return nil, err
	}

	lesson.LessonID = xid.New().String()
	lesson.CreatedAt = time.Now()
	lesson.UpdatedAt = time.Now()
	lesson.DeletedAt = nil
	lesson.ChapterID = "123"

	_, err = db.NamedQuery(
		`INSERT INTO Lessons 
		(LessonID, CreatedAt, UpdatedAt, DeletedAt, LessonName, LessonContent, MinimumScoreToPass, ChapterID)
		VALUES(:LessonID, :CreatedAt, :UpdatedAt, :DeletedAt, :LessonName, :LessonContent, :MinimumScoreToPass, :ChapterID)`,
		lesson)
	if err != nil {
		return nil, err
	}

	return &lesson, nil
}

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

func UserCompletedLesson(lessonComplete LessonComplete) error {
	_, err := db.NamedQuery(
		`INSERT INTO LessonsCompleted(LessonID, UID)
		VALUES(:LessonID, :UID)`,
		lessonComplete)
	return err
}
