package parser

import "errors"

func TryParseData[T any](parser func(subContext string) (T, error), context string) (T, error) {
	var zero T
	// Iterate over all possible substrings of the context
	for i := 0; i < len(context); i++ {
		for j := i + 1; j <= len(context); j++ {
			subContext := context[i:j]
			// Try to parse the current substring
			if parsed, err := parser(subContext); err == nil {
				// If the parser returns nil, it means the substring is parsable, so return the parsed value
				return parsed, nil
			}
		}
	}
	// If no parsable substring is found, return an error and the zero value of T
	return zero, errors.New("no parsable substring found")
}
