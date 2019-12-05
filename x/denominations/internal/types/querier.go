package types

import "strings"

// QueryResultSymbol is a payload for a symbols query
type QueryResultSymbol []string

// String implements fmt.Stringer
func (r QueryResultSymbol) String() string {
	return strings.Join(r[:], "\n")
}
