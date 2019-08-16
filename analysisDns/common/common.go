package common

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type AnalizyData struct {
	TargetAddress string
	Count         int
	UserDevices   []string
}

func ExchangeStringToTime(timeString string) (time.Time, error) {
	timeTime := time.Now()
	fromArray := strings.Split(timeString, " ")
	dateArray := strings.Split(fromArray[0], "-")
	year, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return timeTime, err
	}

	if len(dateArray) < 3 {
		return timeTime, fmt.Errorf("time format error")
	}

	month, err := strconv.Atoi(dateArray[1])
	if err != nil {
		return timeTime, err
	}

	day, err := strconv.Atoi(dateArray[2])
	if err != nil {
		return timeTime, err
	}

	var hour int
	if len(fromArray) > 1 {
		hourArray := strings.Split(fromArray[1], ":")
		hour, _ = strconv.Atoi(hourArray[0])
	}

	timeTime = time.Date(year, time.Month(month), day, hour, 0, 0, 0, timeTime.Location())
	return timeTime, nil
}

func WriteToFile(pathFile string, data interface{}) error {
	_, err := os.Stat(pathFile)
	if err != nil {
		if os.IsExist(err) {
			os.Remove(pathFile)
		}
	}
	f, err := os.Create(pathFile)
	if err != nil {
		return err
	}

	defer f.Close()

	byt, _ := json.Marshal(data)
	f.Write(byt)
	return nil
}

func WriteDataToFile(pathFile string, data []AnalizyData) error {
	_, err := os.Stat(pathFile)
	if err != nil {
		if os.IsExist(err) {
			os.Remove(pathFile)
		}
	}
	f, err := os.Create(pathFile)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(f)
	w.Write([]string{"目标地址", "错误次数", "错误矿机"})
	for _, value := range data {
		record := []string{value.TargetAddress, strconv.Itoa(value.Count), strings.Join(value.UserDevices, ",")}
		w.Write(record)
	}
	w.Flush()
	return nil
}
