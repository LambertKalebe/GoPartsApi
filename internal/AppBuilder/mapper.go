package appbuilder

import (
	"database/sql"
	"fmt"
	"strings"
)

func toAppBuilderSearchResponse(rows *sql.Rows, search string) (appBuilderSearchResponse, error) {
	var resp appBuilderSearchResponse

	defer rows.Close()

	cars := make([]car, 0)

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
		)
		if err != nil {
			return resp, err
		}

		c.ConfigMotor = strings.TrimSpace(
			ccNormalized.String + " " + valves.String,
		)

		cars = append(cars, c)
	}

	if err := rows.Err(); err != nil {
		return resp, err
	}

	resp.Cars = cars
	resp.Search = search

	fmt.Println("resp:", resp)

	return resp, nil
}
