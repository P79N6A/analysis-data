package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/analysis-data/analysisDnsStatic/db"
)

func getTargetAddress(targetAddress string) (string, error) {
	array := strings.Split(targetAddress, ":")
	if len(array) != 2 {
		return targetAddress, fmt.Errorf("unknow address " + targetAddress)
	}

	switch array[1] {
	case "443":
		return "https://" + array[0], nil
	case "80", "8080":
		return "http://" + array[0], nil
	}
	return targetAddress, fmt.Errorf("unknow port " + targetAddress)
}

func isCurrentWebAddress(targetAddress string) bool {
	if strings.Index(targetAddress, ".") > 0 {
		return true
	}
	return false
}

func main() {
	var startTime string
	if len(os.Args) > 1 {
		startTime = os.Args[1]
	}

	if startTime == "" {
		fmt.Println("no input start time")
		return
	}

	var analizyDatas []db.Record

	// datas, err := db.QueryData(startTime)
	// if err != nil {
	// 	fmt.Println("query db fail: ", err.Error())
	// 	return
	// }

	datas, err := db.ReadRecods("./url" + startTime + ".json")
	if err != nil {
		fmt.Println("query db fail: ", err.Error())
		return
	}

	for _, record := range datas.Records {
		if !isCurrentWebAddress(record.TargetAddress) {
			continue
		}
		target, err := getTargetAddress(record.TargetAddress)
		if err != nil {
			continue
		}
		record.TargetAddress = target

		analizyDatas = append(analizyDatas, record)
	}

	WriteDataToFile("analizyData"+startTime+".xls", analizyDatas)
}

func WriteDataToFile(pathFile string, data []db.Record) error {
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
	w.Write([]string{"目标地址", "国家", "包名", "版本号", "用户数"})
	for _, value := range data {
		record := []string{value.TargetAddress, value.ContryZh, value.PkgName, value.AppVersion, value.Count}
		w.Write(record)
	}
	w.Flush()
	return nil
}
