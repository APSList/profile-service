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
