package profile

import (
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
	if u.OrganizationID == uuid.Nil {
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
