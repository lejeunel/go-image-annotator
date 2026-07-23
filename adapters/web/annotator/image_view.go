package annotator

import (
	"encoding/base64"
	"fmt"
	"io"

	"github.com/lejeunel/go-image-annotator/modules/annotator/view"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

type ImageView struct {
	result Node
}

func (p *ImageView) Build(image view.Image) Node {
	if image.Reader == nil {
		return Text("presenting image: got no reader")

	}
	bytes, err := io.ReadAll(image.Reader)

	if err != nil {
		return Text(err.Error())
	}

	b64Image := base64.StdEncoding.EncodeToString(bytes)
	return Img(ID("image"), Src(fmt.Sprintf("data:%v;base64,%s",
		image.MIMEType, b64Image)))

}
