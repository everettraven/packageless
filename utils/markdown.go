package utils

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
)

//RenderMarkdown - renders markdown and outputs it to the console
func RenderMarkdown(input string, colors []string) error {
	// Set up the renderer
	bold := true
	italic := true
	var indent uint = 2
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStyles(
			ansi.StyleConfig{
				Heading: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Color:       &colors[0],
						BlockSuffix: "\n",
						Bold:        &bold,
					},
				},
				H1: ansi.StyleBlock{
					StylePrimitive: ansi.StylePrimitive{
						Prefix:          " ",
						Suffix:          " ",
						Color:           &colors[1],
						BackgroundColor: &colors[0],
						Bold:            &bold,
					},
				},
				Strong: ansi.StylePrimitive{
					Bold: &bold,
				},
				Emph: ansi.StylePrimitive{
					Italic: &italic,
					Color:  &colors[0],
				},
				List: ansi.StyleList{
					StyleBlock: ansi.StyleBlock{
						Indent: &indent,
					},
				},
				HorizontalRule: ansi.StylePrimitive{
					Color:  &colors[2],
					Format: "\n=====================\n",
				},
			},
		),
	)

	out, err := renderer.Render(input)

	if err != nil {
		return err
	}

	fmt.Print(out)

	return nil
}

//RenderInfoMarkdown - Renders markdown and outputs it with a color scheme specific to info messages
func (u *Utility) RenderInfoMarkdown(input string) {
	err := RenderMarkdown(input, []string{"45", "232", "231"})

	if err != nil {
		fmt.Println(input)
	}
}

//RenderErrorMarkdown - Renders markdown and outputs it with a color scheme specific to error messages
func (u *Utility) RenderErrorMarkdown(input string) {
	err := RenderMarkdown(input, []string{"88", "255", "231"})

	if err != nil {
		fmt.Println(input)
	}
}
