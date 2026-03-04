// internal/infrastructure/repository/token_repository.go
package repository

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	appErrors "github.com/overm-app/api-auth/internal/domain/errors"
	"github.com/overm-app/api-auth/internal/domain/models"
	"github.com/overm-app/api-auth/internal/domain/ports"
)

type TokenRepository struct {
	db    *sql.DB
	sq    sq.StatementBuilderType
}

func NewTokenRepository(db *sql.DB) ports.TokenRepository {
	return &TokenRepository{
		db:    db,
		sq:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *TokenRepository) Save(ctx context.Context, token *models.RefreshToken) error {
	query, args, err := r.sq.
		Insert("refresh_tokens").
		Columns("id", "user_id", "token", "expires_at", "created_at").
		Values(token.ID, token.UserID, token.Token, token.ExpiresAt, time.Now()).
		ToSql()
	if err != nil {
		return appErrors.Internal("Failed to build save token query", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return appErrors.Internal("Failed to save refresh token", err)
	}

	return nil
}

func (r *TokenRepository) GetByUserID(ctx context.Context, userID string) (*models.RefreshToken, error) {
	query, args, err := r.sq.
		Select("id", "user_id", "token", "expires_at", "created_at").
		From("refresh_tokens").
		Where(sq.Eq{"user_id": userID}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, appErrors.Internal("Failed to build get token query", err)
	}

	var token models.RefreshToken
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, appErrors.Internal("Failed to scan refresh token", err)
	}

	return &token, nil
}

func (r *TokenRepository) GetByID(ctx context.Context, id string) (*models.RefreshToken, error) {
	query, args, err := r.sq.
		Select("id", "user_id", "token", "expires_at", "created_at").
		From("refresh_tokens").
		Where(sq.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, appErrors.Internal("Failed to build get token by id query", err)
	}

	var token models.RefreshToken
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, appErrors.Internal("Failed to scan refresh token", err)
	}

	return &token, nil
}

func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query, args, err := r.sq.
		Delete("refresh_tokens").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return appErrors.Internal("Failed to build delete token query", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return appErrors.Internal("Failed to delete refresh token", err)
	}

	return nil
}

func (r *TokenRepository) RotateToken(ctx context.Context, oldID string, newToken *models.RefreshToken) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return appErrors.Internal("Failed to begin transaction", err)
	}
	defer tx.Rollback()

	deleteQuery, deleteArgs, err := r.sq.
		Delete("refresh_tokens").
		Where(sq.Eq{"id": oldID}).
		ToSql()
	if err != nil {
		return appErrors.Internal("Failed to build delete query", err)
	}

	_, err = tx.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return appErrors.Internal("Failed to delete old token", err)
	}

	insertQuery, insertArgs, err := r.sq.
		Insert("refresh_tokens").
		Columns("id", "user_id", "token", "expires_at", "created_at").
		Values(newToken.ID, newToken.UserID, newToken.Token, newToken.ExpiresAt, time.Now()).
		ToSql()
	if err != nil {
		return appErrors.Internal("Failed to build insert query", err)
	}

	_, err = tx.ExecContext(ctx, insertQuery, insertArgs...)
	if err != nil {
		return appErrors.Internal("Failed to insert new token", err)
	}

	if err := tx.Commit(); err != nil {
		return appErrors.Internal("Failed to commit token rotation", err)
	}

	return nil
}
