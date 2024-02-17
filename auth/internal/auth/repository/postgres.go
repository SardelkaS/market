package auth_repository

import (
	"auth/internal/auth"
	"auth/internal/auth/model"

	"github.com/jmoiron/sqlx"
)

type postgres struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) auth.Repository {
	return &postgres{
		db: db,
	}
}

func (p postgres) GetUserByName(name string) (*auth_model.User, error) {
	var result auth_model.User
	err := p.db.Get(&result, `SELECT * from "user" WHERE login = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) GetUserById(id int64) (*auth_model.User, error) {
	var result auth_model.User
	err := p.db.Get(&result, `SELECT * from "user" WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) GetUserByInternalId(internalId string) (*auth_model.User, error) {
	var result auth_model.User
	err := p.db.Get(&result, `SELECT * from "user" WHERE internal_id = $1`, internalId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) InsertUser(user auth_model.User) error {
	_, err := p.db.Exec(`INSERT INTO "user"("login", "password", user_role, internal_id, ban, "timezone") values($1, $2, $3, $4, $5, $6)`,
		user.Login, user.Password, user.Role, user.InternalId, user.Ban, user.Timezone)
	return err
}

func (p postgres) UpdatePassword(input auth_model.UpdatePasswordGatewayInput) error {
	_, err := p.db.Exec(`UPDATE "user" SET "password" = $1 WHERE id = $2`, input.HashedPassword, input.UserId)
	return err
}

func (p postgres) UpdateTimezone(input auth_model.UpdateTimezoneGatewayInput) error {
	_, err := p.db.Exec(`UPDATE "user" SET "timezone"=$1 WHERE id=$2`, input.NewTimezone, input.UserId)
	return err
}

func (p postgres) UpdateUserInfo(input auth_model.UpdateUserInfoGatewayInput) error {
	_, err := p.db.Exec(`update "user" set email = $1, "name" = $2, phone_number = $3, birth_date = $4, contact_data = $5 where id = $6`,
		input.Email, input.Name, input.PhoneNumber, input.BirthDate, input.ContactData, input.Id)
	return err
}

func (p postgres) GetUserInfoById(id int64) (*auth_model.UserInfo, error) {
	var result []auth_model.UserInfo
	err := p.db.Select(&result, `
select 
    u.internal_id,
    u.login,
    u.email,
    u."name",
    u.phone_number,
    u.birth_date,
    u.contact_data
	from "user" u
		where u.id = $1`,
		id)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}
