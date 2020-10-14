package packet

import (
	"bytes"
	"encoding/binary"
)

//默认4个字段的包长度
const (
	HEADER_LEN = 4
)

//Pack 封包
func Pack(buf []byte) []byte {

	return append(Int2Bytes(len(buf)), buf...)

}

//UnPack 拆包
func UnPack(buf []byte) (pack []byte, cache []byte) {
	length := len(buf)

	messageLength := BytesToInt(buf[:HEADER_LEN]) //读取期望的数据包大小
	total := HEADER_LEN + messageLength           // 预期的完整的数据包长度
	if total > length {                           //说明已接受的数据包,暂时不够预期的数据包的大小
		pack = []byte{}
		cache = buf
	} else if total == length {
		pack = buf[HEADER_LEN:]
		cache = []byte{}
	} else {
		pack = buf[HEADER_LEN:total]
		cache = buf[total:]
	}
	return

}

func Int2Bytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()

}

func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
