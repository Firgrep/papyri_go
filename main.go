package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const MARKDOWN_TEMPLATE = `---
title: "%s"
pubDate: %s
description: ""
author: "%s"
tags: ["%s"]
---

# %s
`

type Config struct {
	Author       string   `json:"author"`
	DefaultTags  []string `json:"defaultTags"`
	OutputFolder string   `json:"outputFolder"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a filename.")
		os.Exit(1)
	}

	config, err := readConfig("config.json")
	if err != nil {
		log.Fatal("Error reading config file:", err)
	}

	title := os.Args[1]

	fmt.Println("Author:", config.Author)
	fmt.Println("Default Tags:", config.DefaultTags)
	fmt.Println("Output Folder:", config.OutputFolder)

	fileName := strings.ReplaceAll(strings.ToLower(title), " ", "-")
	fileName = strings.ReplaceAll(fileName, "'", "%27")
	fileName = strings.ReplaceAll(fileName, ",", "%2C")

	folderPath := config.OutputFolder

	if err := os.MkdirAll(config.OutputFolder, 0755); err != nil {
		log.Fatal("Error creating directory:", err)
	}

	fileWithDate := fmt.Sprintf("%s/%s_%s.md", folderPath, time.Now().Format("2006-01-02"), fileName)

	file, err := os.Create(fileWithDate)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	capitalizedTitle := cases.Title(language.English, cases.Compact).String(title)

	_, err = fmt.Fprintf(
		file,
		MARKDOWN_TEMPLATE,
		title,
		time.Now().Format("2006-01-02T15:04:05"),
		config.Author,
		strings.Join(config.DefaultTags, "\", \""),
		capitalizedTitle)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(1)
	}

	fmt.Println("Note created:", fileWithDate)
}

func readConfig(filePath string) (Config, error) {
	var config Config

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return config, err
	}

	if config.OutputFolder == "" {
		return config, fmt.Errorf(" outputFolder is empty in the config file")
	}

	return config, nil
}
