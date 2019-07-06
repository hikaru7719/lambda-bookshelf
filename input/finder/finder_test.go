package finder

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	cases := map[string]struct {
		isbn       string
		expectBook *Book
	}{
		"success": {
			isbn: "9784797394481",
			expectBook: &Book{
				ISBN:      "9784797394481",
				Title:     "DNSがよくわかる教科書",
				Author:    "株式会社日本レジストリサービス（JPRS）／渡邉結衣／佐藤新太／藤原和典／",
				Publisher: "SBクリエイティブ",
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			actualBook, _ := Find(tc.isbn)
			if !reflect.DeepEqual(actualBook, tc.expectBook) {
				t.Errorf("want %v, but actual %v", tc.expectBook, actualBook)
			}
		})
	}
}
