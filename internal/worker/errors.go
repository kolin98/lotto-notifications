package worker

import "errors"

var (
	ErrGameNotCheckable = errors.New("workers can be only created for checkable games")
	ErrGameInfoNotSet   = errors.New("game info must be set before creating worker")
	ErrNoResultsFound   = errors.New("no results found")
	ErrResultsNotReady  = errors.New("results not ready")
	ErrMainGameNotFound = errors.New("main game not found")
)
