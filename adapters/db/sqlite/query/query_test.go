package query

import (
	"testing"
	"time"

	pa "github.com/lejeunel/go-image-annotator/shared/pagination"
	st "github.com/lejeunel/go-image-annotator/shared/testing"

	"github.com/stretchr/testify/assert"
)

var queryTests = []QueryTest{
	{"filter by collection",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", time.Now(), nil},
			{*st.IdFromInt(1), "second-collection", time.Now(), nil}},
		"collection:\"first-collection\"",
		"",
		*st.IdFromInt(0),
		1,
	},
	{"order by ingested_at",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", MustParseTime("2026-08-25"), nil},
			{*st.IdFromInt(1), "first-collection", MustParseTime("2026-08-01"), nil}},
		"",
		"ingested_at",
		*st.IdFromInt(1),
		2,
	},
	{"order by ingested_at descending",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", MustParseTime("2026-08-25"), nil},
			{*st.IdFromInt(1), "first-collection", MustParseTime("2026-08-01"), nil}},
		"",
		"ingested_at:desc",
		*st.IdFromInt(0),
		2,
	},
	{"filter by ingestion",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", MustParseTime("2026-08-01"), nil},
			{*st.IdFromInt(1), "first-collection", MustParseTime("2026-08-10"), nil}},
		"ingested_at>\"2026-08-02\"",
		"",
		*st.IdFromInt(1),
		1,
	},
	{"filter by metadata bool value",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", time.Now(), map[string]any{"active": true}},
			{*st.IdFromInt(1), "first-collection", time.Now(), nil}},
		"meta.active",
		"",
		*st.IdFromInt(0),
		1,
	},
	{"filter by metadata integer",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", time.Now(), map[string]any{"score": 1}},
			{*st.IdFromInt(1), "first-collection", time.Now(), map[string]any{"score": 10}}},
		"meta.score > 5",
		"",
		*st.IdFromInt(1),
		1,
	},
	{"filter by metadata date",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", time.Now(), map[string]any{"captured_at": "2026-08-01"}},
			{*st.IdFromInt(1), "first-collection", time.Now(), map[string]any{"captured_at": "2026-08-30"}}},
		"meta.captured_at < \"2026-08-15\"",
		"",
		*st.IdFromInt(0),
		1,
	},
	{"order by metadata date field",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", time.Now(), map[string]any{"captured_at": "2026-08-30"}},
			{*st.IdFromInt(1), "first-collection", time.Now(), map[string]any{"captured_at": "2026-08-01"}}},
		"",
		"meta.captured_at",
		*st.IdFromInt(1),
		2,
	},
	{"filter by metadata bool value with special char in name",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "first-collection", time.Now(), map[string]any{"is-active": true}},
			{*st.IdFromInt(1), "first-collection", time.Now(), nil}},
		"meta.is-active",
		"",
		*st.IdFromInt(0),
		1,
	},
	{"order by collection",
		[]QueryTestPayload{
			{*st.IdFromInt(0), "collection-v1", time.Now(), nil},
			{*st.IdFromInt(1), "collection-v2", time.Now(), nil}},
		"",
		"collection:desc",
		*st.IdFromInt(1),
		2,
	},
}

func TestFiltering(t *testing.T) {
	for _, tt := range queryTests {
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
