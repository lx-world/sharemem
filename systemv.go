package sharemem

import (
	"fmt"
	"reflect"
	"sync"
	"syscall"
	"unsafe"
)

func NewSystemV(key,size, blockSize int)(*Mem, error){
	if blockSize > size{
		return nil, fmt.Errorf("block size over size")
	}
	fmap, _, errCode := syscall.Syscall(syscall.SYS_SHMGET, uintptr(key), uintptr(size), ipcCreate|0600)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %v", errCode)
	}

	addr, _, errCode := syscall.Syscall(syscall.SYS_SHMAT, fmap, 0, 0)
	if errCode != 0 {
		return nil, fmt.Errorf("syscall error, err: %v", errCode)
	}
	var data []byte
	originLen := size/blockSize
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = addr
	sh.Len = size+(keyLen+1+dataLen)*originLen
	sh.Cap = size+(keyLen+1+dataLen)*originLen

	mem := &Mem{
		l:    sync.RWMutex{},
		blockSize: blockSize+1+keyLen+dataLen,
		data: data,
		m: make(map[string]int),
		ch: make(chan int, originLen),
	}
	for i := 0; i < originLen; i++{
		mem.ch <- i
	}
	return mem, nil
}

func (m *Mem)WriteIdx(key string, data []byte)error{
	keyCk, err := checkkey(key)
	if err != nil{
		return err
	}

	dataCk, err := checkData(data, m.blockSize)
	if err != nil{
		return err
	}

	m.l.Lock()
	defer m.l.Unlock()

	dataTp := make([]byte, len(data)+1+keyLen+dataLen)
	dataTp[0] = 1
	copy(dataTp[1:keyLen+1], keyCk[:])
	copy(dataTp[keyLen+1:], dataCk)

	var pre,next int
	if idx, ok := m.m[key];ok{
		pre = idx*m.blockSize
		next = pre+m.blockSize
	}else {
		 i, err := m.getIdx()
		 if err != nil{
		 	return err
		 }
		 m.m[key] = i
		pre = i*m.blockSize
		next = pre+m.blockSize
	}
	if  next > len(m.data){
		return fmt.Errorf("over size")
	}
	copy(m.data[pre:next], dataTp)
	return nil
}

func (m *Mem)DelIdx(key string){
	m.l.Lock()
	defer m.l.Unlock()

	idx, ok := m.m[key]
	index := idx*m.blockSize
	if index < len(m.data){
		m.data[index] = 0
	}
	if ok && len(m.ch) < cap(m.ch){
		m.ch <- idx
	}
}

func (m *Mem)GetAll()(res map[string][]byte){
	res = make(map[string][]byte)
	for i:=0; i<len(m.data)/m.blockSize;i++{
		pre := i*m.blockSize
		next := pre+m.blockSize
		key, data, err := m.dealBlocak(m.data[pre:next])
		if key=="" || data == nil || err != nil{
			continue
		}
		res[key] =  data
	}
	return
}

func (m *Mem)getIdx()(int, error){
	if len(m.ch) == 0{
		return 0, fmt.Errorf("memory is full")
	}
	return <-m.ch, nil
}