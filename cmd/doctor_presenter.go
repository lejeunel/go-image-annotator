package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"

	im "github.com/lejeunel/go-image-annotator/entities/image"
	dr "github.com/lejeunel/go-image-annotator/use-cases/doctor"
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

// progressRefresh is the minimum delay between two redraws of the progress
// block.
const progressRefresh = 100 * time.Millisecond

// Progress redraws a block of lines in place on a terminal.
type Progress struct {
	w     io.Writer
	lines int // lines printed by the previous render
}

// Render replaces the previously rendered block with lines.
func (p *Progress) Render(lines ...string) {
	p.Clear()
	for _, l := range lines {
		fmt.Fprintf(p.w, "%s\n", l)
	}
	p.lines = len(lines)
}

// Clear erases the block, so that other output can take its place.
func (p *Progress) Clear() {
	if p.lines > 0 {
		// Back to the top of the block, then clear to the end of the screen.
		fmt.Fprintf(p.w, "\033[%dA\033[J", p.lines)
	}
	p.lines = 0
}

type DoctorCLIPresenter struct {
	NumErrors        int
	NumExams         int
	NumImagesChecked int
	NumMissingImages int
	progress         Progress
	lastRender       time.Time
	lastCheck        time.Time
}

func NewDoctorCLIPresenter() *DoctorCLIPresenter {
	p := &DoctorCLIPresenter{}
	p.progress = Progress{w: os.Stdout}
	return p
}

func (p *DoctorCLIPresenter) PresentExam(r dr.Response) {
	p.NumExams += 1

	if r.Error == nil {
		fmt.Printf("%s: %s\n", r.Description, green("OK"))
		return
	}
	p.NumErrors += 1
	fmt.Printf("%s: %s %v\n", r.Description, red("ERROR"), r.Error)
}

func (p *DoctorCLIPresenter) SuccessCheckImageIntegrity(id im.ImageId) {
	p.lastCheck = time.Now()
	p.NumImagesChecked += 1
	p.renderImageIntegrity(false)
}

// ErrorCheckImageIntegrity prints the missing image above the progress block,
// so that it stays on screen.
func (p *DoctorCLIPresenter) ErrorCheckImageIntegrity(id im.ImageId) {
	p.lastCheck = time.Now()
	p.NumImagesChecked += 1
	p.NumMissingImages += 1
	p.progress.Clear()
	fmt.Printf("Image %s: %s file not found in image store\n", id, red("ERROR"))
	p.renderImageIntegrity(true)
}

func (p *DoctorCLIPresenter) imageIntegrityLine() string {
	missing := fmt.Sprint(p.NumMissingImages)
	if p.NumMissingImages > 0 {
		missing = red(missing)
	}

	rate := 1 / time.Since(p.lastCheck).Seconds()
	return fmt.Sprintf("Image store integrity: %d image(s) checked, %s missing [%.2f it/sec]",
		p.NumImagesChecked, missing, rate)
}

func (p *DoctorCLIPresenter) renderImageIntegrity(force bool) {
	if !force && time.Since(p.lastRender) < progressRefresh {
		return
	}
	p.progress.Render(p.imageIntegrityLine())
	p.lastRender = time.Now()
}

func (p *DoctorCLIPresenter) DoneCheckImageIntegrity() {
	p.progress.Clear()

	p.NumExams += 1
	status := green("OK")
	if p.NumMissingImages > 0 {
		p.NumErrors += 1
		status = red("ERROR")
	}
	fmt.Printf("%v: %v\n", p.imageIntegrityLine(), status)
}

// Summarize prints a summary of exams and error counts.
func (p *DoctorCLIPresenter) Summarize() {
	fmt.Printf("Performed %v exams and got %v error(s)\n", p.NumExams, p.NumErrors)
}
