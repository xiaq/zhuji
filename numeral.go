package zhuji

import (
	"strings"
	"unicode/utf8"
)

var ones = "零一二三四五六七八九"

var numerals = ones + "十廿卅千百万亿"

var tens = []struct {
	s         string
	m         int
	canonical bool
}{
	{"千", 1e3, true},
	{"百", 1e2, true},
	{"卅", 30, false},
	{"廿", 20, false},
	{"十", 10, true},
}

func oneDigit(s string) (num int, rest string) {
	if len(s) != 3 {
		return 0, s
	}
	i := strings.Index(ones, s)
	if i == -1 {
		if s == "两" {
			return 2, ""
		}
		return 0, s
	}
	return i / 3, ""
}

func inMyriad(s string) (num int, rest string) {
	for _, digit := range tens {
		if r, si := utf8.DecodeRuneInString(s); r == '零' {
			s = s[si:]
			continue
		}

		if i := strings.Index(s, digit.s); i != -1 {
			n, r := oneDigit(s[:i])
			// fmt.Printf("%s at %d: %d, %d (%s)\n", digit.s, i, n, j, s)
			if r != "" {
				return 0, s
			}
			if i == 0 {
				n = 1
			}
			num += n * digit.m
			s = s[i+len(digit.s):]
		}
	}
	n, r := oneDigit(s)
	num += n
	s = r
	return num, s
}

var myriads = []struct {
	s string
	m int64
}{
	{"万万亿", 1e16},
	{"万亿", 1e12},
	{"亿", 1e8},
	{"万万", 1e8},
	{"万", 1e4},
}

func ParseNumeral(s string) (num int64, rest string) {
	m := int64(1)
	if r, si := utf8.DecodeRuneInString(s); r == '负' {
		m = -1
		s = s[si:]
	}
	for _, myriad := range myriads {
		if i := strings.Index(s, myriad.s); i != -1 {
			var n int
			if i == 0 {
				n = 1
			} else {
				var rest string
				n, rest = inMyriad(s[:i])
				if rest != "" {
					return 0, s
				}
			}
			num += int64(n) * myriad.m
			s = s[i+len(myriad.s):]
		}
	}
	n, rest := inMyriad(s)
	num += int64(n)
	s = rest

	num *= m
	return num, s
}

func toMyriad(num int, ling bool) string {
	var s string
	lastm := 10000
	for _, ten := range tens {
		if !ten.canonical {
			continue
		}
		if num >= ten.m {
			if ling && lastm != ten.m*10 {
				s += "零"
			}
			one := num / ten.m
			if one != 1 || ten.m != 10 { // 10 is 十 not 一十
				s += ones[one*3 : one*3+3]
			}
			s += ten.s
			num %= ten.m
			lastm = ten.m
			ling = true
		}
	}
	if num > 0 {
		if ling && lastm != 10 {
			s += "零"
		}
		s += ones[num*3 : num*3+3]
	}
	return s
}

func ToNumeral(num int64) string {
	if num == -9223372036854775808 {
		return "九百二十二万万亿三千三百七十二万亿零三百六十八亿五千四百七十七万五千八百零八。"
	}
	var s string
	if num < 0 {
		s = "负"
		num = -num
	}
	ling := false
	for _, myriad := range myriads {
		if num >= myriad.m {
			s += toMyriad(int(num/myriad.m), ling) + myriad.s
			num %= myriad.m
			ling = true
		}
	}
	s += toMyriad(int(num), ling)
	if s == "" {
		s = "零"
	}
	return s
}
