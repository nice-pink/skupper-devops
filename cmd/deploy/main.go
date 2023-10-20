package main

import (
	"flag"
	"fmt"

	"github.com/nice-pink/skupper-devops/pkg/deploy"
)

func main() {
	// flags
	src := flag.String("src", "", "Path to src folder.")
	dest := flag.String("dest", "", "Optional: Path to dest folder, where src is going to be copied.")
	gitPushDest := flag.Bool("gitPushDest", false, "Push dest folder?")
	flag.Parse()

	fmt.Println("--------")
	fmt.Println("SRC: " + *src)
	fmt.Println("DEST: " + *dest)
	fmt.Println("--------")
	fmt.Println("")

	// prepare
	deploy.Prepare(*src, *dest, *gitPushDest)
}
