package main

import "strings"

func contains(s []string, e string) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func distinct[T comparable](s []T) []T {
	var zero T
	result := make([]T, 0, len(s))
	temp := map[T]struct{}{}
	for _, item := range s {
		if item == zero {
			continue
		}
		if _, ok := temp[item]; !ok {
			result = append(result, item)
			temp[item] = struct{}{}
		}
	}
	return result
}

func split(s string) []string {
	var list []string
	for _, v := range strings.Split(s, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		list = append(list, v)
	}
	return list
}
