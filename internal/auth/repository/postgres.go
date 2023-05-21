package auth_repository

import (
	"market_auth/internal/auth"
	"market_auth/internal/auth/model"

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
	err := p.db.Get(&result, `SELECT * from "auth" WHERE login = $1`, name)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) GetUserById(id int64) (*auth_model.User, error) {
	var result auth_model.User
	err := p.db.Get(&result, `SELECT * from "auth" WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) GetUserByInternalId(internalId string) (*auth_model.User, error) {
	var result auth_model.User
	err := p.db.Get(&result, `SELECT * from "auth" WHERE internal_id = $1`, internalId)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p postgres) InsertUser(user auth_model.User) error {
	res, err := p.db.Query(`INSERT INTO "auth"(login, password, user_role, internal_id, ban) values($1, $2, $3, $4, $5)`,
		user.Login, user.Password, user.Role, user.InternalId, user.Ban, user.Timezone)
	if res != nil {
		_ = res.Close()
	}
	if err != nil {
		return err
	}
	return nil
}

func (p postgres) UpdatePassword(input auth_model.UpdatePasswordGatewayInput) error {
	res, err := p.db.Query(`UPDATE "auth" SET password = $1 WHERE id = $2`, input.HashedPassword, input.UserId)
	if err != nil {
		return err
	}
	return res.Close()
}

func (p postgres) UpdateTimezone(input auth_model.UpdateTimezoneGatewayInput) error {
	res, err := p.db.Query(`UPDATE "auth" SET timezone=$1 WHERE id=$2`, input.NewTimezone, input.UserId)
	if err != nil {
		return err
	}
	return res.Close()
}
