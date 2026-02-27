package labels

import (
	"context"

	e "datahub/errors"
)

type LabelsHTTPController struct {
	LabelsService *Service
}

type LabelResponseBody struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `doc:"Description" json:"description"`
}

type LabelOutputResponse struct {
	Body LabelResponseBody
}

type LabelCreateRequest struct {
	Body struct {
		Name        string `doc:"Name" json:"name"`
		Description string `doc:"Description" json:"description" required:"false"`
	}
}

type LabelUpdateRequest struct {
	Name string `path:"name"`
	Body Updatables
}

func makeLabelResponse(l Label) *LabelOutputResponse {
	return &LabelOutputResponse{Body: LabelResponseBody{
		Id:          l.Id.String(),
		Name:        l.Name,
		Description: l.Description,
	}}
}

type LabelDeleteRequest struct {
	Name string `path:"name"`
}

func (h *LabelsHTTPController) Create(ctx context.Context, input *LabelCreateRequest) (*LabelOutputResponse, error) {
	label, err := New(input.Body.Name, input.Body.Description)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	if err := h.LabelsService.Create(ctx, label); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeLabelResponse(*label), nil

}

func (h *LabelsHTTPController) Update(ctx context.Context, input *LabelUpdateRequest) (*LabelOutputResponse, error) {
	label, err := h.LabelsService.FindByName(ctx, input.Name)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	if err := h.LabelsService.Update(ctx, label, Updatables{Description: input.Body.Description}); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return makeLabelResponse(*label), nil

}

func (h *LabelsHTTPController) Delete(ctx context.Context, input *LabelDeleteRequest) (*struct{}, error) {
	label, err := h.LabelsService.FindByName(ctx, input.Name)
	if err != nil {
		return nil, e.ToHumaStatusError(err)
	}
	if err := h.LabelsService.Delete(ctx, label); err != nil {
		return nil, e.ToHumaStatusError(err)
	}

	return nil, nil

}
