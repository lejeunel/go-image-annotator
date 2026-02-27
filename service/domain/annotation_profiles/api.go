package annotation_profile

import (
	"context"
	e "datahub/errors"
	g "datahub/generic"
)

type ProfileGetRequest struct {
	Name string `path:"name"`
}

type ProfileListingResponseBody struct {
	Pagination g.PaginationMeta      `json:"pagination"`
	Profiles   []ProfileResponseBody `json:"profiles"`
}
type ProfileListingResponse struct {
	Body ProfileListingResponseBody
}

type ProfileDeleteRequest struct {
	Name string `path:"name"`
}

type ProfileCreateRequest struct {
	Body struct {
		Name   string   `json:"name" doc:"Name" required:"true"`
		Labels []string `json:"labels" doc:"List of labels" required:"true"`
	}
}

type ProfileUpdateRequest struct {
	Name string `path:"name"`
	Body ProfileUpdatables
}

type AddLabelRequest struct {
	ProfileName string `path:"profile" doc:"profile to modify"`
	LabelName   string `path:"label" doc:"name of label to add"`
}

type RemoveLabelRequest struct {
	ProfileName string `path:"profile" doc:"profile to modify"`
	LabelName   string `path:"label" doc:"name of label to remove"`
}

type ProfileResponseBody struct {
	Id     string   `json:"id" doc:"Id"`
	Name   string   `json:"name" doc:"Name"`
	Labels []string `json:"labels"`
}

type ProfilePatchRequest struct {
	Name string `path:"name"`
	Body g.JSONPatches
}

type ProfileResponse struct {
	Body ProfileResponseBody
}

func makeProfileResponse(p AnnotationProfile) *ProfileResponse {
	return &ProfileResponse{Body: ProfileResponseBody{
		Id:     p.Id.String(),
		Name:   p.Name,
		Labels: p.LabelNames()},
	}
}

type ProfileHTTPController struct {
	ProfileService *Service
}

func (h *ProfileHTTPController) List(ctx context.Context, pagination *g.PaginationParams) (*ProfileListingResponse, error) {
	profiles, meta, err := h.ProfileService.List(ctx, *pagination)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	var profileResponses []ProfileResponseBody
	for _, profile := range profiles {
		profileResponses = append(profileResponses,
			ProfileResponseBody{
				Id:     profile.Id.String(),
				Name:   profile.Name,
				Labels: profile.LabelNames(),
			})
	}

	return &ProfileListingResponse{Body: ProfileListingResponseBody{Pagination: *meta, Profiles: profileResponses}}, nil
}

func (h *ProfileHTTPController) Get(ctx context.Context, input *ProfileGetRequest) (*ProfileResponse, error) {
	profile, err := h.ProfileService.FindByName(ctx, input.Name)

	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeProfileResponse(*profile), nil
}

func (h *ProfileHTTPController) Put(ctx context.Context, input *ProfileUpdateRequest) (*ProfileResponse, error) {

	profile, err := h.ProfileService.FindByName(ctx, input.Name)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	updatedProfile, err := h.ProfileService.Update(ctx, profile.Id, input.Body)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeProfileResponse(*updatedProfile), nil
}

func (h *ProfileHTTPController) Delete(ctx context.Context, input *ProfileDeleteRequest) (*struct{}, error) {
	profile, err := h.ProfileService.FindByName(ctx, input.Name)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	err = h.ProfileService.Delete(ctx, profile)

	if err != nil {
		return nil, e.ToHumaStatusError(err)

	}

	return nil, nil
}

func (h *ProfileHTTPController) Create(ctx context.Context, input *ProfileCreateRequest) (*ProfileResponse, error) {

	profile := New(input.Body.Name)
	if err := h.ProfileService.Save(ctx, profile); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	for _, labelName := range input.Body.Labels {
		label, err := h.ProfileService.Labels.FindByName(ctx, labelName)
		if err != nil {
			return nil, e.ToHumaStatusError(err)
		}
		if err := h.ProfileService.AddLabel(ctx, profile, label); err != nil {
			return nil, e.ToHumaStatusError(err)
		}
	}

	return makeProfileResponse(*profile), nil

}

func (h *ProfileHTTPController) AddLabel(ctx context.Context, input *AddLabelRequest) (*ProfileResponse, error) {

	profile, err := h.ProfileService.FindByName(ctx, input.ProfileName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	label, err := h.ProfileService.Labels.FindByName(ctx, input.LabelName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	if err := h.ProfileService.AddLabel(ctx, profile, label); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeProfileResponse(*profile), nil

}

func (h *ProfileHTTPController) RemoveLabel(ctx context.Context, input *RemoveLabelRequest) (*ProfileResponse, error) {

	profile, err := h.ProfileService.FindByName(ctx, input.ProfileName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	label, err := h.ProfileService.Labels.FindByName(ctx, input.LabelName)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	if err := h.ProfileService.RemoveLabel(ctx, profile, label); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeProfileResponse(*profile), nil

}
