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
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func hashDir(p string) {
}

func hashFile(p string) string {
	hasher := sha256.New()

	f, err := os.Open(p)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", p, err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	if _, err := io.Copy(hasher, r); err != nil {
		log.Fatalf("Error hashing %s: %s\n", p, err)
	}

	return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
	// Parse CLI arguments for files and directories to hash.
	// Use a goroutine worker pool to calculate hashes concurrently
	// var hm map[hash.Hash][]string
	hm := make(map[string][]string)

	for _, p := range os.Args[1:] {
		h := hashFile(p)
		hm[h] = append(hm[h], p)
	}

	for k, v := range hm {
		fmt.Printf("%s\n", k)
		for _, f := range v {
			fmt.Printf(" * %s\n", f)
		}
	}
}
