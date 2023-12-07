package sorting

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSortMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SortMiddleware())
	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	t.Run("no sort field or order", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("invalid sort order", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/?sort_by=id&sort_order=invalid", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("mismatching length of fields and orders", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/?sort_by=id,name&sort_order=asc", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("correct sorting options", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/?sort_by=id,name&sort_order=asc,desc", nil)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestGetOptionsFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(SortMiddleware())
	router.GET("/", func(c *gin.Context) {
		options := GetOptionsFromContext(c)
		assert.Equal(t, &Options{Fields: []string{"id", "name"}, Orders: []string{"asc", "desc"}}, options)
		c.String(200, "OK")
	})

	request, _ := http.NewRequest(http.MethodGet, "/?sort_by=id,name&sort_order=asc,desc", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
}
