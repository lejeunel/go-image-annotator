package generic

import (
	"fmt"
)

type APIURLBuilder struct {
	APIVersionString string
}

func NewAPIURLBuilder(version string) *APIURLBuilder {
	return &APIURLBuilder{APIVersionString: version}
}

func (b *APIURLBuilder) Build(endpoint string) string {
	return fmt.Sprintf("/api/%v/%v", b.APIVersionString, endpoint)
}
