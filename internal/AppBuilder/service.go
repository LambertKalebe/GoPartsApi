package appbuilder

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Refatorar o código posteriormente, fui adicionando função corrigindo bugs.
// Ainda faltam contextos de "em diante", por exemplo: "Gol 2000...", "Gol 1990 em diante", etc
func serviceAppBuilder(search []string) ([]appBuilderSearchResponse, error) {
	expandedSearches := expanded(search)

	var wg sync.WaitGroup
	var mu sync.Mutex

	res := make([]appBuilderSearchResponse, 0, len(expandedSearches))

	for _, item := range expandedSearches {
		wg.Add(1)

		go func(item expandedSearch) {
			defer wg.Done()

			quoted := quoteTokens([]string{item.Search})[0]

			rows, err := appBuilderSearch(quoted)
			if err != nil {
				return
			}

			response, err := toAppBuilderSearchResponse(rows, item.Search)
			if err != nil {
				return
			}

			// Pós-processamento no response
			response = postProcessYear(response, item.Year)

			mu.Lock()
			res = append(res, response)
			mu.Unlock()

		}(item)
	}

	wg.Wait()

	return res, nil
}

func quoteTokens(search []string) []string {
	result := make([]string, len(search))

	for i, value := range search {
		tokens := strings.Fields(value)

		for j := range tokens {
			tokens[j] = `"` + tokens[j] + `"`
		}

		result[i] = strings.Join(tokens, " ")
	}

	fmt.Println("quoteTokens", result)

	return result
}

// Eu não aguento mais isso
var yearCandidateRegex = regexp.MustCompile(
	`(?:19|20)[0-9]{2}|[0-9]{2}`,
)

func expanded(search []string) []expandedSearch {
	var result []expandedSearch

	for _, query := range search {
		candidates := findYearCandidates(query)
		candidates = filterYearCandidates(query, candidates)

		switch len(candidates) {
		case 0:
			result = append(result, expandedSearch{
				Search: query,
			})

		case 1:
			result = append(
				result,
				expandSingleYear(query, candidates[0])...,
			)

		case 2:
			result = append(
				result,
				expandYearRange(
					query,
					candidates[0],
					candidates[1],
				)...,
			)

		default:
			// Ainda não temos regra para 3+ anos.
			result = append(result, expandedSearch{
				Search: query,
			})
		}
	}

	return result
}

// parseYear interpreta o PRIMEIRO ano do range (sem contexto anterior).
// Ano de 2 dígitos usa um pivô: 00-26 -> 20xx, 27-99 -> 19xx.
func parseYear(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	if len(s) == 2 {
		if n <= 26 {
			n += 2000
		} else {
			n += 1900
		}
	}

	if n < 1900 || n > 2026 {
		return 0, false
	}

	return n, true
}

// parseSecondYear interpreta o SEGUNDO ano do range, herdando o século
// do primeiro ano (start) quando vem com 2 dígitos. Isso resolve casos
// como "1976/78" -> 1978 (e não 2078) e "1995/02" -> 2002.
func parseSecondYear(s string, start int) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}

	if len(s) == 4 {
		if n < 1900 || n > 2026 {
			return 0, false
		}
		return n, true
	}

	// 2 dígitos: parte do século do ano inicial
	century := (start / 100) * 100
	end := century + n

	// Virada de século dentro do range (ex: 1995/02 -> 2002)
	if end < start {
		end += 100
	}

	if end > 2026 {
		return 0, false
	}

	return end, true
}

func postProcessYear(
	response appBuilderSearchResponse,
	expectedYear int,
) appBuilderSearchResponse {

	if expectedYear == 0 {
		return response
	}

	filtered := make([]car, 0, len(response.Cars))

	for _, vehicle := range response.Cars {
		if vehicle.Year == expectedYear {
			filtered = append(filtered, vehicle)
		}
	}

	response.Cars = filtered

	return response
}

