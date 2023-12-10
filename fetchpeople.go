package main

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type FetchResult struct {
	People []Person
	Err    error
}

func (ar FetchResult) Unwrap() ([]Person, error) {
	return ar.People, ar.Err
}

func AsyncFetchPeople(url string) <-chan FetchResult {
	done := make(chan FetchResult)

	go func() {
		data, err := FetchPeople(url)
		done <- FetchResult{data, err}
	}()

	return done
}

func FetchPeople(url string) ([]Person, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty")
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var people []Person
	if err := json.Unmarshal(resp.Body(), &people); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %s", err)
	}

	return people, nil
}
