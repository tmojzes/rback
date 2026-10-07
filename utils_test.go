package main

import (
	"reflect"
	"testing"
)

func TestIff(t *testing.T) {
	if got := iff(true, "yes", "no"); got != "yes" {
		t.Errorf("iff(true) = %v; want 'yes'", got)
	}
	if got := iff(false, 10, 20); got != 20 {
		t.Errorf("iff(false) = %v; want 20", got)
	}
}

func TestGetOrDefault(t *testing.T) {
	m := map[string]any{
		"str":  "hello",
		"num":  42,
		"nil":  nil,
	}

	if got := getOrDefault[string](m, "str", "default"); got != "hello" {
		t.Errorf("got %v, want hello", got)
	}
	if got := getOrDefault[string](m, "missing", "default"); got != "default" {
		t.Errorf("got %v, want default", got)
	}
	if got := getOrDefault[int](m, "num", 0); got != 42 {
		t.Errorf("got %v, want 42", got)
	}
	// Wrong type
	if got := getOrDefault[string](m, "num", "default"); got != "default" {
		t.Errorf("got %v, want default on type mismatch", got)
	}
	// Nil map
	if got := getOrDefault[string](nil, "any", "default"); got != "default" {
		t.Errorf("got %v, want default on nil map", got)
	}
}

func TestFilter(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	evens := filter(nums, func(n int) bool {
		return n%2 == 0
	})
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(evens, expected) {
		t.Errorf("filter() = %v, want %v", evens, expected)
	}

	empty := filter([]string{}, func(s string) bool { return true })
	if len(empty) != 0 {
		t.Errorf("filter() on empty slice = %v, want empty", empty)
	}
}

func TestMapSlice(t *testing.T) {
	nums := []int{1, 2, 3}
	doubled := mapSlice(nums, func(n int) int {
		return n * 2
	})
	expected := []int{2, 4, 6}
	if !reflect.DeepEqual(doubled, expected) {
		t.Errorf("mapSlice() = %v, want %v", doubled, expected)
	}
}
