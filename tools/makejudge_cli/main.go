package main

import (
	"os"
	"share-cat/Nepeta-cataria/pkg/makejudge"
)

func main() {
	makejudge.RunCLI(os.Args[1:])
}
