package sharemem

import (
	"os"
	"sync"
	"syscall"
)

func NewMMap(path string, size, blockSize int)(*Mem, error){
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil{
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil{
		return nil, err
	}

	originLen := size / blockSize
	sz := size + (keyLen+1+dataLen)*originLen
	if stat.Size() != int64(sz) {
		if err := f.Truncate(int64(sz)); err != nil {
			return nil, err
		}
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, sz, syscall.PROT_WRITE|syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, err
	}
	f.Close()
	mem := &Mem{
		l:         sync.RWMutex{},
		blockSize: blockSize + 1 + keyLen + dataLen,
		data:      data,
		m:         make(map[string]int),
		ch:        make(chan int, originLen),
	}
	for i := 0; i < originLen; i++ {
		mem.ch <- i
	}
	return mem, nil
}
