package shared

import (
	"fmt"
	rt "github.com/lejeunel/go-image-annotator/routes"
)

func MakeNewTaskMessage() string {
	return fmt.Sprintf(`Checks its progress in your <a class="underline" href=%v>task logs</a>`, rt.ListTasksUrl)
}
