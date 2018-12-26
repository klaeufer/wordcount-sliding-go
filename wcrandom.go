package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// WordCountConfig The configuration options available.
type WordCountConfig struct {
	uniqueWords, lastNWords, wordsToGenerate, showTop int
}

// newWord() shows how to build an iterator on string
// can make a stateful iterator by putting some variables before the return...

func wordGenerator(wordBase string, uniqueWords int) func() string {
	return func() string {
		n := rand.Intn(uniqueWords)
		return wordBase + strconv.Itoa(n)
	}
}

func showWordCounts(wc map[string]int, showTop int) {
	type WC struct {
		word  string
		count int
	}

	wcSlice := make([]WC, len(wc))
	wcSlice = wcSlice[:0]
	for word, count := range wc {
		wcSlice = append(wcSlice, WC{word, count})
	}
	sort.Slice(wcSlice, func(i, j int) bool {
		return wcSlice[i].count > wcSlice[j].count
	})
	if len(wcSlice) > showTop {
		wcSlice = wcSlice[:showTop]
	}

	//fmt.Printf("wc len %d, show top %d, slice len %d, slice cap %d", len(wc), showTop, len(wcSlice), cap(wcSlice))
	// printing words...

	pretty := make([]string, len(wcSlice))
	pretty = pretty[:0]
	for _, wc := range wcSlice {
		pretty = append(pretty, fmt.Sprintf("%s: %d", wc.word, wc.count))
	}
	fmt.Printf("words { %s }\n", strings.Join(pretty, ", "))
}

func driver(config *WordCountConfig) {
	nextWord := wordGenerator("unique", config.uniqueWords)
	queue := make([]string, config.lastNWords)
	wc := make(map[string]int)

	queue = queue[:0]
	for i := 0; i < config.wordsToGenerate; i++ {
		word := nextWord()
		//fmt.Printf("queue length %d\n", len(queue))
		if len(queue) >= config.lastNWords {
			droppedWord := queue[0]
			queue = queue[1:]
			wc[droppedWord]--
			if wc[droppedWord] <= 0 {
				delete(wc, droppedWord)
			}
		}
		queue = append(queue, word)
		wc[word]++
		showWordCounts(wc, config.showTop)
	}
}

func main() {
	config := WordCountConfig{}
	flag.IntVar(&config.uniqueWords, "unique", config.uniqueWords, "number of unique words")
	flag.IntVar(&config.lastNWords, "last_n_words", config.uniqueWords, "last n words from current word (to count in word cloud)")
	flag.IntVar(&config.wordsToGenerate, "generate", config.wordsToGenerate, "words to generate randomly")
	flag.IntVar(&config.showTop, "show_top", config.showTop, "show top n words")
	flag.Parse()
	fmt.Println(config)
	driver(&config)
}
