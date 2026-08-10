package vehicles

import (
	"database/sql"
	"errors"
	"fmt"
)

func serviceVehicles(qnt int, page int) (vehiclesResponse, error) {
	rows, err := getVehicles(qnt, page)
	if err != nil {
		return vehiclesResponse{}, err
	}
	res, err := toVehicleResponse(rows, page)
	if err != nil {
		return vehiclesResponse{}, err
	}

	return res, nil

}

func serviceVehicleById(id int) (vehicleByIdResponse, error) {
	rows, err := getVehicleById(id)
	if err != nil {
		return vehicleByIdResponse{}, err
	}
	res, err := toVehicleByIdResponse(rows)
	if err != nil {
		return vehicleByIdResponse{}, err
	}
	return res, nil
}

func serviceVehicleDetailsById(id int) (vehicleDetailsByIdResponse, error) {

	if id <= 0 {
		return vehicleDetailsByIdResponse{}, errors.New("invalid query")
	}

	data, err := getVehicleDetailsById(id)
	if errors.Is(err, sql.ErrNoRows) {
		return vehicleDetailsByIdResponse{}, errors.New("invalid query")
	}
	fmt.Println("Service Data\n", data)

	apps, err := getVehicleAppsByProductId(id)
	if err != nil {
		return vehicleDetailsByIdResponse{}, err
	}
	fmt.Println("Service Apps\n", apps)

	res, err := toVehicleDetailsResponse(data, apps)
	if err != nil {
		fmt.Println("Service Error\n", err)
		return vehicleDetailsByIdResponse{}, err
	}
	fmt.Println("Service \n", res)
	return res, nil
}
