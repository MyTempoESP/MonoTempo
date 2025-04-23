package main

import (
	"database/sql"
)

const (
	CRIA_DB = `
PRAGMA synchronous = OFF;
PRAGMA journal_mode = MEMORY;
BEGIN TRANSACTION;

CREATE TABLE athletes
(
   num              INTEGER NOT NULL PRIMARY KEY
,  event_id         INTEGER NOT NULL
,  track_id         INTEGER NOT NULL
,  name             TEXT
,  city             TEXT
,  team             TEXT
,  sex              TEXT
);
CREATE TABLE equipamento
(
   id               INTEGER NOT NULL PRIMARY KEY
,  idequip          INTEGER NOT NULL
,  event_id         INTEGER NOT NULL
,  modelo           TEXT    NOT NULL
,  check            INTEGER NOT NULL
);
CREATE TABLE event_data
(
   id               INTEGER NOT NULL PRIMARY KEY
,  event_date       TEXT DEFAULT NULL
,  event_title      TEXT
);
CREATE TABLE tracks
(
   id               INTEGER NOT NULL PRIMARY KEY
,  event_id         INTEGER DEFAULT NULL
,  inicio           TEXT DEFAULT NULL
,  chegada          TEXT DEFAULT NULL
,  largada          TEXT DEFAULT NULL
,  race_description TEXT
);
CREATE TABLE staffs
(
   id               INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT
,  event_id         INTEGER NOT NULL
,  nome             TEXT DEFAULT NULL
); 

INSERT INTO equipamento(id, idequip, modelo, event_id, check) VALUES (1,0,'',0,0); 
END TRANSACTION;`

	QUERY_ATUALIZA_EQUIP = `
	REPLACE INTO equipamento (
	    id,
	    idequip,
	    modelo,
	    event_id,
	    check
	)
	VALUES(1, ?, ?, ?, ?);
	`

	QUERY_ATUALIZA_PROVA = `
	REPLACE INTO event_data (
		id,
		event_date,
		event_title
	)
	VALUES (?, ?, ?) 
	`

	QUERY_ATUALIZA_PERCURSO = `
	REPLACE INTO tracks (
	    id,
	    race_description,
	    event_id,
	    inicio,
	    chegada,
	    largada
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	QUERY_ATUALIZA_ATLETA = `
	REPLACE INTO athletes (
	    name,
	    sex,
	    team,
	    city,
	    num,
	    event_id,
	    track_id
	)
	VALUES(?, ?, ?, ?, ?, ?, ?);
	`

	QUERY_ATUALIZA_STAFF = `
	REPLACE INTO staffs (
	    id,
	    event_id,
	    nome
	)
	VALUES(?, ?, ?);
	`
)

func IgnorarForeignKey(db *sql.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
}

func AceitarForeignKey(db *sql.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
}
