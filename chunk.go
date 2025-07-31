package wkproto

// ChunkPacket 消息块
type ChunkPacket struct {
	Framer
	MessageID int64  // 消息ID(同个消息多个块的消息ID相同)
	ChunkID   uint64 // 块ID（顺序递增）
	Payload   []byte // 消息内容
}

// GetPacketType 获得包类型
func (c *ChunkPacket) GetFrameType() FrameType {
	return Chunk
}

func (c *ChunkPacket) Size() int {
	return c.SizeWithProtoVersion(LatestVersion)
}

func (c *ChunkPacket) SizeWithProtoVersion(protVersion uint8) int {
	return encodeChunkSize(c, protVersion)
}

func encodeChunk(chunkPacket *ChunkPacket, enc *Encoder, _ uint8) error {
	enc.WriteInt64(chunkPacket.MessageID)
	enc.WriteUint64(chunkPacket.ChunkID)
	enc.WriteBytes(chunkPacket.Payload)
	return nil
}

func decodeChunk(frame Frame, data []byte, _ uint8) (Frame, error) {
	dec := NewDecoder(data)

	chunkPacket := &ChunkPacket{}
	var err error
	if chunkPacket.MessageID, err = dec.Int64(); err != nil {
		return nil, err
	}

	if chunkPacket.ChunkID, err = dec.Uint64(); err != nil {
		return nil, err
	}

	if chunkPacket.Payload, err = dec.BinaryAll(); err != nil {
		return nil, err
	}

	return chunkPacket, nil
}

func encodeChunkSize(packet *ChunkPacket, _ uint8) int {

	size := 0
	size += MessageIDByteSize   // 消息ID
	size += ChunkIDByteSize     // 块ID
	size += len(packet.Payload) // 消息内容
	return size
}
