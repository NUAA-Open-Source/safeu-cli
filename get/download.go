package get

import (
	"fmt"
	"os"
)

func (item *ItemDownload) download() error {
	// download logics
	return nil
}

func (item *ItemDownload) minusDownCount() error {
	// minus download counter
	return nil
}

func Start(recode string, password string) {
	var item ItemDownload

	item.userRecode = recode
	item.userPassword = password

	err := item.download()
	if err != nil {
		fmt.Println("download error")
		os.Exit(0)
	}

	err = item.minusDownCount()
	if err != nil {
		fmt.Println("download minus counter error")
		os.Exit(0)
	}

	fmt.Println("download finished")

}
