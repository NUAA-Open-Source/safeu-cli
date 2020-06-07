package cmd

import (
	"fmt"

	"github.com/NUAA-Open-Source/safeu-cli/util"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "safeu",
		Short: "SafeU CLI is the command line for SafeU",
		Long: `SafeU CLI is the command line tool for SafeU.

You can access SafeU by via website: https://safeu.a2os.club/
Any question please open issue on https://github.com/NUAA-Open-Source/safeu-cli/issues/new`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetHelpCommand(helpCmd)

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(getCmd)
}

// 打印版本
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "The version number of SafeU CLI tool",
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
