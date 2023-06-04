package auth_model

type User struct {
	Id         *int64  `db:"id"`
	Login      *string `db:"login"`
	Password   *string `db:"password"`
	Email      *string `db:"email"`
	Ban        *bool   `db:"ban"`
	IsDeleted  *bool   `db:"is_deleted"`
	Role       *string `db:"user_role"`
	InternalId *string `db:"internal_id"`
	Timezone   *string `db:"timezone"`
}
