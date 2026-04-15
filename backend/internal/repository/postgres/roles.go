package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewRoleRepo(db *pgxpool.Pool, tr Transaction) *RoleRepo {
	return &RoleRepo{
		db:          db,
		Transaction: tr,
	}
}

type Roles interface {
	GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error)
	GetAll(ctx context.Context) ([]*models.Role, error)
	GetUserCount(ctx context.Context, req []string) (map[string]int, error)
	IsExists(ctx context.Context, roleName string) (bool, error)
	IsExistsById(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, tx Tx, dto *models.RoleDTO) error
	Update(ctx context.Context, tx Tx, dto *models.RoleDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteRoleDTO) error

	AssignPermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error
	DeletePermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error
}

func (r *RoleRepo) GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error) {
	condition := ""
	params := []interface{}{}
	if req.ID != uuid.Nil {
		params = append(params, req.ID)
		condition = fmt.Sprintf("WHERE id = $%d", len(params))
	}
	if req.Slug != "" {
		params = append(params, req.Slug)
		condition = fmt.Sprintf("WHERE slug = $%d", len(params))
	}
	if condition == "" {
		return nil, models.ErrInvalidInput
	}

	query := fmt.Sprintf(`SELECT id, slug, name, level, is_system, created_at, updated_at FROM %s %s`,
		Tables.Roles, condition,
	)
	data := &models.Role{}

	err := r.db.QueryRow(ctx, query, params...).Scan(
		&data.ID,
		&data.Slug,
		&data.Name,
		&data.Level,
		&data.IsSystem,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return data, nil
}

func (r *RoleRepo) GetAll(ctx context.Context) ([]*models.Role, error) {
	query := fmt.Sprintf(`SELECT id, slug, name, level, is_active, is_system, is_editable, created_at, updated_at FROM %s 
		ORDER BY level`,
		Tables.Roles,
	)
	data := []*models.Role{}

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.Role{}
		if err := rows.Scan(
			&item.ID, &item.Slug, &item.Name, &item.Level, &item.IsActive, &item.IsSystem, &item.IsEditable,
			&item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		data = append(data, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return data, nil
}

func (r *RoleRepo) GetUserCount(ctx context.Context, req []string) (map[string]int, error) {
	// Если входящий список пуст, сразу возвращаем пустую мапу
	if len(req) == 0 {
		return make(map[string]int), nil
	}

	query := fmt.Sprintf(`SELECT role_id, COUNT(*) FROM %s 
		WHERE role_id = ANY($1)
		GROUP BY role_id`,
		Tables.Users,
	)

	// Передаем слайс IDs (в зависимости от драйвера может понадобиться pq.Array)
	rows, err := r.db.Query(ctx, query, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var roleID string
		var count int
		if err := rows.Scan(&roleID, &count); err != nil {
			return nil, fmt.Errorf("scan count error: %w", err)
		}
		counts[roleID] = count
	}

	return counts, nil
}

func (r *RoleRepo) IsExists(ctx context.Context, roleName string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE name = $1)`, Tables.Roles)
	var exists bool

	err := r.db.QueryRow(ctx, query, roleName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return exists, nil
}
func (r *RoleRepo) IsExistsById(ctx context.Context, id uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1 AND is_active = true)`, Tables.Roles)
	var exists bool

	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return exists, nil
}

func (r *RoleRepo) Create(ctx context.Context, tx Tx, dto *models.RoleDTO) error {
	if dto.Slug == "root" || dto.Slug == "superadmin" {
		return models.ErrReservedRole
	}

	query := fmt.Sprintf(`INSERT INTO %s (slug, name, level, is_system)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		Tables.Roles,
	)

	err := r.getExec(tx).QueryRow(
		ctx, query, dto.Slug, dto.Name, dto.Level, dto.IsSystem,
	).Scan(&dto.ID, &dto.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleRepo) Update(ctx context.Context, tx Tx, dto *models.RoleDTO) error {
	if dto.Slug == "root" || dto.Slug == "superadmin" {
		return models.ErrReservedRole
	}

	query := fmt.Sprintf(`UPDATE %s SET name=$1, level=$2, slug=$3, is_system=$4, updated_at=NOW() WHERE id=$5`,
		Tables.Roles,
	)

	_, err := r.getExec(tx).Exec(ctx, query, dto.Name, dto.Level, dto.Slug, dto.IsSystem, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteRoleDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1 AND NOT is_system`, Tables.Roles)

	_, err := r.getExec(tx).Exec(ctx, query, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleRepo) AssignPermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (role_id, permission_id) VALUES ($1, $2)`, Tables.RolePermissions)

	_, err := r.getExec(tx).Exec(ctx, query, dto.RoleID, dto.PermissionID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleRepo) DeletePermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id=$1 AND permission_id=$2`, Tables.RolePermissions)

	_, err := r.getExec(tx).Exec(ctx, query, dto.RoleID, dto.PermissionID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
