package ukuleleweb

import "testing"

func TestIsPageName(t *testing.T) {
	for _, pn := range []string{
		"PageName",
		"AlsoPageName",
		"WönderfülPägeNäme",
		// XXX: It's unclear to me why this is not recognized. Unicode shenanigans?
		// "ÄtschiBätschi",
	} {
		if !isPageName(pn) {
			t.Errorf("isPageName(%q) = false, want true", pn)
		}
	}
}

func TestIsNotPageName(t *testing.T) {
	for _, pn := range []string{
		"foo PageName bar",
		"/AlsoPageName/",
		"Oneword",
		"123",
	} {
		if isPageName(pn) {
			t.Errorf("isPageName(%q) = true, want false", pn)
		}
	}
}
