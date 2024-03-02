package feedback

import feedback_model "core/internal/feedback/model"

type Repository interface {
	InsertFeedback(input feedback_model.Feedback) (*int64, error)
	RemoveFeedback(id int64) error
	GetFeedbackByInternalId(internalId *string) (*feedback_model.Feedback, error)
	FetchFeedback(productId *int64, userId *int64, limit, offset *int64) ([]feedback_model.Feedback, error)
	GetFeedbackCount(productId *int64, userId *int64) (*int64, error)
	CheckIsFeedbackLiked(feedbackId int64, userId int64) (*bool, error)
	LikeFeedback(feedbackId *int64, userId *int64) (*int64, error)
	UnlikeFeedback(feedbackId *int64, userId *int64) error
	GetFeedbackInfo(ids []int64, userId *int64) ([]feedback_model.FeedbackInfo, error)
}
