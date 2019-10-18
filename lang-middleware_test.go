package langmiddleware

import (
	"testing"
)

type testTuple struct {
	value, expected interface{}
}

func TestParseLangHeader(t *testing.T) {

	i, err := NewCookieAndHeader("en", []string{"en", "fr"}, "cookie-lang")

	if err != nil {
		t.Fatal("Failed to initialize LangMiddleware")
	}

	testValues := []testTuple{
		testTuple{value: "", expected: "en"},
		testTuple{value: "*", expected: "en"},
		testTuple{value: "en", expected: "en"},
		testTuple{value: "en-US", expected: "en"},
		testTuple{value: "en, en-US", expected: "en"},
		testTuple{value: "fr-BE;q=0.8, en-US;q=0.9, fr-BE;q=0.6", expected: "en"},
		testTuple{value: "de;q=0.85, fr;q=0.9, en;q=0.8, *;q=1", expected: "en"},
	}
	var tmp string
	for _, testCase := range testValues {
		tmp = i.fromLangHeader(testCase.value.(string))
		if testCase.expected != tmp {
			t.Errorf("Error while testing \"%v\" got \"%v\" but expected \"%v\"", testCase.value, tmp, testCase.expected)
		}
	}
}

func benchmarkParseLangHeader(i *LangMiddleware, header string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		i.fromLangHeader(header)
	}
}

func BenchmarkParseLangHeaderSimple(b *testing.B) {
	i, _ := NewCookieAndHeader("en", []string{"en", "fr"}, "cookie-lang")
	benchmarkParseLangHeader(i, "en", b)
}

func BenchmarkParseLangHeaderLocale(b *testing.B) {
	i, _ := NewCookieAndHeader("en", []string{"en", "fr"}, "cookie-lang")
	benchmarkParseLangHeader(i, "en-US", b)
}

func BenchmarkParseLangHeaderMultiple(b *testing.B) {
	i, _ := NewCookieAndHeader("en", []string{"en", "fr"}, "cookie-lang")
	benchmarkParseLangHeader(i, "en, en-US", b)
}

func BenchmarkParseLangHeaderPonderated(b *testing.B) {
	i, _ := NewCookieAndHeader("en", []string{"en", "fr"}, "cookie-lang")
	benchmarkParseLangHeader(i, "en-US, fr-BE;q=0.8", b)
}
