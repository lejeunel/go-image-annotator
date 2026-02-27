package labels

import (
	g "datahub/generic"
)

type Repo interface {
	Create(*Label) error
	Delete(*Label) error
	Update(*Label, Updatables) error

	SetParenting(*Label, *Label) error
	FindByName(string) (*Label, error)
	Find(LabelId) (*Label, error)

	List(g.OrderingArg, g.PaginationParams) ([]Label, *g.PaginationMeta, error)
	Count() (int64, error)
}
