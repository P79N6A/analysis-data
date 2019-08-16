package common

import (
	"encoding/csv"
	"os"
	// "strconv"
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
	ConnectTime     string `gorm:"column:connect_time"`
	ConnectCount    string `gorm:"column:connect_count"`
	UseTime         string `gorm:"column:use_time"`
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
		record := []string{value.UserName, value.DeviceID, value.IsGuest, value.CountryZh, value.Location, value.AppVersion, value.PkgName, value.RegisterTime,
			value.InvitationCount, value.AddNU, value.LeftNU, value.ChangeNU, value.ConnectTime, value.ConnectCount, value.UseTime,
		}
		w.Write(record)
	}
	w.Flush()
	return nil
}
