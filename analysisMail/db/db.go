package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"strings"
	"time"

	"github.com/analysis-data/analysisMail/common"
)

var SQLStr = `SELECT
origin.*,
UU.createTime AS registerTime,
( CASE WHEN UU.usageCount IS NULL THEN 0 ELSE UU.usageCount END ) AS invitationCount,
UU.userId,
UU.left_nu,
UU.add_nu,
UU.change_nu
FROM
(
SELECT
	userName,
	deviceid,
	(case when loginType = 'guest' then '是' else '否' end) as isGuest,
	countryZh,
	appVersion,
	pkgName 
FROM
	(
	SELECT
		userName,
		deviceid,
		loginType,
		createTimestamp,
		countryZh,
		appVersion,
		pkgName 
	FROM
		sys_login 
	WHERE
		pkgName = 'cc.coolline.client' 
		AND createTimestamp >= ?
		AND createTimestamp <= ? 
	ORDER BY
		createTimestamp DESC 
		LIMIT 999999 
	) T 
GROUP BY
	userName 
) origin
LEFT JOIN (
SELECT
	createTime,
	U.userId,
	U.userName,
	I.usageCount,
	R.left_nu,
	R.change_nu,
	R.add_nu 
FROM
	( SELECT createTime, userName, userId FROM sys_user WHERE pkgName = 'cc.coolline.client' ) U
	LEFT JOIN ( SELECT userId, usageCount FROM sys_invitation_code ) I ON I.userId = U.userId
	LEFT JOIN (
	SELECT
		userName,
		userId,
		left_nu,
		change_nu,
		add_nu 
	FROM
		( SELECT self_address, sum( change_nu ) AS left_nu FROM ws_user_transaction_record GROUP BY self_address ) T
		LEFT JOIN ( SELECT userName, address, userId FROM sys_user_key ) S ON S.address = T.self_address
		LEFT JOIN ( SELECT self_address as user_address, sum( change_nu ) AS add_nu FROM ws_user_transaction_record WHERE event_code != 'using' GROUP BY self_address ) TP ON TP.user_address = T.self_address
		LEFT JOIN (
		SELECT
			self_address,
			sum( change_nu ) AS change_nu 
		FROM
			ws_user_transaction_record 
		WHERE
			create_time >= FROM_UNIXTIME(?) 
			AND create_time <= FROM_UNIXTIME(?)
			AND event_code = 'using' 
		GROUP BY
			self_address 
		) L ON L.self_address = T.self_address 
	) R ON R.userId = U.userId 
) UU ON UU.userName = origin.userName `

type GormInterface struct {
	gormDB *gorm.DB
}

func RegisterDB(runMode string) (*GormInterface, error) {
	// (可选)设置最大空闲连接
	maxIdle := 60
	// (可选) 设置最大数据库连接 (go >= 1.2)
	maxConn := 300

	dbLink := "cliUsr:CLiE^R#(WW&%A9QEDp201252*92VPusS#$8203t@tcp(climbvpn.cbrhwddmnfax.ap-southeast-1.rds.amazonaws.com:33061)/climb?charset=utf8"
	if runMode == "dev" {
		dbLink = "root:123456@tcp(58.215.139.156:3307)/climb?charset=utf8"
	}

	db, err := gorm.Open("mysql", dbLink)
	if err != nil {
		return nil, fmt.Errorf("RegisterDateBase error: " + err.Error())
	}

	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxConn)

	return &GormInterface{gormDB: db}, nil
}

func (db *GormInterface) Query(startTime time.Time, data interface{}) error {
	queryTime := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
	fmt.Println("query time: ", common.TimeFormat(queryTime))

	preTime := queryTime.Unix() - 60*24*60

	dbTemp := db.gormDB.Raw(SQLStr, preTime, queryTime.Unix(), preTime, queryTime.Unix()).Scan(data)
	if dbTemp.Error != nil {
		return dbTemp.Error
	}
	return nil
}

func (db *GormInterface) Close() {
	db.gormDB.Close()
}

func (db *GormInterface) BuildSQL(startTime time.Time) string {
	nextTime := startTime.Add(time.Minute * 24 * 60)
	nnTime := nextTime.Add(time.Minute * 24 * 60)

	sql := strings.Replace(SQLStr, "?", "%d", -1)
	sql = fmt.Sprintf(sql, startTime.Unix(), nextTime.Unix(), startTime.Unix(), nextTime.Unix(), nextTime.Unix(), nnTime.Unix(), nextTime.Unix(), nnTime.Unix())
	return sql
}
