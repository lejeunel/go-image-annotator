package query

import (
	"testing"
	"time"

	pa "github.com/lejeunel/go-image-annotator/shared/pagination"

	"github.com/stretchr/testify/assert"
)

var tests = []FilterTest{
	{"filter by collection",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", time.Now(), nil},
			{IdFromInt(1), "second-collection", time.Now(), nil}},
		"collection:\"first-collection\"",
		"",
		IdFromInt(0),
		1,
	},
	{"order by ingested_at",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", MustParseTime("2026-08-25"), nil},
			{IdFromInt(1), "first-collection", MustParseTime("2026-08-01"), nil}},
		"",
		"ingested_at",
		IdFromInt(1),
		2,
	},
	{"order by ingested_at descending",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", MustParseTime("2026-08-25"), nil},
			{IdFromInt(1), "first-collection", MustParseTime("2026-08-01"), nil}},
		"",
		"ingested_at:desc",
		IdFromInt(0),
		2,
	},
	{"filter by ingestion",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", MustParseTime("2026-08-01"), nil},
			{IdFromInt(1), "first-collection", MustParseTime("2026-08-10"), nil}},
		"ingested_at>\"2026-08-02\"",
		"",
		IdFromInt(1),
		1,
	},
	{"filter by metadata bool value",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", time.Now(), map[string]any{"active": true}},
			{IdFromInt(1), "first-collection", time.Now(), nil}},
		"meta.active",
		"",
		IdFromInt(0),
		1,
	},
	{"filter by metadata integer",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", time.Now(), map[string]any{"score": 1}},
			{IdFromInt(1), "first-collection", time.Now(), map[string]any{"score": 10}}},
		"meta.score > 5",
		"",
		IdFromInt(1),
		1,
	},
	{"filter by metadata date",
		[]TestIngestionPayload{
			{IdFromInt(0), "first-collection", time.Now(), map[string]any{"captured_at": "2026-08-01"}},
			{IdFromInt(1), "first-collection", time.Now(), map[string]any{"captured_at": "2026-08-30"}}},
		"meta.captured_at < \"2026-08-15\"",
		"",
		IdFromInt(0),
		1,
	},
}

func TestFiltering(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr, imr, mr := Setup()
			InitFilterTest(imr, cr, mr, tt.images)
			count, err := imr.Count(tt.Filter)
			assert.NoError(t, err, tt.name)
			assert.Equal(t, int64(tt.WantCount), *count, tt.name)
			slice, err := imr.Slice(tt.Filter, pa.PaginationParams{Page: 1, PageSize: 10}, tt.Order)
			assert.NoError(t, err)
			assert.Equal(t, tt.WantFirstId.String(), slice[0].ImageId.String(), tt.name)
		})
	}
}
