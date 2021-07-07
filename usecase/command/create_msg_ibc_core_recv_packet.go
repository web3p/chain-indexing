package command

import (
	entity_event "github.com/crypto-com/chain-indexing/entity/event"
	"github.com/crypto-com/chain-indexing/usecase/event"
	ibc_model "github.com/crypto-com/chain-indexing/usecase/model/ibc"
)

type CreateMsgIBCCoreRecvPacket struct {
	msgCommonParams event.MsgCommonParams
	params          ibc_model.MsgRecvPacketParams
}

func NewCreateMsgIBCCoreRecvPacket(
	msgCommonParams event.MsgCommonParams,
	params ibc_model.MsgRecvPacketParams,
) *CreateMsgIBCCoreRecvPacket {
	return &CreateMsgIBCCoreRecvPacket{
		msgCommonParams,
		params,
	}
}

func (*CreateMsgIBCCoreRecvPacket) Name() string {
	return "CreateMsgIBCCoreRecvPacket"
}

func (*CreateMsgIBCCoreRecvPacket) Version() int {
	return 1
}

func (cmd *CreateMsgIBCCoreRecvPacket) Exec() (entity_event.Event, error) {
	event := event.NewMsgIBCCoreRecvPacket(cmd.msgCommonParams, cmd.params)
	return event, nil
}
