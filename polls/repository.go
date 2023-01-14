package polls

import (
	"context"
	"poll-service/models"
)

type PollRepo interface {
	CreatePoll(c context.Context, name string, choices []string) (string, error)
	Vote(c context.Context, v *models.Vote) error
	GetPoll(c context.Context, id string) (*models.Poll, error)
}
