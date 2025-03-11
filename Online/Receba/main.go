package main

import (
	"fmt"
	"log"
	"os"
)

func (r *Receba) FechaDB() {

	r.db.Close()
}

func (r *Receba) ConfiguraAPI(url string) {

	r.AtletasRota = fmt.Sprintf("http://%s/fetch/prova/atletas", url)
	r.DeviceRota = fmt.Sprintf("http://%s/fetch/device", url)
	r.StaffRota = fmt.Sprintf("http://%s/fetch/staffs", url)
	r.ProvaRota = fmt.Sprintf("http://%s/fetch/prova", url)
	r.InfoRota = fmt.Sprintf("http://%s/status/device", url)

	return
}

func (r *Receba) Atualiza() {

	/*
		Ignora os checks de chave estrangeira
		do MySql.
	*/
	IgnorarForeignKey(r.db)

	r.ConfiguraAPI(os.Getenv("MYTEMPO_API_URL"))

	/* debug */
	log.Println("Buscando equip:")

	equip, err := r.BuscaEquip(os.Getenv("MYTEMPO_EQUIP"))

	if err != nil {

		log.Println(err)

		return
	}

	/* debug */
	log.Println("Atualizando equip:")

	err = r.AtualizaEquip(equip)

	if err != nil {

		log.Println(err)

		return
	}

	/* debug */
	log.Println("Buscando a prova:")

	prova, err := r.BuscaProva(equip.ProvaID)

	if err != nil {

		log.Println(err)

		return
	}

	/* debug */
	log.Println("Atualizando a prova:")

	err = r.AtualizaProva(prova)

	if err != nil {

		log.Println(err)

		return
	}

	/*
		TODO: reportar exatamente os erros antes de enviar info

		exemplo:
		• Wi-Fi [X] OK
		• Equip [X] OK
		• Prova [X] OK
		• Staff [ ] ERRO
	*/

	// TODO: Report to website
	/*
		err = r.EnviarInfo(equip.ID)

		if err != nil {

			log.Println(err)

			return
		}
	*/

	/* debug */
	log.Println("Buscando staffs:")

	staff, err := r.BuscaStaff(prova.ID)

	if err != nil {

		log.Println(err)

		return
	}

	/* debug */
	log.Println("Atualizando staffs:")

	err = r.AtualizaStaff(staff, equip.ProvaID)

	if err != nil {

		log.Println(err)
	}

	return
}

func (r *Receba) AtualizarAtletas() {

	IgnorarForeignKey(r.db)

	r.ConfiguraAPI(os.Getenv("MYTEMPO_API_URL"))

	equip, err := r.BuscaEquip(os.Getenv("MYTEMPO_EQUIP"))

	if err != nil {

		log.Println(err)

		return
	}

	atletas, err := r.BuscaAtletas(equip.ProvaID)

	if err != nil {
		log.Println(err)
		return
	}

	err = r.AtualizaAtletas(atletas)

	if err != nil {
		log.Println(err)
	}
}

func main() {

	//Limpar24h()

	var r Receba

	err := r.ConfiguraDB()

	if err != nil {
		log.Println(err)
		return
	}

	defer r.FechaDB()

	r.Atualiza()
	r.AtualizarAtletas()
}
