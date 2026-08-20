package appbuilder

import (
	"database/sql"
	"strings"
)

func toAppBuilderSearchResponse(rows *sql.Rows, search string) (appBuilderSearchResponse, error) {
	var resp appBuilderSearchResponse

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	cars := make([]car, 0)
	carIDs := make([]int, 0)

	for rows.Next() {
		var c car
		var ccNormalized sql.NullString
		var valves sql.NullString

		err := rows.Scan(
			&c.ID,
			&c.Make,
			&c.Model,
			&c.Version,
			&ccNormalized,
			&valves,
			&c.Year,
			&c.FilterData,
		)
		if err != nil {
			return resp, err
		}

		c.ConfigMotor = strings.TrimSpace(
			ccNormalized.String + " " + valves.String,
		)

		cars = append(cars, c)
		carIDs = append(carIDs, c.ID)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}

	resp.Cars = cars
	resp.Search = search

	return resp, nil
}
