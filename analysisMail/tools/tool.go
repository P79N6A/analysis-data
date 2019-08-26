package main

import (
	"fmt"
	"github.com/analysis-data/analysisMail/common"
	"github.com/analysis-data/analysisMail/db"
	"time"
)

func main() {
	dbHandle, err := db.RegisterDB("prod")
	if err != nil {
		fmt.Println("RegisterDB fail: error: ", err)
		return
	}

	nowTime := time.Now()
	var sessionData []common.ClientSessionData
	err = dbHandle.QueryClientSessionData(nowTime.Unix()-60*60*24, nowTime.Unix(), &sessionData)
	if err != nil {
		fmt.Println("QueryClientSessionData fail: error: ", err)
		return
	}
	fmt.Println("query time: ", time.Now().UnixNano()-nowTime.UnixNano())
}
