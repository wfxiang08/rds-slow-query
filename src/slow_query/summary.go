package slow_query

import (
	"fmt"
	log "github.com/wfxiang08/cyutils/utils/rolling_log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	numbersRegex = regexp.MustCompile(`(N, )+N`)
	stringsRegex = regexp.MustCompile(`('S', )+'S'`)

	countTimeRegex = regexp.MustCompile(`Count: (\d+)  Time=([.0-9]+)s`)
)

type Summary struct {
	Count int64
	Total float64
	Meta  string
	Sql   string
}

func NewSummary(meta string, sql string) *Summary {
	return &Summary{
		Meta: meta,
		Sql:  sql,
	}
}
func (s *Summary) HTML() string {
	return fmt.Sprintf("<span style='color:#000;display: inline-block;padding: 5px 10px 5px 10px;'>%s</span>", s.Meta) +
		"<br/>" +
		fmt.Sprintf("<span style='color:#6B6;display: inline-block;padding: 5px 10px 5px 25px;'>%s</span>", s.Sql)
}
func (s *Summary) String() string {
	return s.Meta + "\n" + s.Sql
}

func (s *Summary) Normalize() {
	// Count: 4  Time=1.12s
	items := countTimeRegex.FindStringSubmatch(s.Meta)

	// log.Printf("Meta: %s, Items: %s", s.Meta, strings.Join(items, ", "))

	count, err := strconv.ParseInt(items[1], 10, 64)
	if err != nil {
		log.ErrorErrorf(err, "parse int error: %s", s.Meta)
		return
	}
	s.Count = count

	time, err1 := strconv.ParseFloat(items[2], 64)
	if err1 != nil {
		log.ErrorErrorf(err1, "parse int error: %s", s.Meta)
		return
	}
	s.Total = time * float64(count)

}

type SummarySlice []*Summary

func (p SummarySlice) Len() int {
	return len(p)
}
func (p SummarySlice) Less(i, j int) bool {
	return p[i].Total > p[j].Total
}
func (p SummarySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func NormalizedSummary(input string, html bool) string {
	input = numbersRegex.ReplaceAllStringFunc(input, func(m string) string {
		parts := numbersRegex.FindStringSubmatch(m)
		if html {
			return fmt.Sprintf("<span style='color:#e100ff;'><%d>个N</span>", len(parts[0])/3+1)
		} else {
			return fmt.Sprintf("N..%d..N", len(parts[0])/3+1)
		}
	})

	input = stringsRegex.ReplaceAllStringFunc(input, func(m string) string {
		parts := stringsRegex.FindStringSubmatch(m)
		if html {
			return fmt.Sprintf("<span style='color:#e100ff;'><%d>个S</span>", len(parts[0])/5+1)
		} else {
			return fmt.Sprintf("S..%d..S", len(parts[0])/5+1)
		}
	})
	return input
}

func ParseSummaries(summary string, html bool) []*Summary {
	var results []*Summary
	// 以空行来分割
	summaryLines := strings.Split(summary, "\n\n")

	for i := 0; i < len(summaryLines); i++ {
		lines := strings.Split(summaryLines[i], "\n")
		if len(lines) >= 2 {
			sqlLines := lines[1:]
			for j := 0; j < len(sqlLines); j++ {
				sqlLines[j] = strings.TrimSpace(sqlLines[j])
			}

			meta := lines[0]
			sql := strings.Join(sqlLines, " ")

			if len(meta) > 0 && len(sql) > 0 {
				if strings.HasPrefix(sql, "INSERT") {
					continue
				}
				item := NewSummary(meta, NormalizedSummary(sql, html))
				item.Normalize()
				results = append(results, item)
			}
		}
	}

	sort.Sort(SummarySlice(results))
	return results
}
