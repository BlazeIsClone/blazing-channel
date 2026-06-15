package post

import (
	"context"
	"database/sql"
	"fmt"
)

type Service struct {
	db   *sql.DB
	repo PostRepo
}

func NewPostSvc(db *sql.DB, repo PostRepo) *Service {
	return &Service{db: db, repo: repo}
}

func (svc *Service) Create(ctx context.Context, cmd CreatePostCommand) (*Post, error) {
	post := &Post{Title: cmd.Title, Body: cmd.Body}

	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create post: begin tx: %w", err)
	}

	if err := svc.repo.Create(ctx, tx, post); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("create post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create post: commit: %w", err)
	}

	return post, nil
}

func (svc *Service) GetAll(ctx context.Context) ([]Post, error) {
	return svc.repo.GetAll(ctx)
}

func (svc *Service) FindByID(ctx context.Context, id int) (*Post, error) {
	return svc.repo.FindByID(ctx, id)
}

func (svc *Service) Update(ctx context.Context, id int, cmd UpdatePostCommand) (*Post, error) {
	post := &Post{Title: cmd.Title, Body: cmd.Body}

	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("update post: begin tx: %w", err)
	}

	if err := svc.repo.UpdateByID(ctx, tx, id, post); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("update post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("update post: commit: %w", err)
	}

	return post, nil
}

func (svc *Service) Delete(ctx context.Context, id int) error {
	tx, err := svc.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("delete post: begin tx: %w", err)
	}

	if err := svc.repo.DeleteByID(ctx, tx, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete post: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("delete post: commit: %w", err)
	}

	return nil
}
