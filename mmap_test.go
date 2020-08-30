package sharemem

import "testing"

func Test_Mmap(t *testing.T){
	m, err := NewMMap("test.log",4,2)
	if err != nil{
		t.Errorf("err: %v", err)
		return
	}
	m.WriteIdx("test1",[]byte("31"))
	res,_ := m.GetKey("test1")
	t.Logf("%s", res)
}
