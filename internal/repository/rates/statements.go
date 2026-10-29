package rates

const (
	createUpdateStmt = `
		INSERT INTO rate_updates (id, base_currency, quote_currency, status)
		VALUES ($1, $2, $3, $4)`

	getUpdateStmt = `
		SELECT base_currency, quote_currency, status, price, rate_timestamp, COALESCE(error, ''), updated_at
		FROM rate_updates
		WHERE id = $1`

	getLatestStmt = `
		SELECT price, rate_timestamp, updated_at
		FROM rates
		WHERE base_currency = $1 AND quote_currency = $2`
)
