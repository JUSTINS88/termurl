package main

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/pflag"
)

const version = "1.0.0"

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	fmt.Printf("Welcome To ZAPS. This version is %s\n", version)

	slug := pflag.StringP("slug", "s", "", "Custom slug for your shorten link")
	pflag.Parse()

	originalLink := pflag.Arg(0)

	err := sendLinkToCreateLink(originalLink, *slug)
	if err != nil {
		return err
	}
	return nil
}
