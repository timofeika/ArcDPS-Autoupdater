package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	const cprght = "ArcDPS autoupdater, (c) .devil 2019"
	fileURL := "https://www.deltaconnected.com/arcdps/x64/d3d9.dll"
	fileSum := "https://www.deltaconnected.com/arcdps/x64/d3d9.dll.md5sum"

	fmt.Println(cprght)

	if fileExists("d3d9.dll") != true {
		DownloadFile("d3d9.dll", fileURL)
	} else {
		if err := DownloadFile("d3d9.dll.md5sum", fileSum); err != nil {
			panic(err)
		}
		d3d9str := calculateHash("d3d9.dll")

		file, err := os.Open("d3d9.dll.md5sum")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		b, err := ioutil.ReadFile("d3d9.dll.md5sum")
		if err != nil {
			panic(err)
		}

		strng := string(b[:32])

		if strings.EqualFold(d3d9str, strng) {
			fmt.Println("ArcDPS is up-to-date!")
		} else {
			fmt.Println("ArcDPS need an update.")
			fmt.Printf("Hash of arcdps is %s, hash of md5sum is: %s\n", d3d9str, strng)
			fmt.Println("Downloading ArcDPS...")
			DownloadFile("d3d9.dll", fileURL)
		}
	}
}

// DownloadFile bla bla bla
func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func calculateHash(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()

	_, err = io.Copy(hash, file)
	if err != nil {
		panic(err)
	}

	return hex.EncodeToString(hash.Sum(nil))

}
