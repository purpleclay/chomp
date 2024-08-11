package chomp

// A mapped combinator is a function capable of converting the output from a [Combinator]
// into any given type. Upon success, it will return the unparsed text, along with the
// mapped value. All combinators are strict and must parse its input. Any failure to do
// so should raise a [CombinatorParseError]. It is designed for exclusive use by the
// [Map] function
type mappedCombinator[S any, T Result] func(string) (string, S, error)

// Map the result of a [Combinator] to any other type
//
//	chomp.Map(
//		chomp.While(chomp.IsDigit),
//		func (in string) int { return len(in) })("123456")
//	// ("", 6, nil)
func Map[S any, T Result](c Combinator[T], mapper func(in T) S) mappedCombinator[S, T] {
	return func(s string) (string, S, error) {
		var mapped S

		rem, out, err := c(s)
		if err != nil {
			return rem, mapped, err
		}

		mapped = mapper(out)
		return rem, mapped, nil
	}
}
