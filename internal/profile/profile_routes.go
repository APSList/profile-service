package profile

import (
	"hostflow/profile-service/internal/middlewares"
	"hostflow/profile-service/pkg/lib"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ProfileRoutes struct {
	logger            lib.Logger
	router            *lib.Router
	profileController *ProfileController
	authMiddleware    middlewares.AuthMiddleware
}

func SetProfileRoutes(
	logger lib.Logger,
	router *lib.Router,
	profileController *ProfileController,
	authMiddleware middlewares.AuthMiddleware,
) ProfileRoutes {
	return ProfileRoutes{
		logger:            logger,
		router:            router,
		profileController: profileController,
		authMiddleware:    authMiddleware,
	}
}

func (route ProfileRoutes) Setup() {
	route.logger.Info("Setting up [PROFILE] routes.")

	users := route.router.Group("/users")
	users.Use(route.authMiddleware.Handler())
	{
		users.GET("", route.profileController.GetUsersHandler)
		users.GET("/:id", route.profileController.GetUserByIDHandler)
		users.PUT("/:id/status", route.profileController.DeactivateHandler)
	}

	organizations := route.router.Group("/organization")
	organizations.Use(route.authMiddleware.Handler())
	{
		organizations.GET("/name", route.profileController.GetOrgNameHandler)
	}

	metrics := route.router.Group("/metrics")
	{
		metrics.GET("", gin.WrapH(promhttp.Handler()))
	}

	health := route.router.Group("/health")
	{
		health.GET("/live", route.LivenessHandler)
		health.GET("/ready", route.ReadinessHandler)
	}

	// Swagger documentation
	route.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.logger.Info("Swagger documentation available at: /swagger/index.html")

	route.logger.Info("[PROFILE] routes setup complete.")
}
