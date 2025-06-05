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

func TestCafeWhenOk(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow&count=5", nil)
	resp := httptest.NewRecorder()
	mainHandle(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
	want := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент,Ложка и вилка"
	assert.Equal(t, want, strings.TrimSpace(resp.Body.String()))
}

func TestCafeNegative(t *testing.T) {
	requests := []struct {
		url      string
		wantCode int
		wantBody string
	}{
		{
			url:      "/cafe?city=london",
			wantCode: http.StatusBadRequest,
			wantBody: "unknown city",
		},
		{
			url:      "/cafe?city=moscow&count=abc",
			wantCode: http.StatusBadRequest,
			wantBody: "incorrect count",
		},
		{
			url:      "/cafe?city=moscow&count=-1",
			wantCode: http.StatusBadRequest,
			wantBody: "incorrect count",
		},
	}

	for _, tt := range requests {
		req := httptest.NewRequest(http.MethodGet, tt.url, nil)
		resp := httptest.NewRecorder()
		mainHandle(resp, req)

		require.Equal(t, tt.wantCode, resp.Code, "Код ответа для URL %s должен быть %d", tt.url, tt.wantCode)
		assert.Equal(t, tt.wantBody, strings.TrimSpace(resp.Body.String()), "Тело ответа для URL %s должно быть '%s'", tt.url, tt.wantBody)
	}
}

func TestCafeCount(t *testing.T) {
	requests := []struct {
		count int
		city  string
		want  int
	}{
		{count: 0, city: "moscow", want: 0},
		{count: 1, city: "moscow", want: 1},
		{count: 2, city: "moscow", want: 2},
		{count: 100, city: "moscow", want: len(cafeList["moscow"])},
		{count: 100, city: "tula", want: len(cafeList["tula"])},
		{count: 2, city: "tula", want: 2},
		{count: 0, city: "tula", want: 0},
	}

	for _, tt := range requests {

		req := httptest.NewRequest(http.MethodGet, "/cafe?city="+tt.city+"&count="+strconv.Itoa(tt.count), nil)
		resp := httptest.NewRecorder()
		mainHandle(resp, req)

		require.Equal(t, http.StatusOK, resp.Code, "Код ответа должен быть 200 OK для city=%s, count=%d", tt.city, tt.count)

		body := strings.TrimSpace(resp.Body.String())
		var got []string
		if body != "" {
			got = strings.Split(body, ",")
		} else {
			got = []string{}
		}

		assert.Len(t, got, tt.want, "Ожидалось %d кафе для city=%s, count=%d, получено %d", tt.want, tt.city, tt.count, len(got))
	}
}

func TestCafeSearch(t *testing.T) {
	requests := []struct {
		search    string
		wantCount int
	}{
		{search: "фасоль", wantCount: 0},
		{search: "кофе", wantCount: 2},
		{search: "вилка", wantCount: 1},
		{search: "КОФЕ", wantCount: 2},
		{search: "Студент", wantCount: 1},
		{search: "", wantCount: len(cafeList["moscow"])},
	}

	for _, tt := range requests {
		req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow&search="+tt.search, nil)
		resp := httptest.NewRecorder()
		mainHandle(resp, req)

		require.Equal(t, http.StatusOK, resp.Code, "Код ответа должен быть 200 OK для search=%s", tt.search)

		body := strings.TrimSpace(resp.Body.String())
		var got []string
		if body != "" {
			got = strings.Split(body, ",")
		} else {
			got = []string{}
		}

		assert.Len(t, got, tt.wantCount, "Ожидалось %d кафе для search=%s, получено %d", tt.wantCount, tt.search, len(got))

		if tt.wantCount > 0 {
			for _, cafe := range got {
				assert.Contains(t, strings.ToLower(cafe), strings.ToLower(tt.search), "Кафе '%s' должно содержать строку поиска '%s'", cafe, tt.search)
			}
		}
	}
}
