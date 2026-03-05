package repository

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/google/uuid"
    "github.com/Masterminds/squirrel"
    sq "github.com/Masterminds/squirrel"

    "github.com/overm-app/api-auth/internal/domain/models"
    "github.com/overm-app/api-auth/internal/domain/ports"
)

type UserRepository struct {
	db *sql.DB
	sq sq.StatementBuilderType
}

func NewUserRepository(db *sql.DB) ports.UserRepository {
	return &UserRepository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
    query, args, err := r.sq.
        Select("id", "public_id", "email", "name", "password_hash", "auth_provider", "avatar_url", "created_at", "updated_at").
        From("users").
        Where(squirrel.Eq{"email": email}).
        Limit(1).
        ToSql()
    if err != nil {
        return nil, fmt.Errorf("failed to build SQL query: %w", err)
    }

    var user models.User
    err = r.db.QueryRowContext(ctx, query, args...).Scan(
        &user.ID,
        &user.PublicID,
        &user.Email,
        &user.Name,
        &user.PasswordHash,
        &user.AuthProvider,
        &user.AvatarURL,
        &user.CreatedAt,
        &user.UpdatedAt,
    )
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to execute SQL query: %w", err)
    }

    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
    user.PublicID = uuid.New().String()

    query, args, err := r.sq.
        Insert("users").
        Columns("public_id", "email", "name", "password_hash", "auth_provider", "avatar_url").
        Values(user.PublicID, user.Email, user.Name, user.PasswordHash, user.AuthProvider, user.AvatarURL).
        Suffix("RETURNING id, created_at, updated_at").
        ToSql()
    if err != nil {
        return fmt.Errorf("failed to build query: %w", err)
    }

    return r.db.QueryRowContext(ctx, query, args...).Scan(
        &user.ID,
        &user.CreatedAt,
        &user.UpdatedAt, 
	)
}