package map_converter

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertStructToMap(t *testing.T) {
	now := time.Now()
	desc := "Description"
	empCount := "100"

	type testStruct struct {
		EmployeesCount *string    `json:"employees_count"`
		FoundationDate *time.Time `json:"foundation_date"`
		Description    *string    `json:"description"`
		UserID         string     `json:"user_id"`
	}

	dto := testStruct{
		EmployeesCount: &empCount,
		FoundationDate: &now,
		Description:    &desc,
		UserID:         "user_id",
	}

	testCases := []struct {
		name string
		dto  interface{}
		want map[string]interface{}
	}{
		{
			name: "Test 1",
			dto:  dto,
			want: map[string]interface{}{
				"employees_count": empCount,
				"foundation_date": now,
				"description":     desc,
				"user_id":         "user_id",
			},
		},
		{
			name: "Test 1",
			dto: testStruct{
				EmployeesCount: nil,
				FoundationDate: &now,
				Description:    &desc,
				UserID:         "user_id",
			},
			want: map[string]interface{}{
				"foundation_date": now,
				"description":     desc,
				"user_id":         "user_id",
			},
		}, {
			name: "Test 1",
			dto: testStruct{
				EmployeesCount: nil,
				FoundationDate: nil,
				Description:    nil,
				UserID:         "user_id",
			},
			want: map[string]interface{}{
				"user_id": "user_id",
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertStructToMap(tt.dto)
			assert.Equal(t, tt.want, got)
		})
	}

}
