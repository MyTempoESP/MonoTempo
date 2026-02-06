package dbman

const (
	CREATE_TIME_TABLE = `
CREATE TABLE IF NOT EXISTS athletes_times (
   id            INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT
,  antenna       INTEGER NOT NULL
,  athlete_num   INTEGER NOT NULL
,  staff         INTEGER NOT NULL
,  athlete_time  TEXT
);`

	WRITE_PRAGMAS = `
PRAGMA synchronous = OFF;
PRAGMA journal_mode = MEMORY;`

	INSERT_TIME = `
INSERT INTO athletes_times
	(antenna, athlete_num, staff, athlete_time) VALUES (?, ?, ?, ?)`
)
