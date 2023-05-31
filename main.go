package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var (
	decryptfile string
	hash        string
	searchFile  string
	gemerate    int
	outputDic   string
	online      string
)

func Sha1(data string, salt string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data + salt))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

func Md5(data string, salt string) string {
	md5 := md5.New()
	md5.Write([]byte(data + salt))
	md5Data := md5.Sum([]byte(""))
	return hex.EncodeToString(md5Data)
}

func generateHash(data string, length int, start int, filename string, online string, hash ...string) {
	var fopen *os.File
	if online == "0" {
		fopen, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
		defer fopen.Close()
		if err != nil {
			panic(err)
		}

	}

	if length <= 0 {
		return
	}
	buffer := make([]rune, length)
	for i := range buffer {
		buffer[i] = '0'
	}

	for {
		word := string(buffer)

		pwdhash := Md5(Sha1(word, "tz"), "biz")
		if online != "0" {
			fopen.WriteString(string(buffer) + "\t" + pwdhash + "\n")
		} else {
			if pwdhash == hash[0] {
				fmt.Println("Crack Success:" + pwdhash + "------" + word)
				return
			}
		}
		i := length - 1
		for ; i >= 0; i-- {

			if buffer[i] >= 0 && buffer[i] < '9' {
				buffer[i] = buffer[i] + 1
				break
			}
			if buffer[i] == '9' {
				buffer[i] = 'a'
				break

			}
			if buffer[i] < 'z' && buffer[i] >= 'a' {
				buffer[i] = buffer[i] + 1
				break
			}
			buffer[i] = '0'
			//fmt.Println(word)

		}
		if i < 0 {
			break
		}

	}

}
func findHash(HashValue string, fileName string, hashFile string) {
	var target []string
	if len(hashFile) > 0 {
		data, err := ioutil.ReadFile(hashFile)
		if err != nil {
			panic(err)
		}
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line != "" {
				target = append(target, strings.TrimSpace(line))
			}
		}
	} else {
		target = append(target, HashValue)
	}

	for _, hash := range target {
		fmt.Println("在搜索hash:" + hash)
		f, err := os.Open(fileName)
		if err != nil {
			fmt.Println("cannot able to read the file", err)
			return
		}

		r := bufio.NewReader(f)
		for {
			line, _, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			if strings.Contains(string(line), hash) {
				line_str := strings.Replace(string(line), "\t", "------", -1)
				fmt.Println("破解成功:" + line_str)
				break
			}
		}
		f.Close()
	}

}
func main() {
	flag.StringVar(&decryptfile, "f", "", "输入需要破解的文件名称")
	flag.StringVar(&hash, "h", "", "输入需要破解hash")
	flag.StringVar(&searchFile, "s", "", "需要搜索的字典文件")
	flag.IntVar(&gemerate, "g", 0, "生成字典的位数")
	flag.StringVar(&outputDic, "o", "hash.dic", "生成的字典文件名")
	flag.StringVar(&online, "l", "0", "生成的字典文件名")
	flag.Parse()
	if gemerate != 0 {
		generateHash("0123456789abcdefghijklmnopqrstuvwxyz", gemerate, 0, outputDic, online, hash)

	} else if searchFile != "" {
		findHash(hash, searchFile, decryptfile)
	} else {
		fmt.Println("参数输入错误")
	}

}
