package wkproto

type EndReason uint8

// EndReason constants define why a stream was completed
const (
	// EndReasonSuccess indicates the stream completed successfully (default)
	EndReasonSuccess EndReason = 0
	// EndReasonTimeout indicates the stream ended due to inactivity timeout
	EndReasonTimeout EndReason = 1
	// EndReasonError indicates the stream ended due to an error
	EndReasonError EndReason = 2
	// EndReasonCancelled indicates the stream was manually cancelled
	EndReasonCancelled EndReason = 3
	// EndReasonForce indicates the stream was forcefully ended (e.g., channel closure)
	EndReasonForce EndReason = 4
)

func (e EndReason) String() string {
	switch e {
	case EndReasonSuccess:
		return "success"
	case EndReasonTimeout:
		return "timeout"
	case EndReasonError:
		return "error"
	case EndReasonCancelled:
		return "cancelled"
	case EndReasonForce:
		return "force"
	default:
		return "unknown"
	}
}

func (e EndReason) Value() uint8 {
	return uint8(e)
}

// ChunkPacket 消息块
type ChunkPacket struct {
	Framer
	MessageID int64     // 消息ID(同个消息多个块的消息ID相同)
	ChunkID   uint64    // 块ID（顺序递增）
	EndReason EndReason // 结束原因
	Payload   []byte    // 消息内容
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
	enc.WriteUint8(chunkPacket.EndReason.Value())
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

	var endReason uint8
	if endReason, err = dec.Uint8(); err != nil {
		return nil, err
	}
	chunkPacket.EndReason = EndReason(endReason)

	if chunkPacket.Payload, err = dec.BinaryAll(); err != nil {
		return nil, err
	}

	return chunkPacket, nil
}

func encodeChunkSize(packet *ChunkPacket, _ uint8) int {

	size := 0
	size += MessageIDByteSize   // 消息ID
	size += ChunkIDByteSize     // 块ID
	size += EndReasonByteSize   // 结束原因
	size += len(packet.Payload) // 消息内容
	return size
}
