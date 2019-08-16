package sql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/op/go-logging"

	"github.com/analysis-data/analysisSession/common"
)

type Cluster struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

var (
	log    = logging.MustGetLogger("sql")
	DB     *Cluster
	dbAddr = "cliUsr:CLiE^R#(WW&%A9QEDp201252*92VPusS#$8203t@tcp(climbvpn.cbrhwddmnfax.ap-southeast-1.rds.amazonaws.com:33061)/climb?charset=utf8"
	// dbAddr = "root:fvsh2225@tcp(61.160.47.31:3306)/climb?charset=utf8"
	// dbAddr = "root:123456@(localhost:3306)/climb?charset=utf8"
	// dbAddr = "root:123456@(221.228.197.195:3307)/climb?charset=utf8"
	// dbAddr = "root:123456@(150.109.51.108:3307)/climb?charset=utf8"
)

func init() {
	var err error
	DB, err = NewCluster(dbAddr)
	if err != nil {
		log.Error("sql init error: " + err.Error())
	}
}

func NewCluster(dbUri string) (*Cluster, error) {
	db, err := sql.Open("mysql", dbUri)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.SetConnMaxLifetime(60000)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	cluster := Cluster{
		db:    db,
		stmts: map[string]*sql.Stmt{},
	}

	if err = cluster.initStmts(); err != nil {
		return nil, err
	}
	return &cluster, nil
}

// type ClientSession struct {
// 	ColumnID                int64
// 	ID                      string
// 	Finished                bool
// 	UserName                string
// 	AppVersion              string
// 	State                   int8
// 	Error                   int8
// 	Errstr                  string
// 	EnableDurationMs        int64
// 	DisableCostMs           int64
// 	RouterAllocCount        int32
// 	RouterAllocSuccessCount int32
// 	CreateTime              string
// }

func (cluster *Cluster) QueryClientSessionErrorDetailByCountry(countryEn, startTime, endTime string) ([]common.ClientSession, error) {
	var arrayData []common.ClientSession
	rows, err := cluster.stmts["stmt_query_error_detail_by_country"].Query(countryEn, startTime, endTime, countryEn, startTime, endTime)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var data common.ClientSession
		err = rows.Scan(&data.ColumnID, &data.ID, &data.Finished, &data.UserName, &data.AppVersion, &data.State,
			&data.Error, &data.Errstr, &data.EnableDurationMs, &data.DisableCostMs, &data.RouterAllocCount, &data.RouterAllocSuccessCount, &data.CreateTime)
		if err != nil {
			log.Error("QueryClientSessionErrorDetailByCountry error: " + err.Error())
			continue
		}
		arrayData = append(arrayData, data)
	}

	return arrayData, nil
}

func (cluster *Cluster) QueryClientSessionErrorDetail(startTime, endTime string) ([]common.ClientSession, error) {
	var arrayData []common.ClientSession
	rows, err := cluster.stmts["stmt_query_error_detail"].Query(startTime, endTime, startTime, endTime)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var data common.ClientSession
		err = rows.Scan(&data.ColumnID, &data.ID, &data.Finished, &data.UserName, &data.AppVersion, &data.State,
			&data.Error, &data.Errstr, &data.EnableDurationMs, &data.DisableCostMs, &data.RouterAllocCount, &data.RouterAllocSuccessCount, &data.CreateTime)
		if err != nil {
			log.Error("QueryClientSessionErrorDetail error: " + err.Error())
			continue
		}
		arrayData = append(arrayData, data)
	}

	return arrayData, nil
}

// type ClientSessionConnectionResult struct {
// 	Toatal           int32
// 	SuccessCount     int32
// 	ConnectionCount  int32
// 	DisableCount     int32
// 	ErrorCount       int32
// 	RegistingCount   int32
// 	UnRegistingCount int32
// }
// stmt_query_client_session_cat_by_country
func (cluster *Cluster) QueryClientSessionConnectionResultByCountry(country, startTime, endTime string) (common.ClientSessionConnectionResult, error) {
	var result common.ClientSessionConnectionResult
	err := cluster.stmts["stmt_query_client_session_cat_by_country"].QueryRow(country, startTime, endTime, country, startTime, endTime).
		Scan(&result.Total, &result.SuccessCount, &result.ConnectionCount, &result.DisableCount, &result.ErrorCount, &result.RegistingCount, &result.UnRegistingCount)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (cluster *Cluster) QueryClientSessionConnectionResult(startTime, endTime string) (common.ClientSessionConnectionResult, error) {
	var result common.ClientSessionConnectionResult
	err := cluster.stmts["stmt_query_client_session_cat"].QueryRow(startTime, endTime, startTime, endTime).
		Scan(&result.Total, &result.SuccessCount, &result.ConnectionCount, &result.DisableCount, &result.ErrorCount, &result.RegistingCount, &result.UnRegistingCount)
	if err != nil {
		return result, err
	}

	return result, nil
}
