package builder

import (
	"testing"

	"github.com/rancher/norman/types"
	"github.com/stretchr/testify/assert"
)

func TestEmptyStringWithDefault(t *testing.T) {
	schema := &types.Schema{
		ResourceFields: map[string]types.Field{
			"foo": {
				Default: "foo",
				Type:    "string",
				Create:  true,
			},
		},
	}
	schemas := types.NewSchemas()
	schemas.AddSchema(*schema)

	builder := NewBuilder(&types.APIContext{})

	// Test if no field we set to "foo"
	result, err := builder.Construct(schema, nil, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := result["foo"]
	assert.True(t, ok)
	assert.Equal(t, "foo", value)

	// Test if field is "" we set to "foo"
	result, err = builder.Construct(schema, map[string]interface{}{
		"foo": "",
	}, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok = result["foo"]
	assert.True(t, ok)
	assert.Equal(t, "foo", value)
}

func TestSingleElementPassedExpectedArray(t *testing.T) {
	schema := &types.Schema{
		ResourceFields: map[string]types.Field{
			"foo": {
				Default: []interface{}{"bar"},
				Type:    "array[string]",
				Create:  true,
			},
		},
	}

	schemas := types.NewSchemas()
	schemas.AddSchema(*schema)
	builder := NewBuilder(&types.APIContext{})

	// Test if no field we set to "bar"
	result, err := builder.Construct(schema, nil, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok := result["foo"]
	assert.True(t, ok)
	assert.Equal(t, []interface{}{"bar"}, value)

	// ensure that if a single element is passed to a field
	// which expects an array, the single element is put inside an array
	// this allows for the following syntax variations to be valid
	//
	//	1)
	//	service-account-key-file: key1.pub
	//	-------
	//  2)
	//	service-account-key-file:
	//		- key1.pub
	//		- key2.pub
	//	-------
	// variation 1 will be transformed into 2 upon submission.
	//
	result, err = builder.Construct(schema, map[string]interface{}{
		"foo": "baz",
	}, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok = result["foo"]
	assert.True(t, ok)
	assert.Equal(t, []interface{}{"baz"}, value)

	result, err = builder.Construct(schema, map[string]interface{}{
		"foo": []interface{}{"qux"},
	}, Create)
	if err != nil {
		t.Fatal(err)
	}
	value, ok = result["foo"]
	assert.True(t, ok)
	assert.Equal(t, []interface{}{"qux"}, value)

}
