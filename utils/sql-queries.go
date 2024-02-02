package utils

import (
	"kanban/board"
	"kanban/logger"
	"time"
)

func (handler *DbHandler) CreateDB() (bool, error) {
	createBoardTableSQL := `CREATE TABLE IF NOT EXISTS Boards (
        ID TEXT PRIMARY KEY,
        Title TEXT NOT NULL
    );`

	createCardListTableSQL := `CREATE TABLE IF NOT EXISTS CardLists (
        ID TEXT PRIMARY KEY,
        Title TEXT NOT NULL,
        BoardID TEXT,
        FOREIGN KEY (BoardID) REFERENCES Boards(ID)
    );`

	createCardTableSQL := `CREATE TABLE IF NOT EXISTS Cards (
        ID TEXT PRIMARY KEY,
        Title TEXT NOT NULL,
        Description TEXT,
        StartDate TEXT,
        EndDate TEXT,
        DueDate TEXT,
        Duration INTEGER,
        ListID TEXT,
        FOREIGN KEY (ListID) REFERENCES CardLists(ID)
    );`

	for _, stmt := range []string{createBoardTableSQL, createCardListTableSQL, createCardTableSQL} {
		if _, err := handler.db.Exec(stmt); err != nil {
			logger.Log.Println(err)
			return false, err
		}
	}

	logger.Log.Println("Database created")
	return true, nil
}

func (handler *DbHandler) InsertBoard(board board.Board) error {
	insertBoardSQL := `INSERT INTO Boards (ID, Title) VALUES (?, ?);`
	if _, err := handler.db.Exec(insertBoardSQL, board.ID, board.Title); err != nil {
		logger.Log.Println(err)
		return err
	}
	for _, cardList := range board.CardLists {
		if err := handler.InsertCardList(cardList, board.ID); err != nil {
			logger.Log.Println(err)
			return err
		}
	}
	return nil
}

func (handler *DbHandler) InsertCardList(cardList board.CardList, boardID string) error {
	tx, err := handler.db.Begin()
	if err != nil {
		logger.Log.Println(err)
		return err
	}
	stmt := `INSERT INTO CardLists (ID, Title, BoardID) VALUES (?, ?, ?);`
	_, err = tx.Exec(stmt, cardList.ID, cardList.Title, boardID)
	if err != nil {
		tx.Rollback()
		logger.Log.Println(err)
		return err
	}
	for _, card := range cardList.Cards {
		if err := handler.InsertCard(card, cardList.ID); err != nil {
			logger.Log.Println(err)
			return err
		}
	}
	return nil
}

func (handler *DbHandler) InsertCard(card board.Card, listID string) error {
	stmt := `INSERT INTO Cards (ID, Title, Description, StartDate, EndDate, DueDate, Duration, ListID) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	_, err := handler.db.Exec(stmt, card.ID, card.Title, card.Description, card.StartDate, card.EndDate, card.DueDate, card.Duration, listID)
	logger.Log.Println(err)
	return err
}

func (handler *DbHandler) GetBoard(boardID string) (board.Board, error) {
	kanban := board.Board{}
	kanban.ID = boardID
	stmt := `SELECT Title FROM Boards WHERE ID = ?;`
	row := handler.db.QueryRow(stmt, boardID)
	if err := row.Scan(&kanban.Title); err != nil {
		return kanban, err
	}
	if err := handler.GetCardLists(&kanban); err != nil {
		return kanban, err
	}
	return kanban, nil
}

func (handler *DbHandler) GetCardLists(kanban *board.Board) error {
	stmt := `SELECT ID, Title FROM CardLists WHERE BoardID = ?;`
	rows, err := handler.db.Query(stmt, kanban.ID)
	if err != nil {
		logger.Log.Println(err)
		return err
	}
	defer rows.Close()
	for rows.Next() {
		cardList := board.CardList{}
		if err := rows.Scan(&cardList.ID, &cardList.Title); err != nil {
			logger.Log.Println(err)
			return err
		}
		if err := handler.GetCards(&cardList); err != nil {
			logger.Log.Println(err)
			return err
		}
		kanban.CardLists = append(kanban.CardLists, cardList)
	}
	return nil
}

func (handler *DbHandler) GetCards(cardList *board.CardList) error {
	stmt := `SELECT ID, Title, Description, StartDate, EndDate, DueDate, Duration FROM Cards WHERE ListID = ?;`
	rows, err := handler.db.Query(stmt, cardList.ID)
	if err != nil {
		logger.Log.Println(err)
		return err
	}

	defer rows.Close()
	for rows.Next() {
		card := board.Card{}
		var startDate string
		var endDate string
		var dueDate string
		if err := rows.Scan(&card.ID, &card.Title, &card.Description, &startDate, &endDate, &dueDate, &card.Duration); err != nil {
			logger.Log.Println(err)
			return err
		}
		parsedEndDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			logger.Log.Println(err)
			return err
		}

		parsedDueDate, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			logger.Log.Println(err)
			return err
		}
		parsedStartDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			logger.Log.Println(err)
			return err
		}
		card.EndDate = parsedEndDate
		card.DueDate = parsedDueDate
		card.StartDate = parsedStartDate
		cardList.Cards = append(cardList.Cards, card)
	}
	return nil
}

func (handler *DbHandler) UpdateBoardTitle(boardID string, title string) error {
	stmt := `UPDATE Boards SET Title = ? WHERE ID = ?;`
	if _, err := handler.db.Exec(stmt, title, boardID); err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}

func (handler *DbHandler) UpdateCardList(cardList board.CardList) error {
	stmt := `UPDATE CardLists SET Title = ? WHERE ID = ?;`
	if _, err := handler.db.Exec(stmt, cardList.Title, cardList.ID); err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}

func (handler *DbHandler) UpdateCard(card board.Card) error {
	stmt := `UPDATE Cards SET Title = ?, Description = ?, StartDate = ?, EndDate = ?, DueDate = ?, Duration = ? WHERE ID = ?;`
	if _, err := handler.db.Exec(stmt, card.Title, card.Description, card.StartDate, card.EndDate, card.DueDate, card.Duration, card.ID); err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}

func (handler *DbHandler) DeleteBoard(boardID string) error {
	stmt := `DELETE FROM Boards WHERE ID = ?;`
	if _, err := handler.db.Exec(stmt, boardID); err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}

func (handler *DbHandler) DeleteCardList(cardListID string) error {
	stmt := `DELETE FROM CardLists WHERE ID = ?;`
	if _, err := handler.db.Exec(stmt, cardListID); err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}

func (handler *DbHandler) DeleteCard(cardID string) error {
	stmt := `DELETE FROM Cards WHERE ID = ?;`
	if _, err := handler.db.Exec(stmt, cardID); err != nil {
		logger.Log.Println(err)
		return err
	}
	return nil
}
