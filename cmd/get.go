package cmd

import (
	"github.com/arcosx/Nuwa/get"
	"github.com/spf13/cobra"
)

var userRecode string
var userPassword string
var targetDir string
var isPrint bool

func init() {
	getCmd.Flags().StringVarP(&userPassword, "password", "p", "", "specific password")
	getCmd.Flags().StringVarP(&targetDir, "dir", "d", "", "download to specific directory")
	getCmd.Flags().BoolVarP(&isPrint, "print", "", false, "print the file URL directly, then you can download the file by other download tools (e.g. wget, aria2).")
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Download file(s) from SafeU",
	Long: `Download file(s) by this command.
SafeU is responsible for ensuring download speed and file safety :)
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userRecode = args[0]
		get.Start(userRecode, userPassword, targetDir, isPrint)
	},
}
