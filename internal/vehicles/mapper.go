package vehicles

import (
	"database/sql"
)

func toVehicleResponse(rows *sql.Rows, page int) (vehiclesResponse, error) {

	resp := vehiclesResponse{}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var v vehicle

		_ = rows.Scan(
			&v.ID,
			&v.Make,
			&v.Model,
			&v.Version,
			&v.EngineCode,
			&v.ConfigMotor,
			&v.TransmissionType,
			&v.Year,
			&v.PartsCount,
		)

		resp.Vehicle = append(resp.Vehicle, v)
	}
	resp.Page = page
	return resp, nil
}

func toVehicleByIdResponse(row *sql.Row) (vehicleByIdResponse, error) {
	resp := vehicleByIdResponse{}

	var v vehicle
	err := row.Scan(

		&v.ID,
		&v.Make,
		&v.Model,
		&v.Version,
		&v.EngineCode,
		&v.TransmissionType,
		&v.Year,
		&v.PartsCount,
	)
	if err != nil {
		return resp, err
	}

	resp.Vehicle = append(resp.Vehicle, v)
	return resp, nil
}

func toVehicleDetailsResponse(
	vehicleDetailsRow *sql.Row,
	vehicleAppRows *sql.Rows,
) (vehicleDetailsByIdResponse, error) {
	resp := vehicleDetailsByIdResponse{}
	var v vehicleDetails

	// Fecha as rows ao terminar
	defer func(vehicleAppRows *sql.Rows) {
		err := vehicleAppRows.Close()
		if err != nil {

		}
	}(vehicleAppRows)

	// Produtos aplicáveis
	for vehicleAppRows.Next() {
		var p parts

		if err := vehicleAppRows.Scan(
			&p.ID,
			&p.Name,
			&p.Make,
			&p.Code,
		); err != nil {
			return vehicleDetailsByIdResponse{}, err
		}

		v.Parts = append(v.Parts, p)
	}

	// Importante: Next() pode terminar por erro
	if err := vehicleAppRows.Err(); err != nil {
		return vehicleDetailsByIdResponse{}, err
	}

	// Dados do veículo
	if err := vehicleDetailsRow.Scan(
		&v.ID,
		&v.Make,
		&v.Model,
		&v.Version,
		&v.Year,
		&v.Engine.Code,
		&v.Engine.Valves,
		&v.Engine.CC,
		&v.Engine.Fuel,
		&v.Engine.Aspiration,
		&v.Transmission.Code,
		&v.Transmission.Name,
		&v.InfoData,
		&v.TechData,
		&v.SourceUrl,
	); err != nil {
		return vehicleDetailsByIdResponse{}, err
	}

	resp.Vehicle = append(resp.Vehicle, v)

	return resp, nil
}
