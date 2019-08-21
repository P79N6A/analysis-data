package main

import (
	// "encoding/json"
	"fmt"
	"github.com/analysis-data/analysisMail/common"
	"github.com/analysis-data/analysisMail/db"
	"github.com/analysis-data/analysisMail/mail"
	"os"
	"time"
)

var userSessionData map[string][]common.ClientSessionAnalysisData

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

	queryTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	preTime := queryTime.Unix() - 60*24*60
	fmt.Println("query time: ", common.TimeFormat(queryTime))

	var analysisData []common.AnalysisData
	err = dbHandle.Query(preTime, queryTime.Unix(), &analysisData)
	if err != nil {
		fmt.Println("Query db fail error: ", err)
		return
	}

	var sessionData []common.ClientSessionData
	err = dbHandle.QueryClientSessionData(preTime, queryTime.Unix(), &sessionData)
	if err != nil {
		fmt.Println("QueryClientSessionData db fail error: ", err)
		return
	}

	analysisSeesionData(common.TimeFormatDateFromUnix(preTime), sessionData, dbHandle)

	dbHandle.Close()

	for idx, userData := range analysisData {
		mapUserData, ok := userSessionData[userData.UserName]
		if !ok {
			continue
		}

		analysisData[idx].ConnectCount = len(mapUserData)

		for _, tempData := range mapUserData {
			analysisData[idx].ConnectTime += (tempData.EndTime - tempData.StartTime)

			if analysisData[idx].LastConnectTime < tempData.StartTime {
				analysisData[idx].Location = tempData.SelectRouter
				analysisData[idx].LastConnectTime = tempData.StartTime
			}

			startTime := tempData.Conns[0].ConnCreateTime
			endTime := tempData.Conns[0].ConnCloseTime

			for _, tempConnData := range tempData.Conns {
				if startTime > tempConnData.ConnCreateTime {
					startTime = tempConnData.ConnCreateTime
				}

				if endTime < tempConnData.ConnCloseTime {
					endTime = tempConnData.ConnCloseTime
				}
			}

			analysisData[idx].UseTime = (endTime - startTime) / 1000
		}

	}

	// byt, _ := json.Marshal(analysisData)
	// fmt.Println(string(byt[:]))

	filePath := "userData" + common.TimeFormatDate(now) + ".xls"
	err = common.WriteDataToFile(filePath, analysisData)
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

func analysisSeesionData(date string, datas []common.ClientSessionData, dbHandle *db.GormInterface) {
	sessionAnalysis := map[string]*common.ClientSessionAnalysisData{}

	for _, tempData := range datas {
		storeData, ok := sessionAnalysis[tempData.ID]
		if !ok {
			var connAnalysisData []common.ClientConnAnalysisData
			var connData []common.ClientConnData
			fmt.Println("start query ", "tp_client_conn"+date, "session: ", tempData.ID)
			err := dbHandle.QueryClientConnDataByID("tp_client_conn"+date, tempData.ID, &connData)
			if err == nil {
				connAnalysisData = analysisConnData(connData)
			}

			storeData = &common.ClientSessionAnalysisData{
				UserName:     tempData.UserName,
				UserDevice:   tempData.UserDevice,
				ID:           tempData.ID,
				PkgName:      tempData.PkgName,
				AppVersion:   tempData.AppVersion,
				SelectRouter: tempData.SelectRouter,
				RemoteInput:  tempData.RemoteInput,
				RemoteOutput: tempData.RemoteOutput,
				StartTime:    tempData.CreateTimestamp,
				EndTime:      tempData.CreateTimestamp,
				Conns:        connAnalysisData,
			}
			sessionAnalysis[tempData.ID] = storeData
		} else {
			if storeData.SelectRouter != tempData.SelectRouter {
				storeData.SelectRouter = tempData.SelectRouter
			}

			if storeData.RemoteInput < tempData.RemoteInput {
				storeData.RemoteInput = tempData.RemoteInput
			}

			if storeData.RemoteOutput < tempData.RemoteOutput {
				storeData.RemoteOutput = tempData.RemoteOutput
			}

			if storeData.StartTime > tempData.CreateTimestamp {
				storeData.StartTime = tempData.CreateTimestamp
			}

			if storeData.EndTime < tempData.CreateTimestamp {
				storeData.EndTime = tempData.CreateTimestamp
			}
		}
	}

	userSessionData = map[string][]common.ClientSessionAnalysisData{}
	for _, tempData := range sessionAnalysis {
		_, ok := userSessionData[tempData.UserName]
		if !ok {
			var sessionData []common.ClientSessionAnalysisData
			sessionData = append(sessionData, *tempData)
			userSessionData[tempData.UserName] = sessionData
		} else {
			userSessionData[tempData.UserName] = append(userSessionData[tempData.UserName], *tempData)
		}
	}
}

func analysisConnData(datas []common.ClientConnData) []common.ClientConnAnalysisData {
	var connData []common.ClientConnAnalysisData

	connAnalysis := map[string]*common.ClientConnAnalysisData{}
	for _, tempData := range datas {
		storeData, ok := connAnalysis[tempData.ID]
		if !ok {
			storeData = &common.ClientConnAnalysisData{
				ID:             tempData.ID,
				RemoteInput:    tempData.RemoteInput,
				RemoteOutput:   tempData.RemoteOutput,
				ViaUserName:    tempData.ViaUserName,
				ViaUserDevice:  tempData.ViaUserDevice,
				TargetAddress:  tempData.TargetAddress,
				ConnCreateTime: tempData.ConnCreateTime,
				ConnCloseTime:  tempData.ConnCloseTime,
			}

			connAnalysis[tempData.ID] = storeData
		} else {
			if storeData.RemoteInput < tempData.RemoteInput {
				storeData.RemoteInput = tempData.RemoteInput
			}

			if storeData.RemoteOutput < tempData.RemoteOutput {
				storeData.RemoteOutput = tempData.RemoteOutput
			}

			if storeData.ConnCloseTime < tempData.ConnCloseTime {
				storeData.ConnCloseTime = tempData.ConnCloseTime
			}
		}
	}

	for _, tempData := range connAnalysis {
		connData = append(connData, *tempData)
	}

	return connData
}

func startFunctionTimer(f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			// 计算下一个24h
			next := now.Add(time.Minute * 60 * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 8, 0, 0, 0, next.Location())
			nextIgnore := next.Sub(now)
			t := time.NewTimer(nextIgnore)
			<-t.C
		}
	}()
}
