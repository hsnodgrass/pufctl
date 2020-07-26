package forgeapi

import (
	"fmt"
	"time"
)

func prefix() string {
	pkgname := "forgeapi"
	_time := time.Now()
	return fmt.Sprintf("%s: %s:", pkgname, _time)
}

// GetError provides a wrapper for errors encountered
// during HTTP GET request creation or execution
type GetError struct {
	Err     error
	URL     string
	Headers map[string]string
}

func (r *GetError) Error() string {
	return fmt.Sprintf("%s HTTP GET request failed. URL: %s, Headers: %#v, Error: %#v", prefix(), r.URL, r.Headers, r.Err)
}

// GetNon200Error provides an error wrapper for when
// HTTP GET requests return a non-200 status code
type GetNon200Error struct {
	URL        string
	StatusCode int
}

func (r *GetNon200Error) Error() string {
	return fmt.Sprintf("%s HTTP GET request returned non-200 code. URL: %s, Code: %d", prefix(), r.URL, r.StatusCode)
}

// JSONDecodeError provides wrapper for errors encountered
// while decoding JSON response bodies.
type JSONDecodeError struct {
	Err error
}

func (r *JSONDecodeError) Error() string {
	return fmt.Sprintf("%s Failed to decode HTTP response body: %#v", prefix(), r.Err)
}

// FetchError provides a wrapper for errors encountered
// during Fetch requests
type FetchError struct {
	Err error
}

func (r *FetchError) Error() string {
	return fmt.Sprintf("%s Fetch failed: %#v", prefix(), r.Err)
}

// ListError provides a wrapper for errors encountered
// during List requests
type ListError struct {
	Err error
}

func (r *ListError) Error() string {
	return fmt.Sprintf("%s List failed: %#v", prefix(), r.Err)
}
