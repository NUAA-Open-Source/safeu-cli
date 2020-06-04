package get

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arcosx/Nuwa/util"
)

// https://app.swaggerhub.com/apis-docs/a2os/safeu/1.0.0-rc#/miscellaneous/getCSRFToken
func (item *ItemDownload) getCsrf() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", util.SAFEU_BASE_URL+"/csrf", nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	item.csrf = resp.Header.Get("X-Csrf-Token")
	item.cookie = resp.Header.Get("Set-Cookie")
	return nil
}

// https://app.swaggerhub.com/apis-docs/a2os/safeu/1.0.0-rc#/validation/validateRC
func (item *ItemDownload) validation() error {
	// validate recode & password
	return nil
}

// https://app.swaggerhub.com/apis-docs/a2os/safeu/1.0.0-rc#/download/downloadItem
func (item *ItemDownload) download() error {
	// download logics
	return nil
}

// https://app.swaggerhub.com/apis-docs/a2os/safeu/1.0.0-rc#/download/downloadAwareness
func (item *ItemDownload) minusDownCount() error {
	// minus download counter
	return nil
}

func Start(recode string, password string) {
	var item ItemDownload

	item.userRecode = recode
	item.userPassword = password

	err := item.getCsrf()
	if err != nil {
		fmt.Println("get CSRF token error")
		os.Exit(0)
	}

	err = item.validation()
	if err != nil {
		fmt.Println("validation error")
		os.Exit(0)
	}

	err = item.download()
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
