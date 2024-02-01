package main

import (
	"kanban/logger"
	"kanban/utils"
	"kanban/view"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.InitLogger()
	// if no kanban.db file exists, create one
	if _, err := os.Stat("./kanban.db"); os.IsNotExist(err) {
		logger.Log.Println("Creating kanban.db")
		dbHandler, err := utils.InitializeDBService()
		if err != nil {
			logger.Log.Println(err)
			os.Exit(1)
		}
		_, err = dbHandler.CreateDB()
		if err != nil {
			logger.Log.Println(err)
			os.Exit(1)
		}
	}
	m := view.KanbanModel{}
	p := tea.NewProgram(&m)
	if _, err := p.Run(); err != nil {
		logger.Log.Println(err)
		os.Exit(1)
	}
}
