package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewUserRepo(db *pgxpool.Pool, tr Transaction) *userRepo {
	return &userRepo{
		db:          db,
		Transaction: tr,
	}
}

type Users interface {
	LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.UserRole, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.UserData, error)
	GetAll(ctx context.Context) ([]*models.UserData, error)
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.UserDataDTO) error
	Update(ctx context.Context, tx Tx, dto *models.UserDataDTO) error
	UpdateSeveral(ctx context.Context, tx Tx, dto []*models.UserDataDTO) error
	DeleteSeveral(ctx context.Context, tx Tx, ids []string) error
}

func (r *userRepo) LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.UserRole, error) {
	query := fmt.Sprintf(`SELECT u.sso_id, r.slug
        FROM %s u
        JOIN %s r ON r.id = u.role_id`,
		Tables.Users, Tables.Roles,
	)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	data := make([]*models.UserRole, 0, 50)
	for rows.Next() {
		item := &models.UserRole{}
		if err := rows.Scan(&item.UserID, &item.RoleName); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		data = append(data, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return data, nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.UserData, error) {
	query := fmt.Sprintf(`SELECT u.id, role_id u.username, u.email, u.sso_id, u.first_name, u.last_name, u.is_active,
		FROM %s u
		WHERE u.id = $1
		GROUP BY u.id`,
		Tables.Users,
	)
	data := &models.UserData{}
	if err := r.db.QueryRow(ctx, query, id).Scan(
		&data.ID, &data.RoleID, &data.Username, &data.Email, &data.SSO_ID, &data.FirstName, &data.LastName, &data.IsActive,
	); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *userRepo) GetAll(ctx context.Context) ([]*models.UserData, error) {
	query := fmt.Sprintf(`SELECT u.id, u.username, u.email, u.sso_id, u.first_name, u.last_name, u.is_active, u.created_at,
			r.name as role, COALESCE(MAX(l.login_at), '1970-01-01 00:00:00') as last_login_at
		FROM %s u
		LEFT JOIN %s r ON u.role_id = r.id
		LEFT JOIN %s l ON u.id = l.user_id
		GROUP BY u.id, r.name
		ORDER BY u.last_name`,
		Tables.Users, Tables.Roles, Tables.UserLogins,
	)
	data := []*models.UserData{}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &models.UserData{}
		if err := rows.Scan(
			&tmp.ID, &tmp.Username, &tmp.Email, &tmp.SSO_ID, &tmp.FirstName, &tmp.LastName,
			&tmp.IsActive, &tmp.CreatedAt, &tmp.Role, &tmp.LastVisit,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row. error: %w", err)
		}
		data = append(data, tmp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return data, nil
}

func (r *userRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.UserDataDTO) error {
	if len(dto) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(dto))

	for i, v := range dto {
		if v.ID == uuid.Nil {
			v.ID = uuid.New()
		}

		rows[i] = []interface{}{
			v.ID,
			v.SSO_ID,
			v.RoleID,
			v.Username,
			v.FirstName,
			v.LastName,
			v.Email,
		}
	}

	columns := []string{"id", "sso_id", "role_id", "username", "first_name", "last_name", "email"}
	_, err := r.getExec(tx).CopyFrom(
		ctx,
		pgx.Identifier{Tables.Users},
		columns,
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *userRepo) Update(ctx context.Context, tx Tx, dto *models.UserDataDTO) error {
	query := fmt.Sprintf(`UPDATE %s	SET role_id = $1, is_active = $2
		WHERE id = $3`,
		Tables.Users,
	)

	_, err := r.getExec(tx).Exec(ctx, query, dto.RoleID, dto.IsActive, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *userRepo) UpdateSeveral(ctx context.Context, tx Tx, dto []*models.UserDataDTO) error {
	if len(dto) == 0 {
		return nil
	}

	n := len(dto)
	usernames := make([]string, n)
	emails := make([]string, n)
	ssoIds := make([]string, n)
	roleIds := make([]uuid.UUID, n)
	firstNames := make([]string, n)
	lastNames := make([]string, n)

	for i, v := range dto {
		usernames[i] = v.Username
		emails[i] = v.Email
		ssoIds[i] = v.SSO_ID
		roleIds[i] = v.RoleID
		firstNames[i] = v.FirstName
		lastNames[i] = v.LastName
	}

	query := fmt.Sprintf(`
        UPDATE %s AS t
        SET 
            username = s.username,
            email = s.email,
			role_id = s.role_id,
            first_name = s.first_name,
            last_name = s.last_name
        FROM (
            SELECT * FROM UNNEST(
                $1::text[], 
                $2::text[], 
                $3::uuid[], 
                $4::text[], 
                $5::text[],
				$6::text[]
            ) AS s(username, email, role_id, sso_id, first_name, last_name)
        ) AS s
        WHERE t.sso_id = s.sso_id`,
		Tables.Users,
	)

	_, err := r.getExec(tx).Exec(ctx, query, usernames, emails, roleIds, ssoIds, firstNames, lastNames)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update: %w", err)
	}
	return nil
}

func (r *userRepo) DeleteSeveral(ctx context.Context, tx Tx, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf(`DELETE FROM %s WHERE sso_id=ANY($1)`, Tables.Users)

	if _, err := r.getExec(tx).Exec(ctx, query, ids); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
