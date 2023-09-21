package postgresql

const queryCreateMatch = `
	INSERT INTO
		match
	(
		tournament_names,
		game_id,
		team_a_id,
		team_b_id,
		team_a_odds,
		team_b_odds,
		date,
		match_link,
		status,
		create_time
	) VALUES (
		:tournament_names,
		:game_id,
		:team_a_id,
		:team_b_id,
		:team_a_odds,
		:team_b_odds,
		:date,
		:match_link,
		:status,
		:create_time
	)  RETURNING
		id
`

const queryGetMatchs = `
	SELECT
		m.id,
		m.tournament_names,
		m.game_id,
		g.game_names,
		g.game_icons,
		m.team_a_id,
		t1.team_names AS team_a_names,
		t1.team_icons AS team_a_icons,
		m.team_a_odds,
		m.team_b_id,
		t2.team_names AS team_b_names,
		t2.team_icons AS team_b_icons,
		m.team_b_odds,
		m.date,
		m.match_link,
		m.status,
		COALESCE(m.winner, 0) AS winner,
		m.create_time, 
		m.update_time
	FROM
		match m
	INNER JOIN
		game g
	ON
		m.game_id = g.id
	INNER JOIN
		team t1
	ON
		m.team_a_id = t1.id
	INNER JOIN
		team t2
	ON
		m.team_b_id = t2.id
	%s
`

const queryUpdateMatch = `
	UPDATE
		match
	SET
		tournament_names = :tournament_names,
		game_id = :game_id,
		team_a_id = :team_a_id,
		team_b_id = :team_b_id,
		team_a_odds = :team_a_odds,
		team_b_odds = :team_b_odds,
		date = :date,
		match_link = :match_link,
		status = :status,
		winner = :winner,
		create_time = :create_time,
		update_time = :update_time
	WHERE
		id = :id
`
