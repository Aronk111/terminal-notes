package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	store := &Store{}
	if err := store.Init(); err != nil {
		log.Fatalf("Could not initialize store: %v", err)
	}

	m := NewModel(store)

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Could not start application: %v", err)
	}
}
