package dbman

const (
	attach = `ATTACH DATABASE '/var/monotempo-data/equipamento.db' AS equip_data;`

	// args: HORA_LARGADA
	queryLargada = `
SELECT
	athlete_num,
	antenna,
	track_id,
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
	queryChegada = `
SELECT
	athlete_num,
	antenna,
	track_id,
	MIN(athlete_time)
FROM
	athletes_times
JOIN
	athletes ON athletes.num LIKE athlete_num
JOIN
	tracks ON tracks.id LIKE track_id
WHERE
	athlete_time >= chegada
GROUP BY
	athlete_num;`

	queryCheckpoint = `
SELECT
	athlete_num,
	antenna,
	track_id,
	MAX(athlete_time)
FROM
	athletes_times
JOIN
	athletes ON athletes.num LIKE athlete_num
JOIN
	tracks ON tracks.id LIKE track_id
WHERE
	athlete_time > largada
GROUP BY
	athlete_num;`
)
