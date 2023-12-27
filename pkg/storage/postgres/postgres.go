package postgres

import (
	"GoNews/pkg/storage"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func Init(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Posts() (posts []storage.Post, err error) {
	rows, err := s.db.Query(
		context.Background(),
		`SELECT * FROM posts`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var p storage.Post
		if err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.AuthorID,
			&p.AuthorName,
			&p.CreatedAt,
			&p.PublishedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

func (s *Storage) AddPost(post storage.Post) error {
	var authorName string
	err := s.db.QueryRow(
		context.Background(),
		`SELECT name FROM authors WHERE id = $1`, post.AuthorID).
		Scan(&authorName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("author with ID %d does not exist", post.AuthorID)
		}
		return fmt.Errorf("failed to query author: %w", err)
	}

	_, err = s.db.Exec(
		context.Background(),
		`INSERT INTO posts (title, content, author_id, author_name, published_at) VALUES ($1, $2, $3, $4, $5)`,
		post.Title,
		post.Content,
		post.AuthorID,
		authorName,
		time.Now().Format(time.DateTime))
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	return nil
}

func (s *Storage) UpdatePost(post storage.Post) error {
	exists, err := checkExistsId(s, post.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("post with ID %d does not exist", post.ID)
	}

	var authorName string
	err = s.db.QueryRow(context.Background(),
		`SELECT name FROM authors WHERE id = $1`, post.AuthorID).
		Scan(&authorName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("author with ID %d does not exist", post.AuthorID)
		}
		return fmt.Errorf("failed to query author: %w", err)
	}

	_, err = s.db.Exec(
		context.Background(),
		`UPDATE posts SET title = $1, content = $2, author_id = $3, author_name = $4 WHERE id = $5;`,
		post.Title,
		post.Content,
		post.AuthorID,
		authorName,
		post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}

	return nil
}

func (s *Storage) DeletePost(post storage.Post) error {
	exists, err := checkExistsId(s, post.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("post with ID %d does not exist", post.ID)
	}

	_, err = s.db.Exec(
		context.Background(),
		`DELETE FROM posts WHERE id = $1`, post.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

// проверка id в БД, перед выполнением операций
func checkExistsId(s *Storage, id int) (exists bool, err error) {
	err = s.db.QueryRow(
		context.Background(),
		`SELECT EXISTS(SELECT * FROM posts WHERE id = $1)`, id).
		Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, err
	}
	return false, nil
}
