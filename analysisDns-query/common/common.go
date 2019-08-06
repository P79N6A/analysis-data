package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
