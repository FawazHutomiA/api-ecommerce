package transaction

const (
	FIND_ALL = `
		SELECT 
			r.id, 
			r.product_id, 
			r.user_id,
			r.amount,
			r.status,
			r.code,
			r.payment_url
		FROM 
			transaction r
		WHERE 
			r.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT  
			r.id, 
			r.product_id, 
			r.user_id,
			r.amount,
			r.status,
			r.code,
			r.payment_url
		FROM 
			transaction r
		where 
			r.id = $1
	`

	FIND_BY_PRODUCT_ID = `
		SELECT  
			r.id, 
			r.product_id, 
			r.user_id,
			r.amount,
			r.status,
			r.code,
			r.payment_url
		FROM 
			transaction r
		where 
			r.product_id = $1
	`

	FIND_BY_USER_ID = `
		SELECT  
			r.id, 
			r.product_id, 
			r.user_id,
			r.amount,
			r.status,
			r.code,
			r.payment_url
		FROM 
			transaction r
		where 
			r.user_id = $1
	`

	INSERT = `
		INSERT INTO transaction (
			id, 
			product_id, 
			user_id,
			amount,
			status,
			code,
			payment_url
		) 
		VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6, 
			$7
		)
	`

	UPDATE_BY_ID = `
		UPDATE 
			transaction
		SET 
			product_id = $2,
			user_id = $3,
			amount = $4,
			status = $5,
			code = $6,
			payment_url = $7,
			updated_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_ID = `
		UPDATE 
			transaction
		SET 
			deleted_at = NOW()
		WHERE 
			id = $1
	`
)
