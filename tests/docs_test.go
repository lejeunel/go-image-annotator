package tests

import (
	d "datahub/docs"
	"strings"
	"testing"
)

func TestDocsNoSourceFilesShouldFail(t *testing.T) {

	builder := d.NewBuilder(d.NewGoldMarkParser())
	_, err := builder.Build()
	AssertErrorIs(t, err, d.ErrPageListEmpty)
}

func TestDocsBuildOnePage(t *testing.T) {

	builder := d.NewBuilder(d.NewGoldMarkParser())
	page := strings.NewReader(`---
Title: Introduction
Weight: 0
---
`)
	builder.AddPage(page)
	parsed, err := builder.Build()
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
---
`))
	builder.AddPage(strings.NewReader(`---
Title: Second Page
Weight: 0
---
`))
	parsed, err := builder.Build()
	AssertNoError(t, err)
	if len(parsed) != 2 {
		t.Fatalf("expected to parse two pages, got %v", len(parsed))
	}
}

func TestDocsMetaDataMissing(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
---
`))
	_, err := builder.Build()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsWeightMissingFromMeta(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
---
`))
	_, err := builder.Build()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsWeightWrongTypeShouldFail(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
Weight: this-should-be-an-int
---
`))
	_, err := builder.Build()
	AssertErrorIs(t, err, d.ErrMetaDataParsing)
}

func TestDocsPagesAreSortedByDecreasingWeight(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: 2nd
Weight: 1
---
`))
	builder.AddPage(strings.NewReader(`---
Title: 1st
Weight: 0
---
`))
	parsed, _ := builder.Build()
	if parsed[0].Title != "1st" {
		t.Fatalf("expected to get first page with title 1st, got %v", parsed[0].Title)
	}
}

func TestDocsParseContent(t *testing.T) {
	builder := d.NewBuilder(d.NewGoldMarkParser())
	builder.AddPage(strings.NewReader(`---
Title: The Title
Weight: 0
---
hello
`))
	parsed, _ := builder.Build()
	if !strings.Contains(parsed[0].Content.String(), "hello") {
		t.Fatal("expected to parse a content with string 'hello'")
	}
}
