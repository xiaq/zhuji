package zhuji

import "fmt"

var (
	defs  = map[string][]string{}
	stack []int
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
			} else {
				fmt.Printf("无%s\n", word)
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
	"零": 零, "一": 一, "二": 二, "三": 三, "四": 四, "五": 五, "六": 六, "七": 七, "八": 八, "九": 九, "十": 十, "乘": 乘, "自": 自, "弃": 弃,
}

func first() int {
	return stack[len(stack)-1]
}

func second() int {
	return stack[len(stack)-2]
}

func pop() int {
	i := first()
	stack = stack[:len(stack)-1]
	return i
}

func push(i int) {
	stack = append(stack, i)
}

func 零() { push(0) }
func 一() { push(1) }
func 二() { push(2) }
func 三() { push(3) }
func 四() { push(4) }
func 五() { push(5) }
func 六() { push(6) }
func 七() { push(7) }
func 八() { push(8) }
func 九() { push(9) }
func 十() { push(10) }

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
