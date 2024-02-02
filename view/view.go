package view

import (
	"kanban/board"
	"kanban/logger"
	"kanban/utils"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	boardBackgroundStyle = lipgloss.NewStyle().Background(lipgloss.Color("#00FF"))

	boardTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#007BFF")).
			Padding(0, 2)
	listTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#200000")).
			Background(lipgloss.Color("#FFFF00")).
			Padding(0, 2).
			Margin(0, 2)

	cardStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFFFFF")).
			Padding(0, 2)
)

type ListItem struct {
	title string
}

func (li ListItem) Title() string {
	return li.title
}

func (li ListItem) FilterValue() string {
	return li.title
}

func (li ListItem) Description() string {
	return "" // return an empty string if you don't need a description
}

type KanbanModel struct {
	dbService   utils.DBService
	selected    map[int]bool // selected items
	lists       [][]ListItem // Changed from []list.Model
	board       board.Board
	newCard     board.Card
	addCard     bool
	addList     bool
	createBoard bool
	cursor      int // cursor position
}

// Generate a UUID for new items
func (m *KanbanModel) GenerateUUID() string {
	return utils.GenerateUUID()
}

// InitializeDB initializes the database service.
func (m *KanbanModel) InitializeDB() error {
	var err error
	m.dbService, err = utils.InitializeDBService()
	return err
}

// LoadBoard loads the board data from the database.
func (m *KanbanModel) LoadBoard(boardID string) error {
	loadedBoard, err := m.dbService.GetBoard(boardID)
	if err != nil {
		return err
	}
	m.board = loadedBoard
	return nil
}

// InitializeLists initializes the lists from the loaded board.
func (m *KanbanModel) InitializeLists() {
	for _, cardList := range m.board.CardLists {
		var items []ListItem
		for _, card := range cardList.Cards {
			items = append(items, ListItem{title: card.Title})
		}
		m.lists = append(m.lists, items) // Adjust the width and height as needed
	}
}

func (m *KanbanModel) AddCard(listIndex int, cardTitle string) {
	cardList := m.board.CardLists[listIndex]
	cardList.Cards = append(cardList.Cards, board.Card{Title: cardTitle})
	m.board.CardLists[listIndex] = cardList
}

func (m *KanbanModel) RemoveCard(listIndex int, cardIndex int) {
	cardList := m.board.CardLists[listIndex]
	cardList.Cards = append(cardList.Cards[:cardIndex], cardList.Cards[cardIndex+1:]...)
	m.board.CardLists[listIndex] = cardList
}

func (m *KanbanModel) MoveCard(fromListIndex int, toListIndex int, cardIndex int) {
	card := m.board.CardLists[fromListIndex].Cards[cardIndex]
	m.RemoveCard(fromListIndex, cardIndex)
	m.AddCard(toListIndex, card.Title)
}

func (m *KanbanModel) SaveBoard() error {
	return m.dbService.UpdateBoard(m.board)
}

func (m *KanbanModel) AddList(listTitle string) {
	m.board.CardLists = append(m.board.CardLists, board.CardList{Title: listTitle})
	m.lists = append(m.lists, []ListItem{})
}

func (m *KanbanModel) RemoveList(listIndex int) {
	m.board.CardLists = append(m.board.CardLists[:listIndex], m.board.CardLists[listIndex+1:]...)
	m.lists = append(m.lists[:listIndex], m.lists[listIndex+1:]...)
}

func (m *KanbanModel) Init() tea.Cmd {
	m.selected = make(map[int]bool) // Initialize the selected map

	// Initialize the dbService
	if err := m.InitializeDB(); err != nil {
		logger.Log.Println(err)
		return func() tea.Msg {
			return err
		}
	}

	boardID := "default-board"

	// Load the board data
	if err := m.LoadBoard(boardID); err != nil {
		logger.Log.Println(err)
		return func() tea.Msg {
			return err
		}
	}

	// Initialize the lists
	m.InitializeLists()

	return nil // You can return a command if needed
}

func (m *KanbanModel) handleAddCard(key string) {
	switch key {
	case "/":
		if m.newCard.Title != "" {
			m.addCard = false
			m.AddCard(m.cursor, m.newCard.Title)
			m.newCard = board.Card{}
			// Force a re-render by updating the model
			m = &KanbanModel{
				board:     m.board,
				dbService: m.dbService,
				cursor:    m.cursor,
				selected:  m.selected,
				addCard:   m.addCard,
				newCard:   m.newCard,
			}
		}
	case "backspace":
		if len(m.newCard.Title) > 0 {
			m.newCard.Title = m.newCard.Title[:len(m.newCard.Title)-1]
		}
	case "esc":
		m.addCard = false
	default:
		m.newCard.Title += key
	}
}

func (m *KanbanModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		// handle the error, perhaps by setting an error field in your model
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			if m.cursor > 0 {
				m.cursor--
			}
		case "right", "l":
			if m.cursor < len(m.board.CardLists)-1 {
				m.cursor++
			}
		case "i":
			m.addCard = !m.addCard
		case "enter":
			// Toggle the selection of the currently highlighted list
			m.selected[m.cursor] = !m.selected[m.cursor]
		default:
			if m.addCard {
				m.handleAddCard(msg.String())
			}
		}
	}
	return m, nil
}

// GetPrefix returns the prefix for a list based on its index.
func (m *KanbanModel) GetPrefix(i int) string {
	prefix := ""
	if i == m.cursor {
		prefix += "> "
	}
	if m.selected[i] {
		prefix += "* "
	}
	return prefix
}

// GetListLines returns the lines for a list.
func (m *KanbanModel) GetListLines(list []ListItem, listTitle string) []string {
	var lines []string
	lines = append(lines, listTitleStyle.Render(listTitle))

	for _, item := range list {
		lines = append(lines, cardStyle.Render(item.title))
	}

	return lines
}

// GetMaxHeight returns the maximum height among a slice of line slices.
func GetMaxHeight(listLines [][]string) int {
	maxHeight := 0
	for _, lines := range listLines {
		if len(lines) > maxHeight {
			maxHeight = len(lines)
		}
	}
	return maxHeight
}

func (m *KanbanModel) View() string {
	const listWidth = 20 // Adjust this value as needed
	if m.board.Title == "" {
		return "No board loaded\n"
	}

	var buf strings.Builder
	var listLines [][]string

	totalWidth := listWidth * len(m.lists)
	titleMargin := (totalWidth+len(m.board.Title))/2 + listWidth // I had to add the length of the title to get it centered???? Literally the opposite of what I expected

	boardTitleStyle.Margin(1, titleMargin)
	buf.WriteString(boardTitleStyle.Render(m.board.Title))
	buf.WriteString("\n")

	for i, list := range m.lists {
		prefix := m.GetPrefix(i)
		listTitle := prefix + m.board.CardLists[i].Title
		lines := m.GetListLines(list, listTitle)
		listLines = append(listLines, lines)
	}

	maxHeight := GetMaxHeight(listLines)

	for i := 0; i < maxHeight; i++ {
		buf.WriteString(strings.Repeat(" ", listWidth))
		for _, lines := range listLines {
			if i < len(lines) {
				buf.WriteString(lines[i])
			}
			buf.WriteString(strings.Repeat(" ", listWidth))
		}
		buf.WriteString("\n")
		buf.WriteString("\n")
	}

	if m.addCard {
		buf.WriteString("Add card: ")
		buf.WriteString(m.newCard.Title)
	}

	// Apply the boardBackgroundStyle to the entire output and add a border
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Render(boardBackgroundStyle.Render(buf.String()))
}
