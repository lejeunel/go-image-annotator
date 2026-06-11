package user

import (
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	s "github.com/lejeunel/go-image-annotator/adapters/sqlite/repos"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
)

func appendUnique(slice []string, s string) []string {
	if slices.Contains(slice, s) {
		return slice
	}
	return append(slice, s)
}

func remove(slice []string, s string) []string {
	result := slice[:0]
	for _, existing := range slice {
		if existing != s {
			result = append(result, existing)
		}
	}
	return result
}

type SQLiteUserRepo struct {
	Db *sqlx.DB
}

type Record struct {
	Id           string `db:"id"`
	Roles        string `db:"roles"`
	IsAdmin      bool   `db:"is_admin"`
	ApiTokenHash string `db:"api_token_hash"`
}

func (r *SQLiteUserRepo) Create(usr u.User) error {
	roles, err := json.Marshal(usr.Roles)
	if err != nil {
		return fmt.Errorf("inserting record: %v: %w", err, e.ErrInternal)
	}
	query := "INSERT INTO users (id,roles,is_admin,api_token_hash) VALUES ($1,$2,$3,$4)"
	_, err = r.Db.Exec(query, usr.Id, roles, usr.IsAdmin, hex.EncodeToString(usr.HashPAT))
	if err != nil {
		return fmt.Errorf("inserting record: %v: %w", err, e.ErrInternal)
	}

	for _, g := range usr.Groups {
		if err := r.AssignToGroup(usr.Id, g); err != nil {
			return err
		}
	}

	return nil
}

