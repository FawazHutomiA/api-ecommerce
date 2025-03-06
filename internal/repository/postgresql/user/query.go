package user

const (
	FIND_ALL = `
		SELECT 
			u.id,
			u.role_id,
			u.name, 
			u.email, 
			u.phone, 
			u.gender, 
			u.birth,
			u.is_active,
			u.created_at
		FROM 
			users u
		LEFT JOIN
			access_role ac ON ac.id = u.role_id
		WHERE 
    		u.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT 
			u.id,
			u.role_id,
			u.name, 
			u.email, 
			u.phone, 
			u.gender, 
			u.birth,
			u.is_active
		FROM 
			users u
		WHERE 
			u.id = $1 and u.deleted_at IS NULL
	`

	FIND_BY_EMAIL = `
		SELECT 
			u.id,
			u.role_id,
			u.name, 
			u.email, 
			u.phone, 
			u.gender, 
			u.birth,
			u.is_active
		FROM 
			users u
		WHERE 
			u.email = $1
	`

	FIND_USER_ROLE = `
		SELECT 
			u.id, 
			u.email, 
			ar."name" "role",
			u.password
		FROM 
			users u 
		JOIN 
			access_role ar on ar.id = u.role_id 
		WHERE 
			u.email = $1 and u.deleted_at IS NULL
	`

	INSERT = `
		INSERT INTO 
			users (
				id, 
				role_id,
				name,
				email, 
				password, 
				phone,
				gender,
				birth,
				is_active
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
			$9
		)
	`

	UPDATE_BY_ID = `
		UPDATE 
			users
		SET 
			role_id = $2,
			name = $3,
			email = $4,
			password = $5,
			phone = $6,
			gender = $7,
			birth = $8,
			is_active = $9,
			updated_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_ID = `
		UPDATE 
			users
		SET 
			deleted_at = NOW()
		WHERE 
			id = $1
	`
)
