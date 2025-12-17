package atleta

type Atleta struct {
	Tempo  string `json:"tempo"`
	Antena int    `json:"antena"`
	Numero int    `json:"numero"`

	ProvaID    int
	PercursoID int `json:"percurso"`
}
