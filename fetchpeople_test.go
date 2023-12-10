package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestFetchPeople(t *testing.T) {
	gm := gomonkey.NewPatches()

	t.Run("should return an error when url is empty", func(t *testing.T) {
		_, err := FetchPeople("")
		assert.Error(t, err)
	})

	t.Run("should return an error when it fails", func(t *testing.T) {
		stub := func(_ *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
			return fmt.Errorf("error")
		}

		gm.ApplyMethod(reflect.TypeOf(&fasthttp.Client{}), "Do", stub)
		defer gm.Reset()

		result, err := FetchPeople("http://test-url")
		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("should return a list of people", func(t *testing.T) {
		calledWithUrl := "initial value"

		stub := func(_ *fasthttp.Client, req *fasthttp.Request, resp *fasthttp.Response) error {
			calledWithUrl = string(req.URI().FullURI())
			body := "[{\"id\": 1, \"name\": \"test-name\", \"email\": \"test-email\"}]"
			resp.SetBody([]byte(body))
			return nil
		}

		gm.ApplyMethod(reflect.TypeOf(&fasthttp.Client{}), "Do", stub)
		defer gm.Reset()

		result, err := FetchPeople("http://test-url")
		assert.Nil(t, err)
		assert.Equal(t, "http://test-url/", calledWithUrl)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, 1, result[0].ID)
		assert.Equal(t, "test-name", result[0].Name)
		assert.Equal(t, "test-email", result[0].Email)
	})

	t.Run("should return a list of people (un-mocked)", func(t *testing.T) {
		result, err := FetchPeople("https://jsonplaceholder.typicode.com/users")
		assert.Nil(t, err)
		assert.Equal(t, 10, len(result))
	})
}