package paginator

import "math"

type Page struct {
	Start int64
	End   int64
}

// CalculatePages 根据总数量和每页多少个计算Page数组
func CalculatePages(total, pageSize int64) []Page {
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))

	pages := make([]Page, totalPage)
	for i := int64(0); i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		pages[i] = Page{
			Start: start,
			End:   end,
		}
	}
	return pages
}
