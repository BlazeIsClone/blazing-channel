package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type DBTX interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type PostRepo interface {
	Create(ctx context.Context, tx DBTX, post *Post) error
	GetAll(ctx context.Context) ([]Post, error)
	FindByID(ctx context.Context, id int) (*Post, error)
	UpdateByID(ctx context.Context, tx DBTX, id int, post *Post) error
	DeleteByID(ctx context.Context, tx DBTX, id int) error
}

type PgSQLPostRepo struct {
	db *sql.DB
}

func NewPgSQLPostRepo(db *sql.DB) *PgSQLPostRepo {
	return &PgSQLPostRepo{db: db}
}

func (repo *PgSQLPostRepo) GetAll(ctx context.Context) ([]Post, error) {
	const query = `SELECT id, title, body, created_at, updated_at FROM posts ORDER BY id`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query posts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

func (repo *PgSQLPostRepo) Create(ctx context.Context, tx DBTX, post *Post) error {
	const query = `INSERT INTO posts (title, body) VALUES ($1, $2) RETURNING id, created_at, updated_at`

	return tx.QueryRowContext(ctx, query, post.Title, post.Body).
		Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
}

func (repo *PgSQLPostRepo) FindByID(ctx context.Context, id int) (*Post, error) {
	const query = `SELECT id, title, body, created_at, updated_at FROM posts WHERE id = $1`
	var post Post

	err := repo.db.QueryRowContext(ctx, query, id).
		Scan(&post.ID, &post.Title, &post.Body, &post.CreatedAt, &post.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find post: %w", err)
	}

	return &post, nil
}

func (repo *PgSQLPostRepo) UpdateByID(ctx context.Context, tx DBTX, id int, post *Post) error {
	const query = `UPDATE posts SET title=$1, body=$2, updated_at=NOW() WHERE id=$3 RETURNING updated_at`

	err := tx.QueryRowContext(ctx, query, post.Title, post.Body, id).Scan(&post.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("update post: %w", err)
	}

	post.ID = id
	return nil
}

func (repo *PgSQLPostRepo) DeleteByID(ctx context.Context, tx DBTX, id int) error {
	const query = `DELETE FROM posts WHERE id = $1`

	res, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}

	return nil
}
