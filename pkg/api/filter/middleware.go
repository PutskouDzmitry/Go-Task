package filter

import (
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	LIMIT_KEY           = "limit"
	OFFSET_KEY          = "offset"
	QUERY_KEY           = "q"
	CONTEXT_LIMIT_KEY   = "limit"
	CONTEXT_OFFSET_KEY  = "offset"
	CONTEXT_FILTERS_KEY = "filters"
	CONTEXT_QUERY_KEY   = "query"
)

type FieldType int

const (
	TYPE_STRING FieldType = iota
	TYPE_TIME
	TYPE_INT
	TYPE_BOOL
	TYPE_FLOAT64
	TYPE_STRING_ARRAY
	TYPE_UUID_ARRAY
)

func FiltersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := make([]Field, 0)
		query := c.Request.URL.Query()
		for key := range query {
			for _, s := range c.QueryArray(key) {
				if key == LIMIT_KEY {
					limit, err := strconv.Atoi(s)
					if err == nil && limit > 0 {
						c.Set(CONTEXT_LIMIT_KEY, uint64(limit))
					}
					continue
				}
				if key == OFFSET_KEY {
					offset, err := strconv.Atoi(s)
					if err == nil && offset > 0 {
						c.Set(CONTEXT_OFFSET_KEY, uint64(offset))
					}
					continue
				}
				if key == QUERY_KEY {
					c.Set(CONTEXT_QUERY_KEY, s)
					continue
				}
				index := strings.Index(s, ":")
				operator := OperatorEq
				value := s
				if index != -1 {
					operator = s[:index]
					value = s[index+1:]
				}
				switch operator {
				case OperatorEq, OperatorGt, OperatorLt, OperatorGte, OperatorLte:
					fields = append(fields, Field{
						Name:     key,
						Operator: operator,
						Value:    value,
					})
				default:
					continue
				}
			}
		}
		c.Set(CONTEXT_FILTERS_KEY, fields)
		c.Next()
	}
}

func parseField(field Field, fieldType FieldType) (interface{}, error) {

	switch fieldType {
	case TYPE_STRING:
		return field.Value.(string), nil
	case TYPE_STRING_ARRAY:
		return strings.Split(field.Value.(string), ","), nil
	case TYPE_UUID_ARRAY:
		strs := strings.Split(field.Value.(string), ",")
		uuids := make([]uuid.UUID, 0, len(strs))
		for _, s := range strs {
			parse, err := uuid.Parse(s)
			if err != nil {
				return nil, err
			}
			uuids = append(uuids, parse)
		}
		return uuids, nil
	case TYPE_INT:
		intValue, err := strconv.Atoi(field.Value.(string))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert %s to int, value: \"%s\"", field.Name, field.Value)
		}
		return intValue, nil
	case TYPE_BOOL:
		strValue := field.Value.(string)
		if strValue == "0" || strValue == "false" {
			return false, nil
		} else if strValue == "1" || strValue == "true" {
			return true, nil
		} else {
			return nil, errors.Errorf("cannot convert %s to bool, value: \"%s\", use \"true\" or \"false\" or \"1\" or \"0\"", field.Name, field.Value)
		}
	case TYPE_TIME:
		parsedTime, err := time.Parse(time.RFC3339, field.Value.(string))
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert %s to time, value: \"%s\", use RFC3339 layout", field.Name, field.Value)
		}
		return parsedTime, nil
	case TYPE_FLOAT64:
		floatValue, err := strconv.ParseFloat(field.Value.(string), 64)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot convert %s to float64, value: \"%s\"", field.Name, field.Value)
		}
		return floatValue, nil
	default:
		return nil, errors.Errorf("unknown type for field: %s", field.Name)
	}
}

func GetOptionsFromContext(c *gin.Context, types map[string]FieldType) (*Options, error) {
	var limit, offset uint64
	var query string

	limitAny, exists := c.Get(CONTEXT_LIMIT_KEY)
	if exists {
		limit = limitAny.(uint64)
	}
	offsetAny, exists := c.Get(CONTEXT_OFFSET_KEY)
	if exists {
		offset = offsetAny.(uint64)
	}
	queryAny, exists := c.Get(CONTEXT_QUERY_KEY)
	if exists {
		query = queryAny.(string)
	}
	fields := c.MustGet(CONTEXT_FILTERS_KEY).([]Field)

	typedFields := make([]Field, 0)

	for _, field := range fields {
		if fieldType, ok := types[field.Name]; ok {
			value, err := parseField(field, fieldType)
			if err != nil {
				return nil, err
			}
			typedFields = append(typedFields, Field{
				Name:     field.Name,
				Operator: field.Operator,
				Value:    value,
			})
		}
	}

	return &Options{
		Limit:  limit,
		Offset: offset,
		Fields: typedFields,
		Query:  query,
	}, nil
}
