package cosmos_gov_v1

import (
	"fmt"
	"time"

	"github.com/crypto-com/chain-indexing/entity/command"
	"github.com/crypto-com/chain-indexing/external/tmcosmosutils"
	"github.com/mitchellh/mapstructure"

	"github.com/crypto-com/chain-indexing/usecase/coin"
	command_usecase "github.com/crypto-com/chain-indexing/usecase/command"
	v1_model "github.com/crypto-com/chain-indexing/usecase/model/v1"
	"github.com/crypto-com/chain-indexing/usecase/parser/utils"

	mapstructure_utils "github.com/crypto-com/chain-indexing/usecase/parser/utils/mapstructure"
)

func ParseMsgDeposit(
	parserParams utils.CosmosParserParams,
) ([]command.Command, []string) {
	// Getting possible signer address from Msg
	var possibleSignerAddresses []string
	if parserParams.Msg != nil {
		if depositor, ok := parserParams.Msg["depositor"]; ok {
			possibleSignerAddresses = append(possibleSignerAddresses, utils.AddressParse(depositor.(string)))
		}
	}

	amountInterface := parserParams.Msg["amount"].([]interface{})
	amount, err := tmcosmosutils.NewCoinsFromAmountInterface(amountInterface)
	if err != nil {
		amount = make([]coin.Coin, 0)
		for i := 0; i < len(amountInterface); i++ {
			amount = append(amount, coin.Coin{})
		}
	}

	cmds := []command.Command{command_usecase.NewCreateMsgDepositV1(
		parserParams.MsgCommonParams,

		v1_model.MsgDepositParams{
			ProposalId: parserParams.Msg["proposal_id"].(string),
			Depositor:  utils.AddressParse(parserParams.Msg["depositor"].(string)),
			Amount:     amount,
		},
	)}

	if !parserParams.MsgCommonParams.TxSuccess {
		return cmds, possibleSignerAddresses
	}

	log := utils.NewParsedTxsResultLog(&parserParams.TxsResult.Log[parserParams.MsgIndex])
	logEvents := log.GetEventsByType("proposal_deposit")
	if logEvents == nil {
		panic("missing `proposal_deposit` event in TxsResult log")
	}

	for _, logEvent := range logEvents {
		if logEvent.HasAttribute("voting_period_start") {
			cmds = append(cmds, command_usecase.NewStartProposalVotingPeriod(
				parserParams.MsgCommonParams.BlockHeight, logEvent.MustGetAttributeByKey("voting_period_start"),
			))
			break
		}
	}

	return cmds, possibleSignerAddresses
}

