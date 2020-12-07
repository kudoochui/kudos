package pkg

import "testing"

func TestPackage(t *testing.T){
	_type := TYPE_DATA
	buffer := []byte{1,2,3,4,5,6,7,8,9}
	t.Log(Encode(_type, buffer))

	buffer1 := []byte{4,0,0,9,1,2,3,4,5,6,7,8,9}
	t.Log(Decode(buffer1))
}