package vehicles

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
)

// VehiclesList
// @Summary List
// @Description Realiza uma consulta de veiculos por paginação
// @Tags Vehicles
// @Produce json
// @Param qnt query int false "Quantidade máxima de veiculos por página"
// @Param page query int false "Página desejada"
// @Router /vehicles [get]
// @Success 200 {object} vehiclesResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func vehiclesHandler(c *echo.Context) error {
	req := new(vehiclesRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Corpo da requisição inválido")
	}

	res, err := serviceVehicles(req.Qnt, req.Page)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}
	err = c.JSON(200, res)
	return nil
}

// VehicleId
// @Summary ID
// @Description Realiza uma consulta basica de um veiculo com base no seu ID
// @Tags Vehicles
// @Produce json
// @Param int query int false "ID do veiculo"
// @Router /vehicles/{id} [get]
// @Success 200 {object} vehicleByIdResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func vehicleByIdHandler(c *echo.Context) error {
	var id, err = strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceVehicleById(id)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}
	return c.JSON(http.StatusOK, res)
}

// VehicleDetails
// @Summary Details
// @Description Realiza uma consulta completa de um veiculo com base no seu ID
// @Tags Vehicles
// @Produce json
// @Param int query int false "ID do veiculo"
// @Router /vehicles/{id}/details [get]
// @Success 200 {object} vehicleDetailsByIdResponse
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
func vehicleDetailsByIdHandler(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Requisição inválido")
	}
	res, err := serviceVehicleDetailsById(id)
	if err != nil {
		fmt.Println("Handler", err)
		err = c.String(http.StatusInternalServerError, err.Error())
		return err
	}

	return c.JSON(http.StatusOK, res)
}
