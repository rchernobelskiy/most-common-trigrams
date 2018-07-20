package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

const itemsToReturn = 100

type trigram struct {
	First  string
	Second string
	Third  string
}
type sortableTrigram struct {
	*trigram
	Count int
}

// standard minheap implementation using golang heap interface
type trigramHeap []sortableTrigram

func (h trigramHeap) Len() int           { return len(h) }
func (h trigramHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h trigramHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *trigramHeap) Push(x interface{}) {
	*h = append(*h, x.(sortableTrigram))
}
func (h *trigramHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// lowercase and remove punctuation and line breaks
func cleanWord(word string) string {
	return strings.ToLower(strings.Map(func(r rune) rune {
		if r == '\r' || r == '\n' || unicode.IsPunct(r) {
			return -1
		}
		return r
	}, word))
}

func processInput(counter map[trigram]int, input io.Reader) {
	// split up stream into words by spaces, ignoring newlines
	scanner := bufio.NewScanner(input)
	scanner.Split(scanWords)
	var first, second, third string

	// store counts for every trigram
	for second == "" && scanner.Scan() {
		second = cleanWord(scanner.Text())
	}
	for third == "" && scanner.Scan() {
		third = cleanWord(scanner.Text())
	}
	for scanner.Scan() {
		word := cleanWord(scanner.Text())
		if word == "" {
			continue
		}
		first, second, third = second, third, word
		counter[trigram{first, second, third}]++
	}
}

func computeMostCommon(counter map[trigram]int, maxItems int) []sortableTrigram {
	minHeap := &trigramHeap{}
	if maxItems <= 0 {
		return *minHeap
	}
	heap.Init(minHeap)
	// construct bounded min heap
	for k, v := range counter {
		if minHeap.Len() == maxItems && v > (*minHeap)[0].Count {
			heap.Pop(minHeap)
		}
		if minHeap.Len() < maxItems || v > (*minHeap)[0].Count {
			tmp := k
			heap.Push(minHeap, sortableTrigram{&tmp, v})
		}
	}
	// reverse comparison order and sort
	sort.Sort(sort.Reverse(minHeap))
	return *minHeap
}

func main() {
	counter := make(map[trigram]int)

	// process input, either args or stdin
	args := os.Args[1:]
	if len(args) != 0 {
		for _, filename := range args {
			file, err := os.Open(filename)
			if err != nil {
				log.Println("Cannot open file: " + filename)
				os.Exit(1)
			}
			processInput(counter, file)
		}
	} else {
		processInput(counter, os.Stdin)
	}

	sortedTrigrams := computeMostCommon(counter, itemsToReturn)

	// output results
	delimiter := ""
	for i, item := range sortedTrigrams {
		if i == 1 {
			delimiter = ", "
		}
		fmt.Printf("%s%d - %s %s %s", delimiter, item.Count, item.First, item.Second, item.Third)
	}
	if len(sortedTrigrams) != 0 {
		fmt.Println()
	}
}
