package jsonframer

import (
	"bytes"
	"strings"

	"github.com/noborus/trdsql"
)

func QueryJSONUsingSQLite3(jsonString string, query string, rootSelector string) (string, error) {
	r := bytes.NewBufferString(jsonString)
	options := []trdsql.ReadOpt{}
	options = append(options, trdsql.InFormat(trdsql.JSON))
	if rootSelector != "" {
		if !strings.HasPrefix(rootSelector, ".") {
			rootSelector = "." + rootSelector
		}
		options = append(options, trdsql.InJQ(rootSelector))
	}
	importer, err := trdsql.NewBufferImporter("input", r, options...)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	writer := trdsql.NewJSONWriter(&trdsql.WriteOpts{OutFormat: trdsql.Format(trdsql.JSON), OutStream: &buf})
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	trd.Driver = "sqlite3"
	if err = trd.Exec(query); err != nil {
		a := err.Error()
		if a == "import: invalid names" {
			return "[]", nil
		}
		return "", err
	}
	return buf.String(), nil
}
