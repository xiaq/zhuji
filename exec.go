package zhuji

import (
	"bytes"
	"fmt"
)

var (
	defs  = map[string][]*Word{}
	stack []int64
	conds []bool
)

func exec(words []*Word) bool {
	if len(words) > 1 && words[1].Name == "者" {
		defs[words[0].Name] = words[2:]
		return true
	}

	nconds := len(conds)
	defer func() { conds = conds[:nconds] }()

	for i := range words {
		word := words[i]
		if f := find(word); f != nil {
			if shoulddo() {
				f()
			}
		} else if f, ok := controls[word.Name]; ok {
			ok = f()
			if !ok {
				return false
			}
		} else {
			fmt.Printf("无「%s」。\n", word.Name)
		}
	}
	return true
}

func shoulddo() bool {
	for _, cond := range conds {
		if !cond {
			return false
		}
	}
	return true
}

func find(word *Word) func() {
	if def, ok := defs[word.Name]; ok {
		// Defined word.
		return func() { exec(def) }
	} else if builtin, ok := builtins[word.Name]; ok {
		// Builtin word.
		return builtin
	} else if word.isNumeral() {
		// Numeric word.
		return func() {
			num, rest := ParseNumeral(word.Name)
			if rest != "" {
				fmt.Printf("「%s」似数非数。\n", word.Name)
			} else {
				push(num)
			}
		}
	}
	return nil
}

func ExecArticle(a Article) {
	for _, s := range a.Sentences {
		exec(s.Words)
	}
}

func printWords(words []*Word) string {
	var b bytes.Buffer
	for i, w := range words {
		if i > 0 {
			b.WriteRune('、')
		}
		b.WriteString(w.Name)
	}
	return b.String()
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

var controls = map[string]func() bool{
	"则": 则, "非": 非, "毕": 毕,
}

func 则() bool {
	if len(stack) == 0 {
		fmt.Println("栈上无元，不知其可。")
		return false
	}
	conds = append(conds, pop() != 0)
	return true
}

func 非() bool {
	if len(conds) == 0 {
		fmt.Println("无则有非，不知其可。")
		return false
	}
	c := &conds[len(conds)-1]
	*c = !*c
	return true
}

func 毕() bool {
	if len(conds) == 0 {
		fmt.Println("无则有毕，不知其可。")
		return false
	}
	conds = conds[:len(conds)-1]
	return true
}

var builtins = map[string]func(){
	// Arithmetic operations.
	"加": 加, "和": 加, "减": 减, "负": 负, "乘": 乘, "除": 除,
	"次方": 次方,
	// Arithmetic predicates.
	"等于": 等于, "大于": 大于, "小于": 小于,
	// Stack operations.
	"复": 自, "易": 易, "自": 自, "弃": 弃,
}

func top() int64 {
	return stack[len(stack)-1]
}

func pop() int64 {
	i := top()
	stack = stack[:len(stack)-1]
	return i
}

func push(i int64) {
	stack = append(stack, i)
}

func pushBool(b bool) {
	if b {
		push(-1)
	} else {
		push(0)
	}
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

func 等于() {
	if atleast(2) {
		pushBool(pop() == pop())
	}
}

func 小于() {
	if atleast(2) {
		pushBool(pop() > pop())
	}
}

func 大于() {
	if atleast(2) {
		pushBool(pop() < pop())
	}
}

func 易() {
	if atleast(2) {
		a := pop()
		b := pop()
		push(a)
		push(b)
	}
}

func 自() {
	if atleast(1) {
		push(top())
	}
}

func 弃() {
	if atleast(1) {
		pop()
	}
}
