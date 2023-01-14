package http

import (
	"poll-service/polls"
	"poll-service/utils"

	"github.com/gin-gonic/gin"
)

const correlationID string = "Correlation-ID"

type PollsMiddleware struct {
	uc polls.UseCase
}

func NewPollsMiddleware(uc polls.UseCase) gin.HandlerFunc {
	return (&PollsMiddleware{
		uc: uc,
	}).HandleCorrelation
}

func (m *PollsMiddleware) HandleCorrelation(c *gin.Context) {
	id := utils.CreateId()
	c.Set(correlationID, id)
}
