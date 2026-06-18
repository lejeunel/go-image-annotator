package presenters

type BaseAnnotoriousRequest struct {
	ImageId    string `json:"image_id"`
	Collection string `json:"collection"`
	Label      string `json:"label"`
}

type Bounds struct {
	MinX float32 `json:"minX"`
	MinY float32 `json:"minY"`
	MaxX float32 `json:"maxX"`
	MaxY float32 `json:"maxY"`
}

type Properties struct {
	Color string `json:"color"`
	Label string `json:"label"`
}

type AnnotoriousBody struct {
	Purpose string `json:"purpose"`
	Value   string `json:"value"`
}
