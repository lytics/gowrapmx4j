package gowrapmx4j

import (
	"reflect"
	"testing"
)

func TestBraceRemoval(t *testing.T) {
	teststr := "{10.100.190.140=115.8 GB, 10.100.37.24=112.99 GB, 10.100.168.68=125.47 GB, 10.100.100.117=125.49 GB, 10.100.18.46=127.44 GB}"
	expected := "10.100.190.140=115.8 GB, 10.100.37.24=112.99 GB, 10.100.168.68=125.47 GB, 10.100.100.117=125.49 GB, 10.100.18.46=127.44 GB"
	out := removeBraces(teststr)
	if out != expected {
		t.Errorf("Braces not removed: \"%s\"", out)
	}
}

func TestBracketRemoval(t *testing.T) {
	teststr := "[10.100.168.68, 10.100.100.117, 10.100.37.24, 10.100.190.140, 10.100.18.46]"
	expected := "10.100.168.68, 10.100.100.117, 10.100.37.24, 10.100.190.140, 10.100.18.46"
	out := removeBrackets(teststr)
	if out != expected {
		t.Errorf("Brackets not removed")
	}
}

func TestValueSplitArray(t *testing.T) {
	teststr := "[10.100.168.68, 10.100.100.117, 10.100.37.24, 10.100.190.140, 10.100.18.46]"
	expected := []string{"10.100.168.68", "10.100.100.117", "10.100.37.24", "10.100.190.140", "10.100.18.46"}
	vals := separateValues(removeBrackets(teststr))

	if !reflect.DeepEqual(expected, vals) {
		t.Errorf("Separated values unequal\n%v != %v", expected, vals)
	}
}

func TestMapParse(t *testing.T) {
	teststr := "{10.100.190.140=115.8 GB, 10.100.37.24=112.99 GB, 10.100.168.68=125.47 GB, 10.100.100.117=125.49 GB, 10.100.18.46=127.44 GB}"
	parsed := parseMap(teststr)
	v, ok := parsed["10.100.190.140"]
	vstr := v.(string)
	if len(parsed) != 5 {
		t.Errorf("Incorrect number of elements parsed: %d", len(parsed))
	}
	if !ok {
		t.Errorf("10.100.190.140 key does not exist in parsed map")
	}
	if vstr != "115.8GB" {
		t.Errorf("Map value %v not equal to expected: \"115.8 GB\"", vstr)
	}
}
