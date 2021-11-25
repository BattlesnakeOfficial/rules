package test

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

var updateFixtures = flag.Bool("update-fixtures", false, "Regenerate fixtures in testdata based on current test output")

// RequireJSONMatchesFixture asserts that the JSON text in actual matches the
// JSON read from filename, without taking into account whitespace and
// ordering. Files can be specified relative to the calling test (e.g.
// testdata/example.json). To regenerate the expected test data automatically
// after making a code change, pass the `-update-fixtures` flag to `go test`.
func RequireJSONMatchesFixture(t *testing.T, filename string, actual string) {
	t.Helper()

	if *updateFixtures {
		var indented bytes.Buffer
		err := json.Indent(&indented, []byte(actual), "", "  ")
		require.NoError(t, err, "Failed to indent JSON")
		err = ioutil.WriteFile(filename, indented.Bytes(), 0644)
		require.NoError(t, err, "Failed to update fixture", filename)

		log.Printf("Updating fixture file %#v", filename)
	}

	expectedData, err := ioutil.ReadFile(filename)
	require.NoError(t, err, "Failed to read fixture", filename)

	require.JSONEq(t, string(expectedData), actual)
}
