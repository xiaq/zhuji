package zhuji

import (
	"strings"
	"unicode/utf8"
)

type parser struct {
	source string
	pos    int
}

func newParser(source string) *parser {
	return &parser{source, 0}
}

// Parser primaries.

func (p *parser) eof() bool {
	return p.pos == len(p.source)
}

func (p *parser) rest() string {
	return p.source[p.pos:]
}

func (p *parser) next() rune {
	r, size := utf8.DecodeRuneInString(p.rest())
	p.pos += size
	return r
}

func (p *parser) peek() rune {
	r, _ := utf8.DecodeRuneInString(p.rest())
	return r
}

func (p *parser) from(i int) string {
	return p.source[i:p.pos]
}

func (p *parser) to(i int) string {
	rest := p.rest()
	if i == -1 {
		return rest
	}
	return rest[:i]
}

func (p *parser) uptoFunc(f func(rune) bool) string {
	return p.to(strings.IndexFunc(p.rest(), f))
}

func (p *parser) uptoAny(chars string) string {
	return p.to(strings.IndexAny(p.rest(), chars))
}

func in(r rune, s string) bool {
	return strings.ContainsRune(s, r)
}

// Parsing functions.

type Word struct {
	Name string
}

func (w *Word) isNumeral() bool {
	r, _ := utf8.DecodeRuneInString(w.Name)
	return in(r, numerals)
}

var keywords = "者自"

func (p *parser) word() (w *Word) {
	begin := p.pos

	w = &Word{}
	defer func() {
		w.Name = p.from(begin)
	}()

	if r := p.next(); in(r, keywords) {
		return
	} else if in(r, numerals) {
		for in(p.peek(), numerals) {
			p.next()
		}
		return
	}

	for !p.eof() {
		r := p.peek()

		if in(r, numerals) || in(r, keywords) || in(r, "、，。也\n") {
			return
		}
		p.next()
	}
	return
}

type Sentence struct {
	Words []*Word
}

func (p *parser) sentence() Sentence {
	s := Sentence{}
	jux := false
	for !p.eof() {
		word := p.word()
		s.Words = append(s.Words, word)
		if jux && word.isNumeral() {
			// Numeral juxtaposes another word: swap. This enables infix
			// notation.
			n := len(s.Words)
			s.Words[n-1], s.Words[n-2] = s.Words[n-2], s.Words[n-1]
		}
		jux = true
		for !p.eof() {
			r := p.peek()
			if in(r, "、，") {
				p.next()
				jux = false
			} else if in(r, "。也\n") {
				p.next()
				for in(p.peek(), "。也\n") {
					p.next()
				}
				return s
			} else {
				break
			}
		}
	}
	return s
}

type Article struct {
	Sentences []Sentence
}

func (p *parser) article() Article {
	a := Article{}
	for !p.eof() {
		a.Sentences = append(a.Sentences, p.sentence())
	}
	return a
}

func Parse(source string) Article {
	return newParser(source).article()
}
