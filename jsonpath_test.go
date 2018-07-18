package gojsonfilterbench

import (
	"fmt"
	"testing"

	"github.com/Knetic/govaluate"
	"github.com/buger/jsonparser"
)

type FlexibleMatcher struct{}

func NewFlexibleMatcher() *FlexibleMatcher {
	return &FlexibleMatcher{}
}

// Hard-coded parameter Array for performance test
func NewParameterArray(size int) *parameterArray {
	arr := &parameterArray{
		data: make([]interface{}, size, size),
	}
	return arr
}

type parameterArray struct {
	data []interface{}
}

func (p parameterArray) Get(name string) (interface{}, error) {
	switch name {
	case "firstName":
		return p.data[0], nil
	case "age":
		return p.data[1], nil
	case "isActive":
		return p.data[2], nil
	default:
		return nil, fmt.Errorf("Not found")
	}
}

// Returns true if it's a match. Error if something went wrong.
func (m *FlexibleMatcher) Match(data []byte, expression *govaluate.EvaluableExpression, parameters parameterArray) (bool, error) {
	var err error

	parameters.data[0], err = jsonparser.GetString(data, "name", "first")
	if err != nil {
		fmt.Printf("GetString Error: %v\n", err.Error())
		return false, err
	}

	parameters.data[1], err = jsonparser.GetInt(data, "age")
	if err != nil {
		fmt.Printf("GetInt Error: %v\n", err.Error())
		return false, err
	}

	parameters.data[2], err = jsonparser.GetBoolean(data, "isActive")
	if err != nil {
		fmt.Printf("GetBoolean Error: %v\n", err.Error())
		return false, err
	}

	result, err := expression.Eval(parameters)
	if err != nil {
		fmt.Printf("Evaluate Error: %v\n", err.Error())
		return false, err
	}
	return result.(bool), err
}

func BenchmarkJsonPathWithGoValuate(b *testing.B) {
	data, totalBytes, err := generateRandomData(1)
	if err != nil || len(data) == 0 {
		b.Fatalf("Data generation error: %s", err)
	}

	m := NewFlexibleMatcher()

	// Expression reformatted:
	expression, err := govaluate.NewEvaluableExpression("firstName == 'Neil' || (age < 50 && isActive == true)")
	if err != nil {
		b.Fatalf("NewEvaluableExpression Error: %s", err)
		return
	}

	// Pre-make parameters and re-use
	parameters := NewParameterArray(3)
	b.SetBytes(int64(totalBytes))
	b.ResetTimer()

	for j := 0; j < b.N; j++ {
		for i := 0; i < len(data); i++ {
			_, err := m.Match(data[i], expression, *parameters)

			if err != nil {
				b.Fatalf("Matcher error: %s", err)
			}
		}
	}
}
