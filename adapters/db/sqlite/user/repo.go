package user

import (
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
)

type SQLiteUserRepo struct {
	Db *sqlx.DB
}

type Record struct {
	Id           string `db:"id"`
	Roles        string `db:"roles"`
	IsAdmin      bool   `db:"is_admin"`
	ApiTokenHash string `db:"api_token_hash"`
	PasswordHash string `db:"password_hash"`
}

type ForgottenPasswordStateRecord struct {
	Id        string    `db:"id"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (r SQLiteUserRepo) Create(usr u.User) error {
	query := "INSERT INTO users (id,is_admin,api_token_hash,password_hash) VALUES ($1,$2,$3,$4)"
	_, err := r.Db.Exec(query, usr.Id, usr.IsAdmin, hex.EncodeToString(usr.HashPAT), hex.EncodeToString(usr.HashPassword))
	if err != nil {
		return fmt.Errorf("inserting record: %v: %w", err, e.ErrInternal)
	}
	if err := r.SetGroups(usr.Id, usr.Groups); err != nil {
		return err
	}
	if err := r.SetRoles(usr.Id, usr.Roles); err != nil {
		return err
	}

	return nil
}
func (r SQLiteUserRepo) SetGroups(user u.UserId, groups []string) error {
	_, err := r.Db.Exec(
		"DELETE FROM users_groups WHERE user_id=$1;",
		user,
	)
	if err != nil {
		return fmt.Errorf("deleting user groups: %v: %w", err, e.ErrInternal)
	}

	for _, group := range groups {
		_, err = r.Db.Exec(
			"INSERT INTO users_groups (user_id, group_id) SELECT $1, id FROM groups WHERE name=$2",
			user, group)
		if err != nil {
			return fmt.Errorf("inserting group %q: %v: %w", group, err, e.ErrInternal)
		}
	}
	return nil
}
func (r SQLiteUserRepo) getGroupNames(userId string) ([]string, error) {
	var groups []string
	query := "SELECT name FROM groups WHERE id=(SELECT group_id FROM users_groups WHERE user_id=$1)"
	err := r.Db.Select(&groups, query, userId)
	if err != nil {
		return nil, fmt.Errorf("fetching groups: %v: %w", err, e.ErrInternal)
	}

	return groups, nil
}
func (r SQLiteUserRepo) getRoleNames(userId string) ([]string, error) {
	var roles []string
	query := "SELECT name FROM roles WHERE id=(SELECT id FROM users_roles WHERE user_id=$1)"
	err := r.Db.Select(&roles, query, userId)
	if err != nil {
		return nil, fmt.Errorf("fetching roles: %v: %w", err, e.ErrInternal)
	}

	return roles, nil
}
func (r SQLiteUserRepo) recordToEntity(rec Record) (*u.User, error) {
	groups, err := r.getGroupNames(rec.Id)
	if err != nil {
		return nil, err
	}
	roles, err := r.getRoleNames(rec.Id)
	if err != nil {
		return nil, err
	}

	patHash, err := hex.DecodeString(rec.ApiTokenHash)
	if err != nil {
		return nil, err
	}
	pwHash, err := hex.DecodeString(rec.PasswordHash)
	if err != nil {
		return nil, err
	}

	user := u.NewUser(rec.Id, u.WithRoles(roles), u.WithGroups(groups),
		u.WithAdmin(rec.IsAdmin), u.WithHashedPersonalAccessToken(patHash),
		u.WithHashedPassword(pwHash))
	return &user, nil
}
func (r SQLiteUserRepo) Find(id u.UserId) (*u.User, error) {
	record := Record{}
	err := r.Db.Get(&record,
		"SELECT id,is_admin,api_token_hash,password_hash FROM users WHERE id=$1", id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrNotFound
		default:
			return nil, fmt.Errorf("finding user by id: %v: %w", err, e.ErrInternal)
		}
	}
	user, err := r.recordToEntity(record)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r SQLiteUserRepo) Delete(id string) error {
	_, err := r.Db.Exec("DELETE FROM users WHERE id=$1", id)

	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteUserRepo) Exists(id string) (bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)`, id)
	if err != nil {
		return exists, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return exists, nil

}
func (r SQLiteUserRepo) Count() (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM users"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting records: %v: %w", err, e.ErrInternal)
	}

	return count, nil
}
func (r SQLiteUserRepo) List(m pag.PaginationParams) ([]u.User, error) {
	q := sq.StatementBuilder.Select("id,is_admin,api_token_hash,password_hash").From("users")
	q = q.Limit(uint64(m.PageSize)).Offset((uint64(m.Page-1) * uint64(m.PageSize)))
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Record{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	users := []u.User{}
	for _, rec := range records {
		user, err := r.recordToEntity(rec)
		if err != nil {
			return nil, fmt.Errorf("converting user records to domain objects: %v: %w", err, e.ErrInternal)
		}
		users = append(users, *user)
	}

	return users, nil
}
func (r SQLiteUserRepo) SetRoles(userId string, roles []string) error {
	_, err := r.Db.Exec(
		"DELETE FROM users_roles WHERE user_id=$1;",
		userId,
	)
	if err != nil {
		return fmt.Errorf("deleting user groups: %v: %w", err, e.ErrInternal)
	}

	for _, role := range roles {
		_, err = r.Db.Exec(
			"INSERT INTO users_roles (user_id, role_id) SELECT $1, id FROM roles WHERE name=$2",
			userId, role)
		if err != nil {
			return fmt.Errorf("inserting role %q: %v: %w", role, err, e.ErrInternal)
		}
	}
	return nil
}
func (r SQLiteUserRepo) SetAdmin(userId string, value bool) error {
	query := "UPDATE users SET is_admin=$2 WHERE id=$1"
	_, err := r.Db.Exec(query, userId, value)
	if err != nil {
		return fmt.Errorf("setting admin right to %v : %v: %w", value, err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteUserRepo) SetAccessTokenHash(userId u.UserId, hash []byte) error {
	query := "UPDATE users SET api_token_hash=$2 WHERE id=$1"
	_, err := r.Db.Exec(query, userId, hex.EncodeToString(hash))
	if err != nil {
		return fmt.Errorf("setting access token hash: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteUserRepo) AddForgottenPasswordState(hash []byte, id u.UserId, expiresAt time.Time) error {
	query := "INSERT INTO forgot_password (token_hash,id,expires_at) VALUES ($1,$2,$3)"
	_, err := r.Db.Exec(query, hex.EncodeToString(hash), id, expiresAt)
	if err != nil {
		return fmt.Errorf("inserting record: %v: %w", err, e.ErrInternal)
	}
	return nil

}
func (r SQLiteUserRepo) FindResetPasswordState(hash []byte) (*u.ForgotPasswordState, error) {

	record := ForgottenPasswordStateRecord{}
	err := r.Db.Get(&record,
		"SELECT id,expires_at FROM forgot_password WHERE token_hash=$1", hex.EncodeToString(hash))

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrNotFound
		default:
			return nil, fmt.Errorf("finding forgotten password state: %v: %w", err, e.ErrInternal)
		}

	}
	return &u.ForgotPasswordState{Id: record.Id, ExpiresAt: &record.ExpiresAt}, nil

}
func (r SQLiteUserRepo) UpdatePassword(id u.UserId, hash []byte) error {
	query := "UPDATE users SET password_hash=$1 WHERE id=$2"
	_, err := r.Db.Exec(query, hex.EncodeToString(hash), id)
	if err != nil {
		return fmt.Errorf("updating password hash: %w: %w", err, e.ErrInternal)
	}
	return nil
}
func (r SQLiteUserRepo) DeleteForgottenPasswordTokens(id u.UserId) error {
	query := "DELETE FROM forgot_password WHERE id=$1"
	_, err := r.Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("deleting forgotten password: %v: %w", err, e.ErrInternal)
	}
	return nil

}
func NewSQLiteUserRepo(db *sqlx.DB) SQLiteUserRepo {
	return SQLiteUserRepo{Db: db}
}
