package fetchall

import (
	"errors"
	"slices"
	"testing"

	e "github.com/lejeunel/go-image-annotator/shared/errors"
)

func TestHandleErrOnCount(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnCount: true, Err: e.ErrInternal})
	itr.Execute(p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error, but got none")
	}
}

func TestHandleErrWhenCountExceedsLimit(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{Count_: 2}, WithLimit(1))
	itr.Execute(p)
	if !errors.Is(p.GotErr, e.ErrLabelLimitExceeded) {
		t.Fatalf("expected label limit exceed error, but got %v", p.GotErr)
	}
}

func TestHandleErrOnFetch(t *testing.T) {
	p := &FakePresenter{}
	itr := NewInteractor(&FakeRepo{ErrOnFetch: true, Err: e.ErrInternal})
	itr.Execute(p)
	if !p.GotInternalErr || p.GotSuccess {
		t.Fatal("expected internal error, but got none")
	}
}

func TestFetchLabels(t *testing.T) {
	p := &FakePresenter{}
	labels := []string{"first-label", "second-labels"}
	itr := NewInteractor(&FakeRepo{Labels: labels})
	itr.Execute(p)
	if !p.GotSuccess {
		t.Fatal("expected success")
	}
	if !slices.Equal(p.Got.Labels, labels) {
		t.Fatalf("expected to retrieve %v, got %v", labels, p.Got.Labels)
	}
}

// func TestListLabel(t *testing.T) {
// 	count := 3
// 	pageSize := 2
// 	page := 1
// 	repo := &FakeRepo{Count_: count}
// 	p := &FakePresenter{}
// 	itr := NewInteractor(repo)
// 	itr.Execute(Request{PageSize: pageSize, Page: int64(page)}, p)
// 	if len(p.Got.Labels) != pageSize {
// 		t.Fatalf("expected to retrieve %v labels, got %v", pageSize, len(p.Got.Labels))
// 	}
// 	got := p.Got
// 	if int(got.Pagination.TotalRecords) != count {
// 		t.Fatalf("expected to retrieve count of %v, got %v", count, got.Pagination.TotalRecords)
// 	}
// 	if int(got.Pagination.TotalPages) != 2 {
// 		t.Fatalf("expected to retrieve total pages %v, got %v", 2, got.Pagination.TotalPages)
// 	}
// 	if int(got.Pagination.Page) != page {
// 		t.Fatalf("expected to retrieve page %v, got %v", page, got.Pagination.Page)
// 	}
// 	if int(got.Pagination.PageSize) != pageSize {
// 		t.Fatalf("expected to retrieve page %v, got %v", pageSize, got.Pagination.Page)
// 	}
// }
