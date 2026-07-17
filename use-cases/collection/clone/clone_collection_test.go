package clone

import (
	"context"
	"testing"

	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestSubmitTaskWithoutIdentity(t *testing.T) {
	itr := NewTestingCloner()
	p := &FakePresenter{}
	itr.Execute(t.Context(), Request{}, p)
	assert.NotNil(t, p.GotErr)
	assert.False(t, p.GotSuccess)
}

func CreateCtxWithUserId(ctx context.Context, userId u.UserId) context.Context {
	user := u.NewUser(userId)
	return u.AppendUserToContext(ctx, user)
}
func TestHandleAuthErr(t *testing.T) {
	group := "my-group"
	itr := NewTestingCloner()
	itr.Auth = fk.Auth{Err: e.ErrAuthorization}
	p := &FakePresenter{}
	itr.Execute(CreateCtxWithUserId(t.Context(), "user@mail.com"), Request{DestinationGroup: &group}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}
