package db_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/analysis-data/analysisDns/db"
)

func TestQueryDB(t *testing.T) {
	nowTime := time.Now()
	fromTime := time.Date(nowTime.Year(), nowTime.Month(), 1, 0, 0, 0, 0, nowTime.Location())
	records, err := db.QueryData(fromTime, nowTime)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(records)
}
