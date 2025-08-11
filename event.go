package wkproto

import "fmt"

type EventPacket struct {
	Framer
	// 事件ID(可以为空)
	Id string
	// 事件类型
	Type string
	// 事件时间戳
	Timestamp int64
	// 事件数据
	Data []byte
}

func (e *EventPacket) GetFrameType() FrameType {
	return EVENT
}

func (e *EventPacket) Size() int {
	return encodeEventSize(e, LatestVersion)
}

func (e *EventPacket) String() string {
	return fmt.Sprintf("Id:%s Type:%s Timestamp:%d Data:%s", e.Id, e.Type, e.Timestamp, string(e.Data))
}

func decodeEvent(frame Frame, data []byte, _ uint8) (Frame, error) {
	dec := NewDecoder(data)
	eventPacket := &EventPacket{}
	eventPacket.Framer = frame.(Framer)
	var err error
	if eventPacket.Id, err = dec.String(); err != nil {
		return nil, err
	}

	if eventPacket.Type, err = dec.String(); err != nil {
		return nil, err
	}
	if eventPacket.Timestamp, err = dec.Int64(); err != nil {
		return nil, err
	}
	if eventPacket.Data, err = dec.BinaryAll(); err != nil {
		return nil, err
	}

	return eventPacket, nil
}

func encodeEvent(eventPacket *EventPacket, enc *Encoder, _ uint8) error {
	enc.WriteString(eventPacket.Id)
	enc.WriteString(eventPacket.Type)
	enc.WriteInt64(eventPacket.Timestamp)
	enc.WriteBytes(eventPacket.Data)
	return nil
}

func encodeEventSize(frame Frame, _ uint8) int {
	eventPacket := frame.(*EventPacket)
	size := 0
	size += (len(eventPacket.Id) + StringFixLenByteSize)
	size += (len(eventPacket.Type) + StringFixLenByteSize)
	size += BigTimestampByteSize
	size += len(eventPacket.Data)
	return size
}
