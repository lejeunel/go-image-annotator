package ingestion

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
)

type IngestionHTTPController struct {
	IngestionService *Service
}

type IngestionRequest struct {
	CollectionName string `path:"name" doc:"Name of collection"`
	Body           ImageIngestionPayload
}

type IngestionResponseBody struct {
	Id           string `json:"id"`
	CollectionId string `json:"collection_id"`
	Camera       string `json:"camera"`
	Site         string `json:"site"`
	Transmitter  string `json:"transmitter"`
	CapturedAt   string `json:"captured_at"`
	CreatedAt    string `json:"created_at"`
	Type_        string `json:"type"`
	Group        string `json:"group"`
}
type IngestionResponse struct {
	Body IngestionResponseBody
}

func (h *IngestionHTTPController) Ingest(ctx context.Context, input *IngestionRequest) (*IngestionResponse, error) {
	image, err := h.IngestionService.Ingest(ctx, input.CollectionName, input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error(), err)
	}

	return &IngestionResponse{Body: IngestionResponseBody{
		Id:           image.Id.String(),
		CollectionId: image.Collection.Id.String(),
		Camera:       image.GetCameraName(),
		Site:         image.GetSiteName(),
		Transmitter:  image.GetTransmitter(),
		CapturedAt:   image.CapturedAt.String(),
		CreatedAt:    image.CreatedAt.String(),
		Type_:        image.Type,
		Group:        image.Group,
	},
	}, nil

}
