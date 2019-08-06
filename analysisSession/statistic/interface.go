package statistic

// sfox_client_error_none,
// sfox_client_error_internal,
// sfox_client_error_no_conn_info,
// sfox_client_error_network_error,
// sfox_client_error_user_login_fail,
// sfox_client_error_user_login_timeout,
// sfox_client_error_router_alloc_fail,
// sfox_client_error_router_p2p_fail,

type StatisticClientSessionErrorState struct {
	AllocError    int32
	InternalError int32
	ConnError     int32
	NetWorkError  int32
	LoginError    int32
	LoginTimeOut  int32
	P2PError      int32
}

type StatisticClientSessionState struct {
	Success     int32
	Connecting  int32
	Disable     int32
	Error       StatisticClientSessionErrorState
	Registing   int32
	Unregisting int32
}
