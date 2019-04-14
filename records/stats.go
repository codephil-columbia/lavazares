package records

import (
	"encoding/json"
)

type statsManager interface {
	recordStat(recordID string, uid string) (json.Marshaler, error)
	hollisticStats(uid string)
}

type LessonStats struct {
	WPM      string `db:"wpm" json:"wpm"`
	Accuracy string `db:"accuracy" json:"accuracy"`
	ID       string `db:"lessonid" json:"id"`
	UID      string `db:"uid" json:"uid"`
}
