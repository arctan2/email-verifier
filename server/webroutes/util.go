package webroutes

import (
	"errors"
	"net/http"
	"strconv"
)

func parseInt64PathValue(p string, r *http.Request) (int64, error) {
	fileIdStr := r.PathValue(p)

	if fileIdStr == "" {
		return 0, errors.New(p + " not provided.")
	}

	fileId, err := strconv.ParseInt(fileIdStr, 10, 64)

	if err != nil {
		return 0, err
	}

	return fileId, nil
}

func parseInt64QueryValue(p string, r *http.Request) (int64, error) {
	fileIdStr := r.URL.Query().Get(p)

	if fileIdStr == "" {
		return 0, errors.New(p + " not provided.")
	}

	fileId, err := strconv.ParseInt(fileIdStr, 10, 64)

	if err != nil {
		return 0, err
	}

	return fileId, nil
}
