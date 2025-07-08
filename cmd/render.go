package cmd

import (
	"github.com/grough/echoprint-go/echo"

	"github.com/spf13/cobra"
)

var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Render delay repeats",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		tempo, _ := cmd.Flags().GetInt16("tempo")
		bars, _ := cmd.Flags().GetInt16("bars")
		division, _ := cmd.Flags().GetInt16("division")

		renderer := echo.NewRenderer(file, tempo, bars, division)
		renderer.Render()
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringP("file", "f", "", "Input audio file path")
	renderCmd.Flags().Int16P("tempo", "t", 120, "Base tempo of the delay effect")
	renderCmd.Flags().Int16P("bars", "b", 8, "Duration of the rendered output in bars")
	renderCmd.Flags().Int16P("division", "d", 1, "Delay time as a division of a bar")
}
