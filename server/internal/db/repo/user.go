package repo

import (
	"context"
	"fmt"

	"hata/internal/orm"
	"hata/internal/orm/user"
)

// UserRepo implements the UserRepository interface.
type UserRepo struct {
	client *orm.Client
}

// GetByUsername retrieves a user by username.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*UserData, error) {
	u, err := r.client.User.Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)

	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil // User not found is not an error
		}
		return nil, fmt.Errorf("failed querying user by username %q: %w", username, err)
	}

	return &UserData{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
	}, nil
}

// Create creates a new user with the given credentials.
func (r *UserRepo) Create(ctx context.Context, username, passwordHash, displayName string) (*UserData, error) {
	u, err := r.client.User.Create().
		SetUsername(username).
		SetPasswordHash(passwordHash).
		SetDisplayName(displayName).
		Save(ctx)

	if err != nil {
		// Check for unique constraint violation
		if orm.IsConstraintError(err) {
			return nil, fmt.Errorf("username %q already exists", username)
		}
		return nil, fmt.Errorf("failed creating user: %w", err)
	}

	return &UserData{
		ID:           u.ID,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		DisplayName:  u.DisplayName,
	}, nil
}

// List retrieves users with pagination.
func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]*UserData, int, error) {
	// Get total count
	total, err := r.client.User.Query().Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting users: %w", err)
	}

	// Get paginated users
	users, err := r.client.User.Query().
		Limit(limit).
		Offset(offset).
		All(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("failed querying users: %w", err)
	}

	// Convert to UserData
	result := make([]*UserData, len(users))
	for i, u := range users {
		result[i] = &UserData{
			ID:           u.ID,
			Username:     u.Username,
			PasswordHash: u.PasswordHash,
			DisplayName:  u.DisplayName,
		}
	}

	return result, total, nil
}
