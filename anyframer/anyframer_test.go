package anyframer_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yesoreyeram/grafana-plugins/anyframer"
)

type Framer = anyframer.AnyFramer

type testInputs struct {
	input   string
	framer  Framer
	wantErr error
}

var testToFrame = func(tt testInputs) func(t *testing.T) {
	return func(t *testing.T) {
		got, err := tt.framer.ToFrame(tt.input)
		if tt.wantErr != nil {
			require.NotNil(t, err)
			assert.Equal(t, tt.wantErr, err)
			return
		}
		require.Nil(t, err)
		require.NotNil(t, got)
		experimental.CheckGoldenJSONFrame(t, "testdata/golden", t.Name(), got, updateGoldenFile)
	}
}

func TestToFrame(t *testing.T) {
	t.Run("empty input should throw error", testToFrame(testInputs{
		wantErr: errors.New("invalid/empty input"),
	}))
	t.Run("single string without quotes", testToFrame(testInputs{
		input: `hello world`,
	}))
	t.Run("single string without quotes and numeric values", testToFrame(testInputs{
		input: `hello 123 world`,
	}))
	t.Run("single string", testToFrame(testInputs{
		input: `"hello"`,
	}))
	t.Run("single number", testToFrame(testInputs{
		input: `123.456`,
	}))
	t.Run("single boolean", testToFrame(testInputs{
		input: `true`,
	}))
	t.Run("basic csv", testToFrame(testInputs{
		input: "a,b,c\n1,2,3\n4,5,6",
	}))
	t.Run("basic tsv", testToFrame(testInputs{
		input: "a	b	c\n1	2	3\n4	5	6",
	}))
	t.Run("simple string array", testToFrame(testInputs{
		input: `["foo","bar"]`,
	}))
	t.Run("simple numeric array", testToFrame(testInputs{
		input: `[123, 456.789]`,
	}))
	t.Run("simple boolean array", testToFrame(testInputs{
		input: `[false, true]`,
	}))
	t.Run("simple object array", testToFrame(testInputs{
		input: `[{"name":"foo","salary": 123, "self_employed":false},{"name":"bar","salary": 456.789, "self_employed":true}]`,
	}))
	t.Run("nested json", testToFrame(testInputs{
		input: `{ "users" : [{"name":"foo","salary": 123, "self_employed":false},{"name":"bar","salary": 456.789, "self_employed":true}] }`,
	}))
	t.Run("simple xml", testToFrame(testInputs{
		input: `<?xml version="1.0" encoding="UTF-8"?><root><row><name>foo</name><salary>123</salary><self_employed>false</self_employed></row><row><name>bar</name><salary>456.789</salary><self_employed>true</self_employed></row></root>`,
	}))
	t.Run("simple html", testToFrame(testInputs{
		input: `<html><body><table class="table table-bordered table-hover table-condensed"><thead><tr><th title="Field #1">name</th><th title="Field #2">salary</th><th title="Field #3">self_employed</th></tr></thead><tbody><tr><td>foo</td><td align="right">123</td><td>false</td></tr><tr><td>bar</td><td align="right">456.789</td><td>true</td></tr></tbody></table></body></html>`,
	}))
	t.Run("simple xml with root selector", testToFrame(testInputs{
		input:  `<?xml version="1.0" encoding="UTF-8"?><root><row><name>foo</name><salary>123</salary><self_employed>false</self_employed></row><row><name>bar</name><salary>456.789</salary><self_employed>true</self_employed></row></root>`,
		framer: Framer{RootSelector: "root.row"},
	}))
	t.Run("simple html with root selector", testToFrame(testInputs{
		input:  `<html><body><table class="table table-bordered table-hover table-condensed"><thead><tr><th title="Field #1">name</th><th title="Field #2">salary</th><th title="Field #3">self_employed</th></tr></thead><tbody><tr><td>foo</td><td align="right">123</td><td>false</td></tr><tr><td>bar</td><td align="right">456.789</td><td>true</td></tr></tbody></table></body></html>`,
		framer: Framer{RootSelector: "html.body.table.tbody.tr"},
	}))
}

func TestAnyFile(t *testing.T) {
	files, _ := os.ReadDir("./testdata/all")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b, err := os.ReadFile("./testdata/all/" + file.Name())
		require.Nil(t, err)
		require.NotNil(t, b)
		framer := Framer{RawURL: file.Name()}
		switch file.Name() {
		case "users.json":
		case "users.csv":
		case "users.tsv":
		case "users.xml":
			framer.RootSelector = "root.row"
		case "users.html":
			framer.RootSelector = "html.body.table.tbody.tr"
		default:
		}
		t.Run(strings.ReplaceAll(file.Name(), ".", "_"), testToFrame(testInputs{input: string(b), framer: framer}))
	}
}

func testName(input string) string {
	return strings.ReplaceAll(strings.ReplaceAll(input, ".json", ""), " ", "-")
}

func TestJSONFiles(t *testing.T) {
	files, _ := os.ReadDir("./testdata/json")
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		b, err := os.ReadFile("./testdata/json/" + file.Name())
		require.Nil(t, err)
		require.NotNil(t, b)
		framer := Framer{RawURL: file.Name()}
		switch file.Name() {
		default:
		}
		t.Run(testName(file.Name()), testToFrame(testInputs{input: string(b), framer: framer}))
		switch file.Name() {
		case "org.json":
			t.Run("org emplyees", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: "employees"},
			}))
		case "library.json":
			t.Run("library books", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: "library.books"},
			}))
			t.Run("library loans", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: "library.loans"},
			}))
			t.Run("library customers", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: "library.customers"},
			}))
			t.Run("library books title", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: `library.books.title`},
			}))
			t.Run("library books price list", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: `$map(library.books,function($v){return { "title": $v.title, "price": $v.price }})`},
			}))
			t.Run("library costly books", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: `$map(library.books[price > 30],function($v){return { "title": $v.title, "price": $v.price }})`},
			}))
			t.Run("library books total value", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: `$sum($map(library.books,function($v){return {"value":$v.price*$v.copies}}).value)`},
			}))
			t.Run("library books stats", testToFrame(testInputs{input: string(b),
				framer: Framer{RawURL: file.Name(), RootSelector: `{
					"total books type": $count(library.books.authors),
					"total books count": $sum(library.books.copies),
					"total books value": $sum(library.books.price),
					"max books value": $max(library.books.price),
					"min books value": $min(library.books.price),
					"all books value": $sum($map(library.books,function($v){return {"value":$v.price*$v.copies}}).value)
				}`},
			}))
		}
	}
}
