package parser

import (
	"fmt"
	"strings"

	commandentity "github.com/crypto-com/chain-indexing/entity/command"
	"github.com/crypto-com/chain-indexing/internal/json"
	"github.com/crypto-com/chain-indexing/internal/typeconv"
	"github.com/crypto-com/chain-indexing/usecase/command"
	"github.com/crypto-com/chain-indexing/usecase/model"
	"github.com/crypto-com/chain-indexing/usecase/parser/utils"
)

func ParseBlockResultsTxsResults(
	block *model.Block,
	blockResults *model.BlockResults,
) ([]commandentity.Command, error) {
	cmds := make([]commandentity.Command, 0)

	for i := range blockResults.TxsResults {
		txHex := block.Txs[i]

		parsedCmds := parseCosmosSendToIBC(block.Height, txHex, &blockResults.TxsResults[i])
		cmds = append(cmds, parsedCmds...)
	}

	return cmds, nil
}

func parseCosmosSendToIBC(
	blockHeight int64,
	txHex string,
	txResults *model.BlockResultsTxsResult,
) []commandentity.Command {
	isEthereumTx := false

	var maybeIBCSendPacketEvent *utils.ParsedTxsResultLogEvent
	for i, event := range txResults.Events {
		if event.Type == "ethereum_tx" {
			isEthereumTx = true
			break
		} else if event.Type == "send_packet" {
			maybeIBCSendPacketEvent = utils.NewParsedTxsResultLogEvent(&txResults.Events[i])
		}
	}

	if !isEthereumTx {
		return nil
	}

	if maybeIBCSendPacketEvent == nil {
		return nil
	}

	params := model.RawCosmosSendToIBCParams{
		PacketChannelOrdering:  maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_channel_ordering"),
		PacketConnection:       maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_connection"),
		PacketData:             maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_data"),
		PacketDataHex:          maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_data_hex"),
		PacketDstChannel:       maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_dst_channel"),
		PacketDstPort:          maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_dst_port"),
		PacketSequence:         maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_sequence"),
		PacketSrcChannel:       maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_src_channel"),
		PacketSrcPort:          maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_src_port"),
		PacketTimeoutHeight:    maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_timeout_height"),
		PacketTimeoutTimestamp: maybeIBCSendPacketEvent.MustGetAttributeByKey("packet_timeout_timestamp"),
	}

	var fungibleTokenPacketData model.FungibleTokenPacketData
	json.MustUnmarshalFromString(params.PacketData, &fungibleTokenPacketData)

	timeoutHeight := mustParseCosmosSendToIBCTimeoutHeight(params.PacketTimeoutHeight)

	cmd := command.NewCreateCosmosSendToIBC(blockHeight, model.CosmosSendToIBCParams{
		TxHash:        TxHash(txHex),
		SourcePort:    params.PacketSrcPort,
		SourceChannel: params.PacketSrcChannel,
		Token: model.CosmosSendToIBCToken{
			Denom:  fungibleTokenPacketData.Denom,
			Amount: fungibleTokenPacketData.Amount,
		},
		Sender:   fungibleTokenPacketData.Sender,
		Receiver: fungibleTokenPacketData.Receiver,
		TimeoutHeight: model.CosmosSendToIBCHeight{
			RevisionNumber: timeoutHeight.RevisionNumber,
			RevisionHeight: timeoutHeight.RevisionHeight,
		},
		TimeoutTimestamp:   params.PacketTimeoutTimestamp,
		PacketSequence:     typeconv.MustAtou64(params.PacketSequence),
		PacketDataHex:      params.PacketDataHex,
		DestinationPort:    params.PacketDstPort,
		DestinationChannel: params.PacketDstChannel,
		ChannelOrdering:    params.PacketChannelOrdering,
		ConnectionID:       params.PacketConnection,
	})

	return []commandentity.Command{cmd}
}

func mustParseCosmosSendToIBCTimeoutHeight(height string) model.CosmosSendToIBCHeight {
	heightTokens := strings.Split(height, "-")
	if len(heightTokens) != 2 {
		panic(fmt.Errorf("invalid height: %s", height))
	}

	revisionNumber := typeconv.MustAtou64(heightTokens[0])
	revisionHeight := typeconv.MustAtou64(heightTokens[1])
	return model.CosmosSendToIBCHeight{
		RevisionNumber: revisionNumber,
		RevisionHeight: revisionHeight,
	}
}