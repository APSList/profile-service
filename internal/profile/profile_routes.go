package profile

import (
	"hostflow/profile-service/pkg/lib"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ProfileRoutes struct {
	logger            lib.Logger
	router            *lib.Router
	profileController *ProfileController
}

func SetProfileRoutes(
	logger lib.Logger,
	router *lib.Router,
	profileController *ProfileController,
) ProfileRoutes {
	return ProfileRoutes{
		logger:            logger,
		router:            router,
		profileController: profileController,
	}
}

func (route ProfileRoutes) Setup() {
	route.logger.Info("Setting up [PROFILE] routes.")

	users := route.router.Group("/users")
	{
		users.GET("", route.profileController.GetUsersHandler)
		users.GET("/:id", route.profileController.GetUserByIDHandler)
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
