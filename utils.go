package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Log struct {
	verbose bool
}

func (l *Log) Println(args ...interface{}) {
	if l.verbose {
		fmt.Println(args...)
	}
}

func (l *Log) Print(args ...interface{}) {
	if l.verbose {
		fmt.Print(args...)
	}
}

func printMatrix(matrix [][]bool) {
	if !verboseFlag {
		return
	}
	for i := 0; i < len(matrix); i++ {
		fmt.Print("[ ")
		for j := 0; j < len(matrix[i]); j++ {
			if matrix[i][j] {
				fmt.Print("T")
			} else {
				fmt.Print("F")
			}
			if j < len(matrix[i])-1 {
				fmt.Print(", ")
			}
		}
		fmt.Println(" ]")
	}
}

// Compares two matrices for equality.
func compQuadraticMatrices(m1 [][]bool, m2 [][]bool) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i := 0; i < len(m1); i++ {
		for j := 0; j < len(m1); j++ {
			if m1[i][j] != m2[i][j] {
				return false
			}
		}
	}
	return true
}

func collectGraphFiles(dir string) []string {
	var graphFiles []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while walking through directory!")
			os.Exit(1)
		}
		if !info.IsDir() {
			graphFiles = append(graphFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while reading directory!")
		os.Exit(1)
	}

	return graphFiles
}

func getTimeMS(start time.Time) float64 {
	return float64(time.Since(start).Nanoseconds()) / 1e6
}
