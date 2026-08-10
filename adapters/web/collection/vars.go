package collection

import (
	"fmt"
	d "github.com/lejeunel/go-image-annotator/adapters/web/dashboard"
)

const (
	createCollectionTargetDiv  = "create-collection"
	CollectionUrl              = "/ui/collection"
	CreateCollectionFormUrl    = "/ui/collection/new"
	nameFieldName              = "name"
	descriptionFieldName       = "description"
	groupFieldName             = "group"
	deepFieldName              = "with_annotations"
	resourceUrlFieldName       = "name"
	publicGroupPlaceholderName = "( public )"
)

func MakeNewTaskMessage() string {
	return fmt.Sprintf(`Checks its progress in your <a class="underline" href=%v>task logs</a>`, d.ListTasksUrl)
}
