package role

const (
	FIND_ALL = `
		SELECT 
			u.id,
			u.name,
			u.description
		FROM 
			access_role u
		WHERE 
    		u.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT 
			u.id,
			u.name,
			u.description
		FROM 
			access_role u
		WHERE 
			u.id = $1;
	`

	INSERT = `
		INSERT INTO 
			access_role (
				id, 
				name,
				description
			) 
		VALUES (
			$1, 
			$2, 
			$3
		)
	`

	UPDATE_BY_ID = `
		UPDATE 
			access_role
		SET 
			name = $2,
			description = $3,
			updated_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_ID = `
		UPDATE 
			access_role
		SET 
			deleted_at = NOW()
		WHERE 
			id = $1
	`
)
