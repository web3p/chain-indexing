package rdbparambase

import (
	"fmt"
	"strconv"

	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbparambase/types"
	"github.com/crypto-com/chain-indexing/appinterface/projection/rdbparambase/view"
	"github.com/crypto-com/chain-indexing/appinterface/rdb"
	event_entity "github.com/crypto-com/chain-indexing/entity/event"
	"github.com/crypto-com/chain-indexing/internal/json"
	event_usecase "github.com/crypto-com/chain-indexing/usecase/event"
	"github.com/crypto-com/chain-indexing/usecase/model/genesis"
)

// a generic Param projection. For table schema refer to view/params.go
type ParamBase struct {
	tableName string

	paramList []types.ParamAccessor
}

func NewParamBase(tableName string, paramList []types.ParamAccessor) *ParamBase {
	return &ParamBase{
		tableName,

		paramList,
	}
}

// Handle
func (projection *ParamBase) HandleEvent(conn *rdb.Handle, event event_entity.Event) error {
	// TODO: support ParamChange proposal
	genesisCreatedEvent, ok := event.(*event_usecase.GenesisCreated)
	if !ok {
		return nil
	}

	view := view.NewParams(conn, projection.tableName)
	for _, param := range projection.paramList {
		if err := projection.persistGenesisParam(view, &genesisCreatedEvent.Genesis, param); err != nil {
			return err
		}
	}

	return nil
}

func (projection *ParamBase) persistGenesisParam(
	view *view.Params, genesis *genesis.Genesis, param types.ParamAccessor,
) error {
	var value string
	switch key := param.Module + param.Key; key {
	case "auth.max_memo_characters":
		value = genesis.AppState.Auth.Params.MaxMemoCharacters
	case "auth.tx_sig_limit":
		value = genesis.AppState.Auth.Params.TxSigLimit
	case "auth.tx_size_cost_per_byte":
		value = genesis.AppState.Auth.Params.TxSizeCostPerByte
	case "auth.sig_verify_cost_ed25519":
		value = genesis.AppState.Auth.Params.SigVerifyCostEd25519
	case "auth.sig_verify_cost_secp256k1":
		value = genesis.AppState.Auth.Params.SigVerifyCostSecp256K1

	case "bank.send_enabled":
		value = json.MustMarshalToString(genesis.AppState.Bank.Params.SendEnabled)
	case "bank.default_send_enabled":
		value = strconv.FormatBool(genesis.AppState.Bank.Params.DefaultSendEnabled)

	case "distribution.base_proposer_reward":
		value = genesis.AppState.Distribution.Params.BaseProposerReward
	case "distribution.bonus_proposer_reward":
		value = genesis.AppState.Distribution.Params.BonusProposerReward
	case "distribution.community_tax":
		value = genesis.AppState.Distribution.Params.CommunityTax
	case "distribution.withdraw_addr_enabled":
		value = strconv.FormatBool(genesis.AppState.Distribution.Params.WithdrawAddrEnabled)

	case "gov.min_deposit":
		value = json.MustMarshalToString(genesis.AppState.Gov.DepositParams.MinDeposit)
	case "gov.max_deposit_period":
		value = genesis.AppState.Gov.DepositParams.MaxDepositPeriod
	case "gov.voting_period":
		value = genesis.AppState.Gov.VotingParams.VotingPeriod
	case "gov.quorum":
		value = genesis.AppState.Gov.TallyParams.Quorum
	case "gov.threshold":
		value = genesis.AppState.Gov.TallyParams.Threshold
	case "gov.VetoThreshold":
		value = genesis.AppState.Gov.TallyParams.VetoThreshold

	case "mint.blocks_per_year":
		value = genesis.AppState.Mint.Params.BlocksPerYear
	case "mint.goal_bonded":
		value = genesis.AppState.Mint.Params.GoalBonded
	case "mint.inflation_max":
		value = genesis.AppState.Mint.Params.InflationMax
	case "mint.inflation_min":
		value = genesis.AppState.Mint.Params.InflationMin
	case "mint.inflation_rate_change":
		value = genesis.AppState.Mint.Params.InflationRateChange
	case "mint.mint_denom":
		value = genesis.AppState.Mint.Params.MintDenom

	case "slashing.downtime_jail_duration":
		value = genesis.AppState.Slashing.Params.DowntimeJailDuration
	case "slashing.min_signed_per_window":
		value = genesis.AppState.Slashing.Params.MinSignedPerWindow
	case "slashing.signed_blocks_window":
		value = genesis.AppState.Slashing.Params.SignedBlocksWindow
	case "slashing.slash_fraction_double_sign":
		value = genesis.AppState.Slashing.Params.SlashFractionDoubleSign
	case "slashing.slash_fraction_downtime":
		value = genesis.AppState.Slashing.Params.SlashFractionDowntime

	case "staking.bond_denom":
		value = genesis.AppState.Staking.Params.BondDenom
	case "staking.historical_entries":
		value = strconv.FormatInt(genesis.AppState.Staking.Params.HistoricalEntries, 10)
	case "staking.max_entries":
		value = strconv.FormatInt(genesis.AppState.Staking.Params.MaxEntries, 10)
	case "staking.max_validators":
		value = strconv.FormatInt(genesis.AppState.Staking.Params.MaxValidators, 10)
	case "staking.unbonding_time":
		value = genesis.AppState.Staking.Params.UnbondingTime

	case "ibc_transfer.receive_enabled":
		value = strconv.FormatBool(genesis.AppState.Transfer.Params.ReceiveEnabled)
	case "ibc_transfer.send_enabled":
		value = strconv.FormatBool(genesis.AppState.Transfer.Params.SendEnabled)

	default:
		return fmt.Errorf("unrecognized param: %s.%s", param.Module, param.Key)
	}

	if err := view.Set(param, value); err != nil {
		return fmt.Errorf("error persisting genesis param %s.%s: %v", param.Module, param.Key, err)
	}

	return nil
}
