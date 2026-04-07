package paginator

import (
	"math"

	"github.com/hdget/sdk/common/protobuf"
	"golang.org/x/exp/constraints"
)

// Paginate 分页
func Paginate[T any](sliceVars []T, pageSize int64) [][]T {
	total := int64(len(sliceVars))
	totalPage := int64(math.Ceil(float64(total) / float64(pageSize)))

	results := make([][]T, totalPage)
	for i := int64(0); i < totalPage; i++ {
		start := i * pageSize
		end := (i + 1) * pageSize
		if end > total {
			end = total
		}

		results[i] = sliceVars[start:end]
	}
	return results
}

// GetLimitFromListParam 从protobuf.ListParam转换成Offset, PageSize
func GetLimitFromListParam[T constraints.Integer](list ...*protobuf.ListParam) (T, T) {
	p := NewFromListParam(list...)
	return T(p.Offset), T(p.PageSize)
}

// GetStartEndPosition 如果是按列表slice进行翻页的话， 计算slice的起始index
func GetStartEndPosition(page, pageSize, total int64) (int64, int64) {
	// 处理当前页面
	if page <= 0 {
		page = 1
	}

	// 处理每页大小
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	start := (page - 1) * pageSize
	end := page * pageSize

	if end > total {
		end = total
	}

	if start > end {
		start = end
	}

	return start, end
}
