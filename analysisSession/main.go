package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/analysis-data/analysisSession/statistic"
)

func main() {
	result := statistic.StatisticClientSessionConnectResult("Iran", "2018-12-26 00:00:00", "2018-12-27 00:00:00")
	code, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error)
	}
	fmt.Println(string(code[:]))
}
