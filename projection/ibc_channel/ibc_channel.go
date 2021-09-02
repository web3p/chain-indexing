package ibc_channel

import (
	"fmt"

	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbprojectionbase"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	entity_projection "github.com/crypto-com/chain-indexing/entity/projection"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
	"github.com/crypto-com/chain-indexing/internal/utctime"
	ibc_channel_view "github.com/crypto-com/chain-indexing/projection/ibc_channel/view"
	"github.com/crypto-com/chain-indexing/usecase/coin"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
)

var _ entity_projection.Projection = &IBCChannel{}

type IBCChannel struct {
	*rdbprojectionbase.Base

	rdbConn rdb.Conn
	logger  applogger.Logger
}

func NewIBCChannel(logger applogger.Logger, rdbConn rdb.Conn) *IBCChannel {
	return &IBCChannel{
		rdbprojectionbase.NewRDbBase(rdbConn.ToHandle(), "IBCChannel"),

		rdbConn,
		logger,
	}
}

func (_ *IBCChannel) GetEventsToListen() []string {
	return []string{
		event_usecase.BLOCK_CREATED,

		event_usecase.MSG_IBC_CREATE_CLIENT_CREATED,
		event_usecase.MSG_IBC_CONNECTION_OPEN_ACK_CREATED,
		event_usecase.MSG_IBC_CONNECTION_OPEN_CONFIRM_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_INIT_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_TRY_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_ACK_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_CONFIRM_CREATED,
		event_usecase.MSG_IBC_TRANSFER_TRANSFER_CREATED,
		event_usecase.MSG_IBC_RECV_PACKET_CREATED,
		event_usecase.MSG_IBC_ACKNOWLEDGEMENT_CREATED,
	}
}

func (projection *IBCChannel) OnInit() error {
	return nil
}

