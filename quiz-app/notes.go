package main

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"
)

var ErrNotFound = errors.New("note not found")

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type storeFile struct {
	NextID int    `json:"next_id"`
	Notes  []Note `json:"notes"`
}

type NoteStore struct {
	mu     sync.Mutex
	path   string
	notes  map[int]Note
	nextID int
}

func NewNoteStore(path string) (*NoteStore, error) {
	s := &NoteStore{
		path:   path,
		notes:  make(map[int]Note),
		nextID: 1,
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	} else if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var f storeFile
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}

	var largestID int
	for _, note := range f.Notes {
		s.notes[note.ID] = note
		largestID = max(largestID, note.ID)
	}

	if f.NextID <= largestID {
		return nil, fmt.Errorf("invalid %s: NextID must be greater than the largest note id %d", path, largestID)
	}
	s.nextID = f.NextID

	return s, nil
}

// The caller must already hold s.mu
func (s *NoteStore) save() error {
	f := storeFile{NextID: s.nextID, Notes: s.noteList()}

	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.path, data, 0644)
}

// The caller must already hold s.mu
func (s *NoteStore) noteList() []Note {
	notes := make([]Note, 0, len(s.notes))

	for _, note := range s.notes {
		notes = append(notes, note)
	}

	slices.SortFunc(notes, func(a, b Note) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return notes
}

func (s *NoteStore) Create(title, content string) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	newNote := Note{
		ID:        s.nextID,
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	s.notes[newNote.ID] = newNote
	s.nextID++

	if err := s.save(); err != nil {
		delete(s.notes, newNote.ID)
		s.nextID--
		return Note{}, err
	}

	return newNote, nil
}

func (s *NoteStore) Get(id int) (Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	note, ok := s.notes[id]
	return note, ok
}

func (s *NoteStore) List() []Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.noteList()
}

func (s *NoteStore) Update(id int, title, content string) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	old, ok := s.notes[id]
	if !ok {
		return Note{}, ErrNotFound
	}

	note := old
	note.Title = title
	note.Content = content

	s.notes[id] = note

	if err := s.save(); err != nil {
		s.notes[id] = old
		return Note{}, err
	}

	return note, nil
}

func (s *NoteStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	old, ok := s.notes[id]
	if !ok {
		return ErrNotFound
	}

	delete(s.notes, id)

	if err := s.save(); err != nil {
		s.notes[id] = old
		return err
	}

	return nil
}
