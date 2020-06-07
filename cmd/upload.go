package cmd

import (
	"fmt"

	"github.com/NUAA-Open-Source/safeu-cli/upload"
	"github.com/NUAA-Open-Source/safeu-cli/util"

	"github.com/spf13/cobra"
)

var UserRecode string
var UserPassword string
var UserDownCount int
var UserExpireTime int

func init() {
	uploadCmd.Flags().StringVarP(&UserRecode, "recode", "r", "", "specific recode")
	uploadCmd.Flags().StringVarP(&UserPassword, "password", "p", "", "specific password")
	uploadCmd.Flags().IntVarP(&UserDownCount, "downcount", "d", 0, "specific down count")
	uploadCmd.Flags().IntVarP(&UserExpireTime, "expiretime", "e", 0, "specific expire time")
}

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
		printUserModifyInfo()
		upload.Start(uploadFiles, UserRecode, UserPassword, UserDownCount, UserExpireTime)
	},
}

func printUserModifyInfo() {
	var showNotifyRecode string
	var showNotifyPassword string
	var showNotifyDownCount int
	var showNotifyExpireTime int
	if UserRecode == "" {
		showNotifyRecode = "RANDOM"
	} else {
		showNotifyRecode = UserRecode
	}
	if UserPassword == "" {
		showNotifyPassword = "EMPTY"
	} else {
		showNotifyPassword = UserPassword
	}
	if UserDownCount == 0 {
		showNotifyDownCount = util.DEFAULT_DOWN_COUNT
	} else {
		showNotifyDownCount = UserDownCount
	}
	if UserExpireTime == 0 {
		showNotifyExpireTime = util.DEFAULT_EXPIRE_TIME
	} else {
		showNotifyExpireTime = UserExpireTime
	}
	fmt.Println(fmt.Sprintf("user setting recode:%s , password: %s , downcount: %d , expiretime: %d", showNotifyRecode, showNotifyPassword, showNotifyDownCount, showNotifyExpireTime))
}
