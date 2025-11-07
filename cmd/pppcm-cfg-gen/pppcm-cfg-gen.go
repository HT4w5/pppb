package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/HT4w5/pppcm/internal/config"
)

func printUsage() {
	fmt.Printf("Usage: pppcm-cfg-gen <tag-prefix> <ttyname-prefix> <user> <password> <ifname-prefix> <link-count>\n")
}

func main() {
	if len(os.Args) != 7 {
		printUsage()
		return
	}

	linkCount, err := strconv.Atoi(os.Args[6])
	if err != nil || linkCount <= 0 {
		fmt.Printf("Invalid link-count %s: %v\n", os.Args[6], err)
		return
	}

	cfg := config.New()
	cfg.Daemon.Expected = linkCount

	for i := range linkCount {
		cfg.Links = append(cfg.Links, &config.LinkConfig{
			Tag:      fmt.Sprintf("%s%d", os.Args[1], i),
			TTYName:  fmt.Sprintf("%s%d", os.Args[2], i),
			User:     os.Args[3],
			Password: os.Args[4],
			IFName:   fmt.Sprintf("%s%d", os.Args[5], i),
		})
	}

	jsonBytes, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		fmt.Printf("Failed to build config: %v", err)
		return
	}

	fmt.Println(string(jsonBytes))
}
