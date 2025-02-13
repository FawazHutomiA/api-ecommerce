package warehouse

const (
	FIND_ALL = `
		SELECT 
			w.id,
			w.name,
			w.address,
			w.province,
			w.city,
			w.zip_code
		FROM 
			warehouse w
		WHERE 
    		w.deleted_at IS NULL
	`

	FIND_BY_ID = `
		SELECT 
			w.id,
			w.name,
			w.address,
			w.province,
			w.city,
			w.zip_code
		FROM 
			warehouse w
		WHERE 
			w.id = $1;
	`

	INSERT = `
		INSERT INTO 
			warehouse (
				id,
				name,
				address,
				province,
				city,
				zip_code
			) 
		VALUES (
			$1, 
			$2, 
			$3,
			$4,
			$5,
			$6
		)
	`

	UPDATE_BY_ID = `
		UPDATE 
			warehouse
		SET 
			id = $1,
			name = $2,
			address = $3,
			province = $4,
			city = $5
			zip_code = $6,
			updated_at = NOW()
		WHERE 
			id = $1
	`

	DELETE_BY_ID = `
		UPDATE 
			warehouse
		SET 
			deleted_at = NOW()
		WHERE 
			id = $1
	`
)