func ParseMsgSubmitProposal(
	parserParams utils.CosmosParserParams,
) ([]command.Command, []string) {
	var rawMsg v1_model.RawMsgSubmitProposal
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure_utils.StringToDurationHookFunc(),
			mapstructure_utils.StringToByteSliceHookFunc(),
		),
		Result: &rawMsg,
	}
	decoder, decoderErr := mapstructure.NewDecoder(decoderConfig)
	if decoderErr != nil {
		panic(fmt.Errorf("error creating ParseMsgSubmitProposal decoder: %v", decoderErr))
	}
	if err := decoder.Decode(parserParams.Msg); err != nil {
		panic(fmt.Errorf("error decoding ParseMsgSubmitProposal: %v", err))
	}

	rawMsg.Proposer = utils.AddressParse(rawMsg.Proposer)

	var cmds []command.Command
	var possibleSignerAddresses []string

	blockHeight := parserParams.MsgCommonParams.BlockHeight

	msgs, ok := parserParams.Msg["messages"].([]interface{})
	if !ok {
		panic(fmt.Errorf("error parsing MsgSubmitProposal.msgs to []interface{}: %v", parserParams.Msg["messages"]))
	}

	for innerMsgIndex, innerMsgInterface := range msgs {
		innerMsg, ok := innerMsgInterface.(map[string]interface{})
		if !ok {
			panic(fmt.Errorf("error parsing MsgSubmitProposal.msgs[%v] to map[string]interface{}: %v", innerMsgIndex, innerMsgInterface))
		}

		innerMsgType, ok := innerMsg["@type"].(string)
		if !ok {
			panic(fmt.Errorf("error missing '@type' in MsgSubmitProposal.msgs[%v]: %v", innerMsgIndex, innerMsg))
		}

		innerMsg["initial_deposit"] = rawMsg.InitialDeposit
		innerMsg["proposer"] = rawMsg.Proposer
		innerMsg["metadata"] = rawMsg.Metadata

		parser := parserParams.ParserManager.GetParser(utils.CosmosParserKey(innerMsgType), utils.ParserBlockHeight(blockHeight))

		msgCommands, signers := parser(utils.CosmosParserParams{
			AddressPrefix:   parserParams.AddressPrefix,
			StakingDenom:    parserParams.StakingDenom,
			TxsResult:       parserParams.TxsResult,
			MsgCommonParams: parserParams.MsgCommonParams,
			Msg:             innerMsg,
			MsgIndex:        parserParams.MsgIndex,
			ParserManager:   parserParams.ParserManager,
		})

		possibleSignerAddresses = append(possibleSignerAddresses, signers...)
		cmds = append(cmds, msgCommands...)
	}

	return cmds, possibleSignerAddresses
}

func ParseMsgVote(
	parserParams utils.CosmosParserParams,
) ([]command.Command, []string) {

	// Getting possible signer address from Msg
	var possibleSignerAddresses []string
	if parserParams.Msg != nil {
		if voter, ok := parserParams.Msg["voter"]; ok {
			possibleSignerAddresses = append(possibleSignerAddresses, utils.AddressParse(voter.(string)))
		}
	}

	return []command.Command{command_usecase.NewCreateMsgVoteV1(
		parserParams.MsgCommonParams,

		v1_model.MsgVoteParams{
			ProposalId: parserParams.Msg["proposal_id"].(string),
			Voter:      utils.AddressParse(parserParams.Msg["voter"].(string)),
			Option:     parserParams.Msg["option"].(string),
			Metadata:   parserParams.Msg["metadata"].(string),
		},
	)}, possibleSignerAddresses
}

func ParseMsgVoteWeighted(
	parserParams utils.CosmosParserParams,
) ([]command.Command, []string) {
	var rawMsg v1_model.RawMsgVoteWeight
	decoderConfig := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure_utils.StringToDurationHookFunc(),
			mapstructure_utils.StringToByteSliceHookFunc(),
		),
		Result: &rawMsg,
	}
	decoder, decoderErr := mapstructure.NewDecoder(decoderConfig)
	if decoderErr != nil {
		panic(fmt.Errorf("error creating ParseMsgSubmitProposal decoder: %v", decoderErr))
	}
	if err := decoder.Decode(parserParams.Msg); err != nil {
		panic(fmt.Errorf("error decoding ParseMsgSubmitProposal: %v", err))
	}

	var cmds []command.Command
	var possibleSignerAddresses []string

	if rawMsg.Options != nil {
		// Getting possible signer address from Msg
		if voter, ok := parserParams.Msg["voter"]; ok {
			possibleSignerAddresses = append(possibleSignerAddresses, utils.AddressParse(voter.(string)))
		}

		cmds = append(cmds, command_usecase.NewCreateMsgVoteWeighted(
			parserParams.MsgCommonParams,

			v1_model.MsgVoteWeightedParams{
				ProposalId:  parserParams.Msg["proposal_id"].(string),
				Voter:       utils.AddressParse(parserParams.Msg["voter"].(string)),
				VoteOptions: rawMsg.Options,
				Metadata:    parserParams.Msg["metadata"].(string),
			},
		))

	}

	return cmds, possibleSignerAddresses
}
