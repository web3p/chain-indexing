package cosmosapp

import cosmosapp_interface "github.com/crypto-com/chain-indexing/appinterface/cosmosapp"

type ValidatorsResp struct {
	MaybeValidatorResponse []cosmosapp_interface.Validator `json:"validators"`
	MaybePagination        *Pagination                     `json:"pagination"`
	// On error
	MaybeCode    *int    `json:"code"`
	MaybeMessage *string `json:"message"`
}

type ValidatorResp struct {
	Validator cosmosapp_interface.Validator `json:"validator"`
}

type StakingPoolResp struct {
	Pool Pool `json:"pool"`
}

type Pool struct {
	NotBondedTokens string `json:"not_bonded_tokens"`
	BondedTokens    string `json:"bonded_tokens"`
}
