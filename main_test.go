package main

import (
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeWhenOk(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=moscow&count=5", nil)
	resp := httptest.NewRecorder()
	mainHandle(resp, req)

	require.Equal(t, http.StatusOK, resp.Code, "Код ответа должен быть 200 OK")
	want := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент,Ложка и вилка"
	assert.Equal(t, want, strings.TrimSpace(resp.Body.String()), "Неверное тело ответа для успешного запроса")
}

func TestCafeNegative(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/cafe?city=london", nil)
	resp := httptest.NewRecorder()
	mainHandle(resp, req)

	require.Equal(t, http.StatusBadRequest, resp.Code, "Код ответа должен быть 400 Bad Request для неизвестного города")
	want := "unknown city"
	assert.Equal(t, want, strings.TrimSpace(resp.Body.String()), "Неверное сообщение об ошибке для неизвестного города")
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
