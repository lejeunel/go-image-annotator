package cloner

import (
	im "github.com/lejeunel/go-image-annotator/entities/image"
	task "github.com/lejeunel/go-image-annotator/entities/task"
	"github.com/lejeunel/go-image-annotator/modules/copier"
	logger "github.com/lejeunel/go-image-annotator/modules/event-logger"
	"iter"
)

type TaskFuncBuilder interface {
	Build(task.CloneTask) (func(), error)
}

type ImageRepo interface {
	Iterate(im.FilteringParams, int) iter.Seq2[im.BaseImage, error]
}

type CloneTaskFuncBuilder struct {
	ImageRepo
	*logger.EventLogger
	copier.Copier
}

func (b *CloneTaskFuncBuilder) Build(t task.CloneTask) (func(), error) {
	// TODO iterate on each batch of source image
	return nil, nil
}

func NewCloneTaskFuncBuilder(r ImageRepo, l *logger.EventLogger, c copier.Copier) CloneTaskFuncBuilder {
	return CloneTaskFuncBuilder{r, l, c}
}
