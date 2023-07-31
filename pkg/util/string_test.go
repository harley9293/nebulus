package util

import (
	"testing"
)

func TestStrFirstToUpper(t *testing.T) {
	if StrFirstToUpper("") != "" {
		t.Error("StrFirstToUpper(\"\") != \"\"")
	}

	if StrFirstToUpper("hello") != "Hello" {
		t.Error("StrFirstToUpper(\"hello\") != \"Hello\"")
	}
}

func TestUnicode2Hans(t *testing.T) {
	if string(Unicode2Hans([]byte(""))) != "" {
		t.Error("Unicode2Hans(\"\") != \"\"")
	}

	if string(Unicode2Hans([]byte(`\u4e2d\u56fd`))) != "中国" {
		t.Error("Unicode2Hans(\"\") != \"\"")
	}

	if string(Unicode2Hans([]byte(`\u12`))) != `\u12` {
		t.Error("Unicode2Hans(\"\") != \"\"")
	}
}
