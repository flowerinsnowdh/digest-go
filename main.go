/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var algorithm hash.Hash
var paths []string

func init() {
	algorithm = func() hash.Hash {
		var simple = filepath.Base(os.Args[0])
		var ext = filepath.Ext(simple)
		simple = strings.TrimSuffix(simple, ext)

		switch strings.ToLower(simple) {
		case "md4":
			return md4.New()
		case "md5":
			return md5.New()
		case "sha1":
			return sha1.New()
		case "sha224":
			return sha256.New224()
		case "sha256":
			return sha256.New()
		case "sha384":
			return sha512.New384()
		case "sha512":
			return sha512.New()
		case "ripemd160":
			return ripemd160.New()
		case "sha3-224":
			return sha3.New224()
		case "sha3-256":
			return sha3.New256()
		case "sha3-384":
			return sha3.New384()
		case "sha3-512":
			return sha3.New512()
		case "sha512-224":
			return sha512.New512_224()
		case "sha512-256":
			return sha512.New512_256()
		case "blake2s-128":
			h, _ := blake2s.New128(nil)
			return h
		case "blake2s-256":
			h, _ := blake2s.New256(nil)
			return h
		case "blake2b-256":
			h, _ := blake2b.New256(nil)
			return h
		case "blake2b-384":
			h, _ := blake2b.New384(nil)
			return h
		case "blake2b-512":
			h, _ := blake2b.New512(nil)
			return h
		case "sm3":
			return sm3.New()
		default:
			_, _ = fmt.Fprintln(os.Stderr, "unknown algorithm")
			os.Exit(1)
			return nil
		}
	}()

	paths = make([]string, 0, len(os.Args)-1)
	for i := 1; i < len(os.Args); i++ {
		paths = append(paths, os.Args[i])
	}
}

func main() {
	if len(paths) == 0 { // use stdin
		sum("stdin", os.Stdin)
	} else {
		var streamsToClose []*os.File = make([]*os.File, 0, len(paths))
		defer func() {
			for _, stream := range streamsToClose {
				_ = stream.Close()
			}
		}()

		for _, p := range paths {
			if stream, err := os.Open(p); err != nil {
				panic(err)
			} else {
				streamsToClose = append(streamsToClose, stream)
				sum(p, stream)
			}
		}
	}
}

func sum(path string, f *os.File) {
	algorithm.Reset()
	var buffer []byte = make([]byte, 4096)
	for {
		if n, err := f.Read(buffer); n <= 0 || err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		} else {
			_, _ = algorithm.Write(buffer[:n])
		}
	}
	fmt.Println(hex.EncodeToString(algorithm.Sum(nil)), path)
}
