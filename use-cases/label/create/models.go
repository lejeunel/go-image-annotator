package create

import (
	lbl "github.com/lejeunel/go-image-annotator/entities/label"
)

type Response struct {
	Name        string
	Description string
}

type Request struct {
	Name        string
	Description string
}

type CreateModel struct {
	Id          lbl.LabelId
	Name        string
	Description string
}
