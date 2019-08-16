package main

import (
	"encoding/json"
	"fmt"
	"github.com/analysis-data/analysisMail/common"
	"github.com/analysis-data/analysisMail/db"
	"github.com/analysis-data/analysisMail/mail"
	"os"
	"time"
)

func main() {
	startFunctionTimer(process)
	var stopChan = make(chan string)
	for {
		select {
		case _ = <-stopChan:
			return
		}
	}
}

func process() {
	now := time.Now()
	runMode := "dev"
	if len(os.Args) > 1 {
		runMode = os.Args[1]
	}

	dbHandle, err := db.RegisterDB(runMode)
	if err != nil {
		fmt.Println("RegisterDB error: ", err)
		return
	}

	var records []common.AnalysisData
	err = dbHandle.Query(now, &records)
	if err != nil {
		fmt.Println("Query db fail error: ", err)
		return
	}

	dbHandle.Close()

	byt, _ := json.Marshal(records)
	fmt.Println(string(byt[:]))

	filePath := "userData" + common.TimeFormatDate(now) + ".xls"
	err = common.WriteDataToFile(filePath, records)
	if err != nil {
		fmt.Println("WriteDataToFile error: ", err)
		return
	}

	err = mail.Upload(filePath)
	if err != nil {
		fmt.Println("Upload file error: ", err.Error())
		return
	}

	err = mail.SendMessage(filePath)
	if err != nil {
		fmt.Println("SendMessage error: ", err.Error())
	}
}

func startFunctionTimer(f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			// 计算下一个24h
			next := now.Add(time.Minute * 60 * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 10, 0, 0, 0, next.Location())
			nextIgnore := next.Sub(now)
			t := time.NewTimer(nextIgnore)
			<-t.C
		}
	}()
}
