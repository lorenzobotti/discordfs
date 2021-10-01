package discordfs

import (
	"testing"
)

func TestPathElements(t *testing.T) {
	testCases := []struct {
		path     string
		elements []string
	}{
		{"go/pee/pee/poo/poo", []string{"go", "pee", "pee", "poo", "poo"}},
		{"/go/pee/pee/poo/poo", []string{"go", "pee", "pee", "poo", "poo"}},
	}
	for _, tC := range testCases {
		el := pathElements(tC.path)
		if !equalStringSlice(el, tC.elements) {
			t.Fatalf("input: %s expected \n%v\n got \n%v\n", tC.path, tC.elements, el)
		}

	}
}

func TestCleanPath(t *testing.T) {
	testCases := []struct {
		c, exp string
	}{
		{"/", "/"},
		{"mamma", "/mamma"},
		{"hamburger.png", "/hamburger.png"},
		{"./hamburger.png", "/hamburger.png"},
		{".///", "/"},
		{"/ciao/..", "/"},
		{"ciao/..", "/"},
		{"large/////automobile", "/large/automobile"},
		{"././am/i/ready/yet/hamburger", "/am/i/ready/yet/hamburger"},
	}

	for _, c := range testCases {
		cleaned := CleanPath(c.c)
		if cleaned != c.exp {
			t.Fatalf("expected %s, got %s", c.exp, cleaned)
		}
	}
}
func TestIsInFolder(t *testing.T) {
	testCases := []struct {
		c                  string
		expectedInFolder   bool
		expectedFolderName string
	}{
		{"/ciao.png", false, "/"},
		{"/pecunia non olet.pdf", false, "/"},
		{"/ciao/pecunia non olet.pdf", true, "ciao"},
		{"/ciao/cane magico/mammamia/", true, "ciao"},
		{"./cane magico/mammamia/", true, "cane magico"},
	}

	for _, c := range testCases {
		isIn := isInFolder(c.c)
		if isIn != c.expectedInFolder {
			t.Fatalf("isInFolder: \"%s\" expected %v, got %v", c.c, c.expectedInFolder, isIn)
		}

		top := topFolder(c.c)
		if top != c.expectedFolderName {
			t.Fatalf("topFolder: \"%s\" expected %v, got %v", c.c, c.expectedFolderName, top)
		}
	}
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func equalStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
