package service

import "testing"

func TestCleanUrl(t *testing.T) {
	tcs := [][]string{
		[]string{
			"http://example.com/hello?a=1&b=2#anchor&hiho=3",
			"http://example.com/hello?a=1&b=2",
		},
		[]string{
			"example.com/hello?a=1&b=2#anchor&hiho=3",
			"http://example.com/hello?a=1&b=2",
		},
		[]string{
			"https://example.com/hello?a=1&b=2#anchor&hiho=3",
			"https://example.com/hello?a=1&b=2",
		},
		[]string{
			"https://example.com/hello",
			"https://example.com/hello",
		},
		[]string{
			"https://example.com",
			"https://example.com",
		},
		[]string{
			"https://example.com/",
			"https://example.com/",
		},
		// []string{
		// 	"https://example/",
		// 	"",
		// },
	}
	for _, tc := range tcs {
		u := tc[0]
		ex := tc[1]
		got := cleanUrl(u)
		if got != ex {
			t.Error("url not cleaned:", got)
			t.Fail()
		}
	}
}
