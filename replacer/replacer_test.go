package replacer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testSt struct {
	in  string
	m   map[string]interface{}
	out string
}

func TestReplaceByMap(t *testing.T) {
	lst := []testSt{
		{
			in: "{QUOTE}\naaa\n{QUERY}",
			m: map[string]interface{}{
				"QUOTE": "quote str",
				"QUERY": "query str",
			},
			out: "quote str\naaa\nquery str",
		},
		{
			in:  "{0}",
			out: "{0}",
		},
		{
			in:  "{0}",
			m:   map[string]interface{}{"0": "1"},
			out: "1",
		},
		{
			in: "{{{{{0}, }}} world",
			m: map[string]interface{}{
				"0": "hello",
			},
			out: "{{{{hello, }}} world",
		},
		{
			in: "1{0}2{1}{2}{3}{0}",
			m: map[string]interface{}{
				"0": "-",
				"1": "-",
				"2": "-",
				"3": "-",
			},
			out: "1-2----",
		},
		{
			in: "hello {0}",
			m: map[string]interface{}{
				"0": "world",
			},
			out: "hello world",
		},
		{
			in: "hello {{0}",
			m: map[string]interface{}{
				"0": "world",
			},
			out: "hello {world",
		},
		{
			in:  "",
			m:   map[string]interface{}{},
			out: "",
		},
		{
			in:  "}{0}{",
			m:   map[string]interface{}{},
			out: "}{0}{",
		},
		{
			in: "}}}{m}{{{{{}",
			m: map[string]interface{}{
				"m": "1",
			},
			out: "}}}1{{{{{}",
		},
		{
			in:  "}}}}}}}}",
			m:   nil,
			out: "}}}}}}}}",
		},
		{
			in:  "{",
			out: "{",
		},
		{
			in: "{   0   }",
			m: map[string]interface{}{
				"0": "1",
			},
			out: "{   0   }",
		},
	}
	for _, item := range lst {
		out := ReplaceByMap(item.in, item.m)
		assert.Equal(t, item.out, out)
	}
}

type testInterfaceSt struct {
	in  interface{}
	out string
}

type testString struct {
}

func (t *testString) String() string {
	return "tttt"
}

func TestAsString(t *testing.T) {
	lst := []testInterfaceSt{
		{
			in:  complex(4, 3),
			out: "(4+3i)",
		},
		{
			in:  nil,
			out: "<nil>",
		},
		{
			in:  "a",
			out: "a",
		},
		{
			in:  1,
			out: "1",
		},
		{
			in:  int64(123),
			out: "123",
		},
		{
			in:  3.14,
			out: "3.14",
		},
		{
			in:  float32(3),
			out: "3",
		},
		{
			in:  true,
			out: "true",
		},
		{
			in:  fmt.Errorf("abc"),
			out: "abc",
		},
		{
			in:  uint32(1345),
			out: "1345",
		},
		{
			in:  &testString{},
			out: "tttt",
		},
		{
			in: map[string]int{
				"a": 1,
			},
			out: `{"a":1}`,
		},
	}
	for _, item := range lst {
		out := asStrValue(item.in)
		assert.Equal(t, item.out, out)
	}
}

type testListSt struct {
	in  string
	m   []interface{}
	out string
}

func TestReplaceByList(t *testing.T) {
	lst := []testListSt{
		{
			in:  "{0}",
			m:   []interface{}{"a"},
			out: "a",
		},
		{
			in:  "{1} {0}",
			m:   []interface{}{"world", "hello"},
			out: "hello world",
		},
		{
			in:  "{{0} {1}{}",
			m:   []interface{}{true, 1},
			out: "{true 1{}",
		},
		{
			in:  "}}}}}}}{{{{{{{{0},{1},{2},{3}",
			m:   []interface{}{0, 1},
			out: "}}}}}}}{{{{{{{0,1,{2},{3}",
		},
		{
			in:  "{{{0000}}",
			m:   []interface{}{123},
			out: "{{{0000}}",
		},
	}
	for _, item := range lst {
		out := ReplaceByList(item.in, item.m...)
		assert.Equal(t, item.out, out)
	}
}
