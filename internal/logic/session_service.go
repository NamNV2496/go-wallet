package logic

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
)

type SessionLogic interface {
	CreateSession(ctx context.Context, sessionId uuid.UUID, userId int64, username string, refreshToken string, userAgent string, clientIp string, isBlocked bool, expiredAt time.Time) (db.Session, error)
	GetSession(ctx context.Context, id uuid.UUID) (db.Session, error)
	UpdateExpiredTime(ctx context.Context, id uuid.UUID, expired time.Time) (db.Session, error)
}

var _ SessionLogic = (*sessionLogic)(nil)

type sessionLogic struct {
	database *databaseaccess.Database
}

func NewSessionLogic(
	database *databaseaccess.Database,
) (SessionLogic, error) {
	return &sessionLogic{
		database: database,
	}, nil
}

func (s sessionLogic) CreateSession(
	ctx context.Context,
	sessionId uuid.UUID,
	userId int64,
	username string,
	refreshToken string,
	userAgent string,
	clientIp string,
	isBlocked bool,
	expiredAt time.Time,
) (db.Session, error) {

	arg := db.CreateSessionParams{
		ID:           sessionId,
		UserID:       userId,
		Username:     username,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIp:     clientIp,
		IsBlocked:    isBlocked,
		ExpiresAt:    expiredAt,
	}
	return s.database.CreateSession(ctx, arg)
}

func (s sessionLogic) GetSession(
	ctx context.Context,
	id uuid.UUID,
) (db.Session, error) {

	return s.database.GetSession(ctx, id)
}

func (s sessionLogic) UpdateExpiredTime(
	ctx context.Context,
	id uuid.UUID,
	expired time.Time,
) (db.Session, error) {

	arg := db.UpdateExpiredTimeParams{
		ID:        id,
		ExpiresAt: expired,
	}
	return s.database.UpdateExpiredTime(ctx, arg)
}
