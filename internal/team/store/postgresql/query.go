package postgresql

const queryCreateTeam = `
	INSERT INTO
		team
	(
		team_names,
		game_id,
		create_time
	) VALUES (
		:team_names,
		:game_id,
		:create_time
	)  RETURNING
		id
`

const queryGetTeams = `
	SELECT
		t.id,
		t.team_names,
		t.game_id,
		g.game_names,
		g.game_icons,
		t.create_time, 
		t.update_time
	FROM
		team t
	INNER JOIN
		game g
	ON
		t.game_id = g.id
`
