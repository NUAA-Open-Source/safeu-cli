package cmd

import (
	"fmt"
	"github.com/arcosx/Nuwa/upload"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file to SafeU",
	Long: `Send and Share file by this command.
SafeU is responsible for ensuring upload speed and file safety
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var uploadFiles []string

		for i := 0; i < len(args); i++ {
			fileName := args[i]
			fmt.Println("ready upload file: " + fileName)
			uploadFiles = append(uploadFiles, fileName)
		}
		upload.Start(uploadFiles)
	},
}
