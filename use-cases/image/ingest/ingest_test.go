package ingest

import (
	"bytes"
	"testing"

	an "github.com/lejeunel/go-image-annotator/entities/annotation"
	clc "github.com/lejeunel/go-image-annotator/entities/collection"
	grp "github.com/lejeunel/go-image-annotator/entities/group"
	m "github.com/lejeunel/go-image-annotator/entities/meta"
	u "github.com/lejeunel/go-image-annotator/entities/user"
	fk "github.com/lejeunel/go-image-annotator/fakes"
	ig "github.com/lejeunel/go-image-annotator/modules/image-ingester"
	e "github.com/lejeunel/go-image-annotator/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestHandleAuthError(t *testing.T) {
	group := grp.NewGroup(grp.NewGroupId(), "a-group")
	collection := clc.NewCollection(clc.NewCollectionId(),
		"a-collection",
		clc.WithGroup(group.Name))
	itr := NewTestingInteractor(&fk.CollectionRepo{
		Existing: []clc.Collection{collection},
	},
		WithAuth(&fk.Auth{ErrOnAuth: e.ErrAuthorization}))
	p := &FakePresenter{}
	itr.Execute(t.Context(), ig.Request{Collection: collection.Name}, p)
	assert.True(t, p.GotAuthErr)
	assert.False(t, p.GotSuccess)
}

func TestNonExistingCollectionShouldFail(t *testing.T) {
	p := &FakePresenter{}
	itr := NewTestingInteractor(&fk.CollectionRepo{ErrOnFind: e.ErrNotFound})
	itr.Execute(t.Context(), ig.Request{}, p)
	assert.True(t, p.GotNotFoundErr)
	assert.False(t, p.GotSuccess)
}

func TestIngestion(t *testing.T) {
	collection := clc.NewCollection(clc.NewCollectionId(), "a-collection")
	ingester := NewTestingIngester()
	itr := NewTestingInteractor(&fk.CollectionRepo{
		Existing: []clc.Collection{collection},
	})
	itr.Ingester = ingester

	user := u.NewUser("a-user")
	ctx := u.AppendUserToContext(t.Context(), user)
	reader := bytes.NewReader([]byte("image-data"))
	req := ig.Request{
		Collection: collection.Name,
		Labels:     []string{"a-label"},
		BoundingBoxes: []an.BoundingBoxRequest{
			{Label: "a-box-label", Xc: 1, Yc: 2, Width: 3, Height: 4},
		},
		Polygons: []an.PolygonRequest{
			{
				Label:  "a-polygon-label",
				Points: an.Points{Coordinates: [][2]float32{{0, 0}, {1, 0}, {1, 1}}},
			},
		},
		Meta: []m.MetaData{
			{Key: "a-key", Value: "a-value"},
			{Key: "another-key", Value: 1.5},
		},
		Reader: reader,
	}

	p := &FakePresenter{}
	itr.Execute(ctx, req, p)

	assert.True(t, p.GotSuccess)
	assert.NotNil(t, ingester.Got)
	assert.Equal(t, user.Id, ingester.Got.UserId)
	assert.Equal(t, collection.Name, ingester.Got.Collection)
	assert.Equal(t, req.Labels, ingester.Got.Labels)
	assert.Equal(t, req.BoundingBoxes, ingester.Got.BoundingBoxes)
	assert.Equal(t, req.Polygons, ingester.Got.Polygons)
	assert.Equal(t, req.Meta, ingester.Got.Meta)
	assert.Same(t, reader, ingester.Got.Reader)
}
