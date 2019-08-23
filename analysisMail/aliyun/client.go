package aliyun

import (
	"fmt"
	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/analysis-data/analysisMail/common"
	"reflect"
	"strconv"
)

// Client 阿里云sdk client
type AliyunClient struct {
	Runmode string
	Client  *sls.Client
}

func NewClient(runmode string) *AliyunClient {
	client := &sls.Client{
		Endpoint:        "ap-southeast-1.log.aliyuncs.com",
		AccessKeyID:     "LTAIXDbKKHuwwjo4",
		AccessKeySecret: "xL2pTLXvFpwwvpKOYx7WgfDLmZzB91",
	}

	return &AliyunClient{Client: client, Runmode: runmode}
}

func (c *AliyunClient) GetSessionData(from int64, to int64, maxLineNum int64, offset int64) ([]common.ClientSessionData, error) {
	var datas []common.ClientSessionData

	res, err := c.GetLogs("sfox-"+c.Runmode, "tp_client_session", from, to, "cc.coolline.client", maxLineNum, offset, false)
	if err != nil {
		return datas, err
	}

	for _, tempData := range res.Logs {
		example := common.ClientSessionData{}
		exchangeStruct(tempData, &example)
		datas = append(datas, example)
	}

	return datas, nil
}

func (c *AliyunClient) GetConnData(from int64, to int64, id string, maxLineNum int64, offset int64) ([]common.ClientConnData, error) {
	var datas []common.ClientConnData

	res, err := c.GetLogs("sfox-"+c.Runmode, "tp_client_conn", from, to, id, maxLineNum, offset, false)
	if err != nil {
		return datas, err
	}

	if !res.IsComplete() {
		return datas, fmt.Errorf("not complete")
	}

	for _, tempData := range res.Logs {
		example := common.ClientConnData{}
		exchangeStruct(tempData, &example)
		datas = append(datas, example)
	}
	return datas, nil
}

func exchangeStruct(source map[string]string, data interface{}) {
	rv := reflect.ValueOf(data).Elem()
	rt := reflect.TypeOf(data).Elem()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		newValue := rv.Field(i)

		if !newValue.CanSet() {
			fmt.Printf("field name: %s can not set value\n", field.Name)
			continue
		}

		tempValue, ok := source[field.Tag.Get("json")]
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
}

func (c *AliyunClient) GetLogs(project, logstore string, from int64, to int64,
	queryExp string, maxLineNum int64, offset int64, reverse bool) (*sls.GetLogsResponse, error) {
	res, err := c.Client.GetLogs(project, logstore, "", from, to, queryExp, maxLineNum, offset, reverse)
	if err != nil {
		return res, err
	}

	return res, nil
}
