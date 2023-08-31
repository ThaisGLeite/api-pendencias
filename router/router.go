package router

import (
	"api-pendencias/database"
	"api-pendencias/pendencia"
	"api-pendencias/utils"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/didip/tollbooth_gin"
	"github.com/gin-gonic/gin"
)

// Define constants for logging levels, Gin mode, and server port.
const (
	InfoLogLevel   = "I"
	ErrorLogLevel  = "E"
	GinModeKey     = "GIN_MODE"
	DefaultGinMode = gin.DebugMode
	PortKey        = "PORT"
	DefaultPort    = "8080"
)

// StartServer function initializes and starts the web server.
func StartServer(database *database.Connection) *gin.Engine {

	// Set the gin mode. This can be either debug or release.
	gin.SetMode(DefaultGinMode)

	// Create a new gin engine. Use gin.New() to have more control over the middleware.
	server := gin.New()

	// Use the Recovery middleware to recover from any panics and write a 500 if it happens.
	server.Use(gin.Recovery())

	// Secure headers middleware
	server.Use(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'")
		// Other headers as needed
		c.Next()
	})

	// Create a new rate limiter. This will limit to 1 request/second.
	limiter := tollbooth.NewLimiter(1, nil)

	// Setup the routes for the server.
	setupRoutes(server, limiter, database)

	// Run the server on the defined port. If there is an error, log it.
	if err := server.Run(":" + PortKey); err != nil {
		utils.HandleError(ErrorLogLevel, "Failed to run server", err)
	}

	return server
}

// setupRoutes function sets up all the routes for the server.
func setupRoutes(router *gin.Engine, limiter *limiter.Limiter, db *database.Connection) {
	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	pendencia := &pendencia.Database{
		Database: db,
	}

	// User routes. These routes are wrapped with a rate limiter middleware.
	router.GET("/pendencia", tollbooth_gin.LimitHandler(limiter), pendencia.GetAllPendencias)      //Return all pendencias
	router.GET("/pendencia/:nome", tollbooth_gin.LimitHandler(limiter), pendencia.GetPendencia)    // Get a pendencia por nome
	router.POST("/pendencia", tollbooth_gin.LimitHandler(limiter), pendencia.CreatePendencia)      // Create a new pendencia
	router.PUT("/pendencia/:nome", tollbooth_gin.LimitHandler(limiter), pendencia.UpdatePendencia) // Update a pendencia
}
