package conf

import (
	"slices"
	"strings"
	"unicode"
)

var keyConverter = DefaultKeyConvertFunc

func SetKeyConverter(conv func(string) string) {
	keyConverter = conv
}

func DisableKeyConverter() {
	keyConverter = func(key string) string { return key }
}

func GetKeyConverter() func(string) string {
	return keyConverter
}

func DefaultKeyConvertFunc(key string) string {
	replace := []string{" ", "-", ".", "+", "~", "\t"}

	var (
		splitter    []int
		upperStrick int
	)

	for i := range key {
		if unicode.IsLetter(rune(key[i])) {
			if unicode.IsUpper(rune(key[i])) {
				upperStrick++
			} else {
				upperStrick = 0
			}
		} else {
			upperStrick = 0
		}
		if i == 0 {
			continue
		}
		if unicode.IsUpper(rune(key[i])) && unicode.IsLower(rune(key[i-1])) {
			splitter = append(splitter, i)
			continue
		}
		if i+1 != len(key) && upperStrick > 0 && unicode.IsLower(rune(key[i+1])) {
			splitter = append(splitter, i)
		}
	}

	slices.Sort(splitter)
	slices.Reverse(splitter)

	for _, s := range splitter {
		key = key[:s] + "_" + key[s:]
	}

	key = strings.ToLower(strings.TrimSpace(key))
	for _, r := range replace {
		key = strings.ReplaceAll(key, r, "_")
	}

	return key
}
