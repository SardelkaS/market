package auth_model

type User struct {
	Id          *int64  `db:"id"`
	Login       *string `db:"login"`
	Password    *string `db:"password"`
	Email       *string `db:"email"`
	Ban         *bool   `db:"ban"`
	IsDeleted   *bool   `db:"is_deleted"`
	Role        *string `db:"user_role"`
	InternalId  *string `db:"internal_id"`
	Timezone    *string `db:"timezone"`
	Name        *string `db:"name"`
	PhoneNumber *string `db:"phone_number"`
	BirthDate   *string `db:"birth_date"`
	ContactData *string `db:"contact_data"`
}

type UserInfo struct {
	InternalId  *string `json:"internal_id" db:"internal_id"`
	Login       *string `json:"login" db:"login"`
	Email       *string `json:"email" db:"email"`
	Name        *string `json:"name" db:"name"`
	PhoneNumber *string `json:"phone_number" db:"phone_number"`
	BirthDate   *string `json:"birth_date" db:"birth_date"`
	ContactData *string `json:"contact_data" db:"contact_data"`
}
