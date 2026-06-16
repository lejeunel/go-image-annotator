package annotorious

type AnnotoriousPolygonRequest struct {
	BaseAnnotoriousRequest
	Annotation AnnotoriousPolygonModel `json:"annotation"`
}

type AnnotoriousPolygonModel struct {
	AnnotationId string            `json:"id"`
	Properties   Properties        `json:"properties"`
	Target       PolygonTarget     `json:"target"`
	Bodies       []AnnotoriousBody `json:"bodies"`
}

type PolygonGeometry struct {
	Bounds Bounds       `json:"bounds"`
	Points [][2]float64 `json:"points"`
}

type PolygonSelector struct {
	Type     string          `json:"type"`
	Geometry PolygonGeometry `json:"geometry"`
}

type PolygonTarget struct {
	Selector PolygonSelector `json:"selector"`
}
