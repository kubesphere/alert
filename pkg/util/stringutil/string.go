// Copyright 2019 The KubeSphere Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package stringutil

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// Creates an slice of slice values not included in the other given slice.
func Diff(base, exclude []string) (result []string) {
	excludeMap := make(map[string]bool)
	for _, s := range exclude {
		excludeMap[s] = true
	}
	for _, s := range base {
		if !excludeMap[s] {
			result = append(result, s)
		}
	}
	return result
}

func Unique(ss []string) (result []string) {
	smap := make(map[string]bool)
	for _, s := range ss {
		smap[s] = true
	}
	for s := range smap {
		result = append(result, s)
	}
	return result
}

func FindString(array []string, str string) int {
	for index, s := range array {
		if str == s {
			return index
		}
	}
	return -1
}

func StringIn(str string, array []string) bool {
	return FindString(array, str) > -1
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func Contains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func SimplifyStringList(s []string) []string {
	b := s[:0]
	for _, x := range s {
		if x := SimplifyString(x); x != "" {
			b = append(b, x)
		}
	}
	return b
}

// "\ta  b  c" => "a b c"
func SimplifyString(s string) string {
	return reMoreSpace.ReplaceAllString(strings.TrimSpace(s), " ")
}

var reMoreSpace = regexp.MustCompile(`\s+`)
