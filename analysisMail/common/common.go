package common

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

func TimeFormat(cur_time time.Time) string {
	str := cur_time.Format("2006-01-02 15:04:05")
	return str
}

func TimeFormatDate(cur_time time.Time) string {
	str := cur_time.Format("20060102")
	return str
}

func TimeFormatDateFromUnix(cur_time int64) string {
	curTime := time.Unix(cur_time, 0)
	return TimeFormatDate(curTime)
}

type AnalysisData struct {
	UserName        string `gorm:"column:userName"`
	DeviceID        string `gorm:"column:deviceid"`
	IsGuest         string `gorm:"column:isGuest"`
	CountryZh       string `gorm:"column:countryZh"`
	Location        string `gorm:"column:location"`
	AppVersion      string `gorm:"column:appVersion"`
	PkgName         string `gorm:"column:pkgName"`
	RegisterTime    string `gorm:"column:registerTime"`
	InvitationCount string `gorm:"column:invitationCount"`
	AddNU           string `gorm:"column:add_nu"`
	LeftNU          string `gorm:"column:left_nu"`
	ChangeNU        string `gorm:"column:change_nu"`
	ConnectTime     int64  `gorm:"column:connect_time"`
	ConnectCount    int    `gorm:"column:connect_count"`
	UseTime         int64  `gorm:"column:use_time"`
	LastConnectTime int64
}

type ClientConnAnalysisData struct {
	ID             string
	RemoteInput    int32
	RemoteOutput   int32
	ViaUserName    string
	ViaUserDevice  string
	TargetAddress  string
	ConnCreateTime int64
	ConnCloseTime  int64
}
type ClientSessionAnalysisData struct {
	UserName     string
	UserDevice   string
	ID           string
	PkgName      string
	AppVersion   string
	SelectRouter string
	RemoteInput  int32
	RemoteOutput int32
	StartTime    int64
	EndTime      int64
	Conns        []ClientConnAnalysisData
}

type ClientSessionData struct {
	UserName        string `gorm:"column:userName"`
	UserDevice      string `gorm:"column:userDevice"`
	ID              string `gorm:"column:id"`
	Finished        string `gorm:"column:finished"`
	PkgName         string `gorm:"column:pkgName"`
	AppVersion      string `gorm:"column:appVersion"`
	SelectRouter    string `gorm:"column:selectRouter"`
	RemoteInput     int32  `gorm:"column:remoteInput"`
	RemoteOutput    int32  `gorm:"column:remoteOutput"`
	Established     int8   `gorm:"column:established"`
	CreateTimestamp int64  `gorm:"column:createTimestamp"`
}

func (ClientSessionData) TableName() string {
	return "tp_client_session"
}

type ClientConnData struct {
	UserName        string `gorm:"column:userName"`
	UserDevice      string `gorm:"column:userDevice"`
	ID              string `gorm:"column:id"`
	SessionID       string `gorm:"column:sessionId"`
	Finished        string `gorm:"column:finished"`
	RemoteInput     int32  `gorm:"column:remoteInput"`
	RemoteOutput    int32  `gorm:"column:remoteOutput"`
	ViaUserName     string `gorm:"column:viaUserName"`
	ViaUserDevice   string `gorm:"column:viaUserDevice"`
	TargetAddress   string `gorm:"column:targetAddress"`
	ConnEstablished int8   `gorm:"column:connEstablished"`
	ConnCreateTime  int64  `gorm:"column:connCreateTime"`
	ConnCloseTime   int64  `gorm:"column:connCloseTime"`
	CreateTimestamp int64  `gorm:"column:createTimestamp"`
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
	title := []string{
		"userName",
		"deviceId",
		"isGuest",
		"countryZh",
		"location",
		"appVersion",
		"pkgName",
		"registerTime",
		"invitationCount",
		"add_nu",
		"left_nu",
		"change_nu",
		"connect_time",
		"connect_count",
		"use_time",
	}

	w.Write(title)

	for _, value := range data {
		connectTime := strconv.FormatInt(value.ConnectTime, 10)
		connectCount := strconv.FormatInt(int64(value.ConnectCount), 10)
		useTime := strconv.FormatInt(value.UseTime, 10)
		record := []string{value.UserName, value.DeviceID, value.IsGuest, value.CountryZh, value.Location, value.AppVersion, value.PkgName, value.RegisterTime,
			value.InvitationCount, value.AddNU, value.LeftNU, value.ChangeNU, connectTime, connectCount, useTime,
		}
		w.Write(record)
	}
	w.Flush()
	return nil
}
