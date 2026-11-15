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
	// missingLogPath receives the IDs of missing images. The file is only
	// created once the first missing image is found.
	missingLogPath string
	missingLog     *os.File
	missingLogErr  error
}

func NewDoctorCLIPresenter(missingLogPath string) *DoctorCLIPresenter {
	p := &DoctorCLIPresenter{missingLogPath: missingLogPath}
	p.progress = Progress{w: os.Stdout}
	return p
}

// logMissingImage appends id to the missing images log. On failure, the error
// is kept for Summarize and further images are not logged.
func (p *DoctorCLIPresenter) logMissingImage(id im.ImageId) {
	if p.missingLogErr != nil {
		return
	}
	if p.missingLog == nil {
		p.missingLog, p.missingLogErr = os.Create(p.missingLogPath)
		if p.missingLogErr != nil {
			return
		}
	}
	_, p.missingLogErr = fmt.Fprintln(p.missingLog, id)
}

func (p *DoctorCLIPresenter) closeMissingLog() {
	if p.missingLog == nil {
		return
	}
	if err := p.missingLog.Close(); err != nil && p.missingLogErr == nil {
		p.missingLogErr = err
	}
	p.missingLog = nil
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
	p.logMissingImage(id)
	p.progress.Clear()
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
	p.closeMissingLog()

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
	p.closeMissingLog()
	if p.missingLogErr != nil {
		fmt.Printf(
			"Missing images log %s: %s %v\n",
			p.missingLogPath,
			red("ERROR"),
			p.missingLogErr,
		)
	} else if p.NumMissingImages > 0 {
		fmt.Printf("%d missing image(s) written to %s\n", p.NumMissingImages, p.missingLogPath)
	}
	fmt.Printf("Performed %v exams and got %v error(s)\n", p.NumExams, p.NumErrors)
}
