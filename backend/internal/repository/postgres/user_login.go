package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var idleThreshold = 1 * time.Hour

type userLoginRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewUserLoginRepo(db *pgxpool.Pool, tr Transaction) *userLoginRepo {
	return &userLoginRepo{
		db:          db,
		Transaction: tr,
	}
}

type UserLogins interface {
	GetByUser(ctx context.Context, req *models.GetUserLoginsDTO) ([]*models.UserLogin, error)
	GetByUserCount(ctx context.Context, userID string) (int64, error)
	GetLastByUser(ctx context.Context, userID string) (*models.UserLogin, error)
	Create(ctx context.Context, tx Tx, dto *models.UserLoginDTO) error
	UpdateLastActivity(ctx context.Context, tx Tx, userID string) (bool, error)
}

func (r *userLoginRepo) GetByUser(ctx context.Context, req *models.GetUserLoginsDTO) ([]*models.UserLogin, error) {
	baseQuery := fmt.Sprintf(`SELECT id, user_id, login_at, ip_address::text, user_agent, metadata 
		FROM %s WHERE user_id = $1`,
		Tables.UserLogins,
	)

	qb := NewQueryBuilder(baseQuery)
	qb.AddDateRangeFilter("login_at", req.StartDate, req.EndDate)
	qb.SetSort("login_at", true)

	if req.Limit > 0 {
		qb.SetLimit(req.Limit)
	}
	if req.Offset > 0 {
		qb.SetOffset(req.Offset)
	}

	query, args := qb.Build()

	var data []*models.UserLogin
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.UserLogin{}
		if err := rows.Scan(&tmp.ID, &tmp.UserID, &tmp.LoginAt, &tmp.IPAddress, &tmp.UserAgent, &tmp.Metadata); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		data = append(data, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return data, nil
}

func (r *userLoginRepo) GetByUserCount(ctx context.Context, userID string) (int64, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE user_id = $1`, Tables.UserLogins)

	var count int64
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user logins: %w", err)
	}
	return count, nil
}

func (r *userLoginRepo) GetLastByUser(ctx context.Context, userID string) (*models.UserLogin, error) {
	query := fmt.Sprintf(`SELECT id, user_id, login_at, ip_address::text, user_agent, metadata, last_activity_at 
		FROM %s WHERE user_id = $1 ORDER BY login_at DESC LIMIT 1`, Tables.UserLogins)

	var login models.UserLogin
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&login.ID, &login.UserID, &login.LoginAt, &login.IPAddress, &login.UserAgent, &login.Metadata, &login.LastActivityAt,
	)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last user login: %w", err)
	}
	return &login, nil
}

func (r *userLoginRepo) Create(ctx context.Context, tx Tx, dto *models.UserLoginDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, ip_address, user_agent, metadata) 
		VALUES ($1, $2, $3, $4, $5)`,
		Tables.UserLogins,
	)
	id := uuid.New()

	_, err := r.getExec(tx).Exec(ctx, query, id, dto.UserID, dto.IPAddress, dto.UserAgent, dto.Metadata)
	if err != nil {
		return fmt.Errorf("failed to create user login: %w", err)
	}
	return nil
}

func (r *userLoginRepo) UpdateLastActivity(ctx context.Context, tx Tx, userID string) (bool, error) {
	now := time.Now()

	login, err := r.GetLastByUser(ctx, userID)
	if err != nil {
		return false, err
	}

	if login == nil {
		return true, nil
	}

	timeSinceActivity := now.Sub(login.LastActivityAt)
	if timeSinceActivity > idleThreshold {
		return true, nil
	}

	query := fmt.Sprintf(`UPDATE %s SET last_activity_at = $1 WHERE id = $2`, Tables.UserLogins)
	_, err = r.getExec(tx).Exec(ctx, query, now, login.ID)
	if err != nil {
		return false, fmt.Errorf("failed to update last activity: %w", err)
	}

	return false, nil
}
