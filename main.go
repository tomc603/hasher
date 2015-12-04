// hasher - Crawl files and directories, generating sha256 sums. Output is
//          grouped by identical hash to reveal duplicate files.
//
//   Copyright 2015 Tom Cameron
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
//

package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func hashDir(p string, hm map[[32]byte][]string) error {
	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", p, err)
	}
	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return err
	}
	if s.IsDir() {
		// The current item is a directory. Hash its files and subdirectories.
		list, err := f.Readdir(-1)
		if err != nil {
			return err
		}

		for _, d := range list {
			hashDir(filepath.Join(p, d.Name()), hm)
		}
	} else {
		// Hash the current item, it is a file
		h, err := hashFile(f)
		if err != nil {
			return err
		}
		hm[h] = append(hm[h], p)
	}
	return nil
}

func hashFile(f *os.File) ([32]byte, error) {
	var v [32]byte
	hasher := sha256.New()

	r := bufio.NewReader(f)
	if _, err := io.Copy(hasher, r); err != nil {
		return [32]byte{}, err
	}
	copy(v[:], hasher.Sum(nil))

	return v, nil
}

func main() {
	// TODO:
	// Spawn goroutine to listen for items and place them in a map
	// Spawn goroutine workers to hash items from a limited channel
	// Place file paths into a queue for processing

	hm := make(map[[32]byte][]string)

	for _, p := range os.Args[1:] {
		err := hashDir(p, hm)
		if err != nil {
			log.Fatalf("ERROR: %s\n", err)
		}
	}

	for k, v := range hm {
		if len(v) > 1 {
			// Print the sha256 hash
			fmt.Printf("%x\n", k)
			for _, f := range v {
				// Print each matching file path
				fmt.Printf("%s\n", f)
			}
			// Extra blank line for clean formatting
			fmt.Print("\n")
		}
	}
}
