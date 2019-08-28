package mail

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

var targetURL = "http://test.shadowfox.cc:34500"

func Upload(filePath string) error {
	bodyBuf := &bytes.Buffer{}

	w := multipart.NewWriter(bodyBuf)
	fw, err := w.CreateFormFile("file", filePath)
	if err != nil {
		return fmt.Errorf("CreateFormFile error: %s", err.Error())
	}

	fh, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("Open fail error: %s", err.Error())
	}

	defer fh.Close()

	_, err = io.Copy(fw, fh)
	if err != nil {
		return fmt.Errorf("Copy fail error: %s", err.Error())
	}

	contentType := w.FormDataContentType()
	w.Close()

	resp, err := http.Post(targetURL+"/upload", contentType, bodyBuf)
	if err != nil {
		return fmt.Errorf("Post fail error: %s", err.Error())
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("POST file result: status: %s, body: %s\n", resp.Status, respBody)
	return nil
}

type MessageAttach struct {
	FileName string `json:"file_name"`
	IsUpload bool   `json:"is_upload"`
}
type MessageData struct {
	Topic      string          `json:"topic"`
	BodyFormat string          `json:"body_format"`
	SenderName string          `json:"sender_name"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	From       string          `json:"from"`
	Attaches   []MessageAttach `json:"attaches"`
}

func SendMessage(filePath string) error {
	attach := MessageAttach{
		FileName: filePath,
		IsUpload: true,
	}

	mail := MessageData{
		Topic:      "daily_user_data",
		SenderName: "每日用戶數據",
		Title:      "每日用户数据",
		Content:    "详见附件",
		From:       "sfoxstatsite@163.com",
		Attaches:   []MessageAttach{attach},
	}

	err := sendMessage(targetURL+"/sendmessage", mail)
	if err != nil {
		return fmt.Errorf("sendMessage fail error: %s", err.Error())
	}
	return nil
}

func sendMessage(url string, messageData interface{}) error {
	data, err := json.Marshal(messageData)
	if err != nil {
		return err
	}

	reqNew := bytes.NewBuffer(data)
	request, err := http.NewRequest("POST", url, reqNew)
	if err != nil {
		return err
	}

	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body[:]))
	}

	return nil
}
