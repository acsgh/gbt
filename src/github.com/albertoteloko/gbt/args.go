package main

import (
	"github.com/albertoteloko/gbt/log"
	"os"
)

type Args []string

var args = Args(os.Args[1:])

func configureLogs() {
	if isDebug() {
		log.Level(log.DEBUG)
	} else {
		log.Level(log.INFO)
	}
}

func isDebug() bool {
	return args.indexOfValue("-d") > -1
}

func (args Args) indexOfValue(value string) int {
	return args.indexOfPredicate(func(arg string) bool {
		return arg == value
	})
}

func (args Args) indexOfPredicate(predicate func(string) bool) int {
	for index, element := range args {
		if predicate(element) {
			return index
		}
	}
	return -1
}