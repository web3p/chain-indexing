package rdbprojectionbase

import (
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"

	"github.com/crypto-com/chain-indexing/appinterface/projection"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	"github.com/crypto-com/chain-indexing/internal/primptr"
)

const DEFAULT_TABLE = "projections"

// Table should have the following schema
// | Field                     | Data Type | Constraint  |
// | ------------------------- | --------- | ----------- |
// | id                        | VARCHAR   | PRIMARY KEY |
// | last_handled_event_height | INT64     | NOT NULL    |

// Base is a bas for projection which keeps track of last handled event height using relational
// database. It implements Id() and GetLastHandledEventHeight() of projection interface.
type Base struct {
	rdbHandle *rdb.Handle
	store     *Store
	config    interface{}

	projectionId string
}

// Create a new Base using table name in the RDb to keep the projection handling records
func NewRDbBase(rdbHandle *rdb.Handle, projectionId string) *Base {
	return NewRDbBaseWithOptions(rdbHandle, projectionId, Options{
		MaybeTable:     primptr.String(DEFAULT_TABLE),
		MaybeConfigPtr: nil,
	})
}

func NewRDbBaseWithOptions(rdbHandle *rdb.Handle, projectionId string, options Options) *Base {
	table := DEFAULT_TABLE
	if options.MaybeTable != nil {
		table = *options.MaybeTable
	}

	base := &Base{
		rdbHandle: rdbHandle,
		store:     NewStore(table),

		projectionId: projectionId,
	}
	if options.MaybeConfigPtr != nil {
		decoderConfig := &mapstructure.DecoderConfig{
			WeaklyTypedInput: true,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				mapstructure.StringToTimeHookFunc(time.RFC3339),
			),
			Result: options.MaybeConfigPtr,
		}
		decoder, decoderErr := mapstructure.NewDecoder(decoderConfig)
		if decoderErr != nil {
			panic(fmt.Errorf("error creating projection config decoder: %v", decoderErr))
		}

		projectionConfigId := strcase.ToSnake(projectionId)
		projectionConfig, projectionConfigFound := projection.GlobalConfig.Projection[projectionConfigId]
		if !projectionConfigFound {
			panic(fmt.Sprintf("no projection config found for projection id: %s", projectionConfigId))
		}
		if err := decoder.Decode(projectionConfig); err != nil {
			panic(fmt.Errorf("error decoding projection %s config: %v", projectionId, err))
		}

		base.config = options.MaybeConfigPtr
	}

	return base
}

type Options struct {
	// Customize table name in the RDb too keep the projection handling records
	MaybeTable *string

	// Pointer to the configuration
	MaybeConfigPtr interface{}
}

type TomlConfig struct {
	Projection map[string]interface{}
}

// Implements projection.Id()
func (base *Base) Id() string {
	return base.projectionId
}

func (base *Base) Config() interface{} {
	return base.config
}

func (base *Base) UpdateLastHandledEventHeight(rdbHandle *rdb.Handle, height int64) error {
	if err := base.store.UpdateLastHandledEventHeight(rdbHandle, base.projectionId, height); err != nil {
		return err
	}
	return nil
}

// Implements projection.GetLastHandledEventHeight()
func (base *Base) GetLastHandledEventHeight() (*int64, error) {
	return base.store.GetLastHandledEventHeight(base.rdbHandle, base.projectionId)
}
