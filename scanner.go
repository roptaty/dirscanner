package dirscanner

import (
	"errors"
	"os"
	"path/filepath"

	"regexp"
)

type needle struct {
	allowedFilePatterns *[]regexp.Regexp

	excludeFilePatterns *[]regexp.Regexp
}

var scanNeedles = map[string]needle{}

// PathResult is used...
type PathResult struct {
	Identifier string
	Path       string
}

// DirScanner ...
type DirScanner struct {
	scanNeedles map[string]needle
}

//NewScanner initializes a new scanner
func NewScanner() DirScanner {
	scanner := DirScanner{map[string]needle{}}
	return scanner
}

func createRegexp(patterns []string) (regExps *[]regexp.Regexp, err error) {
	regs := []regexp.Regexp{}

	if patterns == nil {
		return &regs, nil
	}

	for _, pattern := range patterns {
		var re *regexp.Regexp

		if re, err = regexp.Compile(pattern); err != nil {
			return nil, err
		}

		regs = append(regs, *re)
	}

	return &regs, nil
}

// AddNeedle adds a new needle to the dirscanner
func (d *DirScanner) AddNeedle(identifier string, allowedFilePatterns []string, excludeFilePatterns []string) (err error) {
	var allowedRegex *[]regexp.Regexp
	var excludedRegex *[]regexp.Regexp

	if allowedRegex, err = createRegexp(allowedFilePatterns); err != nil {
		return err
	}

	if excludedRegex, err = createRegexp(excludeFilePatterns); err != nil {
		return err
	}

	d.scanNeedles[identifier] = needle{allowedRegex, excludedRegex}
	return nil
}

//Scan asdfdsf
func (d *DirScanner) Scan(srcPath string) (result *[]PathResult, err error) {
	if len(d.scanNeedles) == 0 {
		return nil, errors.New("No needles are defined")
	}

	if dir, err := os.Stat(srcPath); err != nil {
		return nil, err
	} else if !dir.IsDir() {
		return nil, errors.New("Must be a directory")
	}

	result, err = d.internalScan(srcPath)

	return result, err
}

func (d *DirScanner) internalScan(srcPath string) (result *[]PathResult, err error) {
	var dirEntry *os.File
	var files []os.FileInfo

	results := []PathResult{}

	if dirEntry, err = os.Open(srcPath); err != nil {
		return nil, err
	}

	if files, err = dirEntry.Readdir(-1); err != nil {
		return nil, err
	}

	dirEntry.Close()

	for _, file := range files {
		completeFileName := filepath.Join(srcPath, file.Name())
		identifier, match, exclude, err := d.checkNeedles(completeFileName)

		if err != nil {
			return nil, err
		}

		if exclude {
			continue
		}

		if match {
			results = append(results, PathResult{identifier, completeFileName})
		}

		if fileStat, err := os.Stat(completeFileName); err == nil && fileStat.IsDir() {
			var newResults *[]PathResult

			if newResults, err = d.internalScan(completeFileName); err != nil {
				return nil, err
			}

			results = append(results, *newResults...)

		}
	}

	return &results, nil
}

func (d *DirScanner) checkNeedles(filePath string) (identifier string, match bool, exclude bool, err error) {
	for identifier, nedle := range d.scanNeedles {
		for _, excludeRegex := range *nedle.excludeFilePatterns {
			if exclude = excludeRegex.MatchString(filePath); exclude {
				return identifier, false, exclude, nil
			}

			for _, allowedRegex := range *nedle.allowedFilePatterns {
				if match = allowedRegex.MatchString(filePath); match {
					return identifier, match, false, nil
				}
			}
		}
	}
	return identifier, false, false, nil
}