func previousToken(s string, pos int) string {
	tokens := strings.Fields(strings.TrimSpace(s[:pos]))

	if len(tokens) == 0 {
		return ""
	}

	return tokens[len(tokens)-1]
}
func findYearCandidates(s string) []yearCandidate {
	matches := yearCandidateRegex.FindAllStringIndex(s, -1)

	result := make([]yearCandidate, 0, len(matches))

	for _, match := range matches {
		start := match[0]
		end := match[1]

		// Não aceitar número que faça parte de outro número.
		if start > 0 && isDigit(s[start-1]) {
			continue
		}

		if end < len(s) && isDigit(s[end]) {
			continue
		}

		result = append(result, yearCandidate{
			Start: start,
			End:   end,
			Value: s[start:end],
		})
	}

	return result
}
func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
func filterYearCandidates(
	query string,
	candidates []yearCandidate,
) []yearCandidate {

	result := make([]yearCandidate, 0, len(candidates))

	for _, candidate := range candidates {
		if isModelCandidate(query, candidate) {
			continue
		}

		result = append(result, candidate)
	}

	return result
}

func isModelCandidate(
	query string,
	candidate yearCandidate,
) bool {

	previous := previousToken(query, candidate.Start)

	return strings.EqualFold(previous, "peugeot") &&
		candidate.Value == "2008"
}

func expandSingleYear(
	s string,
	candidate yearCandidate,
) []expandedSearch {

	year, ok := parseYear(candidate.Value)

	if !ok {
		return []expandedSearch{
			{Search: s},
		}
	}

	const firstYear = 1970
	const lastYear = 2026

	var suffix = []string{"...", "/...", ">..."}
	var prefix = []string{"...", ".../", "...<"}

	// XXXX...
	for _, p := range prefix {
		if strings.HasPrefix(strings.TrimSpace(s[candidate.End:]), p) {
			return expandYears(
				s,
				candidate,
				year,
				lastYear,
			)
		}
	}

	// ...XXX
	for _, su := range suffix {
		if strings.HasSuffix(strings.TrimSpace(s[:candidate.Start]), su) {
			return expandYears(
				s,
				candidate,
				firstYear,
				year,
			)
		}
	}

	// Ano simples.
	return []expandedSearch{
		{
			Search:  s,
			Year:    year,
			HasYear: true,
		},
	}
}

func expandYears(
	query string,
	candidate yearCandidate,
	startYear int,
	endYear int,
) []expandedSearch {

	base := strings.TrimSpace(query[:candidate.Start])
	suffix := strings.TrimSpace(query[candidate.End:])

	base = strings.TrimSpace(strings.TrimSuffix(base, "..."))
	suffix = strings.TrimSpace(strings.TrimPrefix(suffix, "..."))

	result := make([]expandedSearch, 0, endYear-startYear+1)

	for year := startYear; year <= endYear; year++ {
		search := buildYearSearch(base, year, suffix)

		result = append(result, expandedSearch{
			Search:  search,
			Year:    year,
			HasYear: true,
		})
	}

	return result
}

func buildYearSearch(
	base string,
	year int,
	suffix string,
) string {

	query := strings.TrimSpace(
		fmt.Sprintf("%s %d", base, year),
	)

	if suffix != "" {
		query += " " + suffix
	}

	return query
}

func expandYearRange(
	s string,
	first yearCandidate,
	second yearCandidate,
) []expandedSearch {

	start, ok := parseYear(first.Value)
	if !ok {
		return []expandedSearch{{Search: s}}
	}

	end, ok := parseSecondYear(second.Value, start)
	if !ok || start > end {
		return []expandedSearch{{Search: s}}
	}

	// 00 01 sem separador continua ambíguo.
	if len(first.Value) == 2 &&
		len(second.Value) == 2 &&
		second.Start-first.End == 0 {

		return []expandedSearch{{Search: s}}
	}

	base := strings.TrimSpace(s[:first.Start])
	suffix := strings.TrimSpace(s[second.End:])

	result := make([]expandedSearch, 0, end-start+1)

	for year := start; year <= end; year++ {
		query := fmt.Sprintf("%s %d", base, year)

		if suffix != "" {
			query += " " + suffix
		}

		result = append(result, expandedSearch{
			Search:  query,
			Year:    year,
			HasYear: true,
		})
	}

	return result
}

//Sinceramente, eu queria muito que alguem visse esse codigo e me falasse que tem alguma maneira mais legivel e facil
