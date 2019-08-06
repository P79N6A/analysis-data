package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

var SQLStr = `SELECT
T.userName,
case when C1.cur_flow > 0 then C1.cur_flow else 0 end as cur_flow,
case when T1.next_login > 0 then 1 else 0 end as next_login,
case when C2.next_flow > 0 then C2.next_flow else 0 end as next_flow
FROM
( SELECT DISTINCT ( userName ) FROM sys_user WHERE createTimestamp >= ? AND createTimestamp <= ? AND countryEn = 'United States' ) T
LEFT JOIN (
SELECT
	userName,
	sum( remoteInput + remoteOutput ) / 1024 / 1024 AS cur_flow 
FROM
	(
	SELECT
		* 
	FROM
		(
		SELECT
			* 
		FROM
			tp_client_session_temp06170627 
		WHERE
			createTimestamp >= ?
			AND createTimestamp <= ?
			AND countryEn = 'United States' 
		ORDER BY
			id,
			createTimestamp DESC 
			LIMIT 999999 
		) T 
	GROUP BY
		id 
	) origin 
GROUP BY
	userName 
) C1 ON C1.userName = T.userName
LEFT JOIN (
SELECT
	userName,
	count( * ) AS next_login 
FROM
	sys_login
WHERE
	createTimestamp >= ?
	AND createTimestamp <= ? 
	AND countryEn = 'United States' 
GROUP BY
	userName 
) T1 ON T1.userName = T.userName
LEFT JOIN (
SELECT
	userName,
	sum( remoteInput + remoteOutput ) / 1024 / 1024 AS next_flow 
FROM
	(
	SELECT
		* 
	FROM
		(
		SELECT
			* 
		FROM
			tp_client_session_temp06170627
		WHERE
			createTimestamp >= ?
			AND createTimestamp <= ? 
			AND countryEn = 'United States' 
		ORDER BY
			id,
			createTimestamp DESC 
			LIMIT 999999 
		) T 
	GROUP BY
		id 
	) origin 
GROUP BY
userName 
) C2 ON C2.userName = T.userName`

type GormInterface struct {
	gormDB *gorm.DB
}

func RegisterDB() (*GormInterface, error) {
	// (可选)设置最大空闲连接
	maxIdle := 60
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := 300
	dbLink := "cliUsr:CLiE^R#(WW&%A9QEDp201252*92VPusS#$8203t@tcp(climbvpn.cbrhwddmnfax.ap-southeast-1.rds.amazonaws.com:33061)/climb?charset=utf8"

	db, err := gorm.Open("mysql", dbLink)
	if err != nil {
		return nil, fmt.Errorf("RegisterDateBase error: " + err.Error())
	}

	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxConn)

	return &GormInterface{gormDB: db}, nil
}

func (db *GormInterface) Query(startTime time.Time, data interface{}) error {
	nextTime := startTime.Add(time.Minute * 24 * 60)
	nnTime := nextTime.Add(time.Minute * 24 * 60)

	dbTemp := db.gormDB.Exec(SQLStr, startTime.Unix(), nextTime.Unix(), startTime.Unix(), nextTime.Unix(), nextTime.Unix(), nnTime.Unix(), nextTime.Unix(), nnTime.Unix()).Find(data)
	if dbTemp.Error != nil {
		return dbTemp.Error
	}
	return nil
}

func (db *GormInterface) BuildSQL(startTime time.Time) string {
	nextTime := startTime.Add(time.Minute * 24 * 60)
	nnTime := nextTime.Add(time.Minute * 24 * 60)

	sql := strings.Replace(SQLStr, "?", "%d", -1)
	sql = fmt.Sprintf(sql, startTime.Unix(), nextTime.Unix(), startTime.Unix(), nextTime.Unix(), nextTime.Unix(), nnTime.Unix(), nextTime.Unix(), nnTime.Unix())
	return sql
}
