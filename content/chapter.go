package content

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

// Chapter metadata. Maps directly to SQL definition in DB.
type Chapter struct {
	ChapterID          string `db:"chapterid" json:"chapterID"`
	ChapterName        string `db:"chaptername" json:"chapterName"`
	ChapterDescription string `db:"chapterdescription" json:"chapterDescription"`
	ChapterImage       string `db:"chapterimage" json:"chapterImage"`
}

// DefaultChapterManager provides read only access to Chapter objs
type DefaultChapterManager struct {
	store  *chapterStore
	logger *log.Logger
}

const defaultChapterManagerLoggerName = "DefaultChapterManager"

// NewDefaultChapterManager creates a new ChapterManager
func NewDefaultChapterManager(db *sqlx.DB) *DefaultChapterManager {
	return &DefaultChapterManager{
		store:  newChapterStore(db),
		logger: log.New(os.Stdout, defaultChapterManagerLoggerName, log.Lshortfile),
	}
}

// GetChapter returns a Chapter by id
func (manager *DefaultChapterManager) GetChapter(id string) (*Chapter, error) {
	return manager.store.Query(id)
}

// GetChapters returns a slice to all Chapters
func (manager *DefaultChapterManager) GetChapters() ([]*Chapter, error) {
	return manager.store.QueryAll()
}

type chapterStore struct {
	db *sqlx.DB
}

func newChapterStore(db *sqlx.DB) *chapterStore {
	return &chapterStore{db: db}
}

func (store *chapterStore) Query(id string) (*Chapter, error) {
	var c Chapter
	err := store.db.QueryRowx("SELECT * FROM Chapter WHERE id = $1", id).StructScan(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (store *chapterStore) QueryAll() ([]*Chapter, error) {
	var all []*Chapter
	rows, err := store.db.Queryx("SELECT * FROM Chapter")
	defer rows.Close()

	for rows.Next() {
		var c Chapter
		err = rows.StructScan(&c)
		if err != nil {
			return nil, err
		}
		all = append(all, &c)
	}

	return all, nil
}
