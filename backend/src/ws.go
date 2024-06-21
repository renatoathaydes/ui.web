package src

import (
	"encoding/json"
	"fmt"
	"io"
)

type WsMessage struct {
	Ok      bool        `json:"ok"`
	Id      int64       `json:"id"`
	Data    interface{} `json:"data"`
	MsgType string      `json:"type"`
}

func WsMessageFromJson(reader io.Reader) (WsMessage, error) {
	dec := json.NewDecoder(reader)
	msg := WsMessage{}
	err := dec.Decode(&msg)
	if err == nil {
		fmt.Printf("Received message: %+v\n", msg)
	}
	return msg, err
}

func (msg WsMessage) WriteJson(writer io.Writer) {
	enc := json.NewEncoder(writer)
	err :=  enc.Encode(msg)
	if err != nil {
		fmt.Printf("Failed to write WS response, closing connection: %s\n", err)
	}
}

func (req WsMessage) WsMessageResponse(data interface{}) WsMessage {
	return WsMessage{
		Ok: req.Ok,
		Id: req.Id,
		MsgType: "response",
		Data: data,
	};
}
