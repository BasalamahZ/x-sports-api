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
	%s
`

const queryUpdateGame = `
	UPDATE
		game
	SET
		game_names = :game_names,
		game_icons = :game_icons,
		create_time = :create_time,
		update_time = :update_time
	WHERE
		id = :id
`
