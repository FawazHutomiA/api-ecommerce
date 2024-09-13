package product

const (
	FIND_ALL = `
		SELECT 
			r.id, 
			r.user_id,
			r.name,
			r.short_description,
			r.description,
			r.goal_amount,
			r.current_amount,
			r.slug,
			r.backer_amount
		FROM 
			product r
		WHERE 
			r.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT  
			r.id, 
			r.user_id,
			r.name,
			r.short_description,
			r.description,
			r.goal_amount,
			r.current_amount,
			r.slug,
			r.backer_amount 
		FROM 
			product r
		where r.id = $1
	`

	FIND_BY_NAME = `
		SELECT  
			r.id, 
			r.user_id,
			r.name,
			r.short_description,
			r.description,
			r.goal_amount,
			r.current_amount,
			r.slug,
			r.backer_amount 
		FROM 
			product r
		where r.name = $1
	`

	FIND_BY_SLUG = `
		SELECT 
			r.id, 
			r.user_id,
			r.name,
			r.short_description,
			r.description,
			r.goal_amount,
			r.current_amount,
			r.slug,
			r.backer_amount 
		FROM 
			product r
		where 
			r.slug = $1
	`

	INSERT = `
		INSERT INTO product (
			id, 
			user_id, 
			name, 
			short_description, 
			description, 
			goal_amount, 
			current_amount, 
			slug, 
			backer_amount
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
		UPDATE product
		SET 
			user_id = $2,
			name = $3,
			short_description = $4,
			description = $5,
			goal_amount = $6,
			current_amount = $7,
			slug = $8,
			backer_amount = $9,
			updated_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_ID = `
		UPDATE 
			product
		SET 
			deleted_at = NOW()
		WHERE 
			id = $1
	`
)
