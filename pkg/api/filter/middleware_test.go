package filter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestFiltersMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name   string
		url    string
		limit  uint64
		offset uint64
		query  string
	}{
		{"no params", "/", 0, 0, ""},
		{"limit only", "/?limit=10", 10, 0, ""},
		{"offset only", "/?offset=5", 0, 5, ""},
		{"query only", "/?q=test", 0, 0, "test"},
		{"all params", "/?limit=10&offset=5&q=test", 10, 5, "test"},
		{"negative limit", "/?limit=-10", 0, 0, ""},
		{"negative offset", "/?offset=-5", 0, 0, ""},
		{"non-numeric limit", "/?limit=ten", 0, 0, ""},
		{"non-numeric offset", "/?offset=five", 0, 0, ""},
		{"query and sort", "/?q=PENDING&sort_by=created_at", 0, 0, "PENDING"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(FiltersMiddleware())
			router.GET("/", func(c *gin.Context) {
				limit, exists := c.Get(CONTEXT_LIMIT_KEY)
				if !exists {
					limit = uint64(0)
				}

				offset, exists := c.Get(CONTEXT_OFFSET_KEY)
				if !exists {
					offset = uint64(0)
				}

				query, exists := c.Get(CONTEXT_QUERY_KEY)
				if !exists {
					query = ""
				}

				assert.Equal(t, tt.limit, limit)
				assert.Equal(t, tt.offset, offset)
				assert.Equal(t, tt.query, query)
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestGetOptionsFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name      string
		url       string
		fieldType FieldType
		fieldName string
		fieldVal  interface{}
		wantErr   bool
	}{
		{"no params", "/", TYPE_INT, "", nil, false},
		{"int field", "/?testField=1", TYPE_INT, "testField", 1, false},
		{"invalid int field", "/?testField=one", TYPE_INT, "testField", 0, true},
		{"bool field", "/?testField=true", TYPE_BOOL, "testField", true, false},
		{"invalid bool field", "/?testField=maybe", TYPE_BOOL, "testField", false, true},
		{"time field", "/?testField=2023-07-13T00:00:00Z", TYPE_TIME, "testField", mustParseTime("2023-07-13T00:00:00Z"), false},
		{"invalid time field", "/?testField=yesterday", TYPE_TIME, "testField", nil, true},
		{"string field", "/?testField=test", TYPE_STRING, "testField", "test", false},
		{"float field", "/?testField=1.23", TYPE_FLOAT64, "testField", 1.23, false},
		{"invalid float field", "/?testField=onepointtwothree", TYPE_FLOAT64, "testField", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.Default()
			router.Use(FiltersMiddleware())
			router.GET("/", func(c *gin.Context) {
				options, err := GetOptionsFromContext(c, map[string]FieldType{
					"testField": tt.fieldType,
				})

				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					if len(options.Fields) > 0 {
						assert.Equal(t, tt.fieldName, options.Fields[0].Name)
						assert.Equal(t, tt.fieldVal, options.Fields[0].Value)
					}
				}
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}
