package main

import (
	"fmt"
	"github.com/analysis-data/analysisMail/mail"
	"os"
)

func main() {
	date := os.Args[1]

	filePath := "userData" + date + ".xls"
	err := mail.Upload(filePath)
	if err != nil {
		fmt.Println("Upload file error: ", err.Error())
		return
	}

	err = mail.SendMessage(filePath)
	if err != nil {
		fmt.Println("SendMessage error: ", err.Error())
	}

	fmt.Println("task send mail complete success")
}
