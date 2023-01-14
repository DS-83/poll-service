package models

type Poll struct {
	ID       string   `json:"poll_id"`
	Question string   `json:"poll_question"`
	Choices  []Choice `json:"poll_choices"`
}

type Choice struct {
	ID     string `json:"choice_id"`
	Name   string `json:"choice_name"`
	PollID string `json:"-"`
	Votes  uint   `json:"votes"`
}

type Vote struct {
	PollID   string
	ChoiceID string
}

func NewVote(cID, pID string) *Vote {
	return &Vote{
		PollID:   pID,
		ChoiceID: cID,
	}
}
