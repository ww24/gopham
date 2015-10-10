package client

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

func readFile(filename string) (res string, err error) {
	f, err := os.Open("fixtures/" + filename)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	sc.Split(split())

	buff := make([]string, 0, 10)
	for sc.Scan() {
		buff = append(buff, sc.Text())
	}
	if err = sc.Err(); err != nil {
		return
	}

	res = strings.Join(buff, "|")

	return
}

func TestSplit(t *testing.T) {
	defer func() {
		cause := recover()
		if cause != nil {
			t.Fatal(cause)
		}
	}()

	files, err := ioutil.ReadDir("fixtures")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		filename := file.Name()
		str, err := readFile(filename)
		if err != nil {
			panic(err)
		}

		f, err := os.Open("fixtures/" + filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		bytes, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		expected := strings.Join(strings.Split(strings.Trim(string(bytes), "\r\n"), "\n\n"), "|")
		if str != expected {
			log.Println("actual:", str)
			log.Println("expected:", expected)
			t.Fatal(filename)
		}

		fmt.Println("passed:", filename)
	}
}
