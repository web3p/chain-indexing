package command

import (
	entity_event "github.com/crypto-com/chainindex/entity/event"
	"github.com/crypto-com/chainindex/usecase/event"
	"github.com/crypto-com/chainindex/usecase/model"
)

type CreateMsgSubmitCancelSoftwareUpgradeProposal struct {
	msgCommonParams event.MsgCommonParams
	params          model.MsgSubmitCancelSoftwareUpgradeProposalParams
}

func NewCreateMsgSubmitCancelSoftwareUpgradeProposal(
	msgCommonParams event.MsgCommonParams,
	params model.MsgSubmitCancelSoftwareUpgradeProposalParams,
) *CreateMsgSubmitCancelSoftwareUpgradeProposal {
	return &CreateMsgSubmitCancelSoftwareUpgradeProposal{
		msgCommonParams,
		params,
	}
}

// Name returns name of command
func (*CreateMsgSubmitCancelSoftwareUpgradeProposal) Name() string {
	return "CreateMsgSubmitProposal"
}

// Version returns version of command
func (*CreateMsgSubmitCancelSoftwareUpgradeProposal) Version() int {
	return 1
}

// Exec process the command data and return the event accordingly
func (cmd *CreateMsgSubmitCancelSoftwareUpgradeProposal) Exec() (entity_event.Event, error) {
	event := event.NewMsgSubmitCancelSoftwareUpgradeProposal(cmd.msgCommonParams, cmd.params)
	return event, nil
}
