package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/rosteroster/memdb/db"
)

func main() {
	// this should accept a single argument on the command line
	// pointing to the path of the input file
	if len(os.Args) < 2 {
		fmt.Println("No filepath specified")
		return
	}
	infile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Failed to open file: %s", os.Args[1])
		return
	}
	defer infile.Close()
	fileScanner := bufio.NewScanner(infile)
	myDB := db.New()

	actions := []string{}
	for fileScanner.Scan() {
		actions = append(actions, fileScanner.Text())
	}

	transaction, err := myDB.NewTransaction(actions)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	myDB.Do(transaction)

	return
}
