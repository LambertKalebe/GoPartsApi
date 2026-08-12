package appbuilder

import (
	"fmt"
	"g0/internal/global"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// Refatorar o código posteriormente, fui adicionando função corrigindo bugs.
// Ainda faltam outros contextos de anos.
func serviceAppBuilder(search []string) ([]appBuilderSearchResponse, error) {
	expandedSearches := expanded(search)

	var wg sync.WaitGroup
	var mu sync.Mutex

	res := make([]appBuilderSearchResponse, 0, len(expandedSearches))

	for _, item := range expandedSearches {
		wg.Add(1)

		go func(item expandedSearch) {
			defer wg.Done()

			quoted := global.QuoteTokens([]string{item.Search})[0]

			rows, err := appBuilderSearch(quoted)
			if err != nil {
				return
			}

			response, err := toAppBuilderSearchResponse(rows, item.Search)
			if err != nil {
				return
			}

			response = postProcessYear(response, item.Year)

			mu.Lock()
			res = append(res, response)
			mu.Unlock()

		}(item)
	}

	wg.Wait()

	return res, nil
}

// EXPANSÃO

var yearCandidateRegex = regexp.MustCompile(
	`(?:19|20)[0-9]{2}|[0-9]{2}`,
)

func expanded(search []string) []expandedSearch {
	var result []expandedSearch

	for _, query := range search {
		r := strings.NewReplacer(
			" de ", " ",
			" da ", " ",
			" do ", " ",
			" até ", " ",
			" ate ", " ",
			" TODOS ", " ",
		)
		query = r.Replace(query)
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
			result = append(result, expandedSearch{
				Search: query,
			})
		}
	}

	return result
}

// PARSE DOS ANOS

// parseYear interpreta o primeiro ano.
//
// 00-26 -> 2000-2026
// 27-99 -> 1927-1999
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

// parseSecondYear interpreta o segundo ano,
// herdando o século do primeiro.
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

	century := (start / 100) * 100
	end := century + n

	// Ex:
	// 1995/02 -> 2002
	if end < start {
		end += 100
	}

	if end > 2026 {
		return 0, false
	}

	return end, true
}

// PÓS PROCESSAMENTO

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

// TOKEN / CANDIDATOS

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

// ANO ÚNICO

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

	// Texto antes do ano.
	base := strings.TrimSpace(
		s[:candidate.Start],
	)

	// Texto depois do ano.
	after := strings.TrimSpace(
		s[candidate.End:],
	)

	// 1. "2009 em diante"

	if after, ok := cutPrefixFold(after, "em diante"); ok {
		return expandYears(
			base,
			after,
			year,
			lastYear,
		)
	}

	// 2. "2009..."

	if after, ok := strings.CutPrefix(after, "..."); ok {
		after = strings.TrimSpace(after)

		return expandYears(
			base,
			after,
			year,
			lastYear,
		)
	}

	// 3. "2009/..."

	if after, ok := strings.CutPrefix(after, "/..."); ok {
		after = strings.TrimSpace(after)

		return expandYears(
			base,
			after,
			year,
			lastYear,
		)
	}

	// 4. "2009 >..."

	if after, ok := strings.CutPrefix(after, ">..."); ok {
		after = strings.TrimSpace(after)

		return expandYears(
			base,
			after,
			year+1,
			lastYear,
		)
	}

	// 5. "até 2009"

	if cleanedBase, ok := cutSuffixFold(base, "até"); ok {
		return expandYears(
			cleanedBase,
			after,
			firstYear,
			year,
		)
	}

	// 6. "...2009"

	if cleanedBase, ok := strings.CutSuffix(base, "..."); ok {
		cleanedBase = strings.TrimSpace(cleanedBase)

		return expandYears(
			cleanedBase,
			after,
			firstYear,
			year,
		)
	}

	// 7. ".../2009"

	if cleanedBase, ok := strings.CutSuffix(base, ".../"); ok {
		cleanedBase = strings.TrimSpace(cleanedBase)

		return expandYears(
			cleanedBase,
			after,
			firstYear,
			year,
		)
	}

	// 8. "após 2009"

	if cleanedBase, ok := cutSuffixFold(base, "após"); ok {
		return expandYears(
			cleanedBase,
			after,
			year+1,
			lastYear,
		)
	}

	// 9. "a partir de 2009"

	if cleanedBase, ok := cutSuffixFold(base, "a partir de"); ok {
		return expandYears(
			cleanedBase,
			after,
			year,
			lastYear,
		)
	}

	// 10. "2009" simples

	return []expandedSearch{
		{
			Search:  s,
			Year:    year,
			HasYear: true,
		},
	}
}

// EXPAND YEARS

func expandYears(
	base string,
	suffix string,
	startYear int,
	endYear int,
) []expandedSearch {

	base = strings.TrimSpace(base)
	suffix = strings.TrimSpace(suffix)

	if startYear > endYear {
		return nil
	}

	result := make([]expandedSearch, 0, endYear-startYear+1)

	for year := startYear; year <= endYear; year++ {

		search := buildYearSearch(
			base,
			year,
			suffix,
		)

		result = append(result, expandedSearch{
			Search:  search,
			Year:    year,
			HasYear: true,
		})
	}

	return result
}

// BUILD SEARCH

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

// RANGE: 2009-2015 / 2009/15 / etc.

func expandYearRange(
	s string,
	first yearCandidate,
	second yearCandidate,
) []expandedSearch {

	start, ok := parseYear(first.Value)

	if !ok {
		return []expandedSearch{
			{Search: s},
		}
	}

	end, ok := parseSecondYear(
		second.Value,
		start,
	)

	if !ok || start > end {
		return []expandedSearch{
			{Search: s},
		}
	}

	// 00 01 sem separador continua ambíguo.
	if len(first.Value) == 2 &&
		len(second.Value) == 2 &&
		second.Start-first.End == 0 {

		return []expandedSearch{
			{Search: s},
		}
	}

	base := strings.TrimSpace(
		s[:first.Start],
	)

	suffix := strings.TrimSpace(
		s[second.End:],
	)

	return expandYears(
		base,
		suffix,
		start,
		end,
	)
}

// HELPERS CASE-INSENSITIVE

func cutPrefixFold(
	s string,
	prefix string,
) (string, bool) {

	if len(s) < len(prefix) {
		return s, false
	}

	if !strings.EqualFold(
		s[:len(prefix)],
		prefix,
	) {
		return s, false
	}

	return strings.TrimSpace(
		s[len(prefix):],
	), true
}

func cutSuffixFold(
	s string,
	suffix string,
) (string, bool) {

	if len(s) < len(suffix) {
		return s, false
	}

	start := len(s) - len(suffix)

	if !strings.EqualFold(
		s[start:],
		suffix,
	) {
		return s, false
	}

	return strings.TrimSpace(
		s[:start],
	), true
}

// Acho que foi o unico trecho que recorri a IA, mas tambem, odeio regex
