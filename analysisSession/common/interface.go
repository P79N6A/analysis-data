package common

type ClientSession struct {
	ColumnID                int64
	ID                      string
	Finished                bool
	UserName                string
	AppVersion              string
	State                   string
	Error                   int8
	Errstr                  string
	EnableDurationMs        int64
	DisableCostMs           int64
	RouterAllocCount        int32
	RouterAllocSuccessCount int32
	CreateTime              string
}

type ClientSessionConnectionResult struct {
	Total            int32
	SuccessCount     int32
	ConnectionCount  int32
	DisableCount     int32
	ErrorCount       int32
	RegistingCount   int32
	UnRegistingCount int32
}
