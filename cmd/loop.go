package cmd

import (
	"github.com/grough/echoprint-go/echo"

	"github.com/spf13/cobra"
)

var loopCmd = &cobra.Command{
	Use:   "loop [input] [output]",
	Short: "Render a seamless loop",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		outputFile := args[1]
		tempo, _ := cmd.Flags().GetFloat64("tempo")
		duration, _ := cmd.Flags().GetFloat64("duration")
		delay, _ := cmd.Flags().GetFloat64("delay")

		renderer, err := echo.NewLoopRenderer(inputFile, outputFile, tempo, duration, delay)
		if err != nil {
			cmd.PrintErrln("Error creating renderer:", err)
			return
		}
		renderer.Render()
	},
}

func init() {
	rootCmd.AddCommand(loopCmd)

	loopCmd.Flags().Float64P("tempo", "t", 120.0, "Base tempo of the loop effect")
	loopCmd.Flags().Float64P("duration", "b", 8.0, "Output duration in beats")
	loopCmd.Flags().Float64P("delay", "d", 1.0, "Loop time in beats")
}
