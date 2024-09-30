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

var cafeList = map[string][]string{
	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("count missing"))
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong count value"))
		return
	}

	city := req.URL.Query().Get("city")

	cafe, ok := cafeList[city]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("wrong city value"))
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandler_Success(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	require.NoError(t, err, "Ошибка при создании запроса")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Ожидается статус код 200")

	assert.NotEmpty(t, rr.Body.String(), "Тело ответа не должно быть пустым")

	cafes := strings.Split(rr.Body.String(), ",")
	assert.Len(t, cafes, 2, "Должно вернуться 2 кафе")
}

func TestMainHandler_WrongCity(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=2&city=unknown", nil)
	require.NoError(t, err, "Ошибка при создании запроса")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Ожидается статус код 400")

	assert.Equal(t, "wrong city value", rr.Body.String(), "Ожидается сообщение об ошибке 'wrong city value'")
}

func TestMainHandler_CountMoreThanTotal(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	require.NoError(t, err, "Ошибка при создании запроса")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)

	handler.ServeHTTP(rr, req)

	expectedResponse := strings.Join(cafeList["moscow"], ",")

	assert.Equal(t, http.StatusOK, rr.Code, "Ожидается статус код 200")

	assert.Equal(t, expectedResponse, rr.Body.String(), "Должны вернуться все доступные кафе")
}
