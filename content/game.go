package content

import (
	"log"
	"os"
	"github.com/jmoiron/sqlx"
)

const gameManagerLoggerName = "GameManager"

type Gametext struct{ 
	Id			  string 	`db:"id" 		json:"gametextId"`
	Txt			  string	`db:"txt" 		json:"gameContnent"`
	Gametype      string 	`db:"gametype"  json:"gameName"`
}

type DefaultGameManager struct {
	store  *gameContentStore
	logger *log.Logger   
}

func NewDefaultGameContentManager(db *sqlx.DB) *DefaultGameManager{
	return &DefaultGameManager{
		store:newGameStore(db),
		logger: log.New(os.Stdout, gameManagerLoggerName, log.Lshortfile),
	}
}

type gameContentStore struct{
	db *sqlx.DB
}

func newGameStore(db *sqlx.DB) *gameContentStore{
	return &gameContentStore{db:db}
}

func (manager *DefaultGameManager) ReturnBoatText() ([]*Gametext, error) {
	return manager.store.QueryBoat()
}

func (manager *DefaultGameManager) ReturnCocoText() ([]*Gametext, error) {
	return manager.store.QueryCoco()
}

func (store *gameContentStore) QueryBoat() ([]*Gametext, error) {
	var g []*Gametext
	rows, err := store.db.Queryx("SELECT * FROM gametext where gametype='boatrace'")
	defer rows.Close()

	for rows.Next() {
		var c Gametext
		err = rows.StructScan(&c)
		if err != nil {
			return nil,err
		}
		g = append(g,&c)
	}
	return g, nil
}

func (store *gameContentStore) QueryCoco() ([]*Gametext, error) {
	var g []*Gametext
	rows, err := store.db.Queryx("SELECT * FROM gametext WHERE gametype='coco'")
	defer rows.Close()

	for rows.Next() {
		var c Gametext
		err = rows.StructScan(&c)
		if err != nil {
			return nil,err
		}
		g = append(g,&c)
	}

	return g, nil
}