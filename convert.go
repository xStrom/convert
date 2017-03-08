// Copyright 2017 Kaur Kuut
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	convert(".\\")
}

func convert(path string) {
	// Get the file list for this directory
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("ReadDir failed: %v\n", err)
		panic("")
	}

	for _, fi := range fileInfos {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".m3u") {
			continue
		}

		f, err := os.Open(path + fi.Name())
		if err != nil {
			fmt.Printf("Failed to open file: %v - %v\n", fi.Name(), err)
			panic("")
		}

		rInUTF8 := transform.NewReader(f, charmap.ISO8859_13.NewDecoder())

		b, err := ioutil.ReadAll(rInUTF8)
		if err != nil {
			fmt.Printf("Failed reading file: %v - %v\n", name, err)
			panic("")
		}
		f.Close()

		name8 := name + "8"
		f8, err := os.Create(path + name8)
		f8.Write([]byte{0xEF, 0xBB, 0xBF}) // Byte Order Mark (WinAmp creates this)
		f8.Write(b)
		f8.Close()

		os.Chtimes(path+name8, fi.ModTime(), fi.ModTime())
	}
}
