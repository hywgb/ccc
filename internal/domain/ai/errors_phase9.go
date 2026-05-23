package ai

import "errors"

var (
	ErrDigitalEmployeeNotFound = errors.New("digital employee not found")
	ErrSceneNotFound           = errors.New("digital employee scene not found")
	ErrSceneAlreadyPublished   = errors.New("scene is already published")
	ErrInvalidIntentConfig     = errors.New("invalid intent config JSON")
	ErrQARuleNotFound          = errors.New("qa rule not found")
	ErrInvalidQARuleType       = errors.New("invalid qa rule type")
	ErrQASchemeNotFound        = errors.New("qa scheme not found")
	ErrQAResultNotFound        = errors.New("qa result not found")
	ErrQAResultNotAppealable   = errors.New("qa result cannot be appealed")
	ErrQAResultNotReviewable   = errors.New("qa result cannot be reviewed")
	ErrASRHotwordsNotFound     = errors.New("asr hotwords not found")
	ErrScorecardNotFound       = errors.New("performance scorecard not found")
	ErrEmptyTranscript         = errors.New("transcript is empty")
)
