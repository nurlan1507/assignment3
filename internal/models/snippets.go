package models

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	Db *pgxpool.Pool
}

func (m *SnippetModel) Insert(title string, content string, expires int) (*Snippet, error) {
	stmt := `INSERT INTO snippets(title,content, created,expires) VALUES ($1,$2,current_date, current_date +'$3 year' )`
	query, err := m.Db.Query(context.Background(), stmt, title, content, expires)
	if err != nil {
		return nil, ErrNoRecord
	}
	snippet := &Snippet{}
	err = query.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return snippet, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT * FROM snippets s WHERE id= $1`
	result := m.Db.QueryRow(context.Background(), stmt, id)
	s := &Snippet{}
	err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT * FROM snippets s WHERE current_date < s.expires limit 10`
	result, err := m.Db.Query(context.Background(), stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
	}
	defer result.Close()
	var snippets []*Snippet
	for result.Next() {
		s := &Snippet{}
		err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Expires, &s.Expires)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			}
			return nil, err
		}
		snippets = append(snippets, s)
	}
	return snippets, nil
}

//func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
//	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
//	result, err := m.Db.Exec(stmt, title, content, expires)
//	if err != nil {
//		return 0, err
//	}
//	id, err := result.LastInsertId() //id of inserted obj
//	if err != nil {
//		return 0, err
//	}
//	return int(id), nil
//}
//
//// Get This will return a specific snippet based on its id.
//func (m *SnippetModel) Get(id int) (*Snippet, error) {
//	stmt := `SELECT * FROM snippets  WHERE id = ?`
//	result := m.Db.QueryRow(stmt, id)
//	s := &Snippet{}
//	err := result.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, ErrNoRecord
//		} else {
//			return nil, err
//		}
//	}
//	return s, nil
//}
//
//// Latest This will return the 10 most recently created snippets.
//func (m *SnippetModel) Latest() ([]*Snippet, error) {
//	var stmt = `SELECT id, title, content ,expires, created FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
//	objects, err := m.Db.Query(stmt)
//	if err != nil {
//		if errors.Is(err, sql.ErrNoRows) {
//			return nil, ErrNoRecord
//		} else {
//			return nil, err
//		}
//	}
//	defer objects.Close()
//	var snippets []*Snippet
//	for objects.Next() {
//		object := &Snippet{}
//		err := objects.Scan(&object.ID, &object.Title, &object.Content, &object.Created, &object.Expires)
//		if err != nil {
//			return nil, err
//		}
//		snippets = append(snippets, object)
//	}
//	if err = objects.Err(); err != nil {
//		return nil, err
//	}
//	return snippets, nil
//}