func (projection *IBCChannel) HandleEvents(height int64, events []event_entity.Event) error {
	rdbTx, err := projection.rdbConn.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = rdbTx.Rollback()
		}
	}()

	rdbTxHandle := rdbTx.ToHandle()

	ibcChannelsView := ibc_channel_view.NewIBCChannels(rdbTxHandle)
	ibcClientsView := ibc_channel_view.NewIBCClients(rdbTxHandle)
	ibcConnectionsView := ibc_channel_view.NewIBCConnections(rdbTxHandle)

	// Get the block time of current height
	var blockTime utctime.UTCTime
	for _, event := range events {
		if blockCreatedEvent, ok := event.(*event_usecase.BlockCreated); ok {
			blockTime = blockCreatedEvent.Block.Time
		}
	}

	// NOTES: Why four channel open events are all needed?
	//
	// PacketOrdering info can only be found in MsgIBCChannelOpenInit and MsgIBCChannelOpenTry.
	//
	// Full information of ConnectionID, CounterpartyChannelID, and CounterpartyPortID
	// can only be found in MsgIBCChannelOpenAck and MsgIBCChannelOpenConfirm.

	for _, event := range events {
		if msgIBCCreateClient, ok := event.(*event_usecase.MsgIBCCreateClient); ok {

			client := &ibc_channel_view.IBCClientRow{
				ClientID:            msgIBCCreateClient.Params.ClientID,
				CounterpartyChainID: msgIBCCreateClient.Params.MaybeTendermintLightClient.TendermintClientState.ChainID,
			}
			if err := ibcClientsView.Insert(client); err != nil {
				return fmt.Errorf("error inserting client: %w", err)
			}

		} else if msgIBCConnectionOpenConfirm, ok := event.(*event_usecase.MsgIBCConnectionOpenConfirm); ok {

			clientID := msgIBCConnectionOpenConfirm.Params.ClientID
			counterpartyChainID, err := ibcClientsView.FindCounterpartyChainIDBy(clientID)
			if err != nil {
				return fmt.Errorf("error in finding counterparty_chain_id: %w", err)
			}

			connection := &ibc_channel_view.IBCConnectionRow{
				ConnectionID:             msgIBCConnectionOpenConfirm.Params.ConnectionID,
				ClientID:                 msgIBCConnectionOpenConfirm.Params.ClientID,
				CounterpartyConnectionID: msgIBCConnectionOpenConfirm.Params.CounterpartyConnectionID,
				CounterpartyClientID:     msgIBCConnectionOpenConfirm.Params.CounterpartyClientID,
				CounterpartyChainID:      counterpartyChainID,
			}
			if err := ibcConnectionsView.Insert(connection); err != nil {
				return fmt.Errorf("error inserting connection: %w", err)
			}

		} else if msgIBCConnectionOpenAck, ok := event.(*event_usecase.MsgIBCConnectionOpenAck); ok {

			clientID := msgIBCConnectionOpenAck.Params.ClientID
			counterpartyChainID, err := ibcClientsView.FindCounterpartyChainIDBy(clientID)
			if err != nil {
				return fmt.Errorf("error in finding counterparty_chain_id: %w", err)
			}

			connection := &ibc_channel_view.IBCConnectionRow{
				ConnectionID:             msgIBCConnectionOpenAck.Params.ConnectionID,
				ClientID:                 msgIBCConnectionOpenAck.Params.ClientID,
				CounterpartyConnectionID: msgIBCConnectionOpenAck.Params.CounterpartyConnectionID,
				CounterpartyClientID:     msgIBCConnectionOpenAck.Params.CounterpartyClientID,
				CounterpartyChainID:      counterpartyChainID,
			}
			if err := ibcConnectionsView.Insert(connection); err != nil {
				return fmt.Errorf("error inserting connection: %w", err)
			}

		} else if msgIBCChannelOpenInit, ok := event.(*event_usecase.MsgIBCChannelOpenInit); ok {

			channel := &ibc_channel_view.IBCChannelRow{
				ChannelID:                    msgIBCChannelOpenInit.Params.ChannelID,
				PortID:                       msgIBCChannelOpenInit.Params.PortID,
				ConnectionID:                 "",
				CounterpartyChannelID:        msgIBCChannelOpenInit.Params.Channel.Counterparty.ChannelID,
				CounterpartyPortID:           msgIBCChannelOpenInit.Params.Channel.Counterparty.PortID,
				CounterpartyChainID:          "",
				Status:                       false,
				PacketOrdering:               msgIBCChannelOpenInit.Params.Channel.Ordering,
				LastInPacketSequence:         0,
				LastOutPacketSequence:        0,
				TotalTransferInCount:         0,
				TotalTransferOutCount:        0,
				TotalTransferOutSuccessCount: 0,
				TotalTransferOutSuccessRate:  0,
				CreatedAtBlockTime:           utctime.FromUnixNano(0),
				CreatedAtBlockHeight:         0,
				Verified:                     false,
				Description:                  "",
				LastActivityBlockTime:        utctime.FromUnixNano(0),
				LastActivityBlockHeight:      0,
				BondedTokens:                 *ibc_channel_view.NewEmptyBondedTokens(),
			}
			if err := ibcChannelsView.Insert(channel); err != nil {
				return fmt.Errorf("error inserting channel: %w", err)
			}

		} else if msgIBCChannelOpenTry, ok := event.(*event_usecase.MsgIBCChannelOpenTry); ok {

			channel := &ibc_channel_view.IBCChannelRow{
				ChannelID:                    msgIBCChannelOpenTry.Params.ChannelID,
				PortID:                       msgIBCChannelOpenTry.Params.PortID,
				ConnectionID:                 msgIBCChannelOpenTry.Params.ConnectionID,
				CounterpartyChannelID:        msgIBCChannelOpenTry.Params.Channel.Counterparty.ChannelID,
				CounterpartyPortID:           msgIBCChannelOpenTry.Params.Channel.Counterparty.PortID,
				CounterpartyChainID:          "",
				Status:                       false,
				PacketOrdering:               msgIBCChannelOpenTry.Params.Channel.Ordering,
				LastInPacketSequence:         0,
				LastOutPacketSequence:        0,
				TotalTransferInCount:         0,
				TotalTransferOutCount:        0,
				TotalTransferOutSuccessCount: 0,
				TotalTransferOutSuccessRate:  0,
				CreatedAtBlockTime:           utctime.FromUnixNano(0),
				CreatedAtBlockHeight:         0,
				Verified:                     false,
				Description:                  "",
				LastActivityBlockTime:        utctime.FromUnixNano(0),
				LastActivityBlockHeight:      0,
				BondedTokens:                 *ibc_channel_view.NewEmptyBondedTokens(),
			}
			if err := ibcChannelsView.Insert(channel); err != nil {
				return fmt.Errorf("error inserting channel: %w", err)
			}

		} else if msgIBCChannelOpenAck, ok := event.(*event_usecase.MsgIBCChannelOpenAck); ok {

			connectionID := msgIBCChannelOpenAck.Params.ConnectionID
			counterpartyChainID, err := ibcConnectionsView.FindCounterpartyChainIDBy(connectionID)
			if err != nil {
				return fmt.Errorf("error in finding counterparty_chain_id: %w", err)
			}

			channel := &ibc_channel_view.IBCChannelRow{
				ChannelID:             msgIBCChannelOpenAck.Params.ChannelID,
				PortID:                msgIBCChannelOpenAck.Params.PortID,
				ConnectionID:          msgIBCChannelOpenAck.Params.ConnectionID,
				CounterpartyChannelID: msgIBCChannelOpenAck.Params.CounterpartyChannelID,
				CounterpartyPortID:    msgIBCChannelOpenAck.Params.CounterpartyPortID,
				CounterpartyChainID:   counterpartyChainID,
				CreatedAtBlockTime:    blockTime,
				CreatedAtBlockHeight:  height,
			}
			if err := ibcChannelsView.UpdateFactualColumns(channel); err != nil {
				return fmt.Errorf("error updating channel: %w", err)
			}

			if err := ibcChannelsView.UpdateStatus(channel.ChannelID, true); err != nil {
				return fmt.Errorf("error updating channel status: %w", err)
			}

		} else if msgIBCChannelOpenConfirm, ok := event.(*event_usecase.MsgIBCChannelOpenConfirm); ok {

			connectionID := msgIBCChannelOpenConfirm.Params.ConnectionID
			counterpartyChainID, err := ibcConnectionsView.FindCounterpartyChainIDBy(connectionID)
			if err != nil {
				return fmt.Errorf("error in finding counterparty_chain_id: %w", err)
			}

			channel := &ibc_channel_view.IBCChannelRow{
				ChannelID:             msgIBCChannelOpenConfirm.Params.ChannelID,
				PortID:                msgIBCChannelOpenConfirm.Params.PortID,
				ConnectionID:          msgIBCChannelOpenConfirm.Params.ConnectionID,
				CounterpartyChannelID: msgIBCChannelOpenConfirm.Params.CounterpartyChannelID,
				CounterpartyPortID:    msgIBCChannelOpenConfirm.Params.CounterpartyPortID,
				CounterpartyChainID:   counterpartyChainID,
				CreatedAtBlockTime:    blockTime,
				CreatedAtBlockHeight:  height,
			}
			if err := ibcChannelsView.UpdateFactualColumns(channel); err != nil {
				return fmt.Errorf("error updating channel: %w", err)
			}

			if err := ibcChannelsView.UpdateStatus(channel.ChannelID, true); err != nil {
				return fmt.Errorf("error updating channel status: %w", err)
			}

		} else if msgIBCTransferTransfer, ok := event.(*event_usecase.MsgIBCTransferTransfer); ok {

			// Transfer started by source chain
			channelID := msgIBCTransferTransfer.Params.SourceChannel

			// TotalTransferOutSuccessRate
			if err := ibcChannelsView.Increment(channelID, "total_transfer_out_count", 1); err != nil {
				return fmt.Errorf("error increasing total_transfer_out_count: %w", err)
			}
			if err := ibcChannelsView.UpdateTotalTransferOutSuccessRate(channelID); err != nil {
				return fmt.Errorf("error updating total_transfer_out_success_rate: %w", err)
			}

			lastOutPacketSequence := msgIBCTransferTransfer.Params.PacketSequence
			if err := ibcChannelsView.UpdateSequence(channelID, "last_out_packet_sequence", lastOutPacketSequence); err != nil {
				return fmt.Errorf("error updating last_out_packet_sequence: %w", err)
			}

			if err := ibcChannelsView.UpdateLastActivityTimeAndHeight(channelID, blockTime, height); err != nil {
				return fmt.Errorf("error updating channel last_activity_time: %w", err)
			}

		} else if msgIBCRecvPacket, ok := event.(*event_usecase.MsgIBCRecvPacket); ok {

			// Transfer started by destination chain
			channelID := msgIBCRecvPacket.Params.Packet.DestinationChannel

			if err := ibcChannelsView.Increment(channelID, "total_transfer_in_count", 1); err != nil {
				return fmt.Errorf("error increasing total_transfer_in_count: %w", err)
			}

			lastInPacketSequence := msgIBCRecvPacket.Params.PacketSequence
			if err := ibcChannelsView.UpdateSequence(channelID, "last_in_packet_sequence", lastInPacketSequence); err != nil {
				return fmt.Errorf("error updating last_in_packet_sequence: %w", err)
			}

			if err := ibcChannelsView.UpdateLastActivityTimeAndHeight(channelID, blockTime, height); err != nil {
				return fmt.Errorf("error updating channel last_activity_time: %w", err)
			}

			// TODO:	should use the below checking when we fix the `success value inverted bug`
			//				in `ibcmsg.go -> ParseMsgRecvPacket()`
			// if msgIBCRecvPacket.Params.MaybeFungibleTokenPacketData.Success {
			if msgIBCRecvPacket.Params.PacketAck.MaybeError == nil {

				amount := msgIBCRecvPacket.Params.MaybeFungibleTokenPacketData.Amount.Uint64()
				denom := msgIBCRecvPacket.Params.MaybeFungibleTokenPacketData.Denom
				destinationChannelID := msgIBCRecvPacket.Params.Packet.DestinationChannel
				destinationPortID := msgIBCRecvPacket.Params.Packet.DestinationPort

				if err := projection.updateBondedTokensWhenMsgIBCRecvPacket(
					ibcChannelsView,
					channelID,
					amount,
					denom,
					destinationChannelID,
					destinationPortID,
				); err != nil {
					return fmt.Errorf("error updateChannelBondedTokensWhenMsgIBCRecvPacket: %w", err)
				}

			}

		} else if msgIBCAcknowledgement, ok := event.(*event_usecase.MsgIBCAcknowledgement); ok {

			// Transfer started by source chain
			channelID := msgIBCAcknowledgement.Params.Packet.SourceChannel

			lastOutPacketSequence := msgIBCAcknowledgement.Params.PacketSequence
			if err := ibcChannelsView.UpdateSequence(channelID, "last_out_packet_sequence", lastOutPacketSequence); err != nil {
				return fmt.Errorf("error updating last_out_packet_sequence: %w", err)
			}

			if err := ibcChannelsView.UpdateLastActivityTimeAndHeight(channelID, blockTime, height); err != nil {
				return fmt.Errorf("error updating channel last_activity_time: %w", err)
			}

			if msgIBCAcknowledgement.Params.MaybeFungibleTokenPacketData.Success {

				// TotalTransferOutSuccessRate
				if err := ibcChannelsView.Increment(channelID, "total_transfer_out_success_count", 1); err != nil {
					return fmt.Errorf("error increasing total_transfer_out_success_count: %w", err)
				}
				if err := ibcChannelsView.UpdateTotalTransferOutSuccessRate(channelID); err != nil {
					return fmt.Errorf("error updating total_transfer_out_success_rate: %w", err)
				}

				amount := msgIBCAcknowledgement.Params.MaybeFungibleTokenPacketData.Amount.Uint64()
				denom := msgIBCAcknowledgement.Params.MaybeFungibleTokenPacketData.Denom
				destinationChannelID := msgIBCAcknowledgement.Params.Packet.DestinationChannel
				destinationPortID := msgIBCAcknowledgement.Params.Packet.DestinationPort

				if err := projection.updateBondedTokensWhenMsgIBCAcknowledgement(
					ibcChannelsView,
					channelID,
					amount,
					denom,
					destinationChannelID,
					destinationPortID,
				); err != nil {
					return fmt.Errorf("error updateChannelBondedTokensWhenMsgIBCAcknowledgement: %w", err)
				}

			}

		}
	}

	if err := projection.UpdateLastHandledEventHeight(rdbTxHandle, height); err != nil {
		return fmt.Errorf("error updating last handled event height: %v", err)
	}

	if err := rdbTx.Commit(); err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}
	committed = true

	return nil
}

