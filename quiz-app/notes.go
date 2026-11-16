package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteStore struct {
	db *sql.DB
}

func NewNoteStore(dsn string) (*NoteStore, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	err = migrate(db, "migrations")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &NoteStore{db: db}, nil
}

func (s *NoteStore) Close() error {
	return s.db.Close()
}

func (s *NoteStore) Create(title, content string) (Note, error) {
	note := Note{Title: title, Content: content, CreatedAt: time.Now().UTC()}

	err := s.db.QueryRow(
		`INSERT INTO notes (
			title,
			content,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4) RETURNING id`,
		note.Title, note.Content, note.CreatedAt, note.UpdatedAt,
	).Scan(&note.ID)

	if err != nil {
		return Note{}, err
	}

	note.UpdatedAt = note.CreatedAt

	return note, nil
}

func (s *NoteStore) Get(id int) (Note, error) {
	var note Note

	err := s.db.QueryRow(
		`SELECT id, title, content, created_at, updated_at FROM notes WHERE id = $1`, id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, ErrNotFound
	} else if err != nil {
		return Note{}, err
	}

	note.CreatedAt = note.CreatedAt.UTC()
	note.UpdatedAt = note.UpdatedAt.UTC()

	return note, nil
}

func (s *NoteStore) List() ([]Note, error) {
	rows, err := s.db.Query(`SELECT id, title, content, created_at, updated_at FROM notes ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]Note, 0)

	for rows.Next() {
		var note Note

		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		note.CreatedAt = note.CreatedAt.UTC()
		note.UpdatedAt = note.UpdatedAt.UTC()

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
		UPDATE notes SET title = $1, content = $2, updated_at = $3 WHERE id = $4
		RETURNING id, title, content, created_at, updated_at`, title, content, time.Now().UTC(), id,
	).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return Note{}, ErrNotFound
	} else if err != nil {
		return Note{}, err
	}

	note.CreatedAt = note.CreatedAt.UTC()
	note.UpdatedAt = note.UpdatedAt.UTC()

	return note, nil
}

func (s *NoteStore) Delete(id int) error {
	res, err := s.db.Exec(`DELETE FROM notes WHERE id = $1`, id)
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
