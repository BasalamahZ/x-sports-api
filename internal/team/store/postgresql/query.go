package postgresql

const queryCreateTeam = `
	INSERT INTO
		team
	(
		team_names,
		team_icons,
		game_id,
		create_time
	) VALUES (
		:team_names,
		:team_icons,
		:game_id,
		:create_time
	)  RETURNING
		id
`

const queryGetTeams = `
	SELECT
		t.id,
		t.team_names,
		t.team_icons,
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
	%s
`

const queryUpdateTeam = `
	UPDATE
		team
	SET
		team_names = :team_names,
		team_icons = :team_icons,
		game_id = :game_id,
		create_time = :create_time,
		update_time = :update_time
	WHERE
		id = :id
`
