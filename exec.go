package zhuji

import (
	"fmt"
	"unicode/utf8"
)

var (
	defs  = map[string][]string{}
	stack []int64
)

func exec(words []string) {
	if len(words) > 1 && words[1] == "者" {
		defs[words[0]] = words[2:]
	} else {
		for _, word := range words {
			if def, ok := defs[word]; ok {
				exec(def)
			} else if builtin, ok := builtins[word]; ok {
				builtin()
			} else if r, _ := utf8.DecodeRuneInString(word); in(r, numerals) {
				num, size := ParseNumeral(word)
				if size != len(word) {
					fmt.Printf("「%s」似数非数。\n", word)
				} else {
					push(num)
				}
			} else {
				fmt.Printf("无「%s」。\n", word)
			}
		}
	}
}

func ExecArticle(a Article) {
	for _, s := range a.Sentences {
		exec(s.Words)
	}
}

func ShowIfNonEmpty() {
	if len(stack) > 0 {
		for i, n := range stack {
			if i > 0 {
				fmt.Print("、")
			}
			fmt.Print(ToNumeral(int64(n)))
		}
		fmt.Println("。")
	}
}

var builtins = map[string]func(){
	"加": 加, "和": 加, "减": 减, "乘": 乘, "除": 除,
	"自": 自, "弃": 弃,
}

func first() int64 {
	return stack[len(stack)-1]
}

func second() int64 {
	return stack[len(stack)-2]
}

func pop() int64 {
	i := first()
	stack = stack[:len(stack)-1]
	return i
}

func push(i int64) {
	stack = append(stack, i)
}

func atleast(n int) bool {
	if len(stack) < n {
		if n == 1 {
			fmt.Println("无元。")
		} else {
			fmt.Printf("无%s元。\n", ToNumeral(int64(n)))
		}
		return false
	}
	return true
}

func 加() {
	if atleast(2) {
		push(pop() + pop())
	}
}

func 减() {
	if atleast(2) {
		push(-(pop() - pop()))
	}
}

func 乘() {
	if atleast(2) {
		push(pop() * pop())
	}
}

func 除() {
	if atleast(2) {
		a := pop()
		b := pop()
		push(b / a)
	}
}

func 自() {
	if atleast(1) {
		push(first())
	}
}

func 弃() {
	if atleast(1) {
		pop()
	}
}
