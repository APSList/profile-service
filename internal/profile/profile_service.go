package profile

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type ProfileService struct {
	repo *ProfileRepository
}

func GetProfileService(repo *ProfileRepository) *ProfileService {
	return &ProfileService{
		repo: repo,
	}
}

// GetUsers returns all users
func (s *ProfileService) GetUsers() ([]User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUsers returns all users belonging to the requester's organization.
// It uses the organizationID extracted from the OIDC/JWT token.
func (s *ProfileService) GetUsersProtected(ctx context.Context, organizationID int64) ([]User, error) {
	// We call the specific repository method that filters by Org
	users, err := s.repo.GetUsersByOrganizationID(ctx, organizationID)
	if err != nil {
		// You can add custom logging or error wrapping here
		return nil, err
	}

	// Logic check: If no users are found, pgx returns an empty slice, not an error.
	if users == nil {
		return []User{}, nil
	}

	return users, nil
}

// GetUserByID returns a user by ID
func (s *ProfileService) GetUserByID(id uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// GetUserByOrganizationID returns a user by org ID
func (s *ProfileService) GetUserByOrganizationID(orgID uuid.UUID) (*User, error) {
	user, err := s.repo.GetUserByOrganizationID(orgID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// CreateUser creates a new user
func (s *ProfileService) CreateUser(u *User) (*User, error) {
	// optional: basic validation
	if u.OrganizationID == 0 {
		return nil, errors.New("organization ID is required")
	}
	if u.Email == "" {
		return nil, errors.New("email is required")
	}
	if u.Name == "" {
		return nil, errors.New("name is required")
	}

	created, err := s.repo.CreateUser(u)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *ProfileService) DeactivateUser(ctx context.Context, targetID string, adminID string, orgID int64, adminRole string) error {
	// 1. Authorization check
	if adminRole != "OWNER" {
		return errors.New("unauthorized: only owners can deactivate members")
	}

	// 2. Prevent self-deactivation (Safety)
	if targetID == adminID {
		return errors.New("cannot deactivate your own account")
	}

	// 3. Execute update
	return s.repo.UpdateStatus(ctx, targetID, orgID, "INACTIVE")
}

func (s *ProfileService) GetOrganizationName(ctx context.Context, orgID int64) (string, error) {
	return s.repo.GetNameByID(ctx, orgID)
}
