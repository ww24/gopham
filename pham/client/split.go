package client

import (
	"bufio"
	"strings"
)

func split() (f bufio.SplitFunc) {
	// 読み取り位置
	offset := 0
	f = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// 1 行取得
		d := data[offset:]
		advance, token, err = bufio.ScanLines(d, atEOF)

		// 1 行に満たない場合はそのまま return
		if advance == 0 && token == nil {
			offset = 0
			return
		}

		// 1 行以上の場合
		if token != nil {
			offset += advance

			// 改行が存在した場合
			if advance == 1 && len(token) == 0 {
				// 改行が連続した場合
				if offset > 1 {
					advance = offset
					token = []byte(strings.Trim(string(data[:offset-2]), "\n"))
					offset = 0
					return
				}

				token = nil
				if len(data) <= offset {
					offset = 0
				}
				return
			}

			// ファイルの終端
			if atEOF && advance == len(d) {
				advance = offset
				token = []byte(strings.Trim(string(data[:offset]), "\n"))
				offset = 0
				return
			}

			advance, token, err = f(data, atEOF)
			return
		}

		offset = 0
		return
	}
	return
}
