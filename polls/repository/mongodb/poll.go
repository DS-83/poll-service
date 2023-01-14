package mongodb

import (
	"context"

	"poll-service/utils/logger"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	e "poll-service/err"
	"poll-service/models"
)

const talblePolls = "polls"

type PollRepo struct {
	db     *mongo.Database
	logger *logger.Logger
}

type poll struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Question string             `bson:"question"`
	Choices  []choice           `bson:"choice"`
}

type choice struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Votes uint               `bson:"votes"`
}

type vote struct {
	PollID   primitive.ObjectID
	ChoiceID primitive.ObjectID
}

func NewPollRepo(db *mongo.Database, l *logger.Logger) *PollRepo {
	return &PollRepo{
		db:     db,
		logger: l,
	}
}

func (r *PollRepo) CreatePoll(ctx context.Context, question string, choiceNames []string) (string, error) {
	cur := r.db.Collection(talblePolls)

	choices := []choice{}
	for _, name := range choiceNames {
		choice := newChoice(name)
		choices = append(choices, *choice)
	}

	poll := poll{
		Question: question,
		Choices:  choices,
	}

	res, err := cur.InsertOne(ctx, poll)
	if err != nil {
		r.logger.Error(ctx, err)
		return "", err
	}
	r.logger.Info(ctx, "insert data in database complete")

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Error(ctx, "_id can not be retrieved from the InsertedID field of the returned InsertOneResult")
	}
	r.logger.Info(ctx, "retirieving _id from InsertedID field complete")

	id := oid.Hex()

	return id, nil
}

func (r *PollRepo) Vote(c context.Context, v *models.Vote) error {
	cur := r.db.Collection(talblePolls)

	vote, err := toMongoVote(v)
	if err != nil {
		return err
	}

	filtr := bson.M{"_id": vote.PollID}
	update := bson.M{"$inc": bson.M{"choice.$[x].votes": 1}}
	arrayFilters := options.ArrayFilters{
		Filters: bson.A{bson.M{"x._id": vote.ChoiceID}},
	}
	opt := options.UpdateOptions{
		ArrayFilters: &arrayFilters,
	}

	res, err := cur.UpdateOne(c, filtr, update, &opt)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		r.logger.Error(c, "request does not matched any document")
		return e.ErrNoDocument
	}
	if res.ModifiedCount != 1 {
		r.logger.Error(c, "document have not been updated")
		return e.ErrMongoDocumentNotUpdated
	}
	r.logger.Debug(
		c,
		"MatchedCount: ", res.MatchedCount,
		"ModifiedCount: ", res.ModifiedCount,
	)

	return nil
}

func (r *PollRepo) GetPoll(c context.Context, id string) (*models.Poll, error) {
	cur := r.db.Collection(talblePolls)

	poll := new(poll)

	r.logger.Debug(c, "convert id string to ObjectId")
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error(c, "error converting id string to ObjectId")
		return nil, err
	}

	r.logger.Debug(c, "start request to db, poll id:", id)
	if err := cur.FindOne(c, bson.M{"_id": obj}).Decode(poll); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, e.ErrNoDocument
		}
		return nil, err
	}
	return toModelsPoll(poll), nil
}

func newChoice(name string) *choice {
	return &choice{
		ID:    primitive.NewObjectID(),
		Name:  name,
		Votes: 0,
	}
}

func toMongoVote(v *models.Vote) (*vote, error) {
	cID, err := primitive.ObjectIDFromHex(v.ChoiceID)
	if err != nil {
		return nil, err
	}

	pID, err := primitive.ObjectIDFromHex(v.PollID)
	if err != nil {
		return nil, err
	}

	return &vote{
		PollID:   pID,
		ChoiceID: cID,
	}, nil
}

func toModelsChoice(c *choice, pID primitive.ObjectID) *models.Choice {
	return &models.Choice{
		ID:     c.ID.Hex(),
		Name:   c.Name,
		PollID: pID.Hex(),
		Votes:  c.Votes,
	}
}

func toModelsPoll(p *poll) *models.Poll {
	choices := []models.Choice{}

	for _, choice := range p.Choices {
		choices = append(choices, *toModelsChoice(&choice, p.ID))
	}

	return &models.Poll{
		ID:       p.ID.Hex(),
		Question: p.Question,
		Choices:  choices,
	}
}
