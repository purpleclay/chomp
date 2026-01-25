package chomp

import (
	"strings"
	"testing"
)

const (
	asciiSentence   = "the quick brown fox jumps over the lazy dog"
	unicodeSentence = "素早い茶色のキツネが怠惰な犬を飛び越える"
)

func BenchmarkTag(b *testing.B) {
	tests := []struct {
		name  string
		tag   string
		input string
	}{
		{
			name:  "Ascii",
			tag:   "the quick",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			tag:   "素早い",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Tag(tt.tag)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkTagNoCase(b *testing.B) {
	tests := []struct {
		name  string
		tag   string
		input string
	}{
		{
			name:  "Ascii",
			tag:   "THE QUICK",
			input: asciiSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := TagNoCase(tt.tag)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkChar(b *testing.B) {
	tests := []struct {
		name  string
		char  rune
		input string
	}{
		{
			name:  "Ascii",
			char:  't',
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			char:  '素',
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Char(tt.char)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkAnyChar(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Ascii",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := AnyChar()
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkTake(b *testing.B) {
	tests := []struct {
		name  string
		n     uint
		input string
	}{
		{
			name:  "Ascii",
			n:     10,
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			n:     5,
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Take(tt.n)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkUntil(b *testing.B) {
	tests := []struct {
		name  string
		until string
		input string
	}{
		{
			name:  "Ascii",
			until: "jumps",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			until: "の",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Until(tt.until)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkAny(b *testing.B) {
	tests := []struct {
		name  string
		chars string
		input string
	}{
		{
			name:  "Ascii",
			chars: "the quic",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			chars: "素早い茶色",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Any(tt.chars)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkNot(b *testing.B) {
	tests := []struct {
		name  string
		chars string
		input string
	}{
		{
			name:  "Ascii",
			chars: "xyz",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			chars: "犬猫",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Not(tt.chars)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkOneOf(b *testing.B) {
	tests := []struct {
		name  string
		chars string
		input string
	}{
		{
			name:  "Ascii",
			chars: "abcdefghijklmnopqrst",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			chars: "素早茶色",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := OneOf(tt.chars)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkNoneOf(b *testing.B) {
	tests := []struct {
		name  string
		chars string
		input string
	}{
		{
			name:  "Ascii",
			chars: "xyz",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			chars: "犬猫",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := NoneOf(tt.chars)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkWhile(b *testing.B) {
	tests := []struct {
		name      string
		predicate Predicate
		input     string
	}{
		{
			name:      "Digit",
			predicate: IsDigit,
			input:     "1234567890" + asciiSentence,
		},
		{
			name:      "Letter/Ascii",
			predicate: IsLetter,
			input:     asciiSentence,
		},
		{
			name:      "Letter/Unicode",
			predicate: IsLetter,
			input:     unicodeSentence,
		},
		{
			name:      "Alphanumeric",
			predicate: IsAlphanumeric,
			input:     "Hello123World456",
		},
		{
			name:      "Space",
			predicate: IsSpace,
			input:     "     " + asciiSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := While(tt.predicate)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkWhileNot(b *testing.B) {
	tests := []struct {
		name      string
		predicate Predicate
		input     string
	}{
		{
			name:      "Digit/Ascii",
			predicate: IsDigit,
			input:     asciiSentence + "123",
		},
		{
			name:      "Digit/Unicode",
			predicate: IsDigit,
			input:     unicodeSentence + "123",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := WhileNot(tt.predicate)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkSatisfy(b *testing.B) {
	tests := []struct {
		name  string
		pred  func(rune) bool
		input string
	}{
		{
			name:  "Ascii",
			pred:  func(r rune) bool { return r >= 'a' && r <= 'z' },
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			pred:  func(r rune) bool { return r >= 0x3000 },
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Satisfy(tt.pred)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkPair(b *testing.B) {
	tests := []struct {
		name   string
		first  Combinator[string]
		second Combinator[string]
		input  string
	}{
		{
			name:   "Ascii",
			first:  Tag("the quick"),
			second: Tag(" brown"),
			input:  asciiSentence,
		},
		{
			name:   "Unicode",
			first:  Tag("素早い"),
			second: Tag("茶色"),
			input:  unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Pair(tt.first, tt.second)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkDelimited(b *testing.B) {
	tests := []struct {
		name  string
		open  Combinator[string]
		inner Combinator[string]
		close Combinator[string]
		input string
	}{
		{
			name:  "Parentheses",
			open:  Tag("("),
			inner: Until(")"),
			close: Tag(")"),
			input: "(the quick brown fox)",
		},
		{
			name:  "Quotes",
			open:  Tag("\""),
			inner: Until("\""),
			close: Tag("\""),
			input: "\"the quick brown fox\"",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Delimited(tt.open, tt.inner, tt.close)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkSepPair(b *testing.B) {
	input := "key:value rest"
	parser := SepPair(While(IsLetter), Tag(":"), While(IsLetter))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkAll(b *testing.B) {
	tests := []struct {
		name   string
		parser Combinator[[]string]
		input  string
	}{
		{
			name:   "ThreeTags",
			parser: All(Tag("the"), Tag(" quick"), Tag(" brown")),
			input:  asciiSentence,
		},
		{
			name:   "FiveTags",
			parser: All(Tag("the"), Tag(" "), Tag("quick"), Tag(" "), Tag("brown")),
			input:  asciiSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = tt.parser(tt.input)
			}
		})
	}
}

func BenchmarkOpt(b *testing.B) {
	tests := []struct {
		name  string
		input string
		match bool
	}{
		{
			name:  "Match",
			input: asciiSentence,
			match: true,
		},
		{
			name:  "NoMatch",
			input: "xyz" + asciiSentence,
			match: false,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Opt(Tag("the"))
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkMap(b *testing.B) {
	input := "12345" + asciiSentence
	parser := Map(While(IsDigit), func(s string) int { return len(s) })

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkMany(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Small",
			input: strings.Repeat("ab", 5) + "rest",
		},
		{
			name:  "Medium",
			input: strings.Repeat("ab", 50) + "rest",
		},
		{
			name:  "Large",
			input: strings.Repeat("ab", 500) + "rest",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Many(Tag("ab"))
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkPeek(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Ascii",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Peek(Take(5))
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkFlatten(b *testing.B) {
	input := asciiSentence
	parser := Flatten(All(Tag("the"), Tag(" quick"), Tag(" brown")))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkFirst(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "FirstMatch",
			input: asciiSentence,
		},
		{
			name:  "LastMatch",
			input: "fox jumps over",
		},
	}

	parser := First(Tag("the"), Tag("quick"), Tag("fox"))

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkVerify(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Pass",
			input: "12345" + asciiSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Verify(While(IsDigit), func(s string) bool {
				return len(s) >= 3
			})
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkRecognize(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Ascii",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Recognize(Pair(Take(3), Take(3)))
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkConsumed(b *testing.B) {
	input := asciiSentence
	parser := Consumed(Pair(Tag("the"), Tag(" quick")))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkEof(b *testing.B) {
	parser := Eof()
	input := ""

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkAllConsuming(b *testing.B) {
	input := "the quick brown fox"
	parser := AllConsuming(Tag(input))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkRest(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Ascii",
			input: asciiSentence,
		},
		{
			name:  "Unicode",
			input: unicodeSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Rest()
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkValue(b *testing.B) {
	input := "true" + asciiSentence
	parser := Value(Tag("true"), true)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkCond(b *testing.B) {
	tests := []struct {
		name  string
		cond  bool
		input string
	}{
		{
			name:  "True",
			cond:  true,
			input: asciiSentence,
		},
		{
			name:  "False",
			cond:  false,
			input: asciiSentence,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Cond(tt.cond, Tag("the"))
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkCut(b *testing.B) {
	input := asciiSentence
	parser := Cut(Tag("the"))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkCrlf(b *testing.B) {
	input := "\r\n" + asciiSentence

	parser := Crlf()
	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkEol(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Ascii",
			input: asciiSentence + "\n",
		},
		{
			name:  "Unicode",
			input: unicodeSentence + "\n",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Eol()
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkUntilScaling(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Small",
			input: "start end",
		},
		{
			name:  "Medium",
			input: strings.Repeat("x", 100) + " end",
		},
		{
			name:  "Large",
			input: strings.Repeat("x", 10000) + " end",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := Until("end")
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkWhileScaling(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Small",
			input: "123abc",
		},
		{
			name:  "Medium",
			input: strings.Repeat("1", 100) + "abc",
		},
		{
			name:  "Large",
			input: strings.Repeat("1", 10000) + "abc",
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			parser := While(IsDigit)
			b.ReportAllocs()
			b.SetBytes(int64(len(tt.input)))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = parser(tt.input)
			}
		})
	}
}

func BenchmarkKeyValuePair(b *testing.B) {
	input := "Content-Type: application/json\r\n"
	parser := SepPair(Until(":"), Tag(": "), Until("\r\n"))

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkGitDiffHeader(b *testing.B) {
	input := "@@ -25,10 +30,15 @@"
	parser := Delimited(
		Tag("@@ "),
		SepPair(
			Pair(Tag("-"), While(IsDigit)),
			Tag(" "),
			Pair(Tag("+"), While(IsDigit)),
		),
		Tag(" @@"),
	)

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}

func BenchmarkCSVField(b *testing.B) {
	input := "field_value,next_field,another"
	parser := Not(",\n")

	b.ReportAllocs()
	b.SetBytes(int64(len(input)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = parser(input)
	}
}
