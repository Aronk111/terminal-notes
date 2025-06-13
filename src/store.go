package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Note struct {
	ID    uint
	Title string
	Body  string
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error

	s.conn, err = sql.Open("sqlite3", "../dist/notes.db")
	if err != nil {
		log.Printf("Could not open database: %v", err)
		return err
	}

	createTableStatement := `CREATE TABLE IF NOT EXISTS notes (
		id integer not null primary key,
		title text not null,
		body text not null
	)`

	if _, err := s.conn.Exec(createTableStatement); err != nil {
		log.Printf("Could not execute statement: %v", err)
		return err
	}

	return nil
}

func (s *Store) GetNotes() ([]Note, error) {
	rows, err := s.conn.Query(`SELECT * FROM notes`)
	if err != nil {
		log.Printf("Could not get notes: %v", err)
		return nil, err
	}

	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var note Note
		rows.Scan(&note.ID, &note.Title, &note.Body)
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *Store) SaveNote(note Note) error {
	if note.ID == 0 {
		note.ID = uint(time.Now().UTC().UnixNano())
	}

	upsertQuery := `INSERT INTO notes (id, title, body)
	VALUES (?, ?, ?)
	ON CONFLICT(id) DO UPDATE
	SET title=excluded.title, body=excluded.body;`

	if _, err := s.conn.Exec(upsertQuery, note.ID, note.Title, note.Body); err != nil {
		log.Printf("Could not update or insert note: %v", err)
		return err
	}

	return nil
}
