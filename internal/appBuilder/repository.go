package appbuilder

import (
	"database/sql"
	"g0/internal/database"
)

const appBuilderSearchQuery = `
SELECT
    v.id,
    m.name AS make,
    v.model,
    v.version,
    COALESCE(printf('%.1f', e.cc_normalized), '') AS cc,
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

func appBuilderSearch(search string) (*sql.Rows, error) {

	rows, err := database.DB.Query(appBuilderSearchQuery, search)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
