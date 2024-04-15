package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"github.com/xhd2015/xgo/runtime/mock"
)

func TestFetchPeople(t *testing.T) {
	t.Run("should return an error when url is empty", func(t *testing.T) {
		_, err := FetchPeople("")
		assert.Error(t, err, fmt.Errorf("url is empty"))
	})

	t.Run("should return an error when it fails", func(t *testing.T) {
		stub := func(req *fasthttp.Request, resp *fasthttp.Response) error {
			return fmt.Errorf("test-error")
		}

		mock.Patch(fasthttp.Do, stub)

		people, err := FetchPeople("http://test-url")
		assert.Nil(t, people)
		assert.Error(t, err, "test-error", fmt.Errorf("test-error"))
	})

	t.Run("should return a list of people", func(t *testing.T) {
		calledWithUrl := "initial value"
		responseBody := "[{\"id\": 1, \"name\": \"test-name\", \"email\": \"test-email\"}]"

		stub := func(req *fasthttp.Request, resp *fasthttp.Response) error {
			calledWithUrl = string(req.URI().FullURI())
			resp.SetBody([]byte(responseBody))
			return nil
		}

		mock.Patch(fasthttp.Do, stub)

		people, err := FetchPeople("http://test-url")
		assert.Nil(t, err)
		assert.Equal(t, "http://test-url/", calledWithUrl)
		assert.Equal(t, 1, len(people))
		assert.Equal(t, 1, people[0].ID)
		assert.Equal(t, "test-name", people[0].Name)
		assert.Equal(t, "test-email", people[0].Email)
	})

	t.Run("should return a list of people (un-mocked)", func(t *testing.T) {
		people, err := FetchPeople("https://jsonplaceholder.typicode.com/users")
		assert.Nil(t, err)
		assert.Equal(t, 10, len(people))
	})
}

func TestFetchPeopleAsync(t *testing.T) {
	t.Run("should send error trough a channel when failing", func(t *testing.T) {
		_, err := (<-FetchPeopleAsync("")).Unwrap()
		assert.Error(t, err, "url is empty", fmt.Errorf("url is empty"))
	})

	t.Run("should send results trough a channel when succeeding", func(t *testing.T) {
		calledWithUrl := "initial value"
		responseBody := "[{\"id\": 1, \"name\": \"test-name\", \"email\": \"test-email\"}]"

		stub := func(req *fasthttp.Request, resp *fasthttp.Response) error {
			calledWithUrl = string(req.URI().FullURI())
			resp.SetBody([]byte(responseBody))
			return nil
		}

		fetchPeopleOptions := FetchPeopleOptions{
			FastHTTPDo: stub,
		}

		people, err := (<-FetchPeopleAsync("http://test-url", fetchPeopleOptions)).Unwrap()
		assert.Nil(t, err)
		assert.Equal(t, "http://test-url/", calledWithUrl)
		assert.Equal(t, 1, len(people))
		assert.Equal(t, 1, people[0].ID)
		assert.Equal(t, "test-name", people[0].Name)
		assert.Equal(t, "test-email", people[0].Email)
	})
}
