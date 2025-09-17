package utils

// CalculateOffsetAndLimit calculates the offset and limit for pagination.
func CalculateOffsetAndLimit(page, pageSize int) (offset, limit int) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset = (page - 1) * pageSize
	limit = pageSize
	return offset, limit
}
