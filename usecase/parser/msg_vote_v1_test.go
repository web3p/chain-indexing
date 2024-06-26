package parser_test

import (
	"github.com/crypto-com/chain-indexing/infrastructure/tendermint"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chain-indexing/entity/command"
	command_usecase "github.com/crypto-com/chain-indexing/usecase/command"
	"github.com/crypto-com/chain-indexing/usecase/event"
	"github.com/crypto-com/chain-indexing/usecase/model"
	model_gov_v1 "github.com/crypto-com/chain-indexing/usecase/model/gov/v1"
	"github.com/crypto-com/chain-indexing/usecase/parser"
	usecase_parser_test "github.com/crypto-com/chain-indexing/usecase/parser/test"
)

var _ = Describe("ParseMsgCommands", func() {
	Describe("v1.MsgVote", func() {
		It("should parse gov.v1.MsgVote command in the transaction", func() {
			block, _ := mustParseBlockResp(usecase_parser_test.TX_MSG_VOTE_V1_BLOCK_RESP)
			blockResults := mustParseBlockResultsResp(
				usecase_parser_test.TX_MSG_VOTE_V1_BLOCK_RESULTS_RESP,
				&tendermint.Base64BlockResultEventAttributeDecoder{},
			)

			tx := MustParseTxsResp(usecase_parser_test.TX_MSG_VOTE_V1_TXS_RESP)
			txs := []model.CosmosTxWithHash{*tx}

			accountAddressPrefix := "crc"
			bondingDenom := "basecro"

			pm := usecase_parser_test.InitParserManager()

			cmds, possibleSignerAddresses, err := parser.ParseBlockTxsMsgToCommands(
				pm,
				block.Height,
				blockResults,
				txs,
				accountAddressPrefix,
				bondingDenom,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))

			Expect(cmds).To(Equal([]command.Command{
				command_usecase.NewCreateMsgVoteV1(
					event.MsgCommonParams{
						BlockHeight: int64(5183),
						TxHash:      "D2711F0542407D7D7F4A2E34184D122D68E8E7E207E329E4354F96171793B16F",
						TxSuccess:   true,
						MsgIndex:    0,
					},
					model_gov_v1.MsgVoteParams{
						ProposalId: "3",
						Voter:      "crc18z6q38mhvtsvyr5mak8fj8s8g4gw7kjjtsgrn7",
						Option:     "VOTE_OPTION_NO",
					},
				),
			}))
			Expect(possibleSignerAddresses).To(Equal([]string{"crc18z6q38mhvtsvyr5mak8fj8s8g4gw7kjjtsgrn7"}))
		})
	})
})
