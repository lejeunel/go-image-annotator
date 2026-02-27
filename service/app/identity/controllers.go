package controllers

import (
	"context"
	a "datahub/app/authorizer"
)

type IdentityHTTPController struct {
	Authorizer *a.Authorizer
}
type IdentityOutputBody struct {
	Entitlements []string `json:"entitlements"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	Groups       []string `json:"groups"`
}
type IdentityOutput struct {
	Body IdentityOutputBody
}

func (h *IdentityHTTPController) Get(ctx context.Context, input *struct{}) (*IdentityOutput, error) {
	entitlements, err := h.Authorizer.IdentityProvider.Entitlements(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := h.Authorizer.IdentityProvider.Groups(ctx)
	if err != nil {
		return nil, err
	}
	username, err := h.Authorizer.IdentityProvider.Username(ctx)
	if err != nil {
		return nil, err
	}
	email, err := h.Authorizer.IdentityProvider.Email(ctx)
	if err != nil {
		return nil, err
	}
	body := IdentityOutputBody{Entitlements: entitlements,
		Username: username,
		Email:    email,
		Groups:   groups}
	return &IdentityOutput{Body: body}, nil
}
