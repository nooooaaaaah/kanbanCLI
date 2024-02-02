package utils

import (
	"database/sql"
	"kanban/board"
	"kanban/logger"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type DbHandler struct {
	db *sql.DB
}

func NewDbHandler(db *sql.DB) *DbHandler {
	return &DbHandler{db: db}
}

type DBService interface {
	CreateDB() (bool, error)
	InsertBoard(board board.Board) error
	InsertCardList(cardList board.CardList, boardID string) error
	InsertCard(card board.Card, cardListID string) error
	GetBoard(boardID string) (board.Board, error)
	GetCardLists(board *board.Board) error
	GetCards(cardList *board.CardList) error
	DeleteBoard(boardID string) error
	DeleteCardList(cardListID string) error
	DeleteCard(cardID string) error
	UpdateBoardTitle(boardID string, title string) error
	UpdateCardList(cardList board.CardList) error
	UpdateCard(card board.Card) error
}

func connectToDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./kanban.db")
	if err != nil {
		logger.Log.Println(err)
		return nil, err
	}
	return db, nil
}

func InitializeDBService() (*DbHandler, error) {
	db, err := connectToDB()
	if err != nil {
		logger.Log.Println(err)
		return nil, err
	}
	return NewDbHandler(db), nil
}

// Genereate UUIDs for the board, cardlist, and cardlist
func GenerateUUID() string {
	return uuid.New().String()
}
