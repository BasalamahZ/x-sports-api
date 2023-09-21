package postgresql

const queryCreateNews = `
	INSERT INTO
		news
	(
		title,
		game_id,
		description,
		image_news,
		date,
		create_time
	) VALUES (
		:title,
		:game_id,
		:description,
		:image_news,
		:date,
		:create_time
	)  RETURNING
		id
`

const queryGetNews = `
	SELECT
		n.id,
		n.title,
		n.game_id,
		g.game_names,
		g.game_icons,
		n.description,
		n.image_news,
		n.date,
		n.create_time, 
		n.update_time
	FROM
		news n
	INNER JOIN
		game g
	ON
		n.game_id = g.id
	%s
`

const queryUpdateNews = `
	UPDATE
		news
	SET
		title = :title,
		game_id = :game_id,
		description = :description,
		image_news = :image_news,
		date = :date,
		create_time = :create_time,
		update_time = :update_time
	WHERE
		id = :id
`
