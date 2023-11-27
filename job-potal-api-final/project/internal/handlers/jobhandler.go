package handlers

import (
	"encoding/json"
	"net/http"

	//"project/internal/auth"
	"project/internal/middlewear"
	"project/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func (h *handler) companyCreation(c *gin.Context) {

	ctx := c.Request.Context()

	traceId, ok := ctx.Value(middlewear.TraceIdKey).(string)

	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id not found in userSignin handler")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var companyCreation model.CreateCompany
	body := c.Request.Body
	err := json.NewDecoder(body).Decode(&companyCreation)
	if err != nil {
		log.Error().Err(err).Msg("error in decoding")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	validate := validator.New()
	err = validate.Struct(&companyCreation)
	if err != nil {
		log.Error().Err(err).Msg("error in validating ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "invalid input"})
		return
	}

	us, err := h.us.CompanyCreate(companyCreation)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user signup problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "user signup failed"})
		return
	}
	c.JSON(http.StatusOK, us)

}

func (h *handler) getAllCompany(c *gin.Context) {
	ctx := c.Request.Context()

	traceId, ok := ctx.Value(middlewear.TraceIdKey).(string)

	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id not found in userSignin handler")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	us, err := h.us.GetAllCompanies()
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user signup problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "getallcomapy failed"})
		return
	}
	c.JSON(http.StatusOK, us)

}

func (h *handler) getCompany(c *gin.Context) {

	ctx := c.Request.Context()
	traceId, _ := ctx.Value(middlewear.TraceIdKey).(string)
	id, erro := strconv.Atoi(c.Param("company_id"))
	if erro != nil {
		log.Error().Str("traceId", traceId).Msg("trace id not found in userSignin handler")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": http.StatusText(http.StatusBadRequest)})
		return
	}

	us, err := h.us.GetCompany(id)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("getcompany problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "user getcompany failed"})
		return
	}
	c.JSON(http.StatusOK, us)
}

func (h *handler) postJobByCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewear.TraceIdKey).(string)
	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id not found in  handler")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	id, erro := strconv.ParseUint(c.Param("company_id"), 10, 32)
	if erro != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	var jobCreation model.NewJobRequest
	body := c.Request.Body
	err := json.NewDecoder(body).Decode(&jobCreation)
	if err != nil {
		log.Error().Err(err).Msg("error in decoding")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	validate := validator.New()
	err = validate.Struct(&jobCreation)
	if err != nil {
		log.Error().Err(err).Msg("error in validating ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	us, err := h.us.JobCreate(jobCreation, id)

	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("job creatuion problem in db")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, us)

}

func (h *handler) getJob(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewear.TraceIdKey).(string)
	id, erro := strconv.Atoi(c.Param("company_id"))
	if erro != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": http.StatusText(http.StatusBadRequest)})
		return
	}

	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id not found in userSignin handler")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	us, err := h.us.GetJobs(id)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in getjob")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in getjob"})
		return
	}
	c.JSON(http.StatusOK, us)

}

func (h *handler) getAllJob(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewear.TraceIdKey).(string)
	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id not found in userSignin handler")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	us, err := h.us.GetAllJobs()
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("problem in getalljob")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in getalljob"})
		return
	}
	c.JSON(http.StatusOK, us)

}

func (h *handler) processApplications(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middlewear.TraceIdKey).(string)
	if !ok {
		log.Error().Str("traceId", traceId).Msg("trace id not found ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var appData []model.NewUserApplication

	err := json.NewDecoder(c.Request.Body).Decode(&appData)

	if err != nil {
		log.Error().Err(err).Str("trace id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "please provide proper data"})
		return
	}

	a, err := h.us.ProcessJobApplications(appData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, a)

}
