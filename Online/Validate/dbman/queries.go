package dbman

const (
	ATTACH = `ATTACH DATABASE '/var/monotempo-data/equipamento.db' AS equip_data;`

	// args: HORA_LARGADA
	QUERY_ATLETAS = `
SELECT
	athlete_num
FROM
	athletes_times`
)