func (projection *IBCChannel) updateBondedTokensWhenMsgIBCRecvPacket(
	ibcChannelsView *ibc_channel_view.IBCChannels,
	channelID string,
	amount uint64,
	denom string,
	destinationChannelID string,
	destinationPortID string,
) error {

	bondedTokens, err := ibcChannelsView.FindBondedTokensBy(channelID)
	if err != nil {
		return fmt.Errorf("error finding channel bonded_tokens: %w", err)
	}

	amountInCoinInt := coin.NewIntFromUint64(amount)

	if projection.tokenOriginFromThisChain(bondedTokens, denom) {
		// This token is from this chain, now it is sent back.
		// Subtract it from the bondedTokens.OnCouterpartyChain
		token := ibc_channel_view.NewBondedToken(denom, amountInCoinInt)
		projection.subtractTokenOnCounterpartyChain(bondedTokens, token)
	} else {
		// This token is NOT from this chain.
		// Add it to bondedTokens.OnThisChain
		newDenom := fmt.Sprintf("%s/%s/%s", destinationPortID, destinationChannelID, denom)
		token := ibc_channel_view.NewBondedToken(newDenom, amountInCoinInt)
		projection.addTokenOnThisChain(bondedTokens, token)
	}

	if err := ibcChannelsView.UpdateBondedTokens(channelID, bondedTokens); err != nil {
		return fmt.Errorf("error ibcChannelsView.UpdateBondedTokens: %w", err)
	}

	return nil
}

