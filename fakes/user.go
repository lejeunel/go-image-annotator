package fake

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	pag "github.com/lejeunel/go-image-annotator/shared/pagination"
	"slices"
	"time"
)

type UserRepo struct {
	ErrOnFind                      error
	ErrOnAssignToGroup             error
	ErrOnAssignRole                error
	ErrOnCreate                    error
	ErrOnUpdatePassword            error
	ErrOnDelete                    error
	ErrOnAddForgottenPasswordState error
	ErrOnCount                     error
	ErrOnList                      error
	ErrOnSetAccessTokenHash        error
	ErrOnSetAdmin                  error
	ErrOnUnassignFromGroup         error
	ErrOnUnAssignRole              error
	Missing                        bool
	Return                         *usr.User
	GotNewGroup                    *string
	ExistingIds                    []string
	GotNewRole                     *string
	GotId                          usr.UserId
	GotHash                        []byte
	ReturnPasswordState            *usr.ForgotPasswordState
	Created                        *usr.User
	DeletedPreviousTokens          bool
	GotExpiresAt                   time.Time
	Count_                         int64
	GotSetAdmin                    bool
	GotUnassignedGroup             string
	GotUnassignedRole              string
}

func (r *UserRepo) Find(id string) (*usr.User, error) {
	if r.Missing {
		return nil, e.ErrNotFound
	}
	if r.ErrOnFind != nil {
		return nil, r.ErrOnFind
	}
	return r.Return, nil
}
func (r *UserRepo) AssignToGroup(id string, group string) error {
	if r.ErrOnAssignToGroup != nil {
		return r.ErrOnAssignToGroup
	}
	r.GotNewGroup = &group
	return nil
}
func (r *UserRepo) AssignRole(id string, role string) error {
	if r.ErrOnAssignRole != nil {
		return r.ErrOnAssignRole
	}
	r.GotNewRole = &role
	return nil
}
func (r *UserRepo) FindResetPasswordState(hash []byte) (*usr.ForgotPasswordState, error) {
	if r.Missing {
		return nil, e.ErrNotFound
	}
	r.GotHash = hash
	return r.ReturnPasswordState, nil
}
func (r *UserRepo) UpdatePassword(id usr.UserId, hash []byte) error {
	if r.ErrOnUpdatePassword != nil {
		return r.ErrOnUpdatePassword
	}
	r.GotId = id
	r.GotHash = hash
	return nil
}
func (r *UserRepo) Create(u usr.User) error {
	if r.ErrOnCreate != nil {
		return r.ErrOnCreate
	}
	r.Created = &u
	return nil
}
func (r *UserRepo) Exists(id string) (bool, error) {
	if slices.Contains(r.ExistingIds, id) {
		return true, nil
	}
	return false, nil
}
func (r *UserRepo) Delete(string) error {
	if r.ErrOnDelete != nil {
		return r.ErrOnDelete
	}
	return nil
}
func (r *UserRepo) DeleteForgottenPasswordTokens(usr.UserId) error {
	r.DeletedPreviousTokens = true
	return nil
}
func (r *UserRepo) AddForgottenPasswordState(hash []byte, id usr.UserId, expires time.Time) error {
	if r.ErrOnAddForgottenPasswordState != nil {
		return r.ErrOnAddForgottenPasswordState
	}
	r.GotId = id
	r.GotHash = hash
	r.GotExpiresAt = expires
	return nil
}
func (r *UserRepo) Count() (int64, error) {
	if r.ErrOnCount != nil {
		return 0, r.ErrOnCount
	}
	return int64(r.Count_), nil
}

func (r *UserRepo) List(req pag.PaginationParams) ([]usr.User, error) {
	if r.ErrOnList != nil {
		return nil, r.ErrOnList

	}

	result := []usr.User{}
	for range req.PageSize {
		usr := usr.NewUser("the-id")
		result = append(result, usr)
	}
	return result, nil

}

func (r *UserRepo) SetAccessTokenHash(id usr.UserId, hash []byte) error {
	if r.ErrOnSetAccessTokenHash != nil {
		return r.ErrOnSetAccessTokenHash
	}
	r.GotId = id
	r.GotHash = hash
	return nil
}

func (r *UserRepo) SetAdmin(id usr.UserId, value bool) error {
	if r.ErrOnSetAdmin != nil {
		return r.ErrOnSetAdmin
	}
	r.GotSetAdmin = value
	return nil
}

func (r *UserRepo) UnAssignFromGroup(id string, group string) error {
	if r.ErrOnUnassignFromGroup != nil {
		return r.ErrOnUnassignFromGroup
	}
	r.GotUnassignedGroup = group
	return nil
}

func (r *UserRepo) UnAssignRole(id string, role string) error {
	if r.ErrOnUnAssignRole != nil {
		return r.ErrOnUnAssignRole
	}
	r.GotUnassignedRole = role
	return nil
}
