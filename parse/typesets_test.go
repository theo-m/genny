package parse_test

import (
	"testing"

	"github.com/cheekybits/genny/parse"
	"github.com/stretchr/testify/assert"
)

func TestArgsToTypeset(t *testing.T) {

	args := "Person=george:man,woman Animal=rex:dog,cat Place=london,paris"
	ts, err := parse.TypeSet(args)

	if assert.NoError(t, err) {
		if assert.Equal(t, 8, len(ts)) {

			assert.Equal(t, ts[0]["Person"].Alias, "george")
			assert.Equal(t, ts[0]["Animal"].Alias, "rex")
			assert.Equal(t, ts[0]["Person"].Type, "man")
			assert.Equal(t, ts[0]["Animal"].Type, "dog")
			assert.Equal(t, ts[0]["Place"].Type, "london")

			assert.Equal(t, ts[0]["Person"].Alias, "george")
			assert.Equal(t, ts[0]["Animal"].Alias, "rex")
			assert.Equal(t, ts[1]["Person"].Type, "man")
			assert.Equal(t, ts[1]["Animal"].Type, "dog")
			assert.Equal(t, ts[1]["Place"].Type, "paris")

			assert.Equal(t, ts[2]["Person"].Type, "man")
			assert.Equal(t, ts[2]["Animal"].Type, "cat")
			assert.Equal(t, ts[2]["Place"].Type, "london")

			assert.Equal(t, ts[3]["Person"].Type, "man")
			assert.Equal(t, ts[3]["Animal"].Type, "cat")
			assert.Equal(t, ts[3]["Place"].Type, "paris")

			assert.Equal(t, ts[4]["Person"].Type, "woman")
			assert.Equal(t, ts[4]["Animal"].Type, "dog")
			assert.Equal(t, ts[4]["Place"].Type, "london")

			assert.Equal(t, ts[5]["Person"].Type, "woman")
			assert.Equal(t, ts[5]["Animal"].Type, "dog")
			assert.Equal(t, ts[5]["Place"].Type, "paris")

			assert.Equal(t, ts[6]["Person"].Type, "woman")
			assert.Equal(t, ts[6]["Animal"].Type, "cat")
			assert.Equal(t, ts[6]["Place"].Type, "london")

			assert.Equal(t, ts[7]["Person"].Type, "woman")
			assert.Equal(t, ts[7]["Animal"].Type, "cat")
			assert.Equal(t, ts[7]["Place"].Type, "paris")

		}
	}

	ts, err = parse.TypeSet("Person=man Animal=dog Place=london")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ts))
	}
	ts, err = parse.TypeSet("Person=1,2,3,4,5 Animal=1,2,3,4,5 Place=1,2,3,4,5")
	if assert.NoError(t, err) {
		assert.Equal(t, 125, len(ts))
	}
	ts, err = parse.TypeSet("Person=1 Animal=1,2,3,4,5 Place=1,2")
	if assert.NoError(t, err) {
		assert.Equal(t, 10, len(ts))
	}

	ts, err = parse.TypeSet("Person=interface{} Animal=interface{} Place=interface{}")
	if assert.NoError(t, err) {
		assert.Equal(t, 1, len(ts))
		assert.Equal(t, ts[0]["Animal"].Type, "interface{}")
		assert.Equal(t, ts[0]["Person"].Type, "interface{}")
		assert.Equal(t, ts[0]["Place"].Type, "interface{}")
	}

}
