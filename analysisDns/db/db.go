package db

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/analysis-data/analysisDns/common"
)

type Records struct {
	Records []Record `json:"RECORDS"`
}

type Record struct {
	UserDevice    string `json:"userDevice" orm:"column(userDevice)"`
	TargetAddress string `json:"targetAddress" orm:"column(targetAddress)"`
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "cliUsr:CLiE^R#(WW&%A9QEDp201252*92VPusS#$8203t@tcp(climbvpn.cbrhwddmnfax.ap-southeast-1.rds.amazonaws.com:33061)/climb?charset=utf8")
}

var createTableSQl = "CREATE tp_share_conn_back like tp_share_conn"
var querySQL = "SELECT T.userDevice, T.targetAddress FROM ( SELECT DISTINCT ( userDevice ) FROM `tp_share_conn_back20190424-20190430` WHERE userDevice != '' AND UNIX_TIMESTAMP(createTime) >= ? and UNIX_TIMESTAMP(createTime) <= ?) origin LEFT JOIN ( SELECT * FROM tp_share_conn WHERE closeReason = 'dns-error' AND targetAddress != '' AND UNIX_TIMESTAMP(createTime) >= ? and UNIX_TIMESTAMP(createTime) <= ?) T ON T.userDevice = origin.userDevice order by T.userDevice"

// func CreateDataBack(from, to time.Time) error {
// 	o := orm.NewOrm()
// 	o.Using("climb")

// 	o.
// }

func QueryData(from, to time.Time) (*Records, error) {
	o := orm.NewOrm()
	o.Using("climb")

	var records []Record
	_, err := o.Raw(querySQL, from.Unix(), to.Unix(), from.Unix(), to.Unix()).QueryRows(&records)
	if err != nil {
		return nil, err
	}
	
	re := Records{Records: records}
	err = common.WriteToFile("records.json", re)
	if err != nil {
		return nil, err
	}
	return &re, nil
}

func ReadRecods(filePath string) (*Records, error) {
	var records Records
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return &records, err
	}

	err = json.Unmarshal(fileBytes, &records)
	if err != nil {
		return &records, err
	}

	return &records, nil
}