package sql

var stmtDefins = map[string]string{
	"stmt_query_error_detail_by_country": `SELECT
			columnId,
			id,
			finished,
			userName,
			appVersion,
			state,
			error,
			errstr,
			enableDurationMs,
			disableCostMs,
			routerAllocCount,
			routerAllocSuccessCount,
			createTime 
		FROM
			tp_client_session 
		WHERE
			id IN (
			SELECT
				id 
			FROM
				(
				SELECT
					origin.id AS id,
					success.id AS successId,
					origin.state,
					origin.finished 
				FROM
					(
					SELECT
						max( columnId ) AS columnId 
					FROM
						tp_client_session 
					WHERE
						userName != '' 
						AND pkgName = 'cc.shadowfox.client' 
						AND countryEn = ? 
						AND createTime >= ? 
						AND createTime < ? 
					GROUP BY
						id 
					) framework
					LEFT JOIN tp_client_session AS origin ON framework.columnId = origin.columnId
					LEFT JOIN (
					SELECT
						id 
					FROM
						tp_client_session 
					WHERE
						userName != '' 
						AND pkgName = 'cc.shadowfox.client' 
						AND state = 'sfox-cli-established' 
						AND countryEn = ? 
						AND createTime >= ?
						AND createTime < ?
					GROUP BY
						id 
					) success ON origin.id = success.id 
				) S 
			WHERE
				successId IS NULL 
				AND state = 'sfox-cli-error' 
			) 
		ORDER BY
			id`,
	"stmt_query_error_detail": `SELECT
			columnId,
			id,
			finished,
			userName,
			appVersion,
			state,
			error,
			errstr,
			enableDurationMs,
			disableCostMs,
			routerAllocCount,
			routerAllocSuccessCount,
			createTime 
		FROM
			tp_client_session 
		WHERE
			id IN (
			SELECT
				id 
			FROM
				(
				SELECT
					origin.id AS id,
					success.id AS successId,
					origin.state,
					origin.finished 
				FROM
					(
					SELECT
						max( columnId ) AS columnId 
					FROM
						tp_client_session 
					WHERE
						userName != '' 
						AND pkgName = 'cc.shadowfox.client' 
						AND createTime >= ? 
						AND createTime < ? 
					GROUP BY
						id 
					) framework
					LEFT JOIN tp_client_session AS origin ON framework.columnId = origin.columnId
					LEFT JOIN (
					SELECT
						id 
					FROM
						tp_client_session 
					WHERE
						userName != '' 
						AND pkgName = 'cc.shadowfox.client' 
						AND state = 'sfox-cli-established' 
						AND createTime >= ?
						AND createTime < ?
					GROUP BY
						id 
					) success ON origin.id = success.id 
				) S 
			WHERE
				successId IS NULL 
				AND state = 'sfox-cli-error' 
			) 
		ORDER BY
			id`,
	"stmt_query_client_session_cat_by_country": `SELECT
			count( * ) AS total,
			sum( CASE WHEN successId IS NULL THEN 0 ELSE 1 END ) AS successCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-connecting' THEN 1 ELSE 0 END ) AS connectingCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-disable' THEN 1 ELSE 0 END ) AS disableCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-error' THEN 1 ELSE 0 END ) AS errorCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-registing' THEN 1 ELSE 0 END ) AS registingCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-unregisting' THEN 1 ELSE 0 END ) AS unregistingCount 
		FROM
			(
			SELECT
				origin.id AS id,
				success.id AS successId,
				origin.state,
				origin.finished 
			FROM
				(
				SELECT
					max( columnId ) AS columnId 
				FROM
					tp_client_session 
				WHERE
					userName != '' 
					AND pkgName = 'cc.shadowfox.client' 
					AND countryEn = ? 
					AND createTime >= ? 
					AND createTime < ?
				GROUP BY
					id 
				) framework
				LEFT JOIN tp_client_session AS origin ON framework.columnId = origin.columnId
				LEFT JOIN (
				SELECT
					id 
				FROM
					tp_client_session 
				WHERE
					userName != '' 
					AND pkgName = 'cc.shadowfox.client' 
					AND state = 'sfox-cli-established' 
					AND countryEn = ? 
					AND createTime >= ?
					AND createTime < ?
				GROUP BY
					id 
				) success ON origin.id = success.id 
			) S`,
	"stmt_query_client_session_cat": `SELECT
			count( * ) AS total,
			sum( CASE WHEN successId IS NULL THEN 0 ELSE 1 END ) AS successCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-connecting' THEN 1 ELSE 0 END ) AS connectingCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-disable' THEN 1 ELSE 0 END ) AS disableCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-error' THEN 1 ELSE 0 END ) AS errorCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-registing' THEN 1 ELSE 0 END ) AS registingCount,
			sum( CASE WHEN successId IS NULL AND state = 'sfox-cli-unregisting' THEN 1 ELSE 0 END ) AS unregistingCount 
		FROM
			(
			SELECT
				origin.id AS id,
				success.id AS successId,
				origin.state,
				origin.finished 
			FROM
				(
				SELECT
					max( columnId ) AS columnId 
				FROM
					tp_client_session 
				WHERE
					userName != '' 
					AND pkgName = 'cc.shadowfox.client' 
					AND createTime >= ? 
					AND createTime < ?
				GROUP BY
					id 
				) framework
				LEFT JOIN tp_client_session AS origin ON framework.columnId = origin.columnId
				LEFT JOIN (
				SELECT
					id 
				FROM
					tp_client_session 
				WHERE
					userName != '' 
					AND pkgName = 'cc.shadowfox.client' 
					AND state = 'sfox-cli-established' 
					AND createTime >= ?
					AND createTime < ?
				GROUP BY
					id 
				) success ON origin.id = success.id 
			) S`,
}

func (cluster *Cluster) initStmts() error {
	for name, def := range stmtDefins {
		stmt, err := cluster.db.Prepare(def)
		if err != nil {
			log.Error("Cluster: prepare stmt " + name + " error, " + err.Error() + ", [" + def + "]")
			return err
		}
		cluster.stmts[name] = stmt
	}
	return nil
}
