package collections

import (
	"context"
	e "datahub/errors"

	g "datahub/generic"
)

type CollectionGetRequest struct {
	Name string `path:"name"`
}

type CollectionDeleteRequest struct {
	Name string `path:"name"`
}

type CollectionCreateRequest struct {
	Body struct {
		Name        string `json:"name" doc:"Name"`
		Description string `json:"description" doc:"Description" required:"false"`
		Group       string `json:"group" doc:"Group" required:"true"`
	}
}

type CollectionUpdatables struct {
	Name        string `json:"name" doc:"New collection name"`
	Description string `json:"description" doc:"New description"`
	Group       string `json:"group" doc:"Group"`
}

type CollectionUpdateRequest struct {
	Name string `path:"name"`
	Body CollectionUpdatables
}

type CollectionPatchRequest struct {
	Name string `path:"name"`
	Body g.JSONPatches
}

type CollectionCloneRequest struct {
	Id   string `path:"id"`
	Body struct {
		Name string `doc:"Name of new collection" json:"name"`
		Deep bool   `doc:"Also clone annotations" json:"deep"`
	}
}

type CollectionResponseBody struct {
	Id          string `json:"id" doc:"Id"`
	Name        string `json:"name" doc:"Name"`
	Description string `json:"description"`
	Group       string `json:"group"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
}

type CollectionResponse struct {
	Body CollectionResponseBody
}

type ProfileAssignmentRequest struct {
	ProfileName    string `path:"profilename"`
	CollectionName string `path:"collectionname"`
}

type ProfileUnassignRequest struct {
	CollectionName string `path:"collectionname"`
}

type ProfileAssignmentResponse struct{}

func makeCollectionResponse(c Collection) *CollectionResponse {
	return &CollectionResponse{Body: CollectionResponseBody{
		Id:          c.Id.String(),
		Name:        c.Name,
		Description: c.Description,
		Group:       c.Group,
		UpdatedAt:   c.UpdatedAt.String(),
		CreatedAt:   c.CreatedAt.String()}}
}

type CollectionHTTPController struct {
	CollectionService *Service
}

func (h *CollectionHTTPController) Get(ctx context.Context, input *CollectionGetRequest) (*CollectionResponse, error) {
	collection, err := h.CollectionService.FindByName(ctx, input.Name)

	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeCollectionResponse(*collection), nil
}

func (h *CollectionHTTPController) Put(ctx context.Context, input *CollectionUpdateRequest) (*CollectionResponse, error) {

	updatedCollection, err := h.CollectionService.Update(ctx, input.Name, input.Body)
	if err != nil {
		return nil, err
	}

	return makeCollectionResponse(*updatedCollection), nil
}

func (h *CollectionHTTPController) Patch(ctx context.Context, input *CollectionPatchRequest) (*CollectionResponse, error) {
	patchedCollection, err := h.CollectionService.Patch(ctx, input.Name, input.Body)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	return makeCollectionResponse(*patchedCollection), nil

}

func (h *CollectionHTTPController) Delete(ctx context.Context, input *CollectionDeleteRequest) (*struct{}, error) {
	collection, err := h.CollectionService.FindByName(ctx, input.Name)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	err = h.CollectionService.Delete(ctx, collection.Id)

	if err != nil {
		return nil, e.ToHumaStatusError(err)

	}

	return nil, nil
}

func (h *CollectionHTTPController) Create(ctx context.Context, input *CollectionCreateRequest) (*CollectionResponse, error) {

	collection, err := New(input.Body.Name, input.Body.Description, input.Body.Group)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	if err := h.CollectionService.Create(ctx, collection); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeCollectionResponse(*collection), nil
}

func (h *CollectionHTTPController) AssignProfile(ctx context.Context, input *ProfileAssignmentRequest) (*ProfileAssignmentResponse, error) {
	collection, err := h.CollectionService.FindByName(ctx, input.CollectionName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	profile, err := h.CollectionService.Profiles.FindByName(ctx, input.ProfileName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	if err := h.CollectionService.AssignProfile(ctx, profile, collection); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return &ProfileAssignmentResponse{}, nil
}

func (h *CollectionHTTPController) UnassignProfile(ctx context.Context, input *ProfileUnassignRequest) (*ProfileAssignmentResponse, error) {
	collection, err := h.CollectionService.FindByName(ctx, input.CollectionName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	if err := h.CollectionService.UnassignProfile(ctx, collection); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return &ProfileAssignmentResponse{}, nil
}
