package main

import (
	// "encoding/json"
	"fmt"
	"github.com/analysis-data/analysisMail/aliyun"
	"github.com/analysis-data/analysisMail/common"
	"github.com/analysis-data/analysisMail/db"
	"github.com/analysis-data/analysisMail/mail"
	"os"
	"time"
)

var userSessionData map[string][]common.ClientSessionAnalysisData
var sessionAnalysis map[string]*common.ClientSessionAnalysisData

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
	fmt.Println("query time: ", common.TimeFormat(queryTime), " from: ", preTime, " to: ", queryTime.Unix())

	var analysisData []common.AnalysisData
	err = dbHandle.Query(preTime, queryTime.Unix(), &analysisData)
	if err != nil {
		fmt.Println("Query db fail error: ", err)
		return
	}

	fmt.Println("Query login info cost time ms: ", (time.Now().UnixNano()-now.UnixNano())/1000/1000)
	now = time.Now()

	var sessionData []common.ClientSessionData
	err = dbHandle.QueryClientSessionData(preTime, queryTime.Unix(), &sessionData)
	if err != nil {
		fmt.Println("QueryClientSessionData db fail error: ", err)
		return
	}

	fmt.Println("Query session info cost time ms: ", (time.Now().UnixNano()-now.UnixNano())/1000/1000)

	dbHandle.Close()

	// 处理session数据
	sessionAnalysis = map[string]*common.ClientSessionAnalysisData{}
	analysisSeesionData(sessionData)
	/*aliyun chuli */
	/* 	client := aliyun.NewClient(runMode)

	   	offset := 0

	   	// 处理session数据
	   	sessionAnalysis = map[string]*common.ClientSessionAnalysisData{}
	   	for {
	   		sessionData, err := client.GetSessionData(preTime, queryTime.Unix(), int64(100), int64(offset))
	   		if err != nil {
	   			fmt.Println("GetSessionData error: ", err)
	   			return
	   		}

	   		analysisSeesionData(sessionData)
	   		offset += len(sessionData)
	   		if len(sessionData) < 100 {
	   			break
	   		}
	   	}

	   	// 处理conn数据
	   	queryAndAnalysisConnData(client, preTime, queryTime.Unix())

	*/

	analysisUserData()

	for idx, userData := range analysisData {
		mapUserData, ok := userSessionData[userData.UserName]
		if !ok {
			continue
		}

		analysisData[idx].ConnectCount = len(mapUserData)

		for _, tempData := range mapUserData {
			analysisData[idx].ConnectTime += (tempData.EndTime - tempData.StartTime/1000)

			if analysisData[idx].LastConnectTime < tempData.StartTime {
				analysisData[idx].Location = tempData.SelectRouter
				analysisData[idx].LastConnectTime = tempData.StartTime
			}

			// startTime := tempData.EndTime * 1000
			// endTime := tempData.StartTime

			// if len(tempData.Conns) == 0 {
			// 	startTime = 0
			// 	endTime = 0
			// }

			// for _, tempConnData := range tempData.Conns {
			// 	if startTime > tempConnData.ConnCreateTime {
			// 		startTime = tempConnData.ConnCreateTime
			// 	}

			// 	if endTime < tempConnData.ConnCloseTime {
			// 		endTime = tempConnData.ConnCloseTime
			// 	}
			// }

			// analysisData[idx].UseTime += (endTime - startTime) / 1000
		}

	}

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

	sessionAnalysis = map[string]*common.ClientSessionAnalysisData{}
	userSessionData = map[string][]common.ClientSessionAnalysisData{}

	fmt.Println("task complete success")
}

func analysisSeesionData(datas []common.ClientSessionData) {
	for _, tempData := range datas {
		storeData, ok := sessionAnalysis[tempData.ID]
		if !ok {
			storeData = &common.ClientSessionAnalysisData{
				UserName:     tempData.UserName,
				UserDevice:   tempData.UserDevice,
				ID:           tempData.ID,
				PkgName:      tempData.PkgName,
				AppVersion:   tempData.AppVersion,
				SelectRouter: tempData.SelectRouter,
				RemoteInput:  tempData.RemoteInput,
				RemoteOutput: tempData.RemoteOutput,
				StartTime:    tempData.EnableTime,
				EndTime:      tempData.CreateTimestamp,
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
				storeData.StartTime = tempData.EnableTime
			}

			if storeData.EndTime < tempData.CreateTimestamp {
				storeData.EndTime = tempData.CreateTimestamp
			}
		}
	}
}

func analysisUserData() {
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

func queryAndAnalysisConnData(client *aliyun.AliyunClient, from, to int64) {
	for key, tempSession := range sessionAnalysis {
		fmt.Println("start query ", " tp_client_conn", " user: ", tempSession.UserName, " session: ", key)

		// err := dbHandle.QueryClientConnDataByID("tp_client_conn"+date, tempData.ID, &connData)
		// if err == nil {
		// 	connAnalysisData = analysisConnData(connData)
		// }
		offset := 0
		var connAnalysisData []common.ClientConnAnalysisData
		for {
			connData, err := client.GetConnData(from, to, key, int64(100), int64(offset))
			if err != nil {
				fmt.Println("GetConnData: ", key, " error: ", err)
				return
			}
			tempData := connAnalysisData

			count := len(connData)
			offset += count

			if count > 0 {
				connAnalysisData = analysisConnData(connData, tempData)
			}
			if count < 100 {
				break
			}
		}

		tempSession.Conns = connAnalysisData
	}
}

func analysisConnData(datas []common.ClientConnData, connsData []common.ClientConnAnalysisData) []common.ClientConnAnalysisData {
	var analysisData []common.ClientConnAnalysisData
	for _, tempData := range datas {
		isExist := false
		for _, storeData := range connsData {
			if storeData.ID == tempData.ID {
				if storeData.RemoteInput < tempData.RemoteInput {
					storeData.RemoteInput = tempData.RemoteInput
				}

				if storeData.RemoteOutput < tempData.RemoteOutput {
					storeData.RemoteOutput = tempData.RemoteOutput
				}

				if storeData.ConnCloseTime < tempData.ConnCloseTime {
					storeData.ConnCloseTime = tempData.ConnCloseTime
				}

				analysisData = append(analysisData, storeData)
				isExist = true
				break
			}
		}

		if !isExist {
			storeData := common.ClientConnAnalysisData{
				ID:             tempData.ID,
				RemoteInput:    tempData.RemoteInput,
				RemoteOutput:   tempData.RemoteOutput,
				ViaUserName:    tempData.ViaUserName,
				ViaUserDevice:  tempData.ViaUserDevice,
				TargetAddress:  tempData.TargetAddress,
				ConnCreateTime: tempData.ConnCreateTime,
				ConnCloseTime:  tempData.ConnCloseTime,
			}
			analysisData = append(analysisData, storeData)
		}

	}

	return analysisData
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
