package handlers

import (
	"log"
	"project/internal/auth"
	"project/internal/middlewear"
	"project/internal/services"

	"github.com/gin-gonic/gin"
)

func Api(a *auth.Auth, s services.Serviceinterface) *gin.Engine {
	r := gin.New()
	m, _ := middlewear.NewMiddleWear(a)
	h, err := NewHandler(a, s)
	if err != nil {
		log.Panic("service is not setup")
		return nil
	}
	r.Use(m.Log(), gin.Recovery())
	r.POST("/signup", h.userSignin)
	r.POST("/login", h.userLogin)

	r.POST("/createCompany", m.Auth(h.companyCreation))
	r.GET("/getAllCompany", m.Auth(h.getAllCompany))
	r.GET("/getCompany/:company_id", m.Auth(h.getCompany))

	r.POST("/api/companies/:company_id/jobs", m.Auth(h.postJobByCompany))
	r.GET("/companies/:company_id/jobs", m.Auth(h.getJob))
	r.GET("/jobs", m.Auth(h.getAllJob))
	r.POST("/api/applications", m.Auth(h.processApplications))

	r.POST("/api/process/application", h.ProcessApplication)

	r.POST("/api/forgetpassword", m.Auth(h.forgetpassword))
	r.POST("/requestNewPassword", h.requestNewPassword)
	return r
}
