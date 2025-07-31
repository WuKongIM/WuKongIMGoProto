package wkproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkEncodeAndDecode(t *testing.T) {
	packet := &ChunkPacket{
		MessageID: 1,
		ChunkID:   1,
		Payload:   []byte("hello world"),
		EndReason: 1,
	}

	codec := New()

	// 编码
	packetBytes, err := codec.EncodeFrame(packet, 1)
	assert.NoError(t, err)
	// panic(fmt.Sprintf("%v",packetBytes))
	// 解码
	resultPacket, _, err := codec.DecodeFrame(packetBytes, 1)
	assert.NoError(t, err)
	resultChunkPacketPacket, ok := resultPacket.(*ChunkPacket)
	assert.Equal(t, true, ok)

	assert.Equal(t, packet.MessageID, resultChunkPacketPacket.MessageID)
	assert.Equal(t, packet.ChunkID, resultChunkPacketPacket.ChunkID)
	assert.Equal(t, packet.Payload, resultChunkPacketPacket.Payload)
}
