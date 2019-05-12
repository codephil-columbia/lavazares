package content

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
)

type Gametext struct{ 
	Gamecontent      string    `db:"gamecontent" json:"gameContent"`
}

type DefaultGameManager struct {
	store  *gameContentStore
	logger *log.Logger   
}

func NewDefaultGameContentManager(db *sqlx.DB) *DefaultGameContentManager{
	return &DefaultGameContentManager{
		store:newGameStore(db),
		logger: log.New(os.Stdout, DefaultGameContentManager, log.Lshortfile),
	}
}

type gameContentStore struct{
	db *sqlx.DB
}

func newGameStore(db *sqlx.DB) *gameStore{
	return &gameStore{db:db}
}

func (manager *DefaultGameManager) ReturnBoatText() {
	return manager.store.Query()
}

func (store *gameContentStore) Query() (*Gametext, error) {
	var g Gametext
	err := store.db.QueryRowx().StructScan(&g)
	if err != nil {
		return nil,err
	}
	return &g, nil
}