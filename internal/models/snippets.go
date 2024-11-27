package models

import (
    "database/sql"
    "errors"
    "time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expiress int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
    ID int
    Title string
    Content string
    Created time.Time
    Expires time.Time
}

type SnippetModel struct {
    DB *sql.DB
}

// Insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expireDays int) (int, error) {
    stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    result, err := m.DB.Exec(stmt, title, content, expireDays)
    if err != nil {
        return 0, err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

// Return a specific snippet based on its id
func (m *SnippetModel) Get(id int) (*Snippet, error) {
    stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

    s := &Snippet{}

    err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrNoRecord
        } else {
            return nil, err
        }
    }
    return s, nil
}

// Return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
    stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE  expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

    rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    snippets := []*Snippet{}

    for rows.Next() {
        s := &Snippet{}

        err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
        if err != nil {
            return nil, err
        }

        snippets = append(snippets, s)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return snippets, nil
}
