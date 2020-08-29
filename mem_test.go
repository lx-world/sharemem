package sharemem

import "testing"

func Test_checkkey(t *testing.T){
	str := "hello"
	res,err := checkkey(str)
	if err != nil{
		t.Errorf("%v",err)
		return
	}
	t.Logf("len: %d, %v",len(res), res)
}
