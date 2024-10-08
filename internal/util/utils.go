package util

import (
	"io"
	"log"
	"strings"
)

func PadLeftMin(s string, min int) string {
	remaining := min - len(s)

	if remaining == 0 {
		return s
	}

	if remaining < 0 {
		log.Fatal("The minimum must always be greater than or equal to the length of the string!")
	}

	return strings.Repeat(" ", remaining) + s
}

func PadRightMin(s string, min int) string {
	remaining := min - len(s)

	if remaining == 0 {
		return s
	}

	if remaining < 0 {
		log.Fatalf("The minimum must always be greater than or equal to the length of the string! %v was less than %v", min, len(s))
	}

	return s + strings.Repeat(" ", remaining)
}

func MinLength(strings *[]string) int {
	min := 0

	for _, str := range *strings {
		if len(str) > min {
			min = len(str)
		}
	}

	return min
}

func Map[T any, V any](items *[]T, fn func(T, int) V) []V {
	var mapped []V

	for i, item := range *items {
		mapped = append(mapped, fn(item, i))
	}

	return mapped
}

func MapToSlice[K comparable, V any, T any](m *map[K]V, fn func(K, V) T) []T {
	var slice []T

	for k, v := range *m {
		slice = append(slice, fn(k, v))
	}

	return slice
}

func LPad(str string, padding int) string {
	return strings.Repeat(" ", padding) + str
}

func IsCancel(err error) bool {
	if err == nil {
		return false
	}

	if err == io.EOF {
		return true
	}

	return false
}
