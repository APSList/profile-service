package profile

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProfileController struct {
	service Service
}

func GetProfileController(service Service) *ProfileController {
	return &ProfileController{
		service: service,
	}
}

// LivenessHandler godoc
// @Summary Liveness probe
// @Description Check if the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/liveness [get]
func (route ProfileRoutes) LivenessHandler(c *gin.Context) {
	// If this handler runs, the process is alive
	c.JSON(200, gin.H{
		"status": "alive",
	})
}

// ReadinessHandler godoc
// @Summary Readiness probe
// @Description Check if the service is ready to handle requests
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health/readiness [get]
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

/*// GetUsersHandler godoc
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
}*/

// GetUsersHandler godoc
// @Summary Get organization users
// @Description Returns a list of users belonging to the requester's organization. Requires OWNER role.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} User
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (c *ProfileController) GetUsersHandler(ctx *gin.Context) {
	// 1. Extract claims set by your Auth Middleware
	orgID, orgExists := ctx.Get("organization_id")
	role, roleExists := ctx.Get("role")

	if !orgExists || !roleExists {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "Missing authentication claims",
		})
		return
	}

	// 2. Access Management: Only allow OWNERS to fetch the full list
	if role.(string) != "OWNER" {
		ctx.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "Forbidden",
			Message: "Only organization owners can view the user list",
		})
		return
	}

	// 3. Call the service with the specific Organization ID
	// Note: ensure your service method signature accepts (ctx, orgID)
	users, err := c.service.GetUsersProtected(ctx.Request.Context(), orgID.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch users",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

// DeactivateHandler godoc
// @Summary Deactivate a user
// @Description Deactivates a user account within the organization. Requires OWNER role.
// @Tags users
// @Security ApiKeyAuth
// @Param id path string true "User ID to deactivate"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/{id}/deactivate [delete]
func (c *ProfileController) DeactivateHandler(ctx *gin.Context) {
	fmt.Println("Deactivate handler called")
	targetID := ctx.Param("id")

	// These values are set by your Auth Middleware from the Supabase JWT
	adminID := ctx.GetString("user_id")
	orgID := ctx.GetInt64("organization_id")
	role := ctx.GetString("role")

	err := c.service.DeactivateUser(ctx, targetID, adminID, orgID, role)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(204)
}

// GetOrgNameHandler godoc
// @Summary Get organization name
// @Description Returns the name of the organization associated with the current user's token
// @Tags organization
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string "Example: {"name": "My Org"}"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /org/name [get]
func (c *ProfileController) GetOrgNameHandler(ctx *gin.Context) {
	// 1. Get the organization_id that the Middleware extracted from the JWT
	// Use ctx.Get() because it was stored there by c.Set("organization_id", value)
	val, exists := ctx.Get("organization_id")

	if !exists {
		ctx.JSON(401, gin.H{"error": "Organization ID not found in session"})
		return
	}

	// 2. Type-assert the value to int64
	orgID, ok := val.(int64)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Internal server error: ID format mismatch"})
		return
	}

	// 3. Call your service using the ID from the token
	name, err := c.service.GetOrganizationName(ctx, orgID)
	if err != nil {
		ctx.JSON(404, gin.H{"error": "Organization not found"})
		return
	}

	ctx.JSON(200, gin.H{"name": name})
}

// GetUserByIDHandler godoc
// @Summary Get a user by ID
// @Description Returns a single user's profile information by their UUID
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
