package cmd

import (
	"github.com/arcosx/Nuwa/get"
	"github.com/spf13/cobra"
)

var userRecode string
var userPassword string

func init() {
	getCmd.Flags().StringVarP(&userPassword, "password", "p", "", "specific password")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download file(s) from SafeU",
	Long: `Download file(s) by this command.
SafeU is responsible for ensuring download speed and file safety :)
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		get.Start(userRecode, userPassword)
	},
}
