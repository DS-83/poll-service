package polls

import (
	"context"
	"poll-service/models"
)

type UseCase interface {
	CreatePoll(c context.Context, question string, choices []string) (*models.Poll, error)
	Vote(c context.Context, v *models.Vote) error
	GetResult(c context.Context, id string) (*models.Poll, error)
}
