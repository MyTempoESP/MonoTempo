package dbman

const (
	attach = `ATTACH DATABASE '/var/monotempo-data/equipamento.db' AS equip_data;`

	// args: HORA_LARGADA
	queryAtletas = `
SELECT DISTINCT
	athlete_num
FROM
	athletes_times;`
)
