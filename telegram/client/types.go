package client

import "encoding/json"

type BaseResp struct {
	Ok     bool            `json:"ok"`
	Result json.RawMessage `json:"result"`
}

type Update struct {
	UpdateId int64    `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

type Chat struct {
	Id int64 `json:"id"`
}

type Message struct {
	Chat      *Chat   `json:"chat,omitempty"`
	MessageId int64   `json:"message_id"`
	Text      *string `json:"text,omitempty"`
	Date      int64   `json:"date"`
}
