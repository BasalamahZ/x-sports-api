package postgresql

const queryCreateGame = `
	INSERT INTO
		game
	(
		game_names,
		game_icons,
		create_time
	) VALUES (
		:game_names,
		:game_icons,
		:create_time
	) RETURNING
		id
`
