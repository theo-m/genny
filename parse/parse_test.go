package parse_test

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/cheekybits/genny/parse"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	// input
	filename       string
	outputFilename string
	pkgName        string
	in             string
	tag            string
	types          []map[string]parse.TypeRef

	// expectations
	expectedOut string
	expectedErr error
}{
	{
		filename:    "generic_queue.go",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/queue/int_queue.go`,
	},
	{
		filename:    "generic_queue.go",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{"CrazyNumber", "int"}}},
		expectedOut: `test/queue/int_queue_aliased.go`,
	},
	{
		filename:    "generic_queue.go",
		pkgName:     "changed",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/queue/changed/int_queue.go`,
	},
	{
		filename:    "generic_queue.go",
		in:          `test/queue/generic_queue.go`,
		types:       []map[string]parse.TypeRef{{"Something": parse.TypeRef{"float32", "float32"}}},
		expectedOut: `test/queue/float32_queue.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypes/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{{
			"KeyType":   parse.TypeRef{"string", "string"},
			"ValueType": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/multipletypes/string_int_simplemap.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypes/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{{
			"KeyType":   parse.TypeRef{"interface{}", "interface{}"},
			"ValueType": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/multipletypes/interface_int_simplemap.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypes/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{{
			"KeyType":   parse.TypeRef{"*MyType1", "*MyType1"},
			"ValueType": parse.TypeRef{"*MyOtherType", "*MyOtherType"}}},
		expectedOut: `test/multipletypes/custom_types_simplemap.go`,
	},
	{
		filename:    "generic_internal.go",
		in:          `test/unexported/generic_internal.go`,
		types:       []map[string]parse.TypeRef{{"secret": parse.TypeRef{"*myType", "*myType"}}},
		expectedOut: `test/unexported/mytype_internal.go`,
	},
	{
		filename: "generic_simplemap.go",
		in:       `test/multipletypesets/generic_simplemap.go`,
		types: []map[string]parse.TypeRef{
			{"KeyType": parse.TypeRef{"int", "int"}, "ValueType": parse.TypeRef{"string", "string"}},
			{"KeyType": parse.TypeRef{"float64", "float64"}, "ValueType": parse.TypeRef{"bool", "bool"}},
		},
		expectedOut: `test/multipletypesets/many_simplemaps.go`,
	},
	{
		filename:    "generic_number.go",
		in:          `test/numbers/generic_number.go`,
		types:       []map[string]parse.TypeRef{{"NumberType": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/numbers/int_number.go`,
	},
	{
		filename:    "generic_digraph.go",
		in:          `test/bugreports/generic_digraph.go`,
		types:       []map[string]parse.TypeRef{{"Node": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/bugreports/int_digraph.go`,
	},
	{
		filename:    "generic_new_and_make_slice.go",
		in:          `test/bugreports/generic_new_and_make_slice.go`,
		types:       []map[string]parse.TypeRef{{"NumberType": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/bugreports/int_new_and_make_slice.go`,
	},
	{
		filename:    "cell_x.go",
		in:          `test/bugreports/cell_x.go`,
		types:       []map[string]parse.TypeRef{{"X": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/bugreports/cell_int.go`,
	},
	{
		filename:    "interface_generic_type.go",
		in:          `test/bugreports/interface_generic_type.go`,
		types:       []map[string]parse.TypeRef{{"GenericType": parse.TypeRef{"uint8", "uint8"}}},
		expectedOut: `test/bugreports/interface_uint8.go`,
	},
	{
		filename:    "negation_generic.go",
		in:          `test/bugreports/negation_generic.go`,
		types:       []map[string]parse.TypeRef{{"SomeThing": parse.TypeRef{"string", "string"}}},
		expectedOut: `test/bugreports/negation_string.go`,
	},
	{
		filename:    "buildtags.go",
		in:          `test/buildtags/buildtags.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{"int", "int"}}},
		expectedOut: `test/buildtags/buildtags_expected.go`,
		tag:         "genny",
	},
	{
		filename:    "buildtags.go",
		in:          `test/buildtags/buildtags.go`,
		types:       []map[string]parse.TypeRef{{"_t_": parse.TypeRef{"string", "string"}}},
		expectedOut: `test/buildtags/buildtags_expected_nostrip.go`,
		tag:         "",
	},
}

func TestParse(t *testing.T) {

	for _, test := range tests {

		test.in = contents(test.in)
		test.expectedOut = contents(test.expectedOut)

		bytes, err := parse.Generics(test.filename, test.outputFilename, test.pkgName, test.tag, strings.NewReader(test.in), test.types)

		// check the error
		if test.expectedErr == nil {
			assert.NoError(t, err, "(%s) No error was expected but got: %s", test.filename, err)
		} else {
			assert.NotNil(t, err, "(%s) No error was returned by one was expected: %s", test.filename, test.expectedErr)
			assert.IsType(t, test.expectedErr, err, "(%s) Generate should return object of type %v", test.filename, test.expectedErr)
		}

		// assert the response
		if !assert.Equal(t, string(bytes), test.expectedOut, "Parse didn't generate the expected output.") {
			log.Println("EXPECTED: " + test.expectedOut)
			log.Println("ACTUAL: " + string(bytes))
		}

	}

}

func contents(s string) string {
	if strings.HasSuffix(s, "go") {
		file, err := ioutil.ReadFile(s)
		if err != nil {
			panic(err)
		}
		return string(file)
	}
	return s
}