func (projection *IBCChannel) updateBondedTokensWhenMsgIBCAcknowledgement(
	ibcChannelsView *ibc_channel_view.IBCChannels,
	channelID string,
	amount uint64,
	denom string,
	destinationChannelID string,
	destinationPortID string,
) error {

	bondedTokens, err := ibcChannelsView.FindBondedTokensBy(channelID)
	if err != nil {
		return fmt.Errorf("error ibcChannelsView.FindBondedTokensBy: %w", err)
	}

	amountInCoinInt := coin.NewIntFromUint64(amount)

	if projection.tokenOriginFromCounterpartyChain(bondedTokens, denom) {
		// This token is from the counterparty chain, now it is sent back to counterparty chain.
		// Subtract it from the bondedTokens.OnThisChain
		token := ibc_channel_view.NewBondedToken(denom, amountInCoinInt)
		projection.subtractTokenOnThisChain(bondedTokens, token)
	} else {
		// This token is NOT from the counterparty chain.
		// Add it to bondedTokens.OnCounterpartyChain
		newDenom := fmt.Sprintf("%s/%s/%s", destinationPortID, destinationChannelID, denom)
		token := ibc_channel_view.NewBondedToken(newDenom, amountInCoinInt)
		projection.addOnCounterpartyChain(bondedTokens, token)
	}

	if err := ibcChannelsView.UpdateBondedTokens(channelID, bondedTokens); err != nil {
		return fmt.Errorf("error ibcChannelsView.UpdateBondedTokens: %w", err)
	}

	return nil
}

