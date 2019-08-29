package common

import (
	"encoding/csv"
	// "encoding/json"
	// "io/ioutil"
	"os"
	"strconv"
)

type Record struct {
	UserName  string  `json:"userName"`
	CurFlow   float64 `json:"curFlow"`
	NextLogin int     `json:"nextLogin"`
	NextFlow  float64 `json:"nextFlow"`
}

func WriteDataToFile(pathFile string, data []Record) error {
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
	title := []string{"userName", "curFlow", "nextLogin", "nextFlow"}

	w.Write(title)

	for _, value := range data {
		curFlow := strconv.FormatFloat(value.CurFlow, 'f', -1, 32)
		nextLogin := strconv.Itoa(value.NextLogin)
		nextFlow := strconv.FormatFloat(value.NextFlow, 'f', -1, 32)
		record := []string{value.UserName, curFlow, nextLogin, nextFlow}
		w.Write(record)
	}
	w.Flush()
	return nil
}

// func ReadRecods(filePath string) (*Records, error) {
// 	var records Records
// 	fileBytes, err := ioutil.ReadFile(filePath)
// 	if err != nil {
// 		return &records, err
// 	}

// 	err = json.Unmarshal(fileBytes, &records)
// 	if err != nil {
// 		return &records, err
// 	}

// 	return &records, nil
// }
