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
	// h := handler{a: a, us: s}
	m, err := middlewear.NewMiddleWear(a)
	if err != nil {
		log.Panic("auth err")
		return nil
	}

	h, err := NewHandler(a, s)
	if err != nil {
		log.Panic("service is not setup")
		return nil
	}

	r.Use(m.Log(), gin.Recovery())

	r.POST("/signup", h.userSignin)
	r.POST("/login", h.userLoginin)
	r.POST("/createCompany", m.Auth(h.companyCreation))
	r.GET("/getAllCompany", m.Auth(h.getAllCompany))
	r.GET("/getCompany/:company_id", m.Auth(h.getCompany))
	r.POST("/companies/:company_id/jobs", m.Auth(h.postJob))
	r.GET("/companies/:company_id/jobs", m.Auth(h.getJob))
	r.GET("/jobs", m.Auth(h.getAllJob))

	return r
}
