# Reenvio
Sistema de reenvio automático de atletas baseado em mudanças no ambiente da prova

```

const (
	/*
		Intervalo entre requests para a API.
	*/
	LOTE_TIMEOUT time.Duration = 5 * time.Second
	ENVIAR_TEMPOS_URL string = "http://%s/receive/tempos"
)

		err = envio.AtualizaEquip()

		if err != nil {
			log.Println("(EnviarAtletasLote) Erro ao atualizar dados do equipamento", err)

			continue
		}

		dados := AtletasForm{
			envio.Equip.ID,
			envio.Equip.ProvaID,
			atletas,
		}

		/* NOTE: DEBUG */
		log.Printf("Dados para a request: %+v\n", dados)

		atletasJson, err := json.Marshal(dados)

		if err != nil {
			log.Println("(EnviarAtletasLote) Erro ao converter os dados de atleta para json", err)

			continue
		}

		var url string = os.Getenv("MYTEMPO_API_URL")
		atletasRota := fmt.Sprintf(ENVIAR_TEMPOS_URL, url)

		var res RespostaAPI

		err = JSONRequest(
			atletasRota,
			atletasJson,
			&res,
		)

		if err != nil {
			log.Printf("(EnviarAtletasLote) Erro ao efetuar a request para %s\n", atletasRota)
			log.Println(err)

			envio.NaoEnviados <- atletas

			continue
		}

		if res.Status == "error" {
			log.Printf("(EnviarAtletasLote) Recebido um status de erro da API: %s\n", res.Message)

			/*
				NOTE: Talvez não seja uma boa ideia tentar um reenvio caso
				      a API retorne uma mensagem de erro, é uma boa repensar.
			*/

			envio.NaoEnviados <- atletas
		} 
```
