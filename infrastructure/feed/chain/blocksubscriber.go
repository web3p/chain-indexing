package chain

import (
	"fmt"
	appevent "github.com/crypto-com/chainindex/appinterface/event"
	"github.com/crypto-com/chainindex/infrastructure/feed/chain/parser"
	"github.com/crypto-com/chainindex/infrastructure/notification"
	"github.com/crypto-com/chainindex/usecase/executor"
)

type BlockSubscriber struct {
	Id int64
}

func NewBlockSubscriber(id int64) *BlockSubscriber {
	return &BlockSubscriber{
		Id: id,
	}
}

func (bs *BlockSubscriber) OnNotification(n *notification.BlockNotification, eventStore *appevent.RDbStore) error {
	fmt.Println("processor", bs.Id, "got notification", n.Height)

	// create an executor instance for current height
	executor := executor.NewBlockExecutor(n.Height)
	commands, err := parser.ParseBlockToCommands(n.Block, n.RawBlock, n.BlockResults)
	if err != nil {
		return fmt.Errorf("error parsing block data to commands %v", err)
	}
	executor.AddAllCommands(commands)

	// generate all events, make them persistent
	if err := executor.ExecAllCommands(); err != nil {
		return fmt.Errorf("error generating all events%v", err)
	}
	if err := executor.StoreAllEvents(eventStore); err != nil {
		return fmt.Errorf("error storing all events%v", err)
	}

	return nil
}
