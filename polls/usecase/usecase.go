package usecase

import (
	"context"
	"poll-service/models"
	"poll-service/polls"
	"poll-service/utils/logger"
)

type PollUsecase struct {
	PollRepo polls.PollRepo
	Logger   *logger.Logger
}

func NewPollUsecase(db polls.PollRepo, l *logger.Logger) *PollUsecase {
	return &PollUsecase{
		PollRepo: db,
		Logger:   l,
	}
}

func (uc *PollUsecase) CreatePoll(c context.Context, question string, choices []string) (*models.Poll, error) {
	id, err := uc.PollRepo.CreatePoll(c, question, choices)
	if err != nil {
		uc.Logger.Error(c, err)
		return nil, err
	}
	uc.Logger.Debug(c, "created poll with id:", id)

	uc.Logger.Info(c, "get poll from DB")
	poll, err := uc.PollRepo.GetPoll(c, id)
	if err != nil {
		uc.Logger.Error(c, err)
		return nil, err
	}
	uc.Logger.Info(c, "finish create poll")

	return poll, nil
}

func (uc *PollUsecase) Vote(c context.Context, v *models.Vote) error {
	if err := uc.PollRepo.Vote(c, v); err != nil {
		return err
	}
	return nil
}

func (uc *PollUsecase) GetResult(c context.Context, id string) (*models.Poll, error) {
	poll, err := uc.PollRepo.GetPoll(c, id)
	if err != nil {
		return nil, err
	}
	return poll, nil
}
