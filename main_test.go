package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		city  string
		count int
		want  int
	}{
		{city: "moscow", count: 0, want: 0},
		{city: "moscow", count: 1, want: 1},
		{city: "moscow", count: 2, want: 2},
		{city: "moscow", count: 100, want: min(100, len(cafeList["moscow"]))},
	}

	for _, v := range request {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&count="+strconv.Itoa(v.count), nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		responseString := strings.TrimSpace(response.Body.String())

		searchCafes := []string{}

		if responseString != "" {
			searchCafes = strings.Split(responseString, ",")
		}

		assert.Len(t, searchCafes, v.want)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	request := []struct {
		city      string
		search    string
		wantCount int
	}{
		{city: "moscow", search: "фасоль", wantCount: 0},
		{city: "moscow", search: "кофе", wantCount: 2},
		{city: "moscow", search: "вилка", wantCount: 1},
	}

	for _, v := range request {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/cafe?city="+v.city+"&search="+v.search, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		responseString := strings.TrimSpace(response.Body.String())

		searchCafes := []string{}

		if responseString != "" {
			searchCafes = strings.Split(responseString, ",")
		}

		assert.Len(t, searchCafes, v.wantCount)

		for _, cafe := range searchCafes {
			assert.Contains(t, strings.ToLower(cafe), strings.ToLower(v.search))
		}
	}
}
