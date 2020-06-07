package get

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/NUAA-Open-Source/safeu-cli/util"
)

func (dm *DownloadModel) getCsrf() error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", util.SAFEU_BASE_URL+"/csrf", nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	dm.Csrf = resp.Header.Get("X-Csrf-Token")
	dm.Cookie = resp.Header.Get("Set-Cookie")
	return nil
}

func (dm *DownloadModel) validation() (err error) {
	// validate recode & password
	var valiReq ValidationRequest

	hasher := sha256.New()
	hasher.Write([]byte(dm.UserPassword))
	sha256Pass := hex.EncodeToString(hasher.Sum(nil))
	valiReq.Password = sha256Pass

	jsonStr, err := json.Marshal(valiReq)
	if err != nil {
		fmt.Println("validation json marshal error", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", util.SAFEU_BASE_URL+"/v1/validation/"+dm.UserRecode, strings.NewReader(string(jsonStr)))

	if err != nil {
		fmt.Println("validation request create NewRequest failed", err)
		return err
	}
	req.Header.Set("x-csrf-token", dm.Csrf)
	req.Header.Set("cookie", dm.Cookie)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("validation request call failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			fmt.Println("password needed / invalid password for this recode, validation failed")
			return fmt.Errorf("invalid password")
		}
		fmt.Println("validation requset response is not 200")
		return fmt.Errorf("validation request response code: %d, content: %s", resp.StatusCode, respBody)
	}

	var valiRes ValidationResponse
	err = json.Unmarshal(respBody, &valiRes)
	if err != nil {
		fmt.Println("validation json unmarshal failed", err)
		return err
	}

	dm.Token = valiRes.Token
	dm.Items = valiRes.Items
	return nil
}

func (dm *DownloadModel) getDownloadURL(isPrint bool) error {
	// download logics
	var downReq DownloadRequest
	downReq.Full = true
	downReq.Items = dm.Items

	jsonReq, err := json.Marshal(&downReq)
	if err != nil {
		fmt.Println("download json marshal error", err)
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", util.SAFEU_BASE_URL+"/v1/item/"+dm.UserRecode, strings.NewReader(string(jsonReq)))
	if err != nil {
		fmt.Println("download create NewRequest failed", err)
		return err
	}
	req.Header.Set("x-csrf-token", dm.Csrf)
	req.Header.Set("cookie", dm.Cookie)
	req.Header.Set("token", dm.Token) // auth token for this recode

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("download request call failed", err)
		return err
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Println("download request response status is not 200")
		return fmt.Errorf("download request response return code: %d, content: %s", resp.StatusCode, respBody)
	}

	var downRes DownloadResponse
	err = json.Unmarshal(respBody, &downRes)
	if err != nil {
		fmt.Println("download json unmarshal failed")
		return err
	}

	if isPrint {
		fmt.Println(downRes.URL)
		os.Exit(0)
	}

	dm.URL = downRes.URL
	return nil
}

func (dm *DownloadModel) downloadFile() error {
	// actual download file
	// and rename the file
	if dm.Dir == "" {
		// default: get the current working directory
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("cannot get current directory, please try `--dir` option instead.")
			return fmt.Errorf("cannot get current directory")
		}
		dm.Dir = dir
	}

	tmpFilename := dm.Dir + string(os.PathSeparator) + string(dm.UserRecode) + ".tmp"
	filename := dm.Dir + string(os.PathSeparator) + string(dm.UserRecode) + ".zip"
	if len(dm.Items) == 1 {
		tmpFilename = dm.Dir + string(os.PathSeparator) + dm.Items[0].OriginalName + ".tmp"
		filename = dm.Dir + string(os.PathSeparator) + dm.Items[0].OriginalName
	}

	out, err := os.Create(tmpFilename)
	if err != nil {
		fmt.Println("cannot create file", tmpFilename, ", error: ", err)
		return fmt.Errorf("create file error: %s", err)
	}

	resp, err := http.Get(dm.URL)
	if err != nil {
		out.Close()
		return fmt.Errorf("cannot get file %s, error: %s", dm.URL, err)
	}
	defer resp.Body.Close()

	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Print("\n")

	// Close the file without defer so it can happen before Rename()
	out.Close()

	if err = os.Rename(tmpFilename, filename); err != nil {
		return fmt.Errorf("cannot rename file %s to %s", tmpFilename, filename)
	}

	// print for debug
	// fmt.Printf("tmp file: %s, file: %s\n", tmpFilename, filename)
	dm.Filepath = filename
	return nil
}

func (dm *DownloadModel) minusDownCount() error {
	// minus download counter
	// download logics
	var minusReq MinusDownCountRequest
	client := &http.Client{}

	for _, item := range dm.Items {
		minusReq.Bucket = item.Bucket
		minusReq.Path = item.Path

		jsonReq, err := json.Marshal(&minusReq)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", util.SAFEU_BASE_URL+"/v1/minusDownCount/"+dm.UserRecode, strings.NewReader(string(jsonReq)))
		if err != nil {
			return err
		}
		req.Header.Set("x-csrf-token", dm.Csrf)
		req.Header.Set("cookie", dm.Cookie)

		// minus download counter, ignore any return content
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func Start(recode string, password string, targetDir string, isPrint bool) {
	var dm DownloadModel

	dm.UserRecode = recode
	dm.UserPassword = password
	dm.Dir = targetDir

	// start := time.Now()
	err := dm.getCsrf()
	if err != nil {
		fmt.Println("get CSRF token error", err)
		os.Exit(1)
	}
	// elapsed := time.Since(start)
	// fmt.Println("get csrf cost:", elapsed)

	// start = time.Now()
	err = dm.validation()
	if err != nil {
		fmt.Println("validation error", err)
		os.Exit(1)
	}
	// elapsed = time.Since(start)
	// fmt.Println("validation cost:", elapsed)

	// start = time.Now()
	err = dm.getDownloadURL(isPrint)
	if err != nil {
		fmt.Println("download error", err)
		os.Exit(1)
	}
	// elapsed = time.Since(start)
	// fmt.Println("get download url cost:", elapsed)

	// start = time.Now()
	err = dm.downloadFile()
	if err != nil {
		fmt.Println("download file error", err)
		os.Exit(1)
	}
	// elapsed = time.Since(start)
	// fmt.Println("download file cost:", elapsed)

	// TODO: can be a goroutine
	// start = time.Now()
	err = dm.minusDownCount()
	if err != nil {
		// do nothing...
		fmt.Println("download minus counter error", err)
		os.Exit(1)
	}
	// elapsed = time.Since(start)
	// fmt.Println("minus downcount cost:", elapsed)

	// print for test
	// fmt.Printf("Recode: %s, Token: %s, URL: %s\n", dm.UserRecode, dm.Token, dm.URL)

	fmt.Println("Download complete! You can get your file at: " + dm.Filepath)

}
