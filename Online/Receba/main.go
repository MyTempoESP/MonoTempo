package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
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
}

func (r *Receba) Atualiza(logger *zap.Logger) {

	/*
		Ignora os checks de chave estrangeira
		do MySql.
	*/
	IgnorarForeignKey(r.db)

	r.ConfiguraAPI(os.Getenv("MYTEMPO_API_URL"))

	logger.Debug(fmt.Sprintf("%s;%s;%s;%s;%s", r.AtletasRota, r.DeviceRota, r.StaffRota, r.ProvaRota, r.InfoRota))

	logger = logger.With(zap.String("device", os.Getenv("MYTEMPO_EQUIP")))

	logger.Info("Buscando o equipamento")

	equip, err := r.BuscaEquip(os.Getenv("MYTEMPO_EQUIP"))

	if err != nil {

		logger.Error("Erro ao buscar o equipamento", zap.Error(err))

		return
	}

	equipLogger := logger.With(
		zap.String("equipamento", equip.Nome),
		zap.Int("equip_id", equip.ID),
		zap.Int("prova_id", equip.ProvaID),
		zap.Int("check_id", equip.Check),
	)

	equipLogger.Info("Equipamento encontrado, atualizando")

	err = r.AtualizaEquip(equip)

	if err != nil {

		equipLogger.Error("Erro ao atualizar o equipamento", zap.Error(err))

		return
	}

	equipLogger.Info("Buscando a prova")

	prova, err := r.BuscaProva(equip.ProvaID)

	if err != nil {

		equipLogger.Error("Erro ao buscar a prova", zap.Error(err))

		return
	}

	logger = logger.With(
		zap.Int("equip_id", equip.ID),
		zap.String("prova", prova.Nome),
		zap.Int("prova_id", prova.ID),
		zap.String("data_prova", prova.Data),
		zap.Int("percursos", len(prova.Percursos)),
	)

	logger.Info("Prova encontrada, atualizando")

	err = r.AtualizaProva(prova)

	if err != nil {

		logger.Error("Erro ao atualizar a prova", zap.Error(err))

		return
	}

	logger.Info("Prova atualizada, uscando os staffs")

	staff, err := r.BuscaStaff(prova.ID)

	if err != nil {

		logger.Error("Erro ao buscar os staffs", zap.Error(err))

		return
	}

	logger = logger.With(
		zap.Int("staffs", len(staff)))

	logger.Info("Staffs encontrados, atualizando")

	err = r.AtualizaStaff(staff, equip.ProvaID)

	if err != nil {

		logger.Error("Erro ao atualizar os staffs", zap.Error(err))

		return
	}

	logger.Info("Staffs atualizados")
}

func (r *Receba) AtualizarAtletas(logger *zap.Logger) {

	IgnorarForeignKey(r.db)

	r.ConfiguraAPI(os.Getenv("MYTEMPO_API_URL"))

	logger.Debug(fmt.Sprintf("%s;%s;%s;%s;%s", r.AtletasRota, r.DeviceRota, r.StaffRota, r.ProvaRota, r.InfoRota))

	logger = logger.With(zap.String("device", os.Getenv("MYTEMPO_EQUIP")))

	logger.Info("Buscando o equipamento")

	equip, err := r.BuscaEquip(os.Getenv("MYTEMPO_EQUIP"))

	if err != nil {

		logger.Error("Erro ao buscar o equipamento", zap.Error(err))

		return
	}

	logger = logger.With(
		zap.Int("equip_id", equip.ID),
		zap.Int("prova_id", equip.ProvaID),
	)

	logger.Info("Equipamento encontrado, buscando os atletas")

	atletas, err := r.BuscaAtletas(equip.ProvaID)

	if err != nil {

		logger.Error("Erro ao buscar os atletas", zap.Error(err))

		return
	}

	logger = logger.With(zap.Int("athletes_count", len(atletas)))

	logger.Info("Atletas encontrados, atualizando")

	err = r.AtualizaAtletas(atletas)

	if err != nil {

		logger.Error("Erro ao atualizar os atletas", zap.Error(err))

		return
	}

	logger.Info("Atletas atualizados")
}

func main() {

	//Limpar24h()

	logger, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	var r Receba

	err = r.ConfiguraDB()

	if err != nil {
		logger.Error("Erro ao configurar o banco de dados", zap.Error(err))
		return
	}

	defer r.FechaDB()

	r.Atualiza(logger)
	r.AtualizarAtletas(logger)
}
