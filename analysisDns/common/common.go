package common

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type AnalizyData struct {
	TargetAddress string
	Count         int
	UserDevices   []string
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
