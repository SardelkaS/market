package feedback_model

type CreateFeedbackBody struct {
	UserId    *int64   `json:"-"`
	ProductId *string  `json:"product_id"`
	Stars     *int64   `json:"stars"`
	Message   *string  `json:"message"`
	Pictures  []string `json:"pictures"`
}

type RemoveFeedbackBody struct {
	UserId     *int64  `json:"-"`
	FeedbackId *string `json:"-"`
}

type FetchFeedbackParams struct {
	UserId    *int64
	ProductId *string
	OnlyMy    *bool
	Limit     *int64
	Offset    *int64
}

type FetchFeedbackResponse struct {
	Feedback []FeedbackInfo `json:"feedback"`
	Count    *int64         `json:"count"`
}
