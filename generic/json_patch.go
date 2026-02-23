package generic

import (
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
)

type JSONPatch struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     string `json:"value"`
}

type JSONPatches []JSONPatch

func (p *JSONPatches) Apply(payload []byte) ([]byte, error) {
	patch, err := p.Decode()
	if err != nil {
		return nil, fmt.Errorf("applying patch: %w", err)
	}

	modifiedBytes, err := patch.Apply(payload)
	if err != nil {
		return nil, fmt.Errorf("applying patch: %w", err)
	}
	return modifiedBytes, nil

}

func (p *JSONPatches) Decode() (*jsonpatch.Patch, error) {
	patchesString, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch([]byte(patchesString))
	if err != nil {
		return nil, err
	}
	return &patch, nil

}
