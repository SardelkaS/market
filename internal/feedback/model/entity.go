package feedback_model

import (
	"github.com/lib/pq"
	"time"
)

type Feedback struct {
	Id         *int64         `db:"id"`
	InternalId *string        `db:"internal_id"`
	UserId     *int64         `db:"user_id"`
	ProductId  *int64         `db:"product_id"`
	CreateDate *time.Time     `db:"create_date"`
	UpdateDate *time.Time     `db:"update_date"`
	Stars      *int64         `db:"stars"`
	Message    *string        `db:"message"`
	Pictures   pq.StringArray `db:"pictures"`
	IsRemoved  *bool          `db:"is_removed"`
}

type FeedbackInfo struct {
	InternalId        *string        `json:"internalId" db:"internal_id"`
	UserName          *string        `json:"user_name" db:"user_name"`
	ProductInternalId *string        `json:"product_internal_id" db:"product_internal_id"`
	CreateDate        *time.Time     `json:"create_date" db:"create_date"`
	UpdateDate        *time.Time     `json:"update_date" db:"update_date"`
	Stars             *int64         `json:"stars" db:"stars"`
	Message           *string        `json:"message" db:"message"`
	Pictures          pq.StringArray `json:"pictures" db:"pictures"`
	Likes             *int64         `json:"likes" db:"likes"`
	Liked             *bool          `json:"liked" db:"liked"`
	IsMy              *bool          `json:"is_my" db:"is_my"`
}
