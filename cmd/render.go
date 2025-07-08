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
		tempo, _ := cmd.Flags().GetInt16("tempo")
		bars, _ := cmd.Flags().GetInt16("bars")
		division, _ := cmd.Flags().GetInt16("division")

		renderer, err := echo.NewRenderer(inputFile, outputFile, int(tempo), int(bars), int(division))
		if err != nil {
			cmd.PrintErrln("Error creating renderer:", err)
			return
		}
		renderer.Render()
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().Int16P("tempo", "t", 120, "Base tempo of the delay effect")
	renderCmd.Flags().Int16P("bars", "b", 8, "Duration of the rendered output in bars")
	renderCmd.Flags().Int16P("division", "d", 1, "Delay time as a division of a bar")
}
