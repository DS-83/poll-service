package http

import (
	"net/http"
	"poll-service/models"
	"poll-service/polls"
	"poll-service/utils/logger"

	e "poll-service/err"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase polls.UseCase
	logger  *logger.Logger
}

type createInput struct {
	PollQuestion string   `json:"poll_question"`
	Choices      []string `json:"choices"`
}

type voteInput struct {
	PollID   string `json:"poll_id"`
	ChoiceID string `json:"choice_id"`
}

type getResultInput struct {
	PollID string `json:"poll_id"`
}

type createResponse struct {
	Poll *models.Poll `json:"poll"`
}

type voteResponse struct {
	Resp string `json:"response"`
}

type getResultResponse struct {
	Poll *models.Poll `json:"poll"`
}

func NewHandler(uc polls.UseCase, l *logger.Logger) *Handler {
	return &Handler{
		useCase: uc,
		logger:  l,
	}
}

func (h *Handler) Create(c *gin.Context) {
	input := createInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		h.logger.Error(c, err)
		return
	}

	if len(input.PollQuestion) == 0 || len(input.Choices) == 0 {
		h.logger.Error(c, "incorrect input")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.logger.Info(c, "successfully unmarshal user input")

	poll, err := h.useCase.CreatePoll(c, input.PollQuestion, input.Choices)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		h.logger.Error(c, err)
		return
	}

	resp := createResponse{
		Poll: poll,
	}
	h.logger.Info(c, "successfuly created poll")
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) Vote(c *gin.Context) {
	input := voteInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		h.logger.Error(c, err)
		return
	}

	if len(input.ChoiceID) == 0 || len(input.PollID) == 0 {
		h.logger.Error(c, "incorrect input")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.logger.Info(c, "successfully unmarshal user input")

	vote := models.NewVote(input.ChoiceID, input.PollID)

	if err := h.useCase.Vote(c, vote); err != nil {
		if err == e.ErrNoDocument {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		h.logger.Error(c, err)
		return
	}

	h.logger.Info(c, "successfylly voted:", vote)

	resp := voteResponse{
		Resp: "vote accepted",
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetResult(c *gin.Context) {
	input := getResultInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		h.logger.Error(c, err)
		return
	}

	if len(input.PollID) == 0 {
		h.logger.Error(c, "incorrect input")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	h.logger.Info(c, "successfully unmarshal user input")

	res, err := h.useCase.GetResult(c, input.PollID)
	if err != nil {
		if err == e.ErrNoDocument {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		h.logger.Error(c, err)
		return
	}

	resp := getResultResponse{
		Poll: res,
	}

	c.JSON(http.StatusOK, resp)
}
