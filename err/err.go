package err

import "errors"

var (
	ErrNoDocument              = errors.New("document not found")
	ErrMongoDocumentNotUpdated = errors.New(("document not updated"))
)
