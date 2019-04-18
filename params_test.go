package main

import (
	"testing"
)

const imagePlaceholder = "{IMAGE}"

var prepareCommandTests = []struct {
	input        Meta
	pdf          string
	expectedArgs []string
}{
	{Meta{Page: 1, Resolutions: Dimensions{
		{Width: 64, Height: 100},
	}},
		"portrait.pdf",
		[]string{"portrait.pdf[0]", "-flatten", "-thumbnail", "64x100!", imagePlaceholder},
	},
	{Meta{Page: 7, Resolutions: Dimensions{
		{Width: 800, Height: 600},
		{Width: 1024, Height: 768},
	}},
		"screenshot.pdf",
		[]string{"screenshot.pdf[6]", "-flatten", "-thumbnail", "1024x768!", "-write", imagePlaceholder,
			"-thumbnail", "800x600!", imagePlaceholder},
	},
	{Meta{Page: 1, Resolutions: Dimensions{
		{Width: 100, Height: 150},
		{Width: 300, Height: 450},
		{Width: 200, Height: 300},
	}},
		"input.pdf",
		[]string{"input.pdf[0]", "-flatten", "-thumbnail", "300x450!", "-write", imagePlaceholder,
			"-thumbnail", "200x300!", "-write", imagePlaceholder,
			"-thumbnail", "100x150!", imagePlaceholder},
	},
}

func TestPrepareCommand(t *testing.T) {
	for _, test := range prepareCommandTests {
		args, files, err := test.input.prepareCommand(test.pdf)
		if err != nil {
			t.Errorf("prepare command: %v", err)
			break
		}

		// replace image file name placeholders with real file names
		for _, file := range files {
			for j := range test.expectedArgs {
				if test.expectedArgs[j] == imagePlaceholder {
					test.expectedArgs[j] = file.Name()
					break
				}
			}
		}

		if len(args) != len(test.expectedArgs) {
			t.Errorf("expected arg length %d, was %d", len(test.expectedArgs), len(args))
			break
		}
		for i, expectedArg := range test.expectedArgs {
			if args[i] != expectedArg {
				t.Errorf(`epected arg[%d]=="%s", was "%s"`, i, expectedArg, args[i])
			}
		}
	}
}
