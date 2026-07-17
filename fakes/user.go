package fake

import (
	usr "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"slices"
)

type UserRepo struct {
	ErrOnFind           error
	ErrOnAssignToGroup  error
	ErrOnAssignRole     error
	ErrOnCreate         error
	ErrOnUpdatePassword error
	ErrOnDelete         error
	Missing             bool
	Return              *usr.User
	GotNewGroup         *string
	ExistingIds         []string
	GotNewRole          *string
	GotId               usr.UserId
	GotHash             []byte
	ReturnPasswordState *usr.ForgotPasswordState
	Created             *usr.User
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
