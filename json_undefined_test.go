package unusual_generics_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xakep666/unusual_generics"
)

func TestJSONUndefined(t *testing.T) {
	var testJSON = []byte(`{"field_a": "some text", "field_b": null}`)

	t.Run("With JSONUndefined", func(t *testing.T) {
		var result struct {
			FieldA unusual_generics.JSONUndefined[string] `json:"field_a"`
			FieldB unusual_generics.JSONUndefined[string] `json:"field_b"`
			FieldC unusual_generics.JSONUndefined[string] `json:"field_c"`
		}

		require.NoError(t, json.Unmarshal(testJSON, &result))

		assert.True(t, result.FieldA.Defined)
		assert.Equal(t, "some text", *result.FieldA.Value)

		assert.True(t, result.FieldB.Defined)
		assert.Nil(t, result.FieldB.Value)

		assert.False(t, result.FieldC.Defined)
	})

	t.Run("Without JSONUndefined", func(t *testing.T) {
		var result struct {
			FieldA *string `json:"field_a"`
			FieldB *string `json:"field_b"`
			FieldC *string `json:"field_c"`
		}

		require.NoError(t, json.Unmarshal(testJSON, &result))

		assert.Equal(t, "some text", *result.FieldA)

		assert.Nil(t, result.FieldB)

		assert.Nil(t, result.FieldC)
	})
}

func ExampleJSONUndefined() {
	var testJSON = []byte(`{"field_a": "some text", "field_b": null}`)

	var result struct {
		FieldA unusual_generics.JSONUndefined[string] `json:"field_a"`
		FieldB unusual_generics.JSONUndefined[string] `json:"field_b"`
		FieldC unusual_generics.JSONUndefined[string] `json:"field_c"`
	}

	if err := json.Unmarshal(testJSON, &result); err != nil {
		panic(err)
	}

	fmt.Printf("FieldA: Defined=%t, Value=%s\n", result.FieldA.Defined, *result.FieldA.Value)
	fmt.Printf("FieldB: Defined=%t, Value=%v\n", result.FieldB.Defined, result.FieldB.Value)
	fmt.Printf("FieldC: Defined=%t, Value=%v\n", result.FieldC.Defined, result.FieldC.Value)
	// Output:
	// FieldA: Defined=true, Value=some text
	// FieldB: Defined=true, Value=<nil>
	// FieldC: Defined=false, Value=<nil>
}

func ExampleJSONUndefined_only_pointers() {
	var testJSON = []byte(`{"field_a": "some text", "field_b": null}`)

	var result struct {
		FieldA *string `json:"field_a"`
		FieldB *string `json:"field_b"`
		FieldC *string `json:"field_c"`
	}

	if err := json.Unmarshal(testJSON, &result); err != nil {
		panic(err)
	}

	fmt.Printf("FieldA: %s\n", *result.FieldA)
	fmt.Printf("FieldB: %v\n", result.FieldB)
	fmt.Printf("FieldC: %v\n", result.FieldC)
	// Output:
	// FieldA: some text
	// FieldB: <nil>
	// FieldC: <nil>
}
