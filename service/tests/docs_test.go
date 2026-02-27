package tests

import (
	d "datahub/docs"
	"strings"
	"testing"
)

func TestDocsNoSourceFilesShouldFail(t *testing.T) {

	builder := d.NewBuilder(d.NewGoldMarkParser())
	err := builder.Parse()
	AssertErrorIs(t, err, d.ErrPageListEmpty)
}

func TestDocsBuildOnePage(t *testing.T) {

	builder := d.NewBuilder(d.NewGoldMarkParser())
	page := strings.NewReader(`---
Title: Introduction
Weight: 0
ShortName: intro
---
`)
	builder.AddPage(page)
	err := builder.Parse()
	parsed := builder.Build()
	AssertNoError(t, err)
	if len(parsed) != 1 {
		t.Fatalf("expected to parse one page, got %v", len(parsed))
	}
}

func TestDocsBuildTwoPages(t *testing.T) {

	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: First Page
Weight: 0
ShortName: first
---
`))
	builder.AddPage(strings.NewReader(`---
Title: Second Page
Weight: 0
ShortName: second
---
`))
	err := builder.Parse()
	pages := builder.Build()
	AssertNoError(t, err)
	if len(pages) != 2 {
		t.Fatalf("expected to parse two pages, got %v", len(pages))
	}
}

func TestDocsMetaDataMissing(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
---
`))
	err := builder.Parse()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsWeightMissingFromMeta(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
ShortName: second
---
`))
	err := builder.Parse()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsShortNameMissingFromMeta(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
---
`))
	err := builder.Parse()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsWeightWrongTypeShouldFail(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
Weight: this-should-be-an-int
ShortName: short
---
`))
	err := builder.Parse()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsPagesAreSortedByDecreasingWeight(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: 2nd
Weight: 1
ShortName: second
---
`))
	builder.AddPage(strings.NewReader(`---
Title: 1st
Weight: 0
ShortName: first
---
`))
	builder.Parse()
	pages := builder.Build()
	if pages[0].Title != "1st" {
		t.Fatalf("expected to get first page with title 1st, got %v", pages[0].Title)
	}
}

func TestDocsParseContent(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
Weight: 0
ShortName: page
---
hello
`))
	builder.Parse()
	pages := builder.Build()
	if !strings.Contains(pages[0].Content.String(), "hello") {
		t.Fatal("expected to parse a content with string 'hello'")
	}
}

func TestDocsParseSummary(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
Weight: 0
ShortName: short
Summary: This is a summary
---
`))
	builder.Parse()
	pages := builder.Build()
	if !strings.Contains(pages[0].Summary, "summary") {
		t.Fatal("expected to parse a summary with string 'summary'")
	}
}
