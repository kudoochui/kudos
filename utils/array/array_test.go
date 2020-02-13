package array

import (
	"gotest.tools/assert"
	"testing"
)

func TestArray1(t *testing.T)  {
	a := []int64{1,2,3,4,5}
	a = PullInt64(a, 5)
	t.Log(a)
	assert.DeepEqual(t, a, []int64{1,2,3,4})
}

func TestArray2(t *testing.T)  {
	a := []int64{1,2,3,4,5}
	a = append(a[:2],a[3:]...)
	t.Log(a)
}