package main

import (
	"cmp"
	"slices"
	"sync"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type NoteStore struct {
	mu     sync.Mutex // must hold to read/write
	notes  map[int]Note
	nextID int
}

func NewNoteStore() *NoteStore {
	return &NoteStore{
		notes:  make(map[int]Note),
		nextID: 1,
	}
}

func (s *NoteStore) Create(title, content string) Note {
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

	return newNote
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

	notes := make([]Note, 0, len(s.notes))

	for _, note := range s.notes {
		notes = append(notes, note)
	}

	slices.SortFunc(notes, func(a, b Note) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return notes
}

func (s *NoteStore) Update(id int, title, content string) (Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	note, ok := s.notes[id]
	if !ok {
		return Note{}, false
	}

	note.Title = title
	note.Content = content

	s.notes[id] = note

	return note, true
}

func (s *NoteStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.notes[id]
	if !ok {
		return false
	}

	delete(s.notes, id)

	return true
}
