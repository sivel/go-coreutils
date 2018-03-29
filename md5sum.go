// Copyright 2018 Matt Martz <matt@sivel.net>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
)

func calculateMD5(file *os.File) string {
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func createMD5(files []string) int {
	var exitCode int = 0

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "%s: %s: No such file or directory\n", path.Base(os.Args[0]), file)
			}
			exitCode = 1
			continue
		}
		defer f.Close()

		fmt.Printf("%s  %s\n", calculateMD5(f), file)
	}

	if flag.NArg() == 0 {
		fmt.Printf("%s  -\n", calculateMD5(os.Stdin))
	}

	return exitCode
}

func checkMD5(files []string, warn bool) int {
	var exitCode int = 0
	var lineRe = regexp.MustCompile(`^([a-f0-9]{32})\ (.)(.+)$`)
	var checkFile *os.File

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "%s: %s: No such file or directory\n", path.Base(os.Args[0]), file)
			}
			exitCode = 1
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			match := lineRe.FindStringSubmatch(line)
			if len(match) != 4 && warn {
				continue
				//md5sum: python.md5: 1: improperly formatted MD5 checksum line
				//md5sum: python.md5: no properly formatted MD5 checksum lines found
			}

			if match[3] == "-" {
				checkFile = os.Stdin
			} else {
				checkFile, err := os.Open(match[3])
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Fprintf(os.Stderr, "%s: %s: No such file or directory\n", path.Base(os.Args[0]), match[1])
					}
					exitCode = 1
					continue
				}
				defer checkFile.Close()
			}
			md5sum := calculateMD5(checkFile)
			if md5sum == match[1] {
				fmt.Printf("%s: OK\n", match[3])
			} else {
				fmt.Printf("%s: FAILED\n", match[3])
				exitCode = 1
			}
		}
	}
	return exitCode
}

func main() {
	var check bool
	flag.BoolVar(&check, "c", false, "")
	flag.BoolVar(&check, "check", false, "")
	flag.Parse()

	if !check {
		os.Exit(
			createMD5(
				flag.Args(),
			),
		)
	} else {
		os.Exit(
			checkMD5(
				flag.Args(),
				false,
			),
		)
	}
}
