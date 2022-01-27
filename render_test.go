package ukuleleweb

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRender(t *testing.T) {
	for _, tt := range []struct{ Input, Want string }{
		{
			Input: "Just a paragraph.",
			Want:  "<p>Just a paragraph.</p>\n",
		},
		{
			Input: "Hello *World*!",
			Want:  "<p>Hello <em>World</em>!</p>\n",
		},
		{
			Input: "Hello WikiLink!",
			Want:  `<p>Hello <a href="/WikiLink">WikiLink</a>!</p>` + "\n",
		},
		{
			Input: "WikiLink at start",
			Want:  `<p><a href="/WikiLink">WikiLink</a> at start</p>` + "\n",
		},
		{
			Input: "at the end a WikiLink",
			Want:  `<p>at the end a <a href="/WikiLink">WikiLink</a></p>` + "\n",
		},
		{
			Input: `<a href="http://wiki/">To the wiki!</a>`,
			Want:  `<p><a href="http://wiki/">To the wiki!</a></p>` + "\n",
		},
		{
			Input: "Hello go/go-link!",
			Want:  `<p>Hello <a href="http://go/go-link">go/go-link</a>!</p>` + "\n",
		},
		// {
		// 	// Should not recognize the inner mention of 'ExamplePage'.
		// 	Input: `<a href="http://wiki/ExamplePage">To the wiki!</a>`,
		// 	Want:  `<p><a href="http://wiki/ExamplePage">To the wiki!</a></p>` + "\n",
		// },
	} {
		got := renderHTML(tt.Input)
		if diff := cmp.Diff(got, tt.Want); diff != "" {
			t.Errorf("renderHTML(%q) = %q, want %q", tt.Input, got, tt.Want)
		}
	}
}