func (r *SQLiteUserRepo) AssignToGroup(user string, group string) error {
	query := "INSERT INTO users_groups (user_id,group_id) VALUES ($1,(SELECT id FROM groups WHERE name=$2))"
	_, err := r.Db.Exec(query, user, group)
	if err != nil {
		return fmt.Errorf("inserting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteUserRepo) UnAssignFromGroup(user string, group string) error {
	query := "DELETE FROM users_groups WHERE user_id=$1 AND group_id=(SELECT id FROM groups WHERE name=$2)"
	_, err := r.Db.Exec(query, user, group)
	if err != nil {
		return fmt.Errorf("unassigning user from group: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteUserRepo) getGroupNames(userId string) ([]string, error) {
	var groups []string
	query := "SELECT name FROM groups WHERE id=(SELECT id FROM users_groups WHERE user_id=$1)"
	err := r.Db.Select(&groups, query, userId)
	if err != nil {
		return nil, fmt.Errorf("fetching groups: %v: %w", err, e.ErrInternal)
	}

	return groups, nil

}
func (r *SQLiteUserRepo) Find(id u.UserId) (*u.User, error) {
	record := Record{}
	err := r.Db.Get(&record,
		"SELECT id,roles,is_admin,api_token_hash FROM users WHERE id=$1", id)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, e.ErrNotFound
		default:
			return nil, fmt.Errorf("finding user by id: %v: %w", err, e.ErrInternal)
		}

	}
	var roles []string
	err = json.Unmarshal([]byte(record.Roles), &roles)
	if err != nil {
		return nil, fmt.Errorf("finding user by id: unmarshaling user roles %v: %w: %w", record.Roles, err, e.ErrInternal)
	}

	groups, err := r.getGroupNames(id)
	if err != nil {
		return nil, err
	}

	hash, err := hex.DecodeString(record.ApiTokenHash)
	if err != nil {
		return nil, err
	}

	user := u.NewUser(record.Id, u.WithRoles(roles), u.WithGroups(groups),
		u.WithAdmin(record.IsAdmin), u.WithHashedPersonalAccessToken(hash))
	return &user, nil
}
func (r *SQLiteUserRepo) Delete(id string) error {
	_, err := r.Db.Exec("DELETE FROM users WHERE id=$1", id)

	if err != nil {
		return fmt.Errorf("deleting record: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteUserRepo) Exists(id string) (bool, error) {
	var exists bool

	err := r.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM users WHERE id=$1)`, id)
	if err != nil {
		return exists, fmt.Errorf("%v: %w", err, e.ErrInternal)
	}

	return exists, nil

}
func (r *SQLiteUserRepo) Count() (int64, error) {
	var count int64

	query := "SELECT COUNT(*) FROM users"
	err := r.Db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("counting records: %v: %w", err, e.ErrInternal)
	}

	return count, nil
}
func (r *SQLiteUserRepo) List(m list.Request) ([]u.User, error) {
	q := sq.StatementBuilder.Select("id").From("users")
	q = q.Limit(uint64(m.PageSize)).Offset((uint64(m.Page-1) * uint64(m.PageSize)))
	sql, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("building query: %v: %w", err, e.ErrInternal)
	}
	records := []Record{}
	if err := r.Db.Select(&records, sql, args...); err != nil {
		return nil, fmt.Errorf("applying query: %v: %w", err, e.ErrInternal)
	}

	objects := []u.User{}
	for _, r := range records {
		user := u.NewUser(r.Id)
		objects = append(objects, user)
	}

	return objects, nil
}
func (r *SQLiteUserRepo) getCurrentRoles(userId string) ([]string, error) {
	var currentRoles string
	query := "SELECT roles FROM users WHERE id=$1"
	err := r.Db.QueryRow(query, userId).Scan(&currentRoles)
	if err != nil {
		return nil, fmt.Errorf("fetching currently applied roles: %v: %w", err, e.ErrInternal)
	}
	var roles []string
	err = json.Unmarshal([]byte(currentRoles), &roles)
	if err != nil {
		return nil, fmt.Errorf("finding user by id: unmarshaling user roles %v: %w: %w", currentRoles, err, e.ErrInternal)
	}
	return roles, nil
}
func (r *SQLiteUserRepo) setRoles(userId string, roles []string) error {
	query := "UPDATE users SET roles=$2 WHERE id=$1"
	data, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("assigning roles %v: %w: %w", roles, err, e.ErrInternal)
	}
	_, err = r.Db.Exec(query, userId, data)
	if err != nil {
		return fmt.Errorf("assigning role: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func (r *SQLiteUserRepo) AssignRole(userId string, role string) error {
	currentRoles, err := r.getCurrentRoles(userId)
	if err != nil {
		return fmt.Errorf("assigning role: %w", err)
	}

	if err := r.setRoles(userId, appendUnique(currentRoles, role)); err != nil {
		return err
	}
	return nil
}
func (r *SQLiteUserRepo) UnAssignRole(userId string, role string) error {
	currentRoles, err := r.getCurrentRoles(userId)
	if err != nil {
		return fmt.Errorf("unassigning role: %w", err)
	}
	if !slices.Contains(currentRoles, role) {
		return fmt.Errorf("unassigning role: checking for existence of role %v: %w", role, e.ErrNotFound)
	}
	if err := r.setRoles(userId, remove(currentRoles, role)); err != nil {
		return err
	}
	return nil
}
func (r *SQLiteUserRepo) SetAdmin(userId string, value bool) error {
	query := "UPDATE users SET is_admin=$2 WHERE id=$1"
	_, err := r.Db.Exec(query, userId, value)
	if err != nil {
		return fmt.Errorf("setting admin right to %v : %v: %w", value, err, e.ErrInternal)
	}
	return nil
}

func (r *SQLiteUserRepo) SetAccessTokenHash(userId u.UserId, hash []byte) error {
	hashStr := hex.EncodeToString(hash)
	query := "UPDATE users SET api_token_hash=$2 WHERE id=$1"
	_, err := r.Db.Exec(query, userId, hashStr)
	if err != nil {
		return fmt.Errorf("setting access token hash: %v: %w", err, e.ErrInternal)
	}
	return nil
}
func NewSQLiteUserRepo(db *sqlx.DB) *SQLiteUserRepo {
	return &SQLiteUserRepo{Db: db}
}

func NewTestSQLiteUserRepo() *SQLiteUserRepo {
	return NewSQLiteUserRepo(s.NewSQLiteDB(":memory:"))
}
