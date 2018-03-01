package ahocorasick

import (
	"container/list"
	//"fmt"
	"unicode/utf8"
)

type trieNode struct {
	count int // 字典敏感词的权重
	fail  *trieNode
	child map[rune]*trieNode // 子节点匹配字符
	index int                // 匹配到的字串在字典中的索引
}

func newTrieNode() *trieNode {
	return &trieNode{
		count: 0,
		fail:  nil,
		child: make(map[rune]*trieNode),
		index: -1,
	}
}

type Matcher struct {
	root *trieNode // 字典根节点
	size int       // 字典总数
	mark []bool
}

func NewMatcher() *Matcher {
	return &Matcher{
		root: newTrieNode(),
		size: 0,
		mark: make([]bool, 0),
	}
}

// initialize the ahocorasick
func (this *Matcher) Build(dictionary []string) {
	for i, _ := range dictionary {
		this.insert(dictionary[i])
	}
	this.build()
	this.mark = make([]bool, this.size)
}

// string match search
// return all strings matched as indexes into the original dictionary
func (this *Matcher) Match(s string) []int {
	curNode := this.root
	this.resetMark()
	var p *trieNode = nil

	ret := make([]int, 0)

	for _, v := range s {
		for curNode.child[v] == nil && curNode != this.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = this.root
		}

		p = curNode
		for p != this.root && p.count > 0 && !this.mark[p.index] {
			this.mark[p.index] = true
			for i := 0; i < p.count; i++ {
				ret = append(ret, p.index)
			}
			p = p.fail
		}
	}

	return ret
}

// 替换字串中的敏感词
func (this *Matcher) Replace(s string) (newstr string) {
	if len(s) < 1 {
		return s
	}
	node := this.root
	key := []rune(s)
	var chars []rune = nil
	slen := len(key)
	for i := 0; i < slen; i++ {
		if _, exists := node.child[key[i]]; exists {
			node = node.child[key[i]]
			for j := i + 1; j < slen; j++ {
				if _, exists := node.child[key[j]]; exists {
					node = node.child[key[j]]
					// 判断是否在屏蔽词末尾节点
					if node.index >= 0 {
						if chars == nil {
							chars = key // 查找到匹配的字串 将原串初始化到chars中
						}
						// 替换需要屏蔽的字符
						for t := i; t <= j; t++ {
							c, _ := utf8.DecodeRuneInString("*")
							chars[t] = c
						}
						i = j
						node = this.root // 回溯到根节点搜索
						break
					}
				} else {
					i = j
					node = this.root // 回溯到根节点搜索
					break
				}
			}
			node = this.root // 回溯到根结点搜索
		}
	}
	if chars == nil {
		return s
	} else {
		return string(chars)
	}
}

// 替换字串中的敏感词为给出的字串
func (this *Matcher) RepaceStr(s string, replace string) string {
	if len(s) < 1 {
		return s
	}
	node := this.root
	key := []rune(s)
	var chars []rune = nil

	slen := len(key)
	for i := 0; i < slen; i++ {
		//fmt.Printf("i:%v\n", key[i])
		if _, exists := node.child[key[i]]; exists {
			node = node.child[key[i]]
			j := i + 1
			haskey := false
			for ; j < slen; j++ {
				if _, exists := node.child[key[j]]; exists {
					node = node.child[key[j]]
					// 判断是否在屏蔽词末尾节点
					if node.index >= 0 {
						chars = append(chars, []rune(replace)...)
						i = j
						node = this.root // 回溯到根节点搜索
						haskey = true
						break
					}
				} else {
					break
				}
			}
			if !haskey {
				for t := i; t <= j && t < slen; t++ {
					chars = append(chars, key[t])
					//fmt.Printf("append1:%v i:%v\n", string(key[t]), t)
				}
				i = j
			}
			node = this.root // 回溯到根节点搜索

		} else {
			chars = append(chars, key[i])
			//fmt.Printf("append2:%v i:%v\n", key[i], i)
		}
	}

	return string(chars)
}

// just return the number of len(Match(s))
func (this *Matcher) GetMatchResultSize(s string) int {

	curNode := this.root
	this.resetMark()
	var p *trieNode = nil

	num := 0

	for _, v := range s {
		for curNode.child[v] == nil && curNode != this.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = this.root
		}

		p = curNode
		for p != this.root && p.count > 0 && !this.mark[p.index] {
			this.mark[p.index] = true
			num += p.count
			p = p.fail
		}
	}

	return num
}

func (this *Matcher) build() {
	ll := list.New()
	ll.PushBack(this.root)
	for ll.Len() > 0 {
		temp := ll.Remove(ll.Front()).(*trieNode)
		var p *trieNode = nil

		for i, v := range temp.child {
			if temp == this.root {
				v.fail = this.root
			} else {
				p = temp.fail
				for p != nil {
					if p.child[i] != nil {
						v.fail = p.child[i]
						break
					}
					p = p.fail
				}
				if p == nil {
					v.fail = this.root
				}
			}
			ll.PushBack(v)
		}
	}
}

func (this *Matcher) insert(s string) {
	curNode := this.root
	for _, v := range s {
		if curNode.child[v] == nil {
			curNode.child[v] = newTrieNode()
		}
		curNode = curNode.child[v]
	}
	curNode.count++
	curNode.index = this.size
	this.size++
}

func (this *Matcher) resetMark() {
	for i := 0; i < this.size; i++ {
		this.mark[i] = false
	}
}
