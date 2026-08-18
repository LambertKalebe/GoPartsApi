package vehicles

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// VehiclesList
// @Summary List
// @Description Realiza uma consulta de veiculos por paginação
// @Tags Vehicles
// @Produce json
// @Param limit query int false "Quantidade máxima de veiculos por página" default(100) minimum(1) maximum(500)
// @Param page query int false "Página desejada"
// @Router /vehicles [get]
// @Success 200 {object} vehiclesResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func vehiclesHandler(c *echo.Context) error {
	req := new(vehiclesRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceVehicles(req.Limit, req.Page)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	err = c.JSON(200, res)
	return nil
}

// VehicleId
// @Summary ID
// @Description Realiza uma consulta basica de um veiculo com base no seu ID
// @Tags Vehicles
// @Produce json
// @Param id path int true "ID do veiculo"
// @Router /vehicles/{id} [get]
// @Success 200 {object} vehicleByIdResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func vehicleByIdHandler(c *echo.Context) error {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceVehicleById(id)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, res)
}

// VehicleDetails
// @Summary Details
// @Description Realiza uma consulta completa de um veiculo com base no seu ID
// @Tags Vehicles
// @Produce json
// @Param id path int true "ID do veiculo"
// @Router /vehicles/{id}/details [get]
// @Success 200 {object} vehicleDetailsByIdResponse
// @Failure 400 {object} httpcustom.ErrorResponse "Requisição inválida"
// @Failure 401 {object} httpcustom.ErrorResponse "Não autorizado"
// @Failure 404 {object} httpcustom.ErrorResponse "Recurso não encontrado"
// @Failure 500 {object} httpcustom.ErrorResponse "Erro interno do servidor"
func vehicleDetailsByIdHandler(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(
			http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceVehicleDetailsById(id)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
