package postgresql

const queryCreateThread = `
	INSERT INTO
		thread
	(
		title,
		game_id,
		description,
		image_thread,
		date,
		create_time
	) VALUES (
		:title,
		:game_id,
		:description,
		:image_thread,
		:date,
		:create_time
	)  RETURNING
		id
`

const queryGetThreads = `
	SELECT
		t.id,
		t.title,
		t.game_id,
		g.game_names,
		g.game_icons,
		t.description,
		t.image_thread,
		t.date,
		t.create_time, 
		t.update_time
	FROM
		thread t
	INNER JOIN
		game g
	ON
		t.game_id = g.id
	%s
`

const queryUpdateThread = `
	UPDATE
		thread
	SET
		title = :title,
		game_id = :game_id,
		description = :description,
		image_thread = :image_thread,
		date = :date,
		create_time = :create_time,
		update_time = :update_time
	WHERE
		id = :id
`
