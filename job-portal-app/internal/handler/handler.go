package handler

import (
	"github.com/gin-gonic/gin"
	"job-portal-api/internal/service"
)

type handler struct {
	service service.UserService
}

type Handler interface {
	UserLogin(c *gin.Context)
	RegisterUser(c *gin.Context)

	GetCompany(c *gin.Context)
	ListCompanies(c *gin.Context)
	CreateCompany(c *gin.Context)

	GetJobPostingByID(c *gin.Context)
	GetAllJobPostings(c *gin.Context)
	ListJobsForCompany(c *gin.Context)
	CreateJobPosting(c *gin.Context)
	ProcessJobApplication(c *gin.Context)
	ForgotPasswordHandler(c *gin.Context)
	UpdatePasswordHandler(c *gin.Context)
}

func NewHandler(svc service.UserService) (Handler, error) {
	return &handler{
		service: svc,
	}, nil
}
