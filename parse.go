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

var keywords = "者自"

func (p *parser) word() string {
	begin := p.pos

	if r := p.next(); in(r, keywords) {
		return p.from(begin)
	} else if in(r, numerals) {
		for in(p.peek(), numerals) {
			p.next()
		}
		return p.from(begin)
	}

	for !p.eof() {
		r := p.peek()

		if in(r, keywords) || in(r, "、，。也\n") {
			return p.from(begin)
		}
		p.next()
	}
	return p.from(begin)
}

type Sentence struct {
	Words []string
}

func (p *parser) sentence() Sentence {
	s := Sentence{}
	for !p.eof() {
		s.Words = append(s.Words, p.word())
		for !p.eof() {
			r := p.peek()
			if in(r, "、，") {
				p.next()
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
