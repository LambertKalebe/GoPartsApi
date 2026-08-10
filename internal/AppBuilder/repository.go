package appbuilder

import (
	"database/sql"
	"fmt"
	"g0/internal/database"
	"sync"
)

const appBuilderSearchQuery = `
SELECT
    v.id,
    m.name AS make,
    v.model,
    v.version,
    e.cc_normalized AS cc,
    e.valves AS valves,
    v.year
FROM fts5_vehicles f
JOIN vehicles v
    ON v.id = f.rowid
LEFT JOIN makes m
    ON m.id = v.make_id
LEFT JOIN engine_configs e
    ON e.id = v.engine_id
WHERE fts5_vehicles MATCH ?
`

// Split tokens in "token" "token2", etc
// Dessa forma a pesquisa de 1.0, 1.6 funciona
func appBuilderSearch(search []carSearchRequest) ([]*sql.Rows, error) {
	var wg sync.WaitGroup

	results := make([]*sql.Rows, len(search))
	errCh := make(chan error, len(search))

	for i := range search {
		i := i
		req := search[i]

		wg.Add(1)

		go func() {
			defer wg.Done()

			fmt.Println("search:", req)

			rows, err := database.DB.Query(
				appBuilderSearchQuery,
				req.Search,
			)
			if err != nil {
				errCh <- err
				return
			}

			results[i] = rows
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			// Se alguma query falhar, fecha as que já foram abertas.
			for _, rows := range results {
				if rows != nil {
					_ = rows.Close()
				}
			}

			return nil, err
		}
	}

	return results, nil
}
