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
	)  RETURNING
		id
`

const queryGetGames = `
	SELECT
		g.id,
		g.game_names,
		g.game_icons,
		g.create_time, 
		g.update_time
	FROM
		game g
`
