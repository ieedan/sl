package args

import "slices"

// extremely primitive argument parser that allows you to single or double quote strings
func Parse(args string) []string {
	index := 0

	parsed := []string{}

	isAtEnd := func() bool {
		return index >= len(args)
	}

	isOneOf := func(s rune, chars ...rune) bool {
		return slices.Index(chars, s) != -1
	}

	for !isAtEnd() {
		char := rune(args[index])
		switch char {
		case '"':
			index++
			start := index

			for !isAtEnd() && args[index] != '"' {
				index++
			}

			if isAtEnd() {
				parsed = append(parsed, args[start:])
			} else {
				parsed = append(parsed, args[start:index])
			}
		case '\'':
			index++
			start := index

			for !isAtEnd() && args[index] != '\'' {
				index++
			}

			if isAtEnd() {
				parsed = append(parsed, args[start:])
			} else {
				parsed = append(parsed, args[start:index])
			}
		case ' ':
		case '\n':
		case '\r':
			index++
			continue
		default:
			start := index

			for !isAtEnd() && !isOneOf(rune(args[index]), '\r', '\n', ' ') {
				index++
			}

			if isAtEnd() {
				parsed = append(parsed, args[start:])
			} else {
				parsed = append(parsed, args[start:index])
			}
		}

		index++
	}

	return parsed
}
