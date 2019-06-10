package plugin

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var spacerTests = map[string]struct {
	count  int
	spaces string
}{
	"count: 4": {
		count:  4,
		spaces: "    ",
	},
	"count: 5": {
		count:  5,
		spaces: "     ",
	},
	"count: 2": {
		count:  2,
		spaces: "  ",
	},
	"count: 8": {
		count:  8,
		spaces: "        ",
	},
	"count: 3": {
		count:  3,
		spaces: "   ",
	},
	"count: 9": {
		count:  9,
		spaces: "         ",
	},
}

func Test_indentSpacer(t *testing.T) {
	t.Parallel()

	spacer := &indentSpacer{}

	for name, tt := range spacerTests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			spacer.WriteTo(&buf, tt.count)
			if buf.String() != tt.spaces {
				t.Errorf("indentSpacer.WriteTo => %#v, want %#v", buf.String(), tt.spaces)
			}

			var str strings.Builder
			spacer.WriteTo(&str, tt.count)
			if str.String() != tt.spaces {
				t.Errorf("indentSpacer.WriteTo => %#v, want %#v", str.String(), tt.spaces)
			}
		})
	}
}

var helps = []struct {
	help *Help
	doc  string
}{
	{
		help: &Help{
			Name:        "test01",
			Description: "single line",
		},
		doc: `test01: single line`,
	},
	{
		help: &Help{
			Name:        "test02",
			Description: "single line",
			Commands: []*Command{
				&Command{
					Command:     "command01",
					Description: "single line 01",
				},
				&Command{
					Command:     "command0002",
					Description: "single line 02",
				},
			},
		},
		doc: `test02: single line
    command01=>   single line 01
    command0002=> single line 02`,
	},
	{
		help: &Help{
			Name: "test0003",
			Description: `first line
second line
third line`,
		},
		doc: `test0003: first line
          second line
          third line`,
	},
	{
		help: &Help{
			Name: "test0004",
			Description: `first line
second line
third line`,
			Commands: []*Command{
				&Command{
					Command:     "command01",
					Description: "single line",
				},
				&Command{
					Command: "command0002",
					Description: `first line
second line
third line`,
				},
			},
		},
		doc: `test0004: first line
          second line
          third line
    command01=>   single line
    command0002=> first line
                  second line
                  third line`,
	},
	{
		help: &Help{
			Name: "テスト00005",
			Description: `first line
second line
third line`,
		},
		doc: `テスト00005: first line
             second line
             third line`,
	},
	{
		help: &Help{
			Name: "テスト00006",
			Description: `first line
second line
third line`,
			Commands: []*Command{
				&Command{
					Command:     "コマンド01",
					Description: "single line",
				},
				&Command{
					Command: "command0002",
					Description: `first line
second line
third line`,
				},
			},
		},
		doc: `テスト00006: first line
             second line
             third line
    コマンド01=>  single line
    command0002=> first line
                  second line
                  third line`,
	},
}

func Test_HelpFormatter(t *testing.T) {
	t.Parallel()

	formatter := &HelpFormatter{
		NameSuffix:    ": ",
		CommandSuffix: "=> ",
		Indent:        4,
	}

	for _, tt := range helps {
		tt := tt
		t.Run(tt.help.Name, func(t *testing.T) {
			doc := formatter.Format(tt.help)
			diff := cmp.Diff(doc, tt.doc)
			if diff != "" {
				t.Errorf("failed to formatHelp(%s), differs: (-got +want)\n%s", tt.help.Name, diff)
			}
		})
	}
}
