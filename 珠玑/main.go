package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"github.com/xiaq/zhuji"
)

func main() {
	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("珠玑> ")
		line, err := stdin.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		if !utf8.ValidString(line) {
			fmt.Println("not UTF-8")
			continue
		}
		article := zhuji.Parse(line[:len(line)-1])
		// fmt.Println(article)
		zhuji.ExecArticle(article)
		zhuji.ShowIfNonEmpty()
	}
}
