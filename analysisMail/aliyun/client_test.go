package aliyun_test

import (
	"fmt"
	"github.com/analysis-data/analysisMail/aliyun"
	"github.com/analysis-data/analysisMail/common"
	"reflect"
	"strconv"
	"testing"
)

func TestAliyunClient_GetLogs(t *testing.T) {
	c := aliyun.NewClient("dev")

	if _, err := c.GetLogs("sfox-dev", "tp_client_session", 1566230400, 1566316800, "", 100, 0, false); err != nil {
		fmt.Println("AliyunClient.GetLogs() error :", err)
	}
}

func TestAliyunClient_GetConnData(t *testing.T) {
	c := aliyun.NewClient("dev")
	conns, err := c.GetConnData(1566403200, 1566489600, "88fedf9bfd657b2d9a5a3416d771e580", 100, 0)
	if err != nil {
		fmt.Println(err)
	}
	for _, tempData := range conns {
		fmt.Println(tempData)
	}
}

func TestAliyunClient_Example(t *testing.T) {
	c := aliyun.NewClient("dev")

	res, err := c.GetLogs("sfox-"+c.Runmode, "tp_client_session", 1566439484, 1566525884, "cc.coolline.client", 1, 0, false)
	if err != nil {
		fmt.Println("getlogs err: ", err)
		return
	}

	if !res.IsComplete() {
		fmt.Println("not complete")
	}

	example := &common.ClientSessionData{}
	rv := reflect.ValueOf(example).Elem()
	rt := reflect.TypeOf(example).Elem()

	for _, tempData := range res.Logs {
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			newValue := rv.Field(i)

			if !newValue.CanSet() {
				fmt.Printf("field name: %s can not set value\n", field.Name)
			}

			tempValue, ok := tempData[field.Tag.Get("json")]
			if !ok {
				continue
			}

			switch newValue.Kind() {
			case reflect.String:
				newValue.SetString(tempValue)
			case reflect.Int32, reflect.Int64, reflect.Int8:
				i, err := strconv.ParseInt(tempValue, 10, 64)
				if err == nil {
					newValue.SetInt(i)
				}
			}
		}
		fmt.Println(example)
	}
}
