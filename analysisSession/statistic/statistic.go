package statistic

import (
	"log"
	"time"

	"github.com/analysis-data/analysisSession/sql"

	"github.com/analysis-data/analysisSession/common"
)

// type StatisticClientSessionErrorState struct {
// 	AllocError    int32
// 	InternalError int32
// 	ConnError     int32
// 	NetWorkError  int32
// 	LoginError    int32
// 	LoginTimeOut  int32
// 	P2PError      int32
// }

func StatisticClientSessionError(countryCode, startTime, endTime string) StatisticClientSessionErrorState {
	var result StatisticClientSessionErrorState
	var clientSessions []common.ClientSession
	var err error

	if startTime == "" {
		var zeroTime time.Time
		startTime = zeroTime.Format("2006-01-02 15:04:05")
	}

	if endTime == "" {
		endTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if countryCode != "" {
		clientSessions, err = sql.DB.QueryClientSessionErrorDetailByCountry(countryCode, startTime, endTime)
		if err != nil {
			log.Println("StatisticClientSessionError query db error: " + err.Error())
			return result
		}
	} else {
		clientSessions, err = sql.DB.QueryClientSessionErrorDetail(startTime, endTime)
		if err != nil {
			log.Println("StatisticClientSessionError query db error: " + err.Error())
			return result
		}
	}

	var (
		preSession  common.ClientSession
		isRegisting bool
	)

	for _, temp := range clientSessions {
		if preSession.ID != "" && temp.ID != preSession.ID {
			if isRegisting {
				result.AllocError++
			}

			if preSession.Errstr == "network-error" {
				result.NetWorkError++
			}
			if preSession.Errstr == "no-conn-info" {
				result.ConnError++
			}
			isRegisting = false
		}

		if temp.State == "sfox-cli-registing" {
			isRegisting = true
		}

		preSession = temp
	}
	return result
}

// type StatisticClientSessionState struct {
// 	Connecting  int32
// 	Disable     int32
// 	Error       StatisticClientSessionErrorState
// 	Registing   int32
// 	Unregisting int32
// }

func StatisticClientSessionConnectResult(countryCode, startTime, endTime string) StatisticClientSessionState {
	var result StatisticClientSessionState
	var data common.ClientSessionConnectionResult
	var err error

	if startTime == "" {
		var zeroTime time.Time
		startTime = zeroTime.Format("2006-01-02 15:04:05")
	}

	if endTime == "" {
		endTime = time.Now().Format("2006-01-02 15:04:05")
	}

	if countryCode != "" {
		data, err = sql.DB.QueryClientSessionConnectionResultByCountry(countryCode, startTime, endTime)
		if err != nil {
			log.Println("StatisticClientSessionConnectResult query db error: " + err.Error())
			return result
		}
	} else {
		data, err = sql.DB.QueryClientSessionConnectionResult(startTime, endTime)
		if err != nil {
			log.Println("StatisticClientSessionConnectResult query db error: " + err.Error())
			return result
		}
	}

	errorResult := StatisticClientSessionError(countryCode, startTime, endTime)
	result = StatisticClientSessionState{
		Success:     data.SuccessCount,
		Connecting:  data.ConnectionCount,
		Disable:     data.DisableCount,
		Error:       errorResult,
		Registing:   data.RegistingCount,
		Unregisting: data.UnRegistingCount,
	}
	return result
}
