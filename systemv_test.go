package sharemem

import "testing"

func Test_NewSystemV(t *testing.T) {
	memory, err := NewSystemV(1, 4, 2)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	memory.WriteIdx("lx", []byte("11"))
	memory.WriteIdx("lx1", []byte("21"))
	memory.DelIdx("lx")
	memory.DelIdx("lx1")
	res := memory.GetAll()
	for key, value := range res {
		t.Logf("key: %s, value: %s\n", key, value)
	}
}
