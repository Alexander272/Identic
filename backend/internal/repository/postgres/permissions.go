package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/Identic/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PermissionRepo struct {
	db *pgxpool.Pool
	Transaction
}

func NewPermissionRepo(db *pgxpool.Pool, tr Transaction) *PermissionRepo {
	return &PermissionRepo{
		db:          db,
		Transaction: tr,
	}
}

type Permissions interface {
	LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.Permission, error)
	Sync(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error
	GetById(ctx context.Context, id uuid.UUID) (*models.Permission, error)
	GetByRole(ctx context.Context, req *models.GetPermsByRoleDTO) ([]*models.Permission, error)
	Count(ctx context.Context, req *models.GetPermsCountDTO) (*models.PermsCount, error)
	CountForAll(ctx context.Context, roleToDescendants map[string][]string) (map[string]models.PermsCount, error)
	Create(ctx context.Context, tx Tx, dto *models.PermissionDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeletePermissionDTO) error
	DeleteByKeys(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error
}

func (r *PermissionRepo) LoadPolicy(ctx context.Context, req *models.GetPoliciesDTO) ([]*models.Permission, error) {
	query := fmt.Sprintf(`SELECT r.slug,  p.object, p.action
		FROM %s rp
		JOIN %s r ON r.id = rp.role_id
		JOIN %s p ON p.id = rp.permission_id`,
		Tables.RolePermissions, Tables.Roles, Tables.Permissions,
	)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	permissions := make([]*models.Permission, 0, 50)
	for rows.Next() {
		item := &models.Permission{}
		if err := rows.Scan(&item.Role, &item.Object, &item.Action); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		permissions = append(permissions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return permissions, nil
}

func (r *PermissionRepo) Sync(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error {
	if len(dto) == 0 {
		return nil
	}
	values := []string{}
	args := []interface{}{}

	for _, v := range dto {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d)", len(args)+1, len(args)+2, len(args)+3))
		args = append(args, v.Object, v.Action, v.Description)
	}

	query := fmt.Sprintf(`INSERT INTO %s (object, action, description)
			VALUES %s
			ON CONFLICT (object, action) 
			DO UPDATE SET description = EXCLUDED.description`,
		Tables.Permissions, strings.Join(values, ", "),
	)

	_, err := r.getExec(tx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Permission, error) {
	query := fmt.Sprintf(`SELECT id, p.object, p.action
		FROM %s p WHERE id=$1`,
		Tables.Permissions,
	)
	data := &models.Permission{}
	err := r.db.QueryRow(ctx, query, id).Scan(&data.ID, &data.Object, &data.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return data, nil
}

func (r *PermissionRepo) GetByRole(ctx context.Context, req *models.GetPermsByRoleDTO) ([]*models.Permission, error) {
	query := fmt.Sprintf(`SELECT r.slug, p.object, p.action
		FROM %s rp
		JOIN %s r ON r.id = rp.role_id
		JOIN %s p ON p.id = rp.permission_id
		WHERE r.slug = $1`,
		Tables.RolePermissions, Tables.Roles, Tables.Permissions,
	)

	data := make([]*models.Permission, 0, 50)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.Permission{}
		if err := rows.Scan(&item.ID, &item.Role, &item.Object, &item.Action); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		data = append(data, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return data, nil
}

func (r *PermissionRepo) Count(ctx context.Context, req *models.GetPermsCountDTO) (*models.PermsCount, error) {
	query := fmt.Sprintf(`SELECT 
			COUNT(*) FILTER (WHERE role_id = $1) AS own_permissions_count,
			COUNT(DISTINCT permission_id) FILTER (WHERE role_id = ANY($2)) AS inherited_permissions_count
		FROM %s
		WHERE role_id = $1 OR role_id = ANY($2)`,
		Tables.RolePermissions,
	)

	data := &models.PermsCount{}
	err := r.db.QueryRow(ctx, query, req.Role, req.Inherited).Scan(&data.Own, &data.Inherited)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	data.Total = data.Own + data.Inherited

	return data, nil
}
func (r *PermissionRepo) CountForAll(ctx context.Context, roleToDescendants map[string][]string) (map[string]models.PermsCount, error) {
	if len(roleToDescendants) == 0 {
		return make(map[string]models.PermsCount), nil
	}

	res := make(map[string]models.PermsCount)

	// Для каждой роли считаем её собственные permissions
	for roleSlug := range roleToDescendants {
		c := models.PermsCount{}

		// Считаем собственные permissions роли
		ownQuery := fmt.Sprintf(`
			SELECT COUNT(rp.permission_id)
			FROM %s rp
			JOIN %s r ON rp.role_id = r.id
			WHERE r.slug = $1`,
			Tables.RolePermissions, Tables.Roles,
		)

		err := r.db.QueryRow(ctx, ownQuery, roleSlug).Scan(&c.Own)
		if err != nil {
			return nil, fmt.Errorf("failed to count own perms for role %s: %w", roleSlug, err)
		}

		c.Total = c.Own
		res[roleSlug] = c
	}

	// Собираем все уникальные descendant slug'и
	allDescendants := make([]string, 0, len(roleToDescendants))
	descendantSet := make(map[string]struct{})
	for _, descendants := range roleToDescendants {
		for _, d := range descendants {
			if _, exists := descendantSet[d]; !exists {
				descendantSet[d] = struct{}{}
				allDescendants = append(allDescendants, d)
			}
		}
	}

	// Считаем permissions для каждого descendant
	descendantPerms := make(map[string]int)
	if len(allDescendants) > 0 {
		descQuery := fmt.Sprintf(`
			SELECT r.slug, COUNT(rp.permission_id)
			FROM %s rp
			JOIN %s r ON rp.role_id = r.id
			WHERE r.slug = ANY($1)
			GROUP BY r.slug`,
			Tables.RolePermissions, Tables.Roles,
		)

		rows, err := r.db.Query(ctx, descQuery, allDescendants)
		if err != nil {
			return nil, fmt.Errorf("failed to count descendant perms: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var slug string
			var count int
			if err := rows.Scan(&slug, &count); err != nil {
				return nil, err
			}
			descendantPerms[slug] = count
		}
	}

	// Для каждой роли суммируем permissions всех её descendants
	for roleSlug, descendants := range roleToDescendants {
		c := res[roleSlug]
		for _, d := range descendants {
			c.Inherited += descendantPerms[d]
		}
		c.Total = c.Own + c.Inherited
		res[roleSlug] = c
	}

	return res, nil
}

func (r *PermissionRepo) Create(ctx context.Context, tx Tx, dto *models.PermissionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, object, action) VALUES ($1, $2, $3)`,
		Tables.Permissions,
	)
	dto.ID = uuid.New()

	_, err := r.getExec(tx).Exec(ctx, query, dto.ID, dto.Object, dto.Action)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) Delete(ctx context.Context, tx Tx, dto *models.DeletePermissionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, Tables.Permissions)

	_, err := r.getExec(tx).Exec(ctx, query, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) DeleteByKeys(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(dto)*2)
	args := make([]interface{}, 0, len(dto)*2)
	for _, v := range dto {
		placeholders = append(placeholders, fmt.Sprintf("($%d::text, $%d::text)", len(args)+1, len(args)+2))
		args = append(args, v.Object, v.Action)
	}

	// НО проще и надежнее использовать расширение unnest или values:
	query := fmt.Sprintf(`DELETE FROM %s 
        WHERE (object, action) NOT IN (
            SELECT * FROM (VALUES %s) AS t(obj, act)
        )`,
		Tables.Permissions,
		strings.Join(placeholders, ","),
	)

	_, err := r.getExec(tx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
