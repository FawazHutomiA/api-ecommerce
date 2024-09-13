package user

const (
	FIND_ALL = `
		SELECT 
			u.id, 
			u.name, 
			u.email, 
			u.occupation, 
			u.phone, 
			u.gender, 
			u.role,
			u.is_google,
			u.is_active,
			u.is_verify
		FROM 
			users u
		WHERE 
    		u.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT 
			u.id, 
			u.name, 
			u.email, 
			u.occupation, 
			u.phone, 
			u.gender, 
			u.role,
			u.is_google,
			u.is_active,
			u.is_verify
		FROM 
			users u
		WHERE 
			u.id = $1;
	`

	FIND_BY_EMAIL = `
		SELECT 
			u.id, 
			u.name, 
			u.email, 
			u.occupation, 
			u.password,
			u.phone, 
			u.gender, 
			u.role,
			u.is_google,
			u.is_active,
			u.is_verify
		FROM 
			users u
		WHERE 
			u.email = $1;
	`

	INSERT = `
		INSERT INTO 
			users (
				id, 
				name,
				email, 
				occupation,
				password, 
				phone,
				role, 
				gender,
				is_google,
				token
			) 
		VALUES (
			$1, 
			$2, 
			$3, 
			$4, 
			$5, 
			$6, 
			$7,
			$8,
			$9,
			$10
		)
	`

	UPDATE_TOKEN_USER = `
		UPDATE 
			users
		SET 
			token = $2, 
			updated_at = NOW()
		WHERE id = $1
	`
)
