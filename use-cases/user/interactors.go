package user

import (
	agr "github.com/lejeunel/go-image-annotator/use-cases/user/assign-group"
	aro "github.com/lejeunel/go-image-annotator/use-cases/user/assign-role"
	"github.com/lejeunel/go-image-annotator/use-cases/user/create"
	"github.com/lejeunel/go-image-annotator/use-cases/user/delete"
	"github.com/lejeunel/go-image-annotator/use-cases/user/list"
	"github.com/lejeunel/go-image-annotator/use-cases/user/read"
	rat "github.com/lejeunel/go-image-annotator/use-cases/user/renew-access-token"
	ugr "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-group"
	uro "github.com/lejeunel/go-image-annotator/use-cases/user/unassign-role"
)

type Interactors struct {
	Find            read.Interactor
	Create          create.Interactor
	Delete          delete.Interactor
	List            list.Interactor
	RenewToken      rat.Interactor
	AssignRole      aro.Interactor
	UnAssignRole    uro.Interactor
	AssignGroup     agr.Interactor
	UnAssignGroup   ugr.Interactor
	DefaultPageSize int
}
