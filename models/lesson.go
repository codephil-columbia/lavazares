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

func UserCompletedLesson(lessonComplete LessonComplete) error {
	_, err := db.NamedQuery(
		`INSERT INTO LessonsCompleted(LessonID, UID)
		VALUES(:LessonID, :UID)`,
		lessonComplete)
	return err
}
