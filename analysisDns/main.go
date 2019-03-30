package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/analysis-data/analysisDns/common"
	"github.com/analysis-data/analysisDns/db"
)

func getTargetAddressWithOutPort(targetAddress string) (string, error) {
	array := strings.Split(targetAddress, ":")
	return array[0], nil
}

var analizyDatas = map[string]*common.AnalizyData{}

func addAnalizyData(targetAddress string, userdevice string) error {
	if data, ok := analizyDatas[targetAddress]; ok {
		data.Count++
		for _, device := range data.UserDevices {
			if device == userdevice {
				return nil
			}
		}
		data.UserDevices = append(data.UserDevices, userdevice)
	} else {
		data := common.AnalizyData{
			TargetAddress: targetAddress,
			Count:         1,
			UserDevices:   []string{userdevice},
		}
		analizyDatas[targetAddress] = &data
	}
	return nil
}

func isCurrentWebAddress(targetAddress string) bool {
	if strings.Index(targetAddress, ".") > 0 {
		return true
	}
	return false
}

func findData(from, to time.Time) error {
	return nil
}

func main() {
	analizyDatas = map[string]*common.AnalizyData{}

	cmd := "read"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	var records *db.Records
	var err error
	switch cmd {
	case "read":
		records, err = db.ReadRecods("./records.json")
		if err != nil {
			fmt.Println("readRecords error: " + err.Error())
			return
		}
	case "query":
		if len(os.Args) <= 2 {
			fmt.Println("please input query time, for example: 2019-03-25 2019-03-29 or from time 2019-03-25")
			return
		}

		fromTime := time.Now()
		endTime := time.Now()
		fromTime, _ = common.ExchangeStringToTime(os.Args[2])
		if len(os.Args) > 3 {
			endTime, _ = common.ExchangeStringToTime(os.Args[3])
		}
		records, err = db.QueryData(fromTime, endTime)
		if err != nil {
			fmt.Println("QueryData error: " + err.Error())
			return
		}
	default:
		fmt.Println("please input date source: read or query")
		return
	}

	for _, record := range records.Records {
		target, _ := getTargetAddressWithOutPort(record.TargetAddress)
		addAnalizyData(target, record.UserDevice)
	}

	var data []common.AnalizyData

	total := 0
	errWebAddressCount := 0
	for key, value := range analizyDatas {
		total += value.Count
		if !isCurrentWebAddress(key) {
			errWebAddressCount += value.Count
		}
		data = append(data, *value)
		// fmt.Printf("%-50s %d\t ", key, value.count)
		// for index, userDevice := range value.userDevices {
		// 	fmt.Printf("%s", userDevice)
		// 	if index < len(value.userDevices)-1 {
		// 		fmt.Printf(",")
		// 	}
		// }
		// fmt.Printf("\n")
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Count > data[j].Count
	})

	common.WriteDataToFile("analizyData.xls", data)
	fmt.Printf("\ntotal: %d, errorWebAddressCount: %d, errorRate: %f", total, errWebAddressCount, float64(errWebAddressCount)/float64(total))
}
