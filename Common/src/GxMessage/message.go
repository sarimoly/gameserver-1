package GxMessage

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
)

const (
	MessageMaskDisconn = 0 //是否断开连接
	MessageMaskNotify  = 1 //是否通知
)

const (
	MessageIdBit    = 0  //消息来源或者目的id
	MessageUnlenBit = 4  //未压缩前消息体长度，可为0，和len相等表示没有压缩
	MessageLenBit   = 6  //消息体长度，可为0
	MessageCmdBit   = 8  //消息命令字
	MessageSeqBit   = 10 //消息序号
	MessageRetBit   = 12 //消息返回值，消息返回时使用
	MessageMaskBit  = 14 //一些标志

	MessageHeaderLen = 16 //消息长度
)

type GxMessage struct {
	Header []byte
	Data   []byte
}

func NewGxMessage() *GxMessage {
	msg := new(GxMessage)
	msg.Header = make([]byte, MessageHeaderLen)
	return msg
}

//打包原生字符串
func (msg *GxMessage) Package(buf []byte) error {
	l := len(buf)
	if l == 0 {
		return nil
	}

	var b bytes.Buffer
	c := false

	//小于指定长度不用检查是否需要压缩
	if l > 10 {
		w := zlib.NewWriter(&b)
		w.Write(buf)
		w.Close()
		c = true
	}

	//压缩后长度比原来小，就保存压缩数据
	if c && b.Len() < l {
		msg.SetUnlen(uint16(l))
		msg.SetLen(uint16(b.Len()))
		msg.Data = make([]byte, b.Len())
		copy(msg.Data[:], b.Bytes())
	} else {
		msg.SetUnlen(uint16(l))
		msg.SetLen(uint16(l))
		msg.Data = make([]byte, l)
		copy(msg.Data[:], buf)
	}
	return nil
}

//解包原生字符串
func (msg *GxMessage) Unpackage() ([]byte, error) {
	if msg.GetLen() == 0 {
		return []byte(""), nil
	}

	if msg.GetLen() == msg.GetUnlen() {
		data := make([]byte, msg.GetLen())
		copy(data[:], msg.Data)
		return data, nil
	} else if msg.GetLen() < msg.GetUnlen() {
		var b bytes.Buffer
		b.Write(msg.Data)
		r, err := zlib.NewReader(&b)
		if err != nil {
			return []byte(""), err
		}
		defer r.Close()

		data := make([]byte, msg.GetUnlen())
		l, _ := r.Read(data)
		if l != int(msg.GetUnlen()) {
			return []byte(""), errors.New("uncompress erro")
		}
		return data, nil
	} else {
		return []byte(""), errors.New("message unpackage erro")
	}
}

//打包protobuf消息
func (msg *GxMessage) PackagePbmsg(pb proto.Message) error {
	buff, err := proto.Marshal(pb)
	if err != nil {
		return err
	}

	return msg.Package(buff)
}

//解包protobuf消息
func (msg *GxMessage) UnpackagePbmsg(pb proto.Message) error {
	data, err := msg.Unpackage()
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, pb)
}

//根据消息长度初始化消息体内存
func (msg *GxMessage) InitData() {
	if msg.GetLen() == 0 {
		return
	}
	msg.Data = make([]byte, msg.GetLen())
}

func (msg *GxMessage) get16(t uint32) uint16 {
	buf := bytes.NewBuffer(make([]byte, 0, MessageHeaderLen))
	buf.Write(msg.Header[t : t+2])

	i16 := make([]byte, 2)
	buf.Read(i16)
	return binary.BigEndian.Uint16(i16)
}

func (msg *GxMessage) set16(t uint32, id uint16) {
	binary.BigEndian.PutUint16(msg.Header[t:t+2], id)
}

func (msg *GxMessage) get32(t uint32) uint32 {
	buf := bytes.NewBuffer(make([]byte, 0, MessageHeaderLen))
	buf.Write(msg.Header[t : t+4])

	i32 := make([]byte, 4)
	buf.Read(i32)
	return binary.BigEndian.Uint32(i32)
}

func (msg *GxMessage) set32(t uint32, id uint32) {
	binary.BigEndian.PutUint32(msg.Header[t:t+4], id)
}

func (msg *GxMessage) GetId() uint32 {
	return msg.get32(MessageIdBit)
}

func (msg *GxMessage) SetId(id uint32) {
	msg.set32(MessageIdBit, id)
}

func (msg *GxMessage) GetCmd() uint16 {
	return msg.get16(MessageCmdBit)
}

func (msg *GxMessage) SetCmd(cmd uint16) {
	msg.set16(MessageCmdBit, cmd)
}

func (msg *GxMessage) GetSeq() uint16 {
	return msg.get16(MessageSeqBit)
}

func (msg *GxMessage) SetSeq(seq uint16) {
	msg.set16(MessageSeqBit, seq)
}

func (msg *GxMessage) GetUnlen() uint16 {
	return msg.get16(MessageUnlenBit)
}

func (msg *GxMessage) SetUnlen(len uint16) {
	msg.set16(MessageUnlenBit, len)
}

func (msg *GxMessage) GetLen() uint16 {
	return msg.get16(MessageLenBit)
}

func (msg *GxMessage) SetLen(len uint16) {
	msg.set16(MessageLenBit, len)
}

func (msg *GxMessage) GetRet() uint16 {
	return msg.get16(MessageRetBit)
}

func (msg *GxMessage) SetRet(ret uint16) {
	msg.set16(MessageRetBit, ret)
}

func (msg *GxMessage) GetMask(mask uint16) bool {
	i := msg.get16(MessageMaskBit)
	return (i & (1 << mask)) != 0
}

func (msg *GxMessage) SetMask(mask uint16, b bool) {
	i := msg.get16(MessageMaskBit)

	if b {
		i |= 1 << mask
	} else {
		i &= ^(1 << mask)
	}

	msg.set16(MessageMaskBit, i)
}

func (msg *GxMessage) String() string {
	return fmt.Sprintf("id: %d, unlen: %d, len: %d, cmd: %d, seq: %d, ret: %d", msg.GetId(), msg.GetUnlen(), msg.GetLen(), msg.GetCmd(), msg.GetSeq(), msg.GetRet())
}

func (msg *GxMessage) Copy() *GxMessage {
	newMsg := NewGxMessage()
	newMsg.Data = make([]byte, msg.GetLen())
	copy(newMsg.Header[:], msg.Header)
	copy(newMsg.Data[:], msg.Data)
	return newMsg
}
