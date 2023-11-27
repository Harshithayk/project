package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/internal/middlewear"
	"project/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *handler) ProcessApplication(c *gin.Context) {
	ctx := c.Request.Context()
	traceid, ok := ctx.Value(middlewear.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceid is not found")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	//_, ok = ctx.Value(auth.Key).(jwt.RegisteredClaims)
	// if !ok {
	// 	log.Error().Str("Trace Id", traceid).Msg("login first")
	// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
	// 	return
	// }
	var applicationData []model.NewUserApplication
	err := json.NewDecoder(c.Request.Body).Decode(&applicationData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "please provide valid details",
		})
		return
	}
	applicationData, err = h.us.ProccessApplication(ctx, applicationData)
	fmt.Println("[][][]", applicationData)
	if err != nil {
		log.Error().Err(err).Str("trace id", traceid)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, applicationData)

}
