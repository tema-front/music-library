package utils

import (
	"errors"
	"net/http"
	"strconv"
)

func ParseLimitAndOffset(r *http.Request, defaultLimit, defaultOffset int) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit := defaultLimit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit <= limit {
			limit = parsedLimit
		}
	}
	
	offset := defaultOffset
	if offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil {
			offset = parsedOffset
		}
	}

	return limit, offset
}


func SafeSlice[T any](offset, limit int, slice []T) ([]T, error) {
	if limit < 0 || offset < 0 {
		return []T{}, errors.New("limit or offset can't be less than 0")
	}

	if offset >= len(slice) || len(slice) == 0 {
		return []T{}, nil
	}

	end := offset + limit
	if end > len(slice) {
		end = len(slice)
	}

	return slice[offset:end], nil
}


