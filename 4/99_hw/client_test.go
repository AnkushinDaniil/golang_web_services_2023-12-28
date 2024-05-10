package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/test-go/testify/assert"
)

type TestCase struct {
	Query      string
	OrderField string
	OrderBy    int
	Limit      int
	Token      string
	Offset     int
	Output     string
	Error      error
}

type TestErrorCase struct {
	Query string
	URL   string
	Error error
}

func SearchServerDummy(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	switch query {
	case "bad url":
		w.WriteHeader(http.StatusRequestTimeout)
		io.WriteString(w, `{"status": 200, "balance": 100500}`)
	case "timeout":
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 200, "balance": 100500}`)
	case "StatusInternalServerError":
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"status": 500, "error": "error"}`)
	case "StatusBadRequestWithBadJson":
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"status": 500, nil}}}}`)
	case "ErrorBadOrderField":
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"status": 500, "error": "ErrorBadOrderField"}`)
	case "unknown bad request error":
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"status": 500, "error": "unknown"}`)
	case "__broken_json":
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status": 400`) // broken json
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func TestSearchServer(t *testing.T) {
	cases := []TestCase{
		{
			Query:      "bo",
			OrderField: "Age",
			OrderBy:    -1,
			Limit:      5,
			Token:      token,
			Offset:     0,
			Error:      nil,
		},
		{
			Query:      "Hilda",
			OrderField: "Age",
			OrderBy:    -1,
			Limit:      5,
			Token:      token,
			Offset:     0,
			Error:      nil,
		},
		{
			Query:      "",
			OrderField: "Id",
			OrderBy:    -1,
			Limit:      5,
			Token:      token,
			Offset:     0,
			Error:      nil,
		},
		{
			Query:      "ab",
			OrderField: "Name",
			OrderBy:    -1,
			Limit:      5,
			Token:      token,
			Offset:     0,
			Error:      nil,
		},
		{
			Query:      "",
			OrderField: "Age",
			OrderBy:    -1,
			Limit:      30,
			Token:      token,
			Offset:     0,
			Error:      nil,
		},
		{
			Query:      "",
			OrderField: "Age",
			OrderBy:    -1,
			Limit:      3,
			Token:      "bad_token",
			Offset:     0,
			Error:      errors.New("Bad AccessToken"),
		},
		{
			Query:      "",
			OrderField: "Age",
			OrderBy:    -1,
			Limit:      -1,
			Token:      token,
			Offset:     0,
			Error:      errors.New("limit must be > 0"),
		},
		{
			Query:      "",
			OrderField: "Age",
			OrderBy:    -1,
			Limit:      1,
			Token:      token,
			Offset:     -1,
			Error:      errors.New("offset must be > 0"),
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	testFile, err := os.Open("test.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer testFile.Close()

	outputsBytes, _ := io.ReadAll(testFile)
	outputs := make([]*SearchResponse, len(cases))
	json.Unmarshal(outputsBytes, &outputs)

	t.Run("OK", func(t *testing.T) {
		for i := 0; i < len(cases); i++ {
			t.Run(fmt.Sprintf("№%v", i+1), func(t *testing.T) {
				req := SearchRequest{
					Limit:      cases[i].Limit,
					Offset:     cases[i].Offset,
					Query:      cases[i].Query,
					OrderField: cases[i].OrderField,
					OrderBy:    cases[i].OrderBy,
				}
				srv := &SearchClient{
					AccessToken: cases[i].Token,
					URL:         ts.URL,
				}
				resp, err := srv.FindUsers(req)
				if err != nil {
					assert.Equal(t, cases[i].Error, err)
					return
				}

				assert.Equal(t, outputs[i], resp)
				assert.Equal(t, cases[i].Error, err)
			})
		}
	})

	tsErr := httptest.NewServer(http.HandlerFunc(SearchServerDummy))
	defer tsErr.Close()

	errorCases := []TestErrorCase{
		{
			Query: "bad url",
			URL:   "bad_url",
			Error: errors.New(
				`unknown error Get "bad_url?limit=2&offset=0&order_by=-1&order_field=Name&query=bad+url": unsupported protocol scheme ""`,
			),
		},
		{
			Query: "timeout",
			URL:   tsErr.URL,
			Error: errors.New(
				"timeout for limit=2&offset=0&order_by=-1&order_field=Name&query=timeout",
			),
		},
		{
			Query: "StatusInternalServerError",
			URL:   tsErr.URL,
			Error: errors.New(
				"SearchServer fatal error",
			),
		},
		{
			Query: "StatusBadRequestWithBadJson",
			URL:   tsErr.URL,
			Error: errors.New(
				"cant unpack error json: invalid character 'n' looking for beginning of object key string",
			),
		},
		{
			Query: "ErrorBadOrderField",
			URL:   tsErr.URL,
			Error: errors.New(
				"OrderFeld Name invalid",
			),
		},
		{
			Query: "unknown bad request error",
			URL:   tsErr.URL,
			Error: errors.New(
				"unknown bad request error: unknown",
			),
		},
		{
			Query: "__broken_json",
			URL:   tsErr.URL,
			Error: errors.New(
				"cant unpack result json: unexpected end of JSON input",
			),
		},
	}

	t.Run("Error", func(t *testing.T) {
		for i := 0; i < len(errorCases); i++ {
			t.Run(fmt.Sprintf("№%v", i+1), func(t *testing.T) {
				req := SearchRequest{
					Limit:      1,
					Offset:     0,
					Query:      errorCases[i].Query,
					OrderField: "Name",
					OrderBy:    -1,
				}
				srv := &SearchClient{
					AccessToken: cases[i].Token,
					URL:         errorCases[i].URL,
				}
				_, err := srv.FindUsers(req)
				if err != nil {
					assert.Equal(t, errorCases[i].Error, err)
					return
				}
			})

			// assert.Equal(t, outputs[i], resp)
			// assert.Equal(t, errorCases[i].Error, err)
		}
	})
}
