package main

type Scanner struct {
	input  string
	tokens []Token

	start   int
	current int
	line    int
}

func (s *Scanner) isAtEnd() bool {
	return s.current > len(s.input)
}

func (s *Scanner) scanToken() {

}

func (s *Scanner) advance() rune {
	out := s.input[s.current]

	s.current++

	return out
}
