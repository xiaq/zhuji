package zhuji

import "fmt"

var (
	defs  = map[string][]*Word{}
	stack []int64
)

func exec(words []*Word) {
	if len(words) > 1 && words[1].Name == "者" {
		defs[words[0].Name] = words[2:]
	} else {
		for _, word := range words {
			if def, ok := defs[word.Name]; ok {
				exec(def)
			} else if builtin, ok := builtins[word.Name]; ok {
				builtin()
			} else if word.isNumeral() {
				num, rest := ParseNumeral(word.Name)
				if rest != "" {
					fmt.Printf("「%s」似数非数。\n", word.Name)
				} else {
					push(num)
				}
			} else {
				fmt.Printf("无「%s」。\n", word.Name)
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
	"加": 加, "和": 加, "减": 减, "负": 负, "乘": 乘, "除": 除,
	"次方": 次方,
	"复":  自, "自": 自, "弃": 弃,
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

func 负() {
	if atleast(1) {
		push(-pop())
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

func pow(b, n int64) int64 {
	if n == 0 {
		return 1
	}
	x := pow(b, n/2)
	x *= x
	if n%2 == 1 {
		x *= b
	}
	return x
}

func 次方() {
	if atleast(2) {
		a := pop()
		b := pop()
		push(pow(b, a))
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
