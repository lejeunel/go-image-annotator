package tests

import (
	"embed"
	"errors"
	"fmt"
	"github.com/go-test/deep"
	"testing"
)

//go:embed test-data/sample-image.png
var testPNGImage []byte

//go:embed test-data/sample-image.jpg
var testJPGImage []byte

//go:embed test-data/sample-image.jpg
//go:embed test-data/sample-image.png
var testFS embed.FS

func AssertErrorIs(t testing.TB, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("wanted %v but got %v", target, err)
	}
}

func AssertError(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Error("wanted an error but didn't get one")
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Error(fmt.Printf("did not want an error but got: %v", err))
	}
}

func AssertDeepEqual(t testing.TB, this any, that any, structName string) {
	diff := deep.Equal(this, that)
	if diff != nil {
		t.Fatalf(fmt.Sprintf("expected to retrieve identical %v structs, but got different fields: %v",
			structName, diff))
	}

}