func (projection *IBCChannel) tokenOriginFromCounterpartyChain(bondedTokens *ibc_channel_view.BondedTokens, denom string) bool {
	for _, token := range bondedTokens.OnThisChain {
		if token.Denom == denom {
			return true
		}
	}
	return false
}

func (projection *IBCChannel) tokenOriginFromThisChain(bondedTokens *ibc_channel_view.BondedTokens, denom string) bool {
	for _, token := range bondedTokens.OnCounterpartyChain {
		if token.Denom == denom {
			return true
		}
	}
	return false
}

func (projection *IBCChannel) subtractTokenOnThisChain(
	bondedTokens *ibc_channel_view.BondedTokens,
	newToken *ibc_channel_view.BondedToken,
) {
	for i, token := range bondedTokens.OnThisChain {
		if token.Denom == newToken.Denom {
			bondedTokens.OnThisChain[i].Amount = bondedTokens.OnThisChain[i].Amount.Sub(newToken.Amount)
			return
		}
	}
}

func (projection *IBCChannel) subtractTokenOnCounterpartyChain(
	bondedTokens *ibc_channel_view.BondedTokens,
	newToken *ibc_channel_view.BondedToken,
) {
	for i, token := range bondedTokens.OnCounterpartyChain {
		if token.Denom == newToken.Denom {
			bondedTokens.OnCounterpartyChain[i].Amount = bondedTokens.OnCounterpartyChain[i].Amount.Sub(newToken.Amount)
			return
		}
	}
}

func (projection *IBCChannel) addOnCounterpartyChain(
	bondedTokens *ibc_channel_view.BondedTokens,
	newToken *ibc_channel_view.BondedToken,
) {
	for i, token := range bondedTokens.OnCounterpartyChain {
		if token.Denom == newToken.Denom {
			// This token already has a record on bondedTokens.OnCounterpartyChain
			bondedTokens.OnCounterpartyChain[i].Amount = bondedTokens.OnCounterpartyChain[i].Amount.Add(newToken.Amount)
			return
		}
	}
	// This token has NO record on bondedTokens.OnCounterpartyChain yet
	bondedTokens.OnCounterpartyChain = append(bondedTokens.OnCounterpartyChain, *newToken)
}

func (projection *IBCChannel) addTokenOnThisChain(
	bondedTokens *ibc_channel_view.BondedTokens,
	newToken *ibc_channel_view.BondedToken,
) {
	for i, token := range bondedTokens.OnThisChain {
		if token.Denom == newToken.Denom {
			// This token already has a record on bondedTokens.OnThisChain
			bondedTokens.OnThisChain[i].Amount = bondedTokens.OnThisChain[i].Amount.Add(newToken.Amount)
			return
		}
	}
	// This token has NO record on bondedTokens.OnThisChain yet
	bondedTokens.OnThisChain = append(bondedTokens.OnThisChain, *newToken)
}
