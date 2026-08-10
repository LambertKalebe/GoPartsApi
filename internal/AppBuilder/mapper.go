package appbuilder

import (
	"database/sql"
	"strings"
	"sync"
)

func toCarSearchResponse(rowsList []*sql.Rows) (carSearchResponse, error) {
	var resp carSearchResponse

	var wg sync.WaitGroup
	var mu sync.Mutex

	errCh := make(chan error, len(rowsList))

	for _, rows := range rowsList {
		if rows == nil {
			continue
		}

		wg.Add(1)

		go func(rows *sql.Rows) {
			defer wg.Done()
			defer rows.Close()

			var cars []car

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
					errCh <- err
					return
				}

				c.ConfigMotor = strings.TrimSpace(
					ccNormalized.String + " " + valves.String,
				)

				cars = append(cars, c)
			}

			if err := rows.Err(); err != nil {
				errCh <- err
				return
			}

			mu.Lock()
			resp.Cars = append(resp.Cars, cars...)
			mu.Unlock()

		}(rows)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			return carSearchResponse{}, err
		}
	}

	return resp, nil
}
