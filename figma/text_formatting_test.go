package figma

import (
	"fmt"
	"slices"
	"testing"
)

var tests = [8]string{
	"hello_world_example",
	"hello world_example",
	"Hello world example",
	"hello_World-example",
	"hello-world EXAMPLE",
	"hello-world/ EXAMPLE",
	"helloWorldExample",
	"HelloWorldExample",
}

func TestNormaliseWords(t *testing.T) {
	want := []string{"hello", "world", "example"}

	for _, tt := range tests {
		testname := fmt.Sprintf("normaliseWords: %s", tt)
		t.Run(testname, func(t *testing.T) {
			ans := normaliseWords(tt)
			if !slices.Equal(ans, want) {
				t.Errorf("got %s, want %s", ans, want)
			}
		})
	}
}

func TestToPascalCase(t *testing.T) {
	want := "HelloWorldExample"

	for _, tt := range tests {
		testname := fmt.Sprintf("ToPascalCase: %s", tt)
		t.Run(testname, func(t *testing.T) {
			ans := ToPascalCase(tt)
			if ans != want {
				t.Errorf("got %s, want %s", ans, want)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	want := "helloWorldExample"

	for _, tt := range tests {
		testname := fmt.Sprintf("ToCamelCase: %s", tt)
		t.Run(testname, func(t *testing.T) {
			ans := ToCamelCase(tt)
			if ans != want {
				t.Errorf("got %s, want %s", ans, want)
			}
		})
	}
}

func TestToKebabCase(t *testing.T) {
	want := "hello-world-example"

	for _, tt := range tests {
		testname := fmt.Sprintf("ToKebabCase: %s", tt)
		t.Run(testname, func(t *testing.T) {
			ans := ToKebabCase(tt)
			if ans != want {
				t.Errorf("got %s, want %s", ans, want)
			}
		})
	}
}

func TestTokenValues(t *testing.T) {
	var values = []struct {
		s, wv, wt string
	}{
		{"TEst", "--test", ":root"},
		{"foo Baz bar", "--foo-baz-bar", ":root"},
		{"foo / bar", "--foo-bar", ":root"},
		{"foo theme/bar", "--bar", "foo-theme"},
	}

	for _, tt := range values {
		testname := fmt.Sprintf("TokenValues: %s", tt.s)
		t.Run(testname, func(t *testing.T) {
			variable, theme := TokenValues(tt.s)
			if variable != tt.wv || theme != tt.wt {
				t.Errorf("got (%s, %s), want (%s, %s)", variable, theme, tt.wv, tt.wt)
			}
		})
	}
}
