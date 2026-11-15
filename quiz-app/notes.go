package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteStore struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS notes (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT      NOT NULL,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL
)`

func NewNoteStore(path string) (*NoteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("schema %s: %w", schema, err)
	}

	return &NoteStore{db: db}, nil
}

func (s *NoteStore) Close() error {
	return s.db.Close()
}

func (s *NoteStore) Create(title, content string) (Note, error) {
	note := Note{Title: title, Content: content, CreatedAt: time.Now().UTC()}

	err := s.db.QueryRow(
		`INSERT INTO notes (title, content, created_at) VALUES (?, ?, ?) RETURNING id`,
		note.Title, note.Content, note.CreatedAt,
	).Scan(&note.ID)

	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (s *NoteStore) Get(id int) (Note, error) {
	var note Note

	err := s.db.QueryRow(
		`SELECT id, title, content, created_at FROM notes WHERE id = ?`, id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, ErrNotFound
	} else if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (s *NoteStore) List() ([]Note, error) {
	rows, err := s.db.Query(`SELECT id, title, content, created_at FROM notes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)

	for rows.Next() {
		var note Note

		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (s *NoteStore) Update(id int, title, content string) (Note, error) {
	var note Note

	err := s.db.QueryRow(`
		UPDATE notes SET title = ?, content = ? WHERE id = ?
		RETURNING id, title, content, created_at`, title, content, id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, ErrNotFound
	} else if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (s *NoteStore) Delete(id int) error {
	res, err := s.db.Exec(`DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}

	return nil
}
