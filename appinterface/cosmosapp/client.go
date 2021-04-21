package cosmosapp

import (
	"errors"

	"github.com/crypto-com/chain-indexing/usecase/coin"
)

type Client interface {
	Account(accountAddress string) (*Account, error)
	Balances(accountAddress string) (coin.Coins, error)
	BalanceByDenom(accountAddress string, denom string) (*coin.Coin, error)
	BondedBalance(accountAddress string) (coin.Coins, error)
	RedelegatingBalance(accountAddress string) (coin.Coins, error)
	UnbondingBalance(accountAddress string) (coin.Coins, error)
	TotalRewards(accountAddress string) (coin.DecCoins, error)
	Commission(validatorAddress string) (coin.DecCoins, error)
	Validator(validatorAddress string) (*Validator, error)
	Delegation(delegator string, validator string) (*DelegationResponse, error)
	AnnualProvisions() (coin.DecCoin, error)
	TotalBondedBalance() (coin.Coin, error)
}

var ErrAccountNotFound = errors.New("account not found")
var ErrAccountNoDelegation = errors.New("account has no delegation")
