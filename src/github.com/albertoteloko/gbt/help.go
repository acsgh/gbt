package main

import (
	"fmt"
)

func isHelp() bool {
	return args.indexOfValue("-help") > -1
}

func printHelp() {
	fmt.Printf("GBT %v\n\n", version)
	fmt.Printf("Usage: gbt [flags] [tasks]\n")
	fmt.Printf("\tFlags:\n")
	fmt.Printf("\t\t-d - Debug the output\n")
	fmt.Printf("\tTasks:\n")
	fmt.Printf("\t\tName\t\t\tDependencies\t\tOrder\n")
	fmt.Printf("\t\tclean\t\t\t[]\t\t\t0\n")
	fmt.Printf("\t\tdependencies\t\t[]\t\t\t1\n")
	fmt.Printf("\t\tbuild\t\t\t[]\t\t\t2\n")
	fmt.Printf("\t\ttest\t\t\t[]\t\t\t3\n")
	fmt.Printf("\t\tbench\t\t\t[]\t\t\t4\n")

}
