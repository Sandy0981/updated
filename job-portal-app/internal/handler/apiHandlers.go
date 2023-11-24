package handler

import (
	"log"

	"github.com/gin-gonic/gin"           // Importing the Gin framework for handling HTTP requests and responses.
	"job-portal-api/internal/auth"       // Importing custom authentication package.
	"job-portal-api/internal/middleware" // Importing custom middleware package.
	"job-portal-api/internal/service"    // Importing custom service package for business logic.
)

// SetupApi is a function that sets up the API routes, middleware, and handlers.
func SetupApi(a auth.Authentication, svc service.UserService) *gin.Engine {
	r := gin.New() // Creating a new Gin engine.

	m, err := middleware.NewMiddleware(a)
	// Creating a new instance of the middleware with the provided authentication.
	if err != nil {
		log.Panic("Error setting up middleware")
		// Logging an error if middleware setup fails.
	}

	h, err := NewHandler(svc)
	// Creating a new instance of the handler with the provided user service.
	if err != nil {
		log.Panic("Error setting up handler")
		// Logging an error if handler setup fails.
	}
	r.Use(m.Log(), gin.Recovery())
	r.GET("/check", m.Authenticate(Check))
	r.POST("/api/register", h.RegisterUser)
	r.POST("/api/login", h.UserLogin)
	r.POST("/api/companies", m.Authenticate(h.CreateCompany))
	r.GET("/api/companies", m.Authenticate(h.ListCompanies))
	r.GET("/api/companies/:companyID", m.Authenticate(h.GetCompany))
	r.POST("/api/companies/:companyID/jobs", m.Authenticate(h.CreateJobPosting))
	r.GET("/api/companies/:companyID/jobs", m.Authenticate(h.ListJobsForCompany))
	r.GET("/api/jobs/:jobID", m.Authenticate(h.GetJobPostingByID))
	r.GET("/api/jobs", m.Authenticate(h.GetAllJobPostings))
	r.POST("/api/process", h.ProcessJobApplication)
	r.POST("/api/forget-password", h.ForgotPasswordHandler)
	return r
	// Returning the configured Gin engine.
}

// Check is a simple handler that responds with a JSON message "ok" for a GET request at "/check".
func Check(c *gin.Context) {
	c.JSON(200, gin.H{
		"Message": "ok",
	})
}
