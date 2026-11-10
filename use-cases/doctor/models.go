package doctor

type Request struct {
	InitialAdminEmail    string
	InitialAdminPassword string
}

type ExamResult struct {
	Description string
	Error       error
	Advice      string
}

type ExamName string

const (
	InitialAdminPassword ExamName = "initial-admin-password"
)

type Response struct {
	Exams map[ExamName]ExamResult
}

func NewDoctorResponse() Response {
	return Response{Exams: make(map[ExamName]ExamResult)}
}

func (r *Response) AddExam(examName ExamName, result ExamResult) {
	r.Exams[examName] = result
}
