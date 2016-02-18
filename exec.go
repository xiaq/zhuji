package zhuji

import (
	"bytes"
	"fmt"
)

var (
	defs  = map[string][]*Word{}
	stack []int64
)

func exec(words []*Word) (rest []*Word) {
	if len(words) > 1 && words[1].Name == "者" {
		defs[words[0].Name] = words[2:]
		return nil
	}
	for i := range words {
		word := words[i]
		// for _, word := range words {
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
		} else if word.Name == "若" {
			rest := exec(words[i+1:])
			if len(rest) == 0 || rest[0].Name != "则" {
				fmt.Println("有若无则，不知其可。")
				return nil
			}
			if len(stack) == 0 {
				fmt.Println("栈上无元，不知其可。")
				return nil
			}
			rest = rest[1:]
			if pop() != 0 {
				// Take the 则 branch.
				// fmt.Println("则", printWords(rest))
				rest = exec(rest)
				if len(rest) == 0 || (rest[0].Name != "非" && rest[0].Name != "毕") {
					fmt.Println("有若无非无毕，不知其可。")
					return nil
				}
				name := rest[0].Name
				rest = rest[1:]
				if name == "非" {
					for i, w := range rest {
						if w.Name == "毕" {
							return rest[i+1:]
						}
					}
					fmt.Println("有非无毕，不知其可。")
					return nil
				}
			} else {
				// Take the 非 branch if it exists
				for i, w := range rest {
					if w.Name == "毕" {
						return rest[i+1:]
					} else if w.Name == "非" {
						// fmt.Println("非", printWords(rest[i+1:]))
						rest = exec(rest[i+1:])
						if len(rest) == 0 || rest[0].Name != "毕" {
							fmt.Println("有非无毕，不知其可。")
							return nil
						}
						return rest[1:]
					}
				}
				fmt.Println("有若无毕，不知其可。")
			}
		} else if word.Name == "则" || word.Name == "非" || word.Name == "毕" {
			return words[i:]
		} else {
			fmt.Printf("无「%s」。\n", word.Name)
		}
	}
	return nil
}

func ExecArticle(a Article) {
	for _, s := range a.Sentences {
		rest := exec(s.Words)
		if len(rest) > 0 {
			fmt.Println("尚有「%s」，不知其可。", printWords(rest))
		}
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
