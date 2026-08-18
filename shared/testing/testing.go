package testing

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

//go:embed sample-image.jpg
var TestJPGImage []byte

//go:embed sample-image.png
var TestPNGImage []byte

type TestingErrPresenter struct {
	GotDuplicationErr bool
	GotValidationErr  bool
	GotInternalErr    bool
	GotNotFoundErr    bool
	GotDependencyErr  bool
	GotErr            error
	GotAuthErr        bool
	GotForbiddenErr   bool
}

func (p *TestingErrPresenter) Error(err error) {
	p.GotErr = err
	switch {
	case errors.Is(err, e.ErrDuplicate):
		p.GotDuplicationErr = true
	case errors.Is(err, e.ErrValidation):
		p.GotValidationErr = true
	case errors.Is(err, e.ErrNotFound):
		p.GotNotFoundErr = true
	case errors.Is(err, e.ErrDependency):
		p.GotDependencyErr = true
	case errors.Is(err, e.ErrAuthorization):
		p.GotAuthErr = true
	case errors.Is(err, e.ErrForbiddenOp):
		p.GotForbiddenErr = true

	default:
		p.GotInternalErr = true
	}
}

type FakeAuth struct {
	Fail bool
}

type FakeProvider struct{}

func (p FakeProvider) Provide() (*u.User, error) {
	return &u.User{}, nil
}

func CreateCtxWithUserId(ctx context.Context, userId u.UserId) context.Context {
	user := u.NewUser(userId)
	return u.AppendUserToContext(ctx, user)
}

func FakeUUIDFromInt(n int) string {
	digit := fmt.Sprintf("%d", n)
	full := strings.Repeat(digit, 32)
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		full[0:8],
		full[8:12],
		full[12:16],
		full[16:20],
		full[20:32],
	)
}

func IdFromInt(i int) *im.ImageId {
	id, err := im.NewImageIdFromString(FakeUUIDFromInt(i))
	if err != nil {
		panic(err)
	}
	return &id
}
