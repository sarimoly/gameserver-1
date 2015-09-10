package GxMessage

import (
	"github.com/golang/protobuf/proto"
	"testing"
)

func Test_NewGxMessage(t *testing.T) {
	msg := NewGxMessage()
	if msg == nil {
		t.Error("new nil message")
	}
	if len(msg.Header) != MessageHeaderLen {
		t.Error("new message header length error, length: ", len(msg.Header))
	}
}

func Test_Package(t *testing.T) {
	msg := NewGxMessage()

	buff := "Test_Package"
	err := msg.Package([]byte(buff))
	if err != nil {
		t.Error("package message error")
	}

	buff2, err2 := msg.Unpackage()
	if err2 != nil {
		t.Error("unpackage message error")
	}

	if string(buff2) != buff {
		t.Error("package and unpackage message error", string(buff2), buff)
	}

	buff = "Test_Packageaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	err = msg.Package([]byte(buff))
	if err != nil {
		t.Error("package message error")
	}

	buff2, err2 = msg.Unpackage()
	if err2 != nil {
		t.Error("unpackage message error")
	}

	if string(buff2) != buff {
		t.Error("package and unpackage message error", string(buff2), buff, msg.GetUnlen(), msg.GetLen())
	}
}

func Test_PbPackage(t *testing.T) {
	msg := NewGxMessage()

	var n1 MessageNot
	n1.Id = proto.Uint32(123)
	n1.Text = proto.String("Test_PbPackage")
	err := msg.PackagePbmsg(&n1)
	if err != nil {
		t.Error("package protobuf message error")
	}

	var n2 MessageNot
	err2 := msg.UnpackagePbmsg(&n2)
	if err2 != nil {
		t.Error("unpackage protobuf message error")
	}

	if *n1.Id != *n2.Id || *n1.Text != *n2.Text {
		t.Error("package and unpackage protobuf message error", n1.String(), n2.String())
	}
}

func Test_Set_Get(t *testing.T) {
	msg := NewGxMessage()

	msg.SetLen(111)
	if msg.GetLen() != 111 {
		t.Error("set message data length error, ", msg.GetLen())
	}

	msg.SetId(222)
	if msg.GetId() != 222 {
		t.Error("set message id error, ", msg.GetId())
	}

	msg.SetSeq(333)
	if msg.GetSeq() != 333 {
		t.Error("set message seq error, ", msg.GetSeq())
	}

	msg.SetCmd(444)
	if msg.GetCmd() != 444 {
		t.Error("set message cmd error, ", msg.GetCmd())
	}

	msg.SetMask(3, true)
	if !msg.GetMask(3) {
		t.Error("set message mask 1 error, ", msg.GetMask(3))
	}
	if msg.GetMask(0) {
		t.Error("set message mask 2 error, ", msg.GetMask(0))
	}
	if msg.GetMask(1) {
		t.Error("set message mask 3 error, ", msg.GetMask(1))
	}
}

func Test_InitData(t *testing.T) {
	msg := NewGxMessage()
	msg.SetLen(12)
	msg.InitData()

	if len(msg.Data) != int(msg.GetLen()) {
		t.Error("new message data length error, ", len(msg.Data), msg.GetLen())
	}

	msg1 := NewGxMessage()
	msg1.SetLen(0)
	msg1.InitData()

	if len(msg1.Data) != int(msg1.GetLen()) {
		t.Error("new message data length error, ", len(msg1.Data), msg1.GetLen())
	}
}
