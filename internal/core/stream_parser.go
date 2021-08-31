package core

import (
	"fmt"
	"io"

	"github.com/yomorun/y3"
	"github.com/yomorun/yomo/internal/frame"
	"github.com/yomorun/yomo/logger"
)

// ParseFrame parses the frame from QUIC stream.
func ParseFrame(stream io.Reader) (frame.Frame, error) {
	buf, err := y3.ReadPacket(stream)
	if err != nil {
		logger.Error("[ParseFrame] read first byte", "err", err)
		return nil, err
	}
	if len(buf) > 512 {
		logger.Debugf("🔗 parsed out total %d bytes: \n\thead 64 bytes are: [%# x], \n\ttail 64 bytes are: [%#x]", len(buf), buf[0:64], buf[len(buf)-64:])
	} else {
		logger.Debugf("🔗 parsed out: [%# x]", buf)
	}

	frameType := buf[0]
	// determine the frame type
	logger.Debugf("[ParseFrame] type=%v", frameType)
	switch frameType {
	case 0x80 | byte(frame.TagOfHandshakeFrame):
		handshakeFrame := readHandshakeFrame(buf)
		logger.Debugf("[ParseFrame] HandshakeFrame name=%s, type=%s", handshakeFrame.Name, handshakeFrame.Type())
		return handshakeFrame, nil
	case 0x80 | byte(frame.TagOfDataFrame):
		data := readDataFrame(buf)
		logger.Debugf("[ParseFrame] DataFrame tid=%s, issuer=%s, data-tag=%#x, len(carriage)=%d", data.TransactionID(), data.Issuer(), data.GetDataTagID(), len(data.GetCarriage()))
		return data, nil
	case 0x80 | byte(frame.TagOfPingFrame):
		return frame.DecodeToPingFrame(buf)
	case 0x80 | byte(frame.TagOfPongFrame):
		return frame.DecodeToPongFrame(buf)
	case 0x80 | byte(frame.TagOfAcceptedFrame):
		return frame.DecodeToAcceptedFrame(buf)
	case 0x80 | byte(frame.TagOfRejectedFrame):
		return frame.DecodeToRejectedFrame(buf)
	default:
		return nil, fmt.Errorf("unknown frame type, buf[0]=%#x", buf[0])
	}
}

func readHandshakeFrame(buf []byte) *frame.HandshakeFrame {
	// parse to HandshakeFrame
	handshake, err := frame.DecodeToHandshakeFrame(buf)
	if err != nil {
		panic(err)
	}
	return handshake
}

func readDataFrame(buf []byte) *frame.DataFrame {
	// parse to DataFrame
	data, err := frame.DecodeToDataFrame(buf)
	if err != nil {
		panic(err)
	}
	return data
}
