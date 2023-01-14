package http

import (
	"poll-service/polls"
	"poll-service/utils/logger"

	"github.com/gin-gonic/gin"
)

func RegisterMidRoutes(r *gin.RouterGroup, uc polls.UseCase, l *logger.Logger) {
	h := NewHandler(uc, l)

	r.POST("createPoll", h.Create)
	r.POST("poll", h.Vote)
	r.POST("getResult", h.GetResult)
}
