package event

import (
	"bytes"

	"github.com/crypto-com/chain-indexing/usecase/model"

	entity_event "github.com/crypto-com/chain-indexing/entity/event"
	jsoniter "github.com/json-iterator/go"
	"github.com/luci/go-render/render"
)

const MSG_MULTI_SEND = "/cosmos.bank.v1beta1.MsgMultiSend"
const MSG_MULTI_SEND_CREATED = "MsgMultiSendCreated"
const MSG_MULTI_SEND_FAILED = "MsgMultiSendFailed"

type MsgMultiSend struct {
	MsgBase

	Inputs  []model.MsgMultiSendInput  `json:"inputs"`
	Outputs []model.MsgMultiSendOutput `json:"outputs"`
}

func NewMsgMultiSend(msgCommonParams MsgCommonParams, params model.MsgMultiSendParams) *MsgMultiSend {
	return &MsgMultiSend{
		NewMsgBase(MsgBaseParams{
			MsgName:         MSG_MULTI_SEND,
			MsgSuccess:      MSG_MULTI_SEND_CREATED,
			MsgFailed:       MSG_MULTI_SEND_FAILED,
			Version:         1,
			MsgCommonParams: msgCommonParams,
		}),

		params.Inputs,
		params.Outputs,
	}
}

func (event *MsgMultiSend) ToJSON() (string, error) {
	encoded, err := jsoniter.Marshal(event)
	if err != nil {
		return "", err
	}

	return string(encoded), nil
}

func (event *MsgMultiSend) String() string {
	return render.Render(event)
}

func DecodeMsgMultiSend(encoded []byte) (entity_event.Event, error) {
	jsonDecoder := jsoniter.NewDecoder(bytes.NewReader(encoded))
	jsonDecoder.DisallowUnknownFields()

	var event *MsgMultiSend
	if err := jsonDecoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}
