package main

import (
	"fmt"
	"github.com/analysis-data/analysisRetain/common"
	"github.com/analysis-data/analysisRetain/db"
)

func main() {
	curTime := "2019-06-25 08:00:00"
	startTime := common.StringToTime(curTime)

	gorDB, err := db.RegisterDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gorDB.BuildSQL(startTime))

	// var recods []common.Record
	// gorDB.Query(startTime, &recods)

	// if err = common.WriteDataToFile("./analysis.xls", recods); err != nil {
	// 	fmt.Println(err)
	// }
}
