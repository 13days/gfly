package codec

import (
	"fmt"
	"github.com/13days/gfly/protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testType struct {
	a int64
	b string
}

func TestPbSerializationMarshal(t *testing.T) {
	pbSer := &pbSerialization{}
	data, err := pbSer.Marshal(nil)
	assert.NotNil(t, err)
	fmt.Println(string(data), err)
	err = pbSer.Unmarshal(data, &protocol.Response{})
	assert.NotNil(t, err)
	fmt.Println(err)
	err = pbSer.Unmarshal(nil, &protocol.Response{})
	assert.NotNil(t, err)
	fmt.Println(err)
	data, err = pbSer.Marshal(&protocol.Response{RetCode: 0, RetMsg:"Success"})
	assert.Nil(t, err)
	fmt.Println(data)
	resp := &protocol.Response{}
	err = pbSer.Unmarshal(data, resp)
	assert.Nil(t, err)
	fmt.Printf("%+v\n", resp)
}
