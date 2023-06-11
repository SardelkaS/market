package feedback_repository

import (
	"github.com/jmoiron/sqlx"
	"market_auth/internal/feedback"
	feedback_model "market_auth/internal/feedback/model"
)

type postgres struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) feedback.Repository {
	return &postgres{
		db: db,
	}
}

func (p postgres) InsertFeedback(input feedback_model.Feedback) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`insert into feedback(user_id, product_id, stars, "message", pictures) values($1, $2, $3, $4, $5) returning id`,
		input.UserId, input.ProductId, input.Stars, input.Message, input.Pictures).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p postgres) RemoveFeedback(id int64) error {
	_, err := p.db.Exec(`update feedback set is_removed = true where id = $1`, id)
	return err
}

func (p postgres) GetFeedbackByInternalId(internalId *string) (*feedback_model.Feedback, error) {
	var result []feedback_model.Feedback
	err := p.db.Select(&result, `select * from feedback where internal_id = $1 and is_removed = false`, internalId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}

func (p postgres) FetchFeedback(productId *int64, userId *int64, limit, offset *int64) ([]feedback_model.Feedback, error) {
	var result []feedback_model.Feedback
	err := p.db.Select(&result, `
	select * from feedback 
	         where (product_id = $1 or $1 is null)
	         	and (user_id = $2 or $2 is null)
	         	and not ($1 is null and $2 is null)
	         	and is_removed = false
	         limit $3 offset $4`,
		productId, userId, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []feedback_model.Feedback{}, nil
	}
	return result, nil
}

func (p postgres) GetFeedbackCount(productId *int64, userId *int64) (*int64, error) {
	var result int64
	err := p.db.Get(&result, `
	select count(*) from feedback 
	         where (product_id = $1 or $1 is null)
	         	and (user_id = $2 or $2 is null)
	         	and not ($1 is null and $2 is null)
	         	and is_removed = false`,
		productId, userId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) LikeFeedback(feedbackId *int64, userId *int64) (*int64, error) {
	var id int64
	err := p.db.QueryRowx(`insert into feedback_like(feedback_id, user_id) values($1, $2) returning id`,
		feedbackId, userId).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (p postgres) UnlikeFeedback(feedbackId *int64, userId *int64) error {
	_, err := p.db.Exec(`delete from feedback_like where feedback_id = $1 and user_id = $2`, feedbackId, userId)
	return err
}

func (p postgres) GetFeedbackInfo(ids []int64, userId *int64) ([]feedback_model.FeedbackInfo, error) {
	var result []feedback_model.FeedbackInfo
	err := p.db.Select(&result, `
select
    f.internal_id,
    u.login as user_name,
    p.internal_id as product_internal_id,
    f.create_date,
    f.update_date,
    f.stars,
    f.message,
    f.pictures,
    (select count(*) from feedback_like fl where fl.feedback_id = f.id) as likes,
    (select count(*) from feedback_like fl where fl.feedback_id = f.id and fl.user_id = $2) > 0 as liked
	from feedback f 
		left join "user" u on u.id = f.user_id
		left outer join product p on p.id = f.product_id
			where f.id = any($1)
				and is_removed = false
		order by (select count(*) from feedback_like fl where fl.feedback_id = f.id), f.create_date`,
		ids, userId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return []feedback_model.FeedbackInfo{}, nil
	}
	return result, nil
}
