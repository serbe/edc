SELECT
		e.id,
		e.start_date,
		c.name
	FROM
		educations AS e
	LEFT JOIN
		contacts AS c ON c.id = e.contact_id
	WHERE
		e.start_date > TIMESTAMP 'now'::timestamp - '1 month'::interval
	ORDER BY
		start_date DESC
	LIMIT 10