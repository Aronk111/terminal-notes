package main

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView uint = iota
	titleView
	bodyView
)

type model struct {
	state       uint
	store       *Store
	notes       []Note
	currentNote Note
	listIndex   int
	textarea    textarea.Model
	textinput   textinput.Model
}

func NewModel(store *Store) model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("Could not get notes: %v", err)
	}

	return model{
		state:     listView,
		store:     store,
		notes:     notes,
		textarea:  textarea.New(),
		textinput: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit
			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.currentNote = Note{}

				m.state = titleView
			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				} else if m.listIndex == 0 {
					m.listIndex = len(m.notes) - 1
				}
			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				} else if m.listIndex == len(m.notes)-1 {
					m.listIndex = 0
				}
			case "enter":
				if len(m.notes) == 0 {
					break
				}

				m.currentNote = m.notes[m.listIndex]

				m.textarea.SetValue(m.currentNote.Body)
				m.textarea.Focus()
				m.textarea.CursorEnd()

				m.state = bodyView
			case "d":
				m.currentNote = m.notes[m.listIndex]

				m.store.DeleteNote(m.currentNote)

				var err error
				m.notes, err = m.store.GetNotes()
				if err != nil {
					return m, tea.Quit
				}

				if len(m.notes) != 0 {
					m.listIndex = len(m.notes) - 1
				}
			}
		case titleView:
			switch key {
			case "esc":
				m.state = listView
			case "enter":
				title := m.textinput.Value()
				if title != "" {
					m.currentNote.Title = title

					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()

					m.state = bodyView
				}
			}
		case bodyView:
			switch key {
			case "esc":
				m.state = listView
			case "ctrl+s":
				body := m.textarea.Value()
				if body != "" {
					m.currentNote.Body = body

					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textarea.CursorEnd()

					var err error
					if err = m.store.SaveNote(m.currentNote); err != nil {
						return m, tea.Quit
					}

					m.notes, err = m.store.GetNotes()
					if err != nil {
						return m, tea.Quit
					}

					m.currentNote = Note{}
					m.state = listView
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}
