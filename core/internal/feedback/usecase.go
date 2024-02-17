package feedback

import feedback_model "market_auth/internal/feedback/model"

type UC interface {
	CreateFeedback(input feedback_model.CreateFeedbackBody) error
	RemoveFeedback(input feedback_model.RemoveFeedbackBody) error
	GetFeedbackByInternalId(internalId string) (*feedback_model.Feedback, error)
	FetchFeedback(input feedback_model.FetchFeedbackParams) (*feedback_model.FetchFeedbackLogicOutput, error)
	LikeFeedback(feedbackId string, userId int64) error
	UnlikeFeedback(feedbackId string, userId int64) error
	GetFeedbackInfo(feedbacks []feedback_model.Feedback, userId int64) ([]feedback_model.FeedbackInfo, error)
}
