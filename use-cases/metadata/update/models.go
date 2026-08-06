package update

type Request struct {
	ImageId    string
	Collection string
	Key        string
	Value      any
}

type Response struct{}
