package profile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileController struct {
	service *ProfileService
}

func GetProfileController(service *ProfileService) *ProfileController {
	return &ProfileController{
		service: service,
	}
}

func (route ProfileRoutes) LivenessHandler(c *gin.Context) {
	// If this handler runs, the process is alive
	c.JSON(200, gin.H{
		"status": "alive",
	})
}

func (route ProfileRoutes) ReadinessHandler(c *gin.Context) {
	// Put real readiness checks here if you have them
	// e.g. database ping, message broker connection, etc.

	// Example:
	// if err := route.db.PingContext(c); err != nil {
	//     c.JSON(503, gin.H{"status": "not ready"})
	//     return
	// }

	c.JSON(200, gin.H{
		"status": "ready",
	})
}

// GetUsersHandler godoc
// @Summary Get all users
// @Description Returns a list of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} profile.User
// @Failure 500 {object} profile.ErrorResponse
// @Router /users [get]
func (c *ProfileController) GetUsersHandler(ctx *gin.Context) {
	users, err := c.service.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch users",
			Message: err.Error(),
		})
		return
	}

	// Return directly
	ctx.JSON(http.StatusOK, users)
}

// GetUserByIDHandler godoc
// @Summary Get a user by ID
// @Description Returns a single user by UUID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} profile.User
// @Failure 400 {object} profile.ErrorResponse
// @Failure 404 {object} profile.ErrorResponse
// @Failure 500 {object} profile.ErrorResponse
// @Router /users/{id} [get]
func (c *ProfileController) GetUserByIDHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid ID format",
			Message: "ID must be a valid UUID",
		})
		return
	}

	user, err := c.service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch user",
			Message: err.Error(),
		})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: "No user exists with the given ID",
		})
		return
	}

	// Return directly
	ctx.JSON(http.StatusOK, user)
}
