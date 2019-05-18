package records

type (
	statQueryer interface {
		recordQueryer
	}

	statsQueryer interface {
		recordsQueryer
	}
)

type LessonStats struct {
	WPM      string `db:"wpm" json:"wpm"`
	Accuracy string `db:"accuracy" json:"accuracy"`
	LessonID string `db:"lessonid" json:"id"`
	UID      string `db:"uid" json:"uid"`
}

type avgTutorialStats struct {
	WPM      string `json:"wpm"`
	Accuracy string `json:"accuracy"`
	UID      string `json:"uid"`
}

type chapterProgress struct {
	PercentComplete string `db:"percentComplete" json:"percentComplete"`
}
