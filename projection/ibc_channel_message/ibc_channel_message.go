package ibc_channel_message

import (
	"fmt"
	"strconv"

	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbprojectionbase"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	projection_entity "github.com/crypto-com/chain-indexing/entity/projection"
	applogger "github.com/crypto-com/chain-indexing/internal/logger"
	"github.com/crypto-com/chain-indexing/internal/utctime"
	"github.com/crypto-com/chain-indexing/projection/ibc_channel_message/view"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
)

var _ projection_entity.Projection = &IBCChannelMessage{}

type IBCChannelMessage struct {
	*rdbprojectionbase.Base

	rdbConn rdb.Conn
	logger  applogger.Logger
}

func NewIBCChannelMessage(
	logger applogger.Logger,
	rdbConn rdb.Conn,
) *IBCChannelMessage {
	return &IBCChannelMessage{
		rdbprojectionbase.NewRDbBase(rdbConn.ToHandle(), "IBCChannelMessage"),

		rdbConn,
		logger,
	}
}

func (_ *IBCChannelMessage) GetEventsToListen() []string {
	return []string{
		event_usecase.BLOCK_CREATED,

		event_usecase.MSG_IBC_CHANNEL_OPEN_INIT_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_TRY_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_ACK_CREATED,
		event_usecase.MSG_IBC_CHANNEL_OPEN_CONFIRM_CREATED,
		event_usecase.MSG_IBC_TRANSFER_TRANSFER_CREATED,
		event_usecase.MSG_IBC_RECV_PACKET_CREATED,
		event_usecase.MSG_IBC_ACKNOWLEDGEMENT_CREATED,
		event_usecase.MSG_IBC_TIMEOUT_CREATED,
		event_usecase.MSG_IBC_TIMEOUT_ON_CLOSE_CREATED,
	}
}

func (projection *IBCChannelMessage) OnInit() error {
	return nil
}

func (projection *IBCChannelMessage) HandleEvents(height int64, events []event_entity.Event) error {
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

	ibcChannelMessagesView := view.NewIBCChannelMessages(rdbTxHandle)
	ibcChannelMessagesTotalView := view.NewIBCChannelMessagesTotal(rdbTxHandle)

	// Get the block time of current height
	var blockTime utctime.UTCTime
	for _, event := range events {
		if blockCreatedEvent, ok := event.(*event_usecase.BlockCreated); ok {
			blockTime = blockCreatedEvent.Block.Time
		}
	}

	var messages []view.IBCChannelMessageRow
	for _, event := range events {
		if typedEvent, ok := event.(*event_usecase.MsgIBCChannelOpenInit); ok {
			channelID := typedEvent.Params.ChannelID

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          "",
				Receiver:        "",
				Denom:           "",
				Amount:          "",
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCChannelOpenTry); ok {

			channelID := typedEvent.Params.ChannelID

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          "",
				Receiver:        "",
				Denom:           "",
				Amount:          "",
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCChannelOpenAck); ok {

			channelID := typedEvent.Params.ChannelID

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          "",
				Receiver:        "",
				Denom:           "",
				Amount:          "",
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCChannelOpenConfirm); ok {

			channelID := typedEvent.Params.ChannelID

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          "",
				Receiver:        "",
				Denom:           "",
				Amount:          "",
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCTransferTransfer); ok {

			channelID := typedEvent.Params.SourceChannel

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          typedEvent.Params.Sender,
				Receiver:        typedEvent.Params.Receiver,
				Denom:           typedEvent.Params.Token.Denom,
				Amount:          strconv.FormatUint(typedEvent.Params.Token.Amount, 10),
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCRecvPacket); ok {

			channelID := typedEvent.Params.Packet.DestinationChannel

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         typedEvent.Params.MaybeFungibleTokenPacketData.Success,
				Sender:          typedEvent.Params.MaybeFungibleTokenPacketData.Sender,
				Receiver:        typedEvent.Params.MaybeFungibleTokenPacketData.Receiver,
				Denom:           typedEvent.Params.MaybeFungibleTokenPacketData.Denom,
				Amount:          typedEvent.Params.MaybeFungibleTokenPacketData.Amount.String(),
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			if typedEvent.Params.PacketAck.MaybeError != nil {
				message.Error = *typedEvent.Params.PacketAck.MaybeError
			}

			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCAcknowledgement); ok {

			channelID := typedEvent.Params.Packet.SourceChannel

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         typedEvent.Params.MaybeFungibleTokenPacketData.Success,
				Sender:          typedEvent.Params.MaybeFungibleTokenPacketData.Sender,
				Receiver:        typedEvent.Params.MaybeFungibleTokenPacketData.Receiver,
				Denom:           typedEvent.Params.MaybeFungibleTokenPacketData.Denom,
				Amount:          typedEvent.Params.MaybeFungibleTokenPacketData.Amount.String(),
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			if typedEvent.Params.MaybeFungibleTokenPacketData.MaybeError != nil {
				message.Error = *typedEvent.Params.MaybeFungibleTokenPacketData.MaybeError
			}

			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCTimeout); ok {

			channelID := typedEvent.Params.Packet.SourceChannel

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          "",
				Receiver:        typedEvent.Params.MaybeMsgTransfer.RefundReceiver,
				Denom:           typedEvent.Params.MaybeMsgTransfer.RefundDenom,
				Amount:          strconv.FormatUint(typedEvent.Params.MaybeMsgTransfer.RefundAmount, 10),
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		} else if typedEvent, ok := event.(*event_usecase.MsgIBCTimeoutOnClose); ok {

			channelID := typedEvent.Params.Packet.SourceChannel

			message := view.IBCChannelMessageRow{
				ChannelID:       channelID,
				BlockHeight:     height,
				BlockTime:       blockTime,
				TransactionHash: typedEvent.TxHash(),
				Success:         true,
				Error:           "",
				Sender:          "",
				Receiver:        typedEvent.Params.MaybeMsgTransfer.RefundReceiver,
				Denom:           typedEvent.Params.MaybeMsgTransfer.RefundDenom,
				Amount:          strconv.FormatUint(typedEvent.Params.MaybeMsgTransfer.RefundAmount, 10),
				MessageType:     typedEvent.MsgName,
				Message:         typedEvent,
			}
			messages = append(messages, message)

		}
	}

	for i, message := range messages {
		if err := ibcChannelMessagesView.Insert(&messages[i]); err != nil {
			return fmt.Errorf("error inserting IBCChannelMessage: %v", err)
		}

		if err := projection.updateIBCChannelMessagesTotal(
			ibcChannelMessagesTotalView,
			message.ChannelID,
			message.MessageType,
		); err != nil {
			return fmt.Errorf("error updating IBCChannelMessageTotal: %v", err)
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

func (projection *IBCChannelMessage) updateIBCChannelMessagesTotal(
	totalView *view.IBCChannelMessagesTotal,
	channelID string,
	messageType string,
) error {
	if err := totalView.Increment(
		fmt.Sprintf("%s:-", channelID), 1,
	); err != nil {
		return fmt.Errorf("error incrementing total message of IBCChannel: %v", err)
	}
	if err := totalView.Increment(
		fmt.Sprintf("%s:%s", channelID, messageType), 1,
	); err != nil {
		return fmt.Errorf("error incrementing total message of IBCChannel: %v", err)
	}
	return nil
}
