package repo

import (
	"context"
	"fmt"
	"time"

	"hata/internal/orm"
	"hata/internal/orm/session"
)

// SessionRepo implements the SessionRepository interface.
type SessionRepo struct {
	client *orm.Client
}

// Create creates a new session for the given user.
func (r *SessionRepo) Create(ctx context.Context, userID int, token string, validUntil time.Time) (*SessionData, error) {
	s, err := r.client.Session.Create().
		SetUserID(userID).
		SetToken(token).
		SetValidUntil(validUntil).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed creating session: %w", err)
	}

	return &SessionData{
		ID:           s.ID,
		UserID:       s.UserID,
		CreatedAt:    s.CreatedAt,
		ValidUntil:   s.ValidUntil,
		Token:        s.Token,
		InvalidSince: s.InvalidSince,
	}, nil
}

// GetByToken retrieves a session by its token.
func (r *SessionRepo) GetByToken(ctx context.Context, token string) (*SessionData, error) {
	s, err := r.client.Session.Query().
		Where(session.TokenEQ(token)).
		Only(ctx)

	if err != nil {
		if orm.IsNotFound(err) {
			return nil, nil // Session not found is not an error
		}
		return nil, fmt.Errorf("failed querying session by token: %w", err)
	}

	return &SessionData{
		ID:           s.ID,
		UserID:       s.UserID,
		CreatedAt:    s.CreatedAt,
		ValidUntil:   s.ValidUntil,
		Token:        s.Token,
		InvalidSince: s.InvalidSince,
	}, nil
}

// Invalidate marks a session as invalid.
func (r *SessionRepo) Invalidate(ctx context.Context, token string) error {
	now := time.Now()

	count, err := r.client.Session.Update().
		Where(session.TokenEQ(token)).
		SetInvalidSince(now).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed invalidating session: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("session with token not found")
	}

	return nil
}
