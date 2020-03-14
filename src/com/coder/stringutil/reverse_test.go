package stringutil

import "testing"

func TestReverse(t *testing.T) {
    cases := []struct {
        input          string
        expectedOutput string
    }{
        {"", ""},
        {"a", "a"},
        {"ab", "ba"},
        {"abc", "cba"},
        {"abcd", "dcba"},
        {"aibohphobia", "aibohphobia"},
    }

    for _, c := range cases {
        if output := Reverse(c.input); output != c.expectedOutput {
            t.Errorf("incorrect output for `%s`: expected `%s` but got `%s`", c.input, c.expectedOutput, output)
        }
    }
}