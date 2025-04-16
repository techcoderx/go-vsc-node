package gqlgen

import "fmt"

// Helper function for handling offset-limit pagination. Returns offset, limit and error.
func Paginate(offset *int, limit *int) (int, int, error) {
	var limit_result int
	var offset_result int
	if limit == nil {
		limit_result = 50
	} else if *limit < 1 || *limit > 100 {
		return 0, 0, fmt.Errorf("limit must be between 1 and 100")
	} else {
		limit_result = *limit
	}
	if offset == nil {
		offset_result = 0
	} else if *offset < 0 || *offset > 10000 {
		return 0, 0, fmt.Errorf("offset must be between 0 and 10000")
	} else {
		offset_result = *offset
	}
	return offset_result, limit_result, nil
}
