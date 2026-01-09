package profile

import (
	"context"
	"errors"
	_ "time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func GetProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

// GetUsersByOrganizationID returns all users belonging to a specific organization.
func (r *ProfileRepository) GetUsersByOrganizationID(ctx context.Context, organizationId int64) ([]User, error) {
	query := `
        SELECT id, organization_id, full_name, role, email, status, created_at, updated_at
        FROM "profiles"
        WHERE organization_id = $1
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(ctx, query, organizationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// pgx.CollectRows handles scanning the entire result set into a slice of structs
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ProfileRepository) GetUsers() ([]User, error) {
	query := `
        SELECT *
        FROM "user"
        ORDER BY created_at DESC
    `

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ProfileRepository) UpdateStatus(ctx context.Context, userID string, orgID int64, status string) error {
	query := `UPDATE "profiles" SET status = $1 WHERE id = $2 AND organization_id = $3`

	result, err := r.db.Exec(ctx, query, status, userID, orgID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("no user found or unauthorized")
	}
	return nil
}

func (r *ProfileRepository) GetNameByID(ctx context.Context, orgID int64) (string, error) {
	var name string
	query := `SELECT name FROM organization WHERE id = $1`

	err := r.db.QueryRow(ctx, query, orgID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (r *ProfileRepository) GetUserByID(id uuid.UUID) (*User, error) {
	query := `
        SELECT id, organization_id, name, role, email, status, created_at, updated_at
        FROM "user"
        WHERE id = $1
    `

	rows, err := r.db.Query(context.Background(), query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) GetUserByOrganizationID(organizationId uuid.UUID) (*User, error) {
	query := `
        SELECT id, organization_id, name, role, email, status, created_at, updated_at
        FROM "user"
        WHERE organization_id = $1
    `

	rows, err := r.db.Query(context.Background(), query, organizationId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *ProfileRepository) CreateUser(u *User) (*User, error) {
	query := `
        INSERT INTO "user" (
            id, organization_id, name, role, email, status, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, organization_id, name, role, email, status, created_at, updated_at
    `

	rows, err := r.db.Query(
		context.Background(),
		query,
		u.ID,
		u.OrganizationID,
		u.Name,
		u.Role,
		u.Email,
		u.Status,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		return nil, err
	}

	return &created, nil
}
