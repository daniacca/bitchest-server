package parser

import (
	"errors"
	"strings"
	"unicode"
)

func Tokenize(input string) ([]string, error) {
	tokens := []string{}
	var current strings.Builder
	var inQuote bool
	var quoteChar rune // can be ' or "
	
	escaped := false
	previousSpace := false

	for _, r := range input {
		switch {
			case escaped: // if the current char is escaped, add it to the current token
				current.WriteRune(r)
				escaped = false // reset the escaped flag
				previousSpace = false // reset the previous space flag, if it was true
			case r == '\\':
				escaped = true // escape the next char
			case r == '"' || r == '\'': // we hit a quote
				if inQuote { // if we're already inside a quote block, check if it's the closing quote
					if r == quoteChar { // it's the closing quote, add the current token and reset the state
						if current.Len() > 0 { // add the current token if it's not empty
							tokens = append(tokens, strings.TrimSpace(current.String()))
						}
						current.Reset()
						inQuote = false
						previousSpace = false // reset the previous space flag, if it was true
					} else { // different quote inside string, keep it
						current.WriteRune(r)
					}
				} else if !previousSpace { // we're not in a quote block, but the previous char was not a space, so add the quote to the current token
					current.WriteRune(r)
				} else { // we're not in a quote block, and the previous char was a space, so start a new quote block
					inQuote = true
					quoteChar = r
					previousSpace = false
				}
			case unicode.IsSpace(r):
				if inQuote { // keep spaces inside the same token
					current.WriteRune(r)
				} else if !previousSpace { // we're not in a quote and the previous char was not a space, so close the current token
					if current.Len() > 0 { // add the current token if it's not empty
						tokens = append(tokens, strings.TrimSpace(current.String()))
					}
					current.Reset()
				}
				previousSpace = true
			default:
				current.WriteRune(r)
				previousSpace = false
		}
	}

	// malformed input: we left characters and there is an unterminated quote, return error
	if inQuote {
		return nil, errors.New("unterminated quote")
	}

	// fallback: we had some characters left in the current token, add them
	if current.Len() > 0 {
		tokens = append(tokens, strings.TrimSpace(current.String()))
	}

	return tokens, nil
}