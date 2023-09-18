package postgresql

const queryGetUserByEmail = `
	SELECT 
		id,
		email,
		password,
		create_time,
		update_time
	FROM
		admin
	WHERE
		email = $1
`
