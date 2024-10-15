package main

import (
	"ems/api/routes"
	"ems/infrastructure/config"
	"ems/infrastructure/database"
	"ems/scheduler"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/**
 * @Function: main
 * @Description: Entry point of the application. It loads the configuration and sets up the server.
 *
 * @Params:
 *    - None
 *
 * @Returns:
 *    - None
 */
func main() {
	loadConfiguration()
	setupServer()
}

/**
 * @Function: loadConfiguration
 * @Description: Loads the application configuration from .env file.
 *
 * @Params:
 *    - None
 *
 * @Returns:
 *    - None
 */
func loadConfiguration() {
	// Load configuration
	err := config.Load()
	if err != nil {
		panic(err)
	}
}

/**
 * @Function: setupRateLimiter
 * @Description: Configures rate limiting for the given routes.
 *
 * @Params:
 *    - routes: The Gin Engine where the rate limiter will be applied.
 *
 * @Returns:
 *    - None
 */
func setupRateLimiter(router *gin.Engine) {
	// Set up rate limiter to 5 requests per second per client
	limiter := tollbooth.NewLimiter(5, nil)
	router.Use(tollbooth_gin.LimitHandler(limiter))
}

/**
 * @Function: setupCors
 * @Description: Configures CORS settings for the given routes.
 *
 * @Params:
 *    - routes: The Gin Engine where the CORS configuration will be applied.
 *
 * @Returns:
 *    - None
 */
func setupCors(router *gin.Engine) {
	// Set up CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	router.Use(cors.New(corsConfig))
}

/**
 * @Function: setupRoutes
 * @Description: Sets up application routes for the given routes.
 *
 * @Params:
 *    - routes: The Gin Engine where the routes will be configured.
 *
 * @Returns:
 *    - None
 */
func setupRoutes(router *gin.Engine, db *gorm.DB) {
	routes.SetupRoutes(router, db)
}

/**
 * @Function: setupServer
 * @Description: Initializes the Gin routes and configures rate limiting, CORS, and routes.
 *
 * @Params:
 *    - None
 *
 * @Returns:
 *    - None
 */
func setupServer() {
	//Initialize Database
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	//Initialize Schedular
	scheduler := scheduler.Scheduler{DB: db}
	scheduler.InitScheduler()

	// Initialize Gin routes
	router := gin.Default()

	// Setup various server configurations
	setupRateLimiter(router)
	setupCors(router)
	setupRoutes(router, db)

	// Start the server
	err = router.Run(config.Config.Port)
	if err != nil {
		panic(err)
	}
}
