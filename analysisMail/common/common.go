package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

func TimeFormatDate(cur_time time.Time) string {
	str := cur_time.Format("20060102")
	return str
}

func TimeFormatDateFromUnix(cur_time int64) string {
	curTime := time.Unix(cur_time, 0)
	return TimeFormatDate(curTime)
}

type AnalysisData struct {
	UserName          string `json:"userName" gorm:"column:userName"`
	DeviceID          string `json:"deviceId" gorm:"column:deviceid"`
	IsGuest           string `json:"isGuest" gorm:"column:isGuest"`
	CountryZh         string `json:"countryZh" gorm:"column:countryZh"`
	Location          string `json:"location" gorm:"column:location"`
	AppVersion        string `json:"appVersion" gorm:"column:appVersion"`
	PkgName           string `json:"pkgName" gorm:"column:pkgName"`
	RegisterTime      string `json:"registerTime" gorm:"column:registerTime"`
	InvitationCount   string `json:"invitationCount" gorm:"column:invitationCount"`
	AddNU             string `json:"add_nu" gorm:"column:add_nu"`
	LeftNU            string `json:"left_nu" gorm:"column:left_nu"`
	ChangeNU          string `json:"change_nu" gorm:"column:change_nu"`
	ConnectTime       int64  `json:"connect_time" gorm:"column:connect_time"`
	ConnectCount      int    `json:"connect_count" gorm:"column:connect_count"`
	UseTime           int64  `json:"use_time" gorm:"column:use_time"`
	IsVip             string `json:"isVip" gorm:"column:isVip"`
	VipStartDate      string `json:"vipStartDate" gorm:"column:vipStartDate"`
	VipExpirationDate string `json:"vipExpireDate" gorm:"column:vipExpireDate"`
	LastConnectTime   int64
}

type ClientConnAnalysisData struct {
	ID             string
	RemoteInput    int64
	RemoteOutput   int64
	ViaUserName    string
	ViaUserDevice  string
	TargetAddress  string
	ConnCreateTime int64
	ConnCloseTime  int64
}
type ClientSessionAnalysisData struct {
	UserName      string
	UserDevice    string
	ID            string
	PkgName       string
	AppVersion    string
	SelectRouter  string
	RemoteInput   int64
	RemoteOutput  int64
	StartTime     int64
	EndTime       int64
	ConnBeginTime int64
	ConnEndTime   int64
	Conns         []ClientConnAnalysisData
}

type ClientSessionData struct {
	UserName        string `json:"userName" gorm:"column:userName"`
	UserDevice      string `json:"userDevice" gorm:"column:userDevice"`
	ID              string `json:"id" gorm:"column:id"`
	Finished        string `json:"finished" gorm:"column:finished"`
	PkgName         string `json:"pkgName" gorm:"column:pkgName"`
	AppVersion      string `json:"appVersion" gorm:"column:appVersion"`
	SelectRouter    string `json:"selectRouter" gorm:"column:selectRouter"`
	RemoteInput     int64  `json:"remoteInput" gorm:"column:remoteInput"`
	RemoteOutput    int64  `json:"remoteOutput" gorm:"column:remoteOutput"`
	Established     int8   `json:"established" gorm:"column:established"`
	EnableTime      int64  `json:"enableTime" gorm:"column:enableTime"`
	CreateTimestamp int64  `json:"__tag__:__receive_time__" gorm:"column:createTimestamp"`
}

func (ClientSessionData) TableName() string {
	return "tp_client_session"
}

type ClientConnData struct {
	UserName        string `json:"userName" gorm:"column:userName"`
	UserDevice      string `json:"userDevice" gorm:"column:userDevice"`
	ID              string `json:"id" gorm:"column:id"`
	SessionID       string `json:"sessionId" gorm:"column:sessionId"`
	Finished        string `json:"finished" gorm:"column:finished"`
	RemoteInput     int64  `json:"remoteInput" gorm:"column:remoteInput"`
	RemoteOutput    int64  `json:"remoteOutput" gorm:"column:remoteOutput"`
	ViaUserName     string `json:"viaUserName" gorm:"column:viaUserName"`
	ViaUserDevice   string `json:"viaUserDevice" gorm:"column:viaUserDevice"`
	TargetAddress   string `json:"targetAddress" gorm:"column:targetAddress"`
	ConnEstablished int8   `json:"connEstablished" gorm:"column:connEstablished"`
	ConnCreateTime  int64  `json:"connCreateTime" gorm:"column:connCreateTime"`
	ConnCloseTime   int64  `json:"connCloseTime" gorm:"column:connCloseTime"`
	CreateTimestamp int64  `json:"__tag__:__receive_time__" gorm:"column:createTimestamp"`
}

func WriteDataToFile(pathFile string, data []AnalysisData) error {
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

	title := []string{}

	example := &AnalysisData{}
	rt := reflect.TypeOf(example).Elem()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		tempValue := field.Tag.Get("json")
		if tempValue == "" {
			continue
		}
		title = append(title, tempValue)
	}

	w.Write(title)

	for _, value := range data {
		var record []string
		newValue := reflect.ValueOf(value)
		for i := 0; i < newValue.NumField(); i++ {
			field := newValue.Field(i)
			var changeValue string
			switch field.Kind() {
			case reflect.Int, reflect.Int32, reflect.Int64, reflect.Int8:
				changeValue = strconv.FormatInt(field.Int(), 10)
			case reflect.String:
				changeValue = field.String()
			default:
				fmt.Println("write file type unKnow: ", field.Kind())
			}

			record = append(record, changeValue)
		}
		w.Write(record)
	}
	w.Flush()
	return nil
}
