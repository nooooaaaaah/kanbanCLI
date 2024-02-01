package board

import "time"

type Card struct {
	ID          string
	DueDate     time.Time
	StartDate   time.Time
	EndDate     time.Time
	Title       string
	Description string
	Duration    time.Duration
}

type CardList struct {
	ID    string
	Title string
	Cards []Card
}

type Board struct {
	ID        string
	Title     string
	CardLists []CardList
}

func NewBoard(title string) Board {
	return Board{Title: title}
}

func NewCardList(title string) CardList {
	return CardList{Title: title}
}

func NewCard(title string) Card {
	return Card{Title: title}
}

func (b *Board) AddCardList(cardList CardList) {
	b.CardLists = append(b.CardLists, cardList)
}

func (cl *CardList) AddCard(card Card) {
	cl.Cards = append(cl.Cards, card)
}

func (b *Board) GetCardList(id string) *CardList {
	for _, cardList := range b.CardLists {
		if cardList.ID == id {
			return &cardList
		}
	}
	return nil
}

func (cl *CardList) GetCard(id string) *Card {
	for _, card := range cl.Cards {
		if card.ID == id {
			return &card
		}
	}
	return nil
}

func (b *Board) RemoveCardList(id string) {
	for i, cardList := range b.CardLists {
		if cardList.ID == id {
			b.CardLists = append(b.CardLists[:i], b.CardLists[i+1:]...)
			return
		}
	}
}

func (cl *CardList) RemoveCard(id string) {
	for i, card := range cl.Cards {
		if card.ID == id {
			cl.Cards = append(cl.Cards[:i], cl.Cards[i+1:]...)
			return
		}
	}
}

func (b *Board) UpdateCardList(cardList CardList) {
	for i, c := range b.CardLists {
		if c.ID == cardList.ID {
			b.CardLists[i] = cardList
			return
		}
	}
}

func (cl *CardList) UpdateCard(card Card) {
	for i, c := range cl.Cards {
		if c.ID == card.ID {
			cl.Cards[i] = card
			return
		}
	}
}

func (b *Board) UpdateCard(card Card) {
	for _, cardList := range b.CardLists {
		for i, c := range cardList.Cards {
			if c.ID == card.ID {
				cardList.Cards[i] = card
				return
			}
		}
	}
}
