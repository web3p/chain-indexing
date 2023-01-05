package v1_model

import (
	"encoding/json"
	"time"

	"github.com/crypto-com/chain-indexing/external/utctime"
	"github.com/crypto-com/chain-indexing/usecase/coin"
)

type RawMsgSubmitProposal struct {
	Type           string                 `json:"@type"`
	Messages       []MsgSubmitProposalMsg `json:"messages"`
	Proposer       string                 `json:"proposer"`
	InitialDeposit coin.Coins             `json:"initial_deposit"`
	Metadata       string                 `json:"metadata"`
}
type MsgSubmitProposalMsg struct {
	Type string `mapstructure:"@type" json:"@type"`
}

type MsgSubmitCommunityPoolSpendProposalParams struct {
	MaybeProposalId *string                                    `json:"proposalId"`
	Content         MsgSubmitCommunityPoolSpendProposalContent `json:"content"`
	ProposerAddress string                                     `json:"proposerAddress"`
	InitialDeposit  coin.Coins                                 `json:"initialDeposit"`
}
type MsgSubmitCommunityPoolSpendProposalContent struct {
	Type             string     `json:"@type"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	RecipientAddress string     `json:"recipientAddress"`
	Amount           coin.Coins `json:"amount"`
}
type RawMsgSubmitCommunityPoolSpendProposalContent struct {
	Type             string        `json:"@type"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	RecipientAddress string        `json:"recipient"`
	Amount           []interface{} `json:"amount"`
}

type MsgSubmitParamChangeProposalParams struct {
	MaybeProposalId *string                             `json:"proposalId"`
	Content         MsgSubmitParamChangeProposalContent `json:"content"`
	ProposerAddress string                              `json:"proposerAddress"`
	InitialDeposit  coin.Coins                          `json:"initialDeposit"`
}
type MsgSubmitParamChangeProposalContent struct {
	Type        string                               `json:"@type"`
	Title       string                               `json:"title"`
	Description string                               `json:"description"`
	Changes     []MsgSubmitParamChangeProposalChange `json:"changes"`
}
type MsgSubmitParamChangeProposalChange struct {
	Subspace string          `json:"subspace"`
	Key      string          `json:"key"`
	Value    json.RawMessage `json:"value"`
}

// MsgSoftwareUpgrade
type MsgSoftwareUpgradeParams struct {
	MaybeProposalId *string                `json:"proposalId"`
	Authority       string                 `json:"authority"`
	Plan            MsgSoftwareUpgradePlan `json:"plan"`
	Proposer        string                 `json:"proposer"`
	InitialDeposit  coin.Coins             `json:"initial_deposit"`
	Metadata        string                 `json:"metadata"`
}

type RawMsgSoftwareUpgrade struct {
	Type      string                                  `json:"@type"`
	Authority string                                  `json:"authority"`
	Plan      RawMsgSubmitSoftwareUpgradeProposalPlan `json:"plan"`
}

type MsgSoftwareUpgrade struct {
	Type      string                 `json:"@type"`
	Authority string                 `json:"authority"`
	Plan      MsgSoftwareUpgradePlan `json:"plan"`
}

type MsgSoftwareUpgradePlan struct {
	Name                string          `json:"name"`
	Time                utctime.UTCTime `json:"time"`
	Height              int64           `json:"height"`
	Info                string          `json:"info"`
	UpgradedClientState string          `json:"upgraded_client_state"`
}

type RawMsgSubmitSoftwareUpgradeProposalContent struct {
	Type        string                                  `json:"@type"`
	Title       string                                  `json:"title"`
	Description string                                  `json:"description"`
	Plan        RawMsgSubmitSoftwareUpgradeProposalPlan `json:"plan"`
}
type RawMsgSubmitSoftwareUpgradeProposalPlan struct {
	Name                string    `json:"name"`
	Time                time.Time `json:"time"`
	Height              string    `json:"height"`
	Info                string    `json:"info"`
	UpgradedClientState string    `json:"upgraded_client_state"`
}

// MsgSoftwareCancel

type MsgCancelUpgradeParams struct {
	MaybeProposalId *string    `json:"proposalId"`
	Authority       string     `json:"authority"`
	Proposer        string     `json:"proposer"`
	InitialDeposit  coin.Coins `json:"initial_deposit"`
	Metadata        string     `json:"metadata"`
}
type RawMsgCancelUpgrade struct {
	Type      string `json:"@type"`
	Authority string `json:"authority"`
}
type MsgCancelUpgrade struct {
	Type      string `json:"@type"`
	Authority string `json:"authority"`
}

type MsgSubmitCancelSoftwareUpgradeMessages struct {
	Type string `json:"@type"`
}

type MsgSubmitTextProposalParams struct {
	MaybeProposalId *string                      `json:"proposalId"`
	Content         MsgSubmitTextProposalContent `json:"content"`
	ProposerAddress string                       `json:"proposerAddress"`
	InitialDeposit  coin.Coins                   `json:"initialDeposit"`
}
type MsgSubmitTextProposalContent struct {
	Type        string `json:"@type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type MsgSubmitUnknownProposalParams struct {
	MaybeProposalId *string                         `json:"proposalId"`
	Content         MsgSubmitUnknownProposalContent `json:"content"`
	ProposerAddress string                          `json:"proposerAddress"`
	InitialDeposit  coin.Coins                      `json:"initialDeposit"`
}

type MsgSubmitUnknownProposalContent struct {
	Type        string      `json:"@type"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	RawContent  interface{} `json:"rawContent"`
}
