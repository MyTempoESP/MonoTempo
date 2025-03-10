package dbman

const (
	ATTACH = `ATTACH DATABASE '/var/monotempo-data/equipamento.db' AS equip_data;`

	// args: HORA_LARGADA
	QUERY_LARGADA = `
SELECT
	athlete_num,
	antenna,
	staff,
	MAX(athlete_time)
FROM
	athletes_times
JOIN
	athletes ON athletes.num LIKE athlete_num
JOIN
	tracks ON tracks.id LIKE track_id
WHERE
	athlete_time <= largada AND
	athlete_time >= inicio
GROUP BY
	athlete_num;`

	// args: HORA_CHEGADA
	QUERY_CHEGADA = `
SELECT
	athlete_num,
	antenna,
	staff,
	MIN(athlete_time)
FROM
	athletes_times
JOIN
	athletes ON athletes.num LIKE athlete_num
JOIN
	tracks ON tracks.id LIKE track_id
WHERE
	athlete_time > largada AND
	athlete_time <= chegada
GROUP BY
	athlete_num;`
)
