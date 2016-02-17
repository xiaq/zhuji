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
					fmt.Println("「%s」似数非数。", word)
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
	"乘": 乘, "自": 自, "弃": 弃,
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

func 自() {
	if len(stack) < 1 {
		fmt.Println("无元。")
	} else {
		push(first())
	}
}

func 乘() {
	if len(stack) < 2 {
		fmt.Println("无二元。")
	} else {
		push(pop() * pop())
	}
}

func 弃() {
	if len(stack) < 1 {
		fmt.Println("无元。")
	} else {
		pop()
	}
}
