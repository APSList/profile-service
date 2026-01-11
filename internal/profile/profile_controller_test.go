package profile

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) GetUsersProtected(ctx context.Context, orgID int64) ([]User, error) {
	args := m.Called(ctx, orgID)
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockProfileService) GetOrganizationName(ctx context.Context, orgID int64) (string, error) {
	args := m.Called(ctx, orgID)
	return args.String(0), args.Error(1)
}

func (m *MockProfileService) GetUserByID(id uuid.UUID) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockProfileService) DeactivateUser(ctx context.Context, tID, aID string, oID int64, r string) error {
	return m.Called(ctx, tID, aID, oID, r).Error(0)
}

// --- TESTI ---

func TestGetUsersHandler_OwnerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProfileService)
	controller := GetProfileController(mockSvc)

	r := gin.Default()
	r.GET("/users", func(c *gin.Context) {
		c.Set("organization_id", int64(1))
		c.Set("role", "OWNER") // Simuliramo OWNER dostop
		controller.GetUsersHandler(c)
	})

	mockSvc.On("GetUsersProtected", mock.Anything, int64(1)).Return([]User{{Name: "Leon"}}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Leon")
}

func TestGetUsersHandler_ForbiddenForMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(MockProfileService)
	controller := GetProfileController(mockSvc)

	r := gin.Default()
	r.GET("/users", func(c *gin.Context) {
		c.Set("organization_id", int64(1))
		c.Set("role", "MEMBER") // Simuliramo navadnega člana
		controller.GetUsersHandler(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "Only organization owners can view")
}

func TestGetUserByIDHandler_InvalidUUID(t *testing.T) {
	controller := GetProfileController(new(MockProfileService))
	r := gin.Default()
	r.GET("/users/:id", controller.GetUserByIDHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/123-ni-uuid", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID must be a valid UUID")
}
