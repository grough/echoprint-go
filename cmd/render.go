package cmd

import (
	"github.com/grough/echoprint-go/echo"

	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render [input] [output]",
	Short: "Render delay repeats",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile := args[1]
		tempo, _ := cmd.Flags().GetFloat64("tempo")
		duration, _ := cmd.Flags().GetFloat64("duration")
		delay, _ := cmd.Flags().GetFloat64("delay")

		renderer, err := echo.NewRenderer(inputFile, outputFile, tempo, duration, delay)
		if err != nil {
			cmd.PrintErrln("Error creating renderer:", err)
			return
		}
		renderer.Render()
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().Float64P("tempo", "t", 120, "Base tempo of the delay effect")
	renderCmd.Flags().Float64P("duration", "b", 8*4, "Output duration in beats")
	renderCmd.Flags().Float64P("delay", "d", 1, "Delay time in beats")
}
