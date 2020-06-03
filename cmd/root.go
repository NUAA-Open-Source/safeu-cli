package cmd

import (
	"fmt"
	"github.com/arcosx/Nuwa/util"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "safeu",
		Short: "Nuwa is Command Line For SafeU",
		Long: `Nuwa is Command Line Tool for SafeU.
You can access SafeU by via website: https://safeu.a2os.club/
Any question please open issue on https://github.com/arcosx/Nuwa/issues/new`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetHelpCommand(helpCmd)

	rootCmd.AddCommand(uploadCmd)
}

// 打印版本
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version number of Nuwa",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(util.VERSION)
	},
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "help",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`upload one file : safeu upload filename
upload more file : safeu upload filename1 filename2 filename3 ....
	`)
	},
}
