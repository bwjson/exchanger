package rates

const (
	createUpdateStmt = `
		INSERT INTO rate_updates (id, base_currency, quote_currency, status, idempotency_key)
		VALUES ($1, $2, $3, $4, NULLIF($5, ''))
		ON CONFLICT (idempotency_key) DO NOTHING
		RETURNING id`

	getUpdateByKeyStmt = `
		SELECT id, base_currency, quote_currency
		FROM rate_updates
		WHERE idempotency_key = $1`

	getUpdateStmt = `
		SELECT base_currency, quote_currency, status, price, rate_timestamp, COALESCE(error, ''), updated_at
		FROM rate_updates
		WHERE id = $1`

	getLatestStmt = `
		SELECT price, rate_timestamp, updated_at
		FROM rates
		WHERE base_currency = $1 AND quote_currency = $2`

	selectPendingStmt = `
		SELECT id, base_currency, quote_currency, attempts, next_attempt_at
		FROM rate_updates
		WHERE status = 'pending' AND next_attempt_at <= now()
		ORDER BY next_attempt_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED`

	saveUpdateStmt = `
		UPDATE rate_updates
		SET status = $2, price = $3, rate_timestamp = $4, error = NULLIF($5, ''),
			attempts = $6, next_attempt_at = $7, updated_at = now()
		WHERE id = $1`

	upsertRateStmt = `
		INSERT INTO rates (base_currency, quote_currency, price, rate_timestamp, updated_at)
		VALUES ($1, $2, $3, $4, now())
		ON CONFLICT (base_currency, quote_currency) DO UPDATE
		SET price = EXCLUDED.price, rate_timestamp = EXCLUDED.rate_timestamp, updated_at = now()
		WHERE rates.rate_timestamp <= EXCLUDED.rate_timestamp`
)
