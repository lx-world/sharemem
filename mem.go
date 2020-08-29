package sharemem

import (
	"fmt"
	"sync"
)

const (
	ipcCreate     = 00001000 // 不存在则创建
	keyOriginSize = 16
	keyLen        = 17
	dataLen       = 4
)

type Mem struct {
	l         sync.RWMutex
	blockSize int
	m         map[string]int
	data      []byte
	ch        chan int
}

func checkkey(key string) ([keyLen]byte, error) {
	var res [keyLen]byte
	if len(key) > keyOriginSize {
		return res, fmt.Errorf("key must <= 16")
	}
	res[0] = byte(len(key))
	for idx, v := range key {
		res[idx+1] = byte(v)
	}
	return res, nil
}

func checkData(data []byte, size int) ([]byte, error) {
	if len(data) > size-1-keyLen-dataLen {
		return nil, fmt.Errorf("data len over range")
	}
	res := make([]byte, len(data)+dataLen)
	copy(res[:4], datalen(len(data)))
	copy(res[4:], data)
	return res, nil
}

func (m *Mem) dealBlocak(data []byte) (string, []byte, error) {
	if len(data) != m.blockSize {
		return "", nil, fmt.Errorf("block size error")
	}
	if data[0] == 0 {
		return "", nil, nil
	}
	key := data[1:18]

	return dealkey(key), dealData(data[18:]), nil
}

func dealkey(data []byte) string {
	l := data[0]
	res := make([]byte, l)
	for i := 0; i < int(l); i++ {
		res[i] = data[i+1]
	}
	return string(res)
}

func dealData(data []byte) []byte {
	l := dataLenToInt(data[:5])
	res := make([]byte, l)
	for i := 0; i < l; i++ {
		res[i] = data[i+dataLen]
	}
	return res
}

func datalen(l int) []byte {
	res := make([]byte, 4)
	res[0], res[1], res[2], res[3] = byte(l>>24), byte(l>>16), byte(l>>8), byte(l)
	return res
}

func dataLenToInt(data []byte) int {
	return int(data[3]) | int(data[2])<<8 | int(data[1])<<16 | int(data[0])<<24
}
