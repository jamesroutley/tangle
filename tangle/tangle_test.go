package tangle_test

import (
	"regexp"
	"testing"

	"github.com/jamesroutley/tangle/tangle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTangle_DefaultConfig(t *testing.T) {
	cases := []struct {
		name        string
		sourceFiles []string
		expected    string
	}{
		{
			name:        "a basic example - one file using the default config",
			sourceFiles: []string{"testdata/basic_example.md"},
			expected:    "let numbers = [1, 2, 3];\n\nconsole.log(numbers);\n",
		},

		{
			name: "Two source files. We expect the code block contents to be concatenated",
			sourceFiles: []string{
				"testdata/two_files_a.md",
				"testdata/two_files_b.md",
			},
			expected: "print('hello')\n\nprint('world')\n",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tangler := tangle.NewTangler()
			output, err := tangler.Tangle(tc.sourceFiles...)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(output))
		})
	}
}

func TestTangle_LangugageFilter(t *testing.T) {
	cases := []struct {
		name        string
		sourceFiles []string
		language    string
		expected    string
	}{
		{
			name:        "One source file with an HTML language filter",
			sourceFiles: []string{"testdata/two_languages.md"},
			language:    "html",
			expected:    "<html></html>\n",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tangler := tangle.NewTangler(
				tangle.LanguageFilterOption("html"),
			)
			output, err := tangler.Tangle(tc.sourceFiles...)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(output))
		})
	}
}

func TestTangle_RegexFilter(t *testing.T) {
	cases := []struct {
		name        string
		sourceFiles []string
		match       string
		expected    string
	}{
		{
			name:        "One source file with an HTML language filter",
			sourceFiles: []string{"testdata/regex_example.md"},
			match:       `^\/\/ hash_table.c`,
			expected: `// hash_table.c
#include "hash_table.h"
`,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			re, err := regexp.Compile(tc.match)
			require.NoError(t, err)
			tangler := tangle.NewTangler(
				tangle.RegexFilterOption(re),
			)
			output, err := tangler.Tangle(tc.sourceFiles...)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, string(output))
		})
	}
}
