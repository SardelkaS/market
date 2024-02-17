package feedback_usecase

import (
	"fmt"
	"github.com/google/uuid"
	"market_auth/internal/failure"
	"market_auth/internal/feedback"
	feedback_model "market_auth/internal/feedback/model"
	"market_auth/internal/product"
	"market_auth/pkg/secure"
)

type uc struct {
	repo      feedback.Repository
	productUC product.UC
}

func New(repo feedback.Repository, productUC product.UC) feedback.UC {
	return &uc{
		repo:      repo,
		productUC: productUC,
	}
}

func (u *uc) CreateFeedback(input feedback_model.CreateFeedbackBody) error {
	if input.UserId == nil || input.Stars == nil || *input.Stars < 1 || *input.Stars > 5 || input.ProductId == nil {
		return failure.ErrInput
	}

	productData, err := u.productUC.GetProduct(*input.ProductId)
	if err != nil {
		return err
	}

	internalId := secure.CalcInternalId(uuid.New().String())
	_, err = u.repo.InsertFeedback(feedback_model.Feedback{
		InternalId: &internalId,
		UserId:     input.UserId,
		ProductId:  productData.Id,
		Stars:      input.Stars,
		Message:    input.Message,
		Pictures:   input.Pictures,
	})
	if err != nil {
		fmt.Printf("Error to create feedback for product %s: %s\n", *productData.InternalId, err.Error())
		return fmt.Errorf("error to save feedback")
	}

	return nil
}

func (u *uc) RemoveFeedback(input feedback_model.RemoveFeedbackBody) error {
	if input.FeedbackId == nil || input.UserId == nil {
		return failure.ErrInput
	}

	feedbackData, err := u.repo.GetFeedbackByInternalId(input.FeedbackId)
	if err != nil {
		fmt.Printf("Error to get feedback %s: %s\n", *input.FeedbackId, err.Error())
		return fmt.Errorf("error to get feedback")
	}
	if feedbackData == nil {
		fmt.Printf("Error to get feedback %s: not found\n", *input.FeedbackId)
		return fmt.Errorf("error to get feedback")
	}

	if *feedbackData.UserId != *input.UserId {
		return fmt.Errorf("feedback belongs to another user")
	}

	err = u.repo.RemoveFeedback(*feedbackData.Id)
	if err != nil {
		fmt.Printf("Error to remove feedback %s: %s", *input.FeedbackId, err.Error())
		return fmt.Errorf("error to remove feedback")
	}

	return nil
}

func (u *uc) GetFeedbackByInternalId(internalId string) (*feedback_model.Feedback, error) {
	feedbackData, err := u.repo.GetFeedbackByInternalId(&internalId)
	if err != nil {
		fmt.Printf("Error to get feedback %s: %s\n", internalId, err.Error())
		return nil, fmt.Errorf("error to get feedback")
	}
	if feedbackData == nil {
		fmt.Printf("Error to get feedback %s: not found\n", internalId)
		return nil, fmt.Errorf("error to get feedback")
	}

	return feedbackData, nil
}

func (u *uc) FetchFeedback(input feedback_model.FetchFeedbackParams) (*feedback_model.FetchFeedbackLogicOutput, error) {
	if input.UserId == nil && input.ProductId == nil {
		return nil, failure.ErrInput
	}

	var productId *int64
	if input.ProductId != nil {
		productData, err := u.productUC.GetProduct(*input.ProductId)
		if err != nil {
			return nil, err
		}

		productId = productData.Id
	}

	var userId *int64
	if input.OnlyMy != nil && *input.OnlyMy {
		userId = input.UserId
	}

	feedbacks, err := u.repo.FetchFeedback(productId, userId, input.Limit, input.Offset)
	if err != nil {
		fmt.Printf("Error to fetch feedback: %s\n", err.Error())
		return nil, fmt.Errorf("error to fetch feedback")
	}

	count, err := u.repo.GetFeedbackCount(productId, userId)
	if err != nil {
		fmt.Printf("Error to get feedback count: %s\n", err.Error())
		return nil, fmt.Errorf("error to get feedback count")
	}

	return &feedback_model.FetchFeedbackLogicOutput{
		Feedback: feedbacks,
		Count:    count,
	}, nil
}

func (u *uc) LikeFeedback(feedbackId string, userId int64) error {
	feedbackData, err := u.GetFeedbackByInternalId(feedbackId)
	if err != nil {
		return err
	}

	if *feedbackData.UserId == userId {
		return fmt.Errorf("it's your feedback")
	}

	isLiked, err := u.repo.CheckIsFeedbackLiked(*feedbackData.Id, userId)
	if err != nil {
		fmt.Printf("Error to check is feedback %s liked by %d: %s", feedbackId, userId, err.Error())
		return fmt.Errorf("error to like feedback")
	}

	if isLiked == nil || *isLiked {
		return fmt.Errorf("feedback already liked")
	}

	_, err = u.repo.LikeFeedback(feedbackData.Id, &userId)
	if err != nil {
		fmt.Printf("Error to like feedback %s: %s\n", feedbackId, err.Error())
		return fmt.Errorf("error to like feedback")
	}

	return nil
}

func (u *uc) UnlikeFeedback(feedbackId string, userId int64) error {
	feedbackData, err := u.GetFeedbackByInternalId(feedbackId)
	if err != nil {
		return err
	}

	if *feedbackData.UserId == userId {
		return fmt.Errorf("it's your feedback")
	}

	err = u.repo.UnlikeFeedback(feedbackData.Id, &userId)
	if err != nil {
		fmt.Printf("Error to unlike feedback %s: %s\n", feedbackId, err.Error())
		return fmt.Errorf("error to unlike feedback")
	}

	return nil
}

func (u *uc) GetFeedbackInfo(feedbacks []feedback_model.Feedback, userId int64) ([]feedback_model.FeedbackInfo, error) {
	var ids []int64
	for _, feedbackData := range feedbacks {
		ids = append(ids, *feedbackData.Id)
	}

	result, err := u.repo.GetFeedbackInfo(ids, &userId)
	if err != nil {
		fmt.Printf("Error to get feedbacks info (%v): %s\n", ids, err.Error())
		return nil, fmt.Errorf("error to get feedbacks info")
	}

	return result, nil
}
