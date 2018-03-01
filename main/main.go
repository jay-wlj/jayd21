package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jayd21/ahocorasick"
)

var (
	ahocorasick_matcher *ahocorasick.Matcher
	comment_replace     string
)

func get_Matcher() *ahocorasick.Matcher {
	if ahocorasick_matcher == nil {
		ahocorasick_matcher = ahocorasick.NewMatcher()
	}
	return ahocorasick_matcher
}

func ReadLine(file string, handler func(string)) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		handler(string(line))
	}
	return
}

// 初始化敏感词库及替换字串
func InitKeyword(file string, replace string) (err error) {
	now := time.Now().Unix()
	comment_replace = replace
	dictionary := []string{}
	ReadLine(file, func(key string) {
		dictionary = append(dictionary, key)
	})
	fmt.Printf("keyWords:\n\"%v\"\n", dictionary)
	get_Matcher().Build(dictionary)

	fmt.Printf("load keyword file success! tick=%v\n", time.Now().Unix()-now)
	return
}

func FilterText(s string) string {
	return get_Matcher().RepaceStr(s, comment_replace)
}

func FindMatch(s string) []int {
	return get_Matcher().Match(s)
}

func main() {
	InitKeyword("/opt/yf-api/conf/keywords_comment.txt", "**")
	str := FilterText("哦啦啦666ga.ga 下来了")
	fmt.Printf(str)
}
