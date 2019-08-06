package db

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/analysis-data/analysisDns-query/common"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Records struct {
	Records []Record `json:"RECORDS"`
}

type Record struct {
	TargetAddress string `json:"targetAddress" orm:"column(targetAddress)"`
	ContryZh      string `json:"countryZh" orm:"column(countryZh)"`
	PkgName       string `json:"pkgName" orm:"column(pkgName)"`
	AppVersion    string `json:"appVersion" orm:"column(appVersion)"`
	Count         string `json:"count" orm:"column(count)"`
}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "cliUsr:CLiE^R#(WW&%A9QEDp201252*92VPusS#$8203t@tcp(climbvpn.cbrhwddmnfax.ap-southeast-1.rds.amazonaws.com:33061)/climb?charset=utf8")
}

var querySQL = "select targetAddress, countryZh, pkgName, appVersion, count(userName) as count from tp_client_conn%s group by targetAddress, countryZh, pkgName, userName"

func QueryData(startTime string) (*Records, error) {
	o := orm.NewOrm()
	o.Using("climb")

	var sql string
	sql = fmt.Sprintf(sql, querySQL, startTime)

	var records []Record
	_, err := o.Raw(sql).QueryRows(&records)
	if err != nil {
		return nil, err
	}

	re := Records{Records: records}
	err = common.WriteToFile("records"+startTime+".json", re)
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
