package codec

import (
	"fmt"
	"testing"

	"github.com/lubanproj/gorpc/protocol"
	"github.com/stretchr/testify/assert"
)

func TestMsgpackSerializationMarshal(t *testing.T) {
	msgSer := &MsgpackSerialization{}
	data, err := msgSer.Marshal(nil)
	assert.NotNil(t, err)
	fmt.Println(string(data), err)
	err = msgSer.Unmarshal(data, &protocol.Response{})
	assert.NotNil(t,err)
	err = msgSer.Unmarshal(nil, &protocol.Response{})
	assert.NotNil(t,err)
	fmt.Println(err)
	data, err = msgSer.Marshal(&protocol.Response{RetCode: 0, RetMsg:"Success"})
	assert.Nil(t, err)
	fmt.Println(data)
	resp := &protocol.Response{}
	err = msgSer.Unmarshal(data, resp)
	assert.Nil(t, err)
	fmt.Printf("%+v\n", resp)

	tt := testType{
		a: 1,
		b: "2",
	}
	data, err = msgSer.Marshal(&tt)
	assert.Nil(t, err)
	fmt.Println(data)
	rt := &testType{}
	err = msgSer.Unmarshal(data, rt)
	assert.Nil(t, err)
	fmt.Printf("%+v\n", rt)
}
