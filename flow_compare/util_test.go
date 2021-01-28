package flow_compare

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestExecFlowCompare(t *testing.T) {

}

func TestGenFlowComparePath(t *testing.T) {
	path, err := GenFlowComparePath("P.S.M", "127.0.0.1:8888", []string{"a","b","c"}, []int{1,2,3}, )
	assert.Nil(t, err)
	fmt.Println(path)
}

func TestGetFlowComparor(t *testing.T) {

}

func TestParseFlowComparePath(t *testing.T) {
	path, err := GenFlowComparePath("P.S.M", "127.0.0.1:8888", []string{"a", "b", "c"}, []int{1, 2, 3})
	assert.Nil(t, err)
	strs := strings.Split(path, "/")
	newPath := strs[len(strs)-1]
	fmt.Println(newPath)
	methods, svrAddr := ParseFlowComparePath(newPath)
	for method, rate := range methods {
		fmt.Println(method, rate)
	}
	fmt.Println(svrAddr)
}
func TestRegisterFlowComparor(t *testing.T) {

}