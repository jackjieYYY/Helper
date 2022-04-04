package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var ConnectStr = "mongodb+srv://closure:k83nF3ea9SWyQgs5@cluster0.ujkjx.mongodb.net/test?authSource=admin&replicaSet=atlas-gnh73i-shard-0&readPreference=primary&appname=MongoDB%20Compass&ssl=true"

func PostToMe(message string) {
	url := "http://mc.mesord.com:8098/api/send_msg_auto"
	method := "POST"

	payload := strings.NewReader(`{"token": "ADMIN_TOKEN@FEXLI_2022" , "msg": "` + message + `", "uid": 563255057, "toImg": false}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func Post(message string) {
	url := "http://mc.mesord.com:8098/api/send_msg_auto"
	method := "POST"

	payload := strings.NewReader(`{"token": "ADMIN_TOKEN@FEXLI_2022" , "msg": "` + message + `", "uid": 913468406, "toImg": false}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func WriteResVersionFile(res *ResVersionInfo) {
	b, err := json.Marshal(res)
	if err != nil {
		Post("An error occurred while creating the file In Redroid.MongoDB sync Func")
		fmt.Println("error:", err)
	}
	_ = ioutil.WriteFile("version.info", b, 0644)
}

func SurveillanceIsLock() bool {
	return isFileExist("Lock")
}

func PingIsLock() bool {
	return isFileExist("PingLock")
}
func PingLock() {
	fmt.Println("start to lock")
	_ = ioutil.WriteFile("PingLock", []byte{}, 0644)
}
func PingUnLock() {
	_ = os.Remove("PingLock")
}
func Lock() {
	fmt.Println("start to lock")
	_ = ioutil.WriteFile("Lock", []byte{}, 0644)
}

func UnLock() {
	_ = os.Remove("Lock")
}

func isFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func Cmd(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		Post(err.Error())
		return "", err
	}
	return string(out), nil
}
