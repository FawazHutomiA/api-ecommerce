package token

const (
	FIND_ALL = `
		SELECT 
			u.id,
			u.user_id,
			u.token,
			u.expired_at
		FROM 
			token u
		WHERE 
    		u.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT 
			u.id,
			u.user_id,
			u.token,
			u.expired_at
		FROM 
			token u
		WHERE 
			u.id = $1;
	`

	FIND_BY_USER_ID = `
		SELECT 
			u.id,
			u.user_id,
			u.token,
			u.expired_at
		FROM 
			token u
		WHERE 
			u.user_id = $1;
	`

	INSERT = `
		INSERT INTO 
			token (
				id, 
				user_id,
				token,
				expired_at
			) 
		VALUES (
			$1, 
			$2, 
			$3,
			$4
		)
	`

	UPDATE_BY_ID = `
		UPDATE 
			token
		SET 
			user_id = $2,
			token = $3,
			expired_at = $4,
			updated_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_ID = `
		UPDATE 
			token
		SET 
			deleted_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_USER_ID = `
		DELETE FROM 
    		token
		WHERE 
			user_id = $1;
	`
)
