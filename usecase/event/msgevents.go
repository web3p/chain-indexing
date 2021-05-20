package event

import "github.com/crypto-com/chain-indexing/entity/event"

type MsgEvent interface {
	event.Event

	MsgType() string
	TxHash() string
	TxSuccess() bool
}

var MSG_EVENTS = []string{
	MSG_SEND_CREATED,
	MSG_SEND_FAILED,
	MSG_MULTI_SEND_CREATED,
	MSG_MULTI_SEND_FAILED,

	MSG_SET_WITHDRAW_ADDRESS_CREATED,
	MSG_SET_WITHDRAW_ADDRESS_FAILED,
	MSG_WITHDRAW_DELEGATOR_REWARD_CREATED,
	MSG_WITHDRAW_DELEGATOR_REWARD_FAILED,
	MSG_WITHDRAW_VALIDATOR_COMMISSION_CREATED,
	MSG_WITHDRAW_VALIDATOR_COMMISSION_FAILED,
	MSG_FUND_COMMUNITY_POOL_CREATED,
	MSG_FUND_COMMUNITY_POOL_FAILED,

	MSG_SUBMIT_PARAM_CHANGE_PROPOSAL_CREATED,
	MSG_SUBMIT_PARAM_CHANGE_PROPOSAL_FAILED,
	MSG_SUBMIT_COMMUNITY_POOL_SPEND_PROPOSAL_CREATED,
	MSG_SUBMIT_COMMUNITY_POOL_SPEND_PROPOSAL_FAILED,
	MSG_SUBMIT_SOFTWARE_UPGRADE_PROPOSAL_CREATED,
	MSG_SUBMIT_SOFTWARE_UPGRADE_PROPOSAL_FAILED,
	MSG_SUBMIT_CANCEL_SOFTWARE_UPGRADE_PROPOSAL_CREATED,
	MSG_SUBMIT_CANCEL_SOFTWARE_UPGRADE_PROPOSAL_FAILED,
	MSG_SUBMIT_TEXT_PROPOSAL_CREATED,
	MSG_SUBMIT_TEXT_PROPOSAL_FAILED,
	MSG_DEPOSIT_CREATED,
	MSG_DEPOSIT_FAILED,
	MSG_VOTE_CREATED,
	MSG_VOTE_FAILED,

	MSG_CREATE_VALIDATOR_CREATED,
	MSG_CREATE_VALIDATOR_FAILED,
	MSG_EDIT_VALIDATOR_CREATED,
	MSG_EDIT_VALIDATOR_FAILED,
	MSG_DELEGATE_CREATED,
	MSG_DELEGATE_FAILED,
	MSG_UNDELEGATE_CREATED,
	MSG_UNDELEGATE_FAILED,
	MSG_BEGIN_REDELEGATE_CREATED,
	MSG_BEGIN_REDELEGATE_FAILED,

	MSG_UNJAIL_CREATED,
	MSG_UNJAIL_FAILED,

	MSG_NFT_ISSUE_DENOM_CREATED,
	MSG_NFT_ISSUE_DENOM_FAILED,
	MSG_NFT_MINT_NFT_CREATED,
	MSG_NFT_MINT_NFT_FAILED,
	MSG_NFT_TRANSFER_NFT_CREATED,
	MSG_NFT_TRANSFER_NFT_FAILED,
	MSG_NFT_EDIT_NFT_CREATED,
	MSG_NFT_EDIT_NFT_FAILED,
	MSG_NFT_BURN_NFT_CREATED,
	MSG_NFT_BURN_NFT_FAILED,
}
