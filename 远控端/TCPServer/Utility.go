package main

import (
	"MyOperatePacket4Server"
	"bufio"
	"golang.org/x/text/encoding/charmap"
	"os"
	"unicode/utf8"
)

func ReadFileContents(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, err := file.Stat()
	FileSize := stats.Size()

	bytes := make([]byte, FileSize)

	buffer := bufio.NewReader(file)

	_, err = buffer.Read(bytes)

	return bytes, err
}

func WriteFileContents(fileName string,b []byte) error {
	file, err := os.Create(fileName)
	if err != nil {return err}
	defer file.Close()

	buffer := bufio.NewWriter(file)

	_, err = buffer.Write(b)

	buffer.Flush()//让缓冲区的内容立即写入文件
	return err
}

func IsLetter(s byte) bool {
	if (s < 'a' || s > 'z') && (s < 'A' || s > 'Z') {
			return false
	}
	return true
}

func IsFullPath(s string) bool {
	if len(s)>=3 && IsLetter(s[0]) && s[1]==':' && (s[2]=='\\' || s[2]=='/'){
		return true
	}
	return false
}

func isSpace(r rune) bool {
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

var q bool=true
func ScanWordsAndQuotes(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			if r=='"' {q=false;start++/*跳过引号*/;break}else{q=true;break}
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if q{
			if isSpace(r) {
				return i + width, data[start:i], nil
			}
		}else{
			if r=='"' {
				return i + width, data[start:i], nil
			}
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	if !q{start--}//之前检测出了引号,但是后面没有检测到的情况
	return start, nil, nil
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func DecodeWindows1250(enc []byte) string {
	dec := charmap.Windows1250.NewDecoder()
	out, _ := dec.Bytes(enc)
	return string(out)
}

func EncodeWindows1250(inp string) []byte {
	enc := charmap.Windows1250.NewEncoder()
	out, _ := enc.Bytes(MyOperatePacket4Server.String2bytes(inp))
	return out
}