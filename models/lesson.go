package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/rs/xid"
	"github.com/jmoiron/sqlx"
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
	LessonID  string  `db:"lessonid" json:"lessonID"`
	ChapterID string  `db:"chapterid" json:"chapterID"`
	UID       *string `db:"uid"`
	WPM       float64 `db:"wpm"`
	Accuracy  float64 `db:"accuracy"`
}

// ChaptersComplete holds the ChapterID and User ID for all of the chapters
// that a User has completed
type ChaptersComplete struct {
	ChapterID string `json:"chapterid" db:"chapterid"`
	UID       string `json:"uid" db:"uid"`
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

type TutorialLessonResponse struct {
	ChapterDescription string `json:"chapterdescription"`
	ChapterImage string `json:"chapterimage"`
	ChapterName string `json:"chaptername"`
	ChapterID string `json:"chapterid"`
	LessonDescriptions pq.StringArray `json:"lessondescriptions"`
	LessonId string `json:"lessonid"`
	LessonName string `json:"lessonname"`
	LessonText pq.StringArray `json:"lessontext"`
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

// UserCompletedChapter takes UID of User and ChapterID of the chapter that
// they completed and inserts the pair into the DB.
func UserCompletedChapter(chapterComplete ChaptersComplete) error {
	_, err := db.NamedQuery(
		`INSERT INTO ChaptersCompleted(ChapterID, UID)
		VALUES(:ChapterID, :UID)`,
		chapterComplete)
	return err
}

// UserCompletedLesson takes UID of User and LessonID of the Lesson that
// they completed and inserts the pair into DB.
func UserCompletedLesson(lessonComplete LessonsComplete) error {
	if hasCompletedLesson(lessonComplete) {
		_, err := db.Query("UPDATE LessonsCompleted SET wpm=$1, accuracy=$2 WHERE uid=$3 and lessonid=$4", lessonComplete.WPM, lessonComplete.Accuracy, lessonComplete.UID, lessonComplete.LessonID)
		return err
	} else {
		_, err := db.NamedQuery(
			`INSERT INTO LessonsCompleted(LessonID, UID, WPM, Accuracy, ChapterID)
		VALUES(:lessonID, :uid, :wpm, :accuracy, :chapterid)`,
			lessonComplete)
		return err
	}
}

func hasCompletedLesson(lessonComplete LessonsComplete) bool {
	err := db.QueryRow("select count(1) from LessonsCompleted where lessonid=$1 and uid=$2", lessonComplete.LessonID, lessonComplete.UID).Scan()
	return (err == sql.ErrNoRows)
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

func NextLessonForStudent(uid string) (*TutorialLessonResponse, error) {
	s, err := GetStudent(uid)
	if err != nil {
		return nil, err
	}
	var resp TutorialLessonResponse
	if err := db.QueryRowx(`
		select chapterimage, lessonname, chaptername, L.lessontext, L.lessondescriptions, L.lessonid, C.chapterdescription, C.chapterid
		from lessons l, lessonscompleted lc, chapters c 
		where l.lessonid != $1 and lc.lessonid != l.lessonid and l.chapterid = c.chapterid and lc.uid = $2
		order by createdat limit 1`,
		s.CurrentLessonID, s.UID).StructScan(&resp); err != nil {
			return nil, err
	}
	return &resp, err
}

func GetAllChapterNames() (*[]string, error) {
	chapters := []string{}
	err := db.Select(&chapters, "select chaptername from chapters order by chaptername asc")
	if err != nil {
		return nil, err
	}
	return &chapters, nil
}

func GetCurrentLessonForStudent(uid string) (*TutorialLessonResponse, error) {
	s, err := GetStudent(uid)
	if err != nil {
		return nil , err
	}

	var resp TutorialLessonResponse
	if err := db.QueryRowx(`
		select chapterimage, lessonname, chaptername, L.lessontext, L.lessondescriptions, L.lessonid, C.chapterdescription, C.chapterid
		from lessons l, chapters c 
		where l.chapterid = c.chapterid and l.lessonid not in (
			select lessonid
    		from lessonscompleted lc
			where lc.uid = $1
		)
		order by createdat limit 1`, s.UID).StructScan(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
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
	var lcs []LessonsComplete
	rows, err := db.Queryx("select lessonid, wpm, accuracy, uid, chapterid from lessonscompleted where uid = $1", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var lc LessonsComplete
		err := rows.StructScan(&lc)
		if err != nil {
			return nil, err
		}
		lcs = append(lcs, lc)
	}
	return &lcs, nil
}

func GetOverallWPMAndAccuracy(uid string) (*map[string]interface{}, error) {
	stats := make(map[string]interface{})

	err := db.QueryRowx(
		`select AVG(accuracy) as avgAccuracy, AVG(wpm) as avgWPM 
		from lessonscompleted 
		where uid=$1`, uid).MapScan(stats)

	// In the case of a new User
	if stats["avgAccuracy"] == nil || stats["avgWPM"] == nil {
		stats = map[string]interface{}{
			"avgAccuracy": "0",
			"avgWPM": "0",
		}
	}

	return &stats, err
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
		and S.currentchaptername != C.chaptername
		and C.chapterid = LC.chapterid
		and S.uid=$1`, uid).MapScan(progress)

	if err != nil {
		return nil, err
	}

	return &progress, nil
}

func hasCompletedLessonInTx(tx *sqlx.Tx, lessonid, uid string) bool {
	err := tx.QueryRow("select count(1) from LessonsCompleted where lessonid=$1 and uid=$2", lessonid, uid).Scan()
	return err == sql.ErrNoRows
}

func hasCompletedChapter(tx *sqlx.Tx, chapterid, uid string) bool {
	var count string
	err := tx.QueryRow("select count(1) from ChaptersCompleted where chapterid=$1 and uid=$2", chapterid, uid).Scan(&count)
	return err == sql.ErrNoRows
}

 func UserDidFinishLesson(lc LessonsComplete) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// Add lesson to lesson completed list, or update the previous record if has already completed lesson
	if hasCompletedLessonInTx(tx, lc.LessonID, *lc.UID) {
		_, err = tx.Exec(`
			UPDATE LessonsCompleted 
			SET wpm=$1, accuracy=$2 
			WHERE uid=$3 and lessonid=$4`,
			lc.WPM, lc.Accuracy, lc.LessonID)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec(`
			INSERT INTO LessonsCompleted(LessonID, UID, WPM, Accuracy, ChapterID) 
			VALUES($1, $2, $3, $4, $5)`,
			lc.LessonID, lc.UID, lc.WPM, lc.Accuracy, lc.ChapterID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	nextLesson := make(map[string]interface{})
	userDidCompleteChapter := false

	// Check whether user has completed chapter with this lesson, if so, make to update that users current chapter
	// if this query is empty, then there are no more unmatched lessons for the chapter == chapter is completed
	err = tx.QueryRowx(`
		SELECT L.lessonname, L.lessonid, L.chapterid 
		FROM Lessons L 
		WHERE L.chapterid=$1 AND L.lessonid 
		NOT IN (
			SELECT lessonid 
			FROM LessonsCompleted 
			WHERE uid=$2 and lessonid = L.lessonid) 
		ORDER BY createdat LIMIT 1`,
		lc.ChapterID, lc.UID).MapScan(nextLesson)
	if err == sql.ErrNoRows {
		userDidCompleteChapter = true
		if !hasCompletedChapter(tx, lc.ChapterID, *lc.UID) {
			_, err := tx.Exec(`INSERT INTO ChaptersCompleted(chapterid, uid) VALUES($1, $2)`, lc.ChapterID, lc.UID)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if userDidCompleteChapter {
		var chapterName string
		// Find the next chapter that needs to be completed
		nextChapter := make(map[string]interface{})
		err := tx.QueryRowx(`SELECT C.chaptername, C.chapterid FROM Chapters C, ChaptersCompleted CC WHERE C.chapterid != CC.chapterid AND CC.uid = $1 ORDER BY C.chaptername LIMIT 1`, lc.UID).MapScan(nextChapter)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.QueryRow(`SELECT chaptername FROM Chapters WHERE chapterid=$1`, nextChapter["chapterid"]).Scan(&chapterName)
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.QueryRowx(`SELECT L.lessonid FROM Chapters C, Lessons L WHERE C.chapterid = $1 AND L.chapterid = $1 ORDER BY L.createdat LIMIT 1`, nextChapter["chapterid"]).MapScan(nextLesson)
		if err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.Exec(`UPDATE Students SET currentlessonid=$1, currentchapterid=$2, currentchaptername=$3`, nextLesson["lessonid"], nextChapter["chapterid"], chapterName)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		_, err = tx.Exec(`UPDATE Students SET currentlessonid=$1`, nextLesson["lessonid"])
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func GetCurrent(uid string ) (*TutorialLessonResponse, error) {
	var resp TutorialLessonResponse
	err := db.QueryRowx(`SELECT chapterimage, lessonname, chaptername, L.lessontext, L.lessondescriptions, L.lessonid, C.chapterdescription, C.chapterid FROM Lessons L, Students S, Chapters C WHERE S.uid=$1 AND S.currentlessonid = L.lessonid and L.chapterid = C.chapterid LIMIT 1`, uid).StructScan(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, err
}