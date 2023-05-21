package auth_model

type User struct {
	Id         *int64  `db:"id"`
	Login      *string `db:"login"`
	Password   *string `db:"password"`
	Role       *string `db:"user_role"`
	InternalId *string `db:"internal_id"`
	Ban        *bool   `db:"ban"`
	Timezone   *string `db:"timezone"`
}
