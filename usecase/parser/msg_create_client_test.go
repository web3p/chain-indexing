package parser_test

import (
	"regexp"
	"strings"

	"github.com/crypto-com/chain-indexing/internal/json"

	"github.com/crypto-com/chain-indexing/usecase/parser/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/crypto-com/chain-indexing/infrastructure/tendermint"
	"github.com/crypto-com/chain-indexing/usecase/event"
	"github.com/crypto-com/chain-indexing/usecase/parser"
	usecase_parser_test "github.com/crypto-com/chain-indexing/usecase/parser/test"
)

var _ = Describe("ParseMsgCommands", func() {
	Describe("MsgIBCCreateClient", func() {
		It("should parse Msg commands when there is MsgIBCCreateClient in the transaction", func() {
			expected := `{
  "name": "MsgCreateClientCreated",
  "version": 1,
  "height": 5,
  "uuid": "{UUID}",
  "msgName": "MsgIBCCreateClient",
  "txHash": "7E34A75D8063BADF7B93538C23C88DEEF1FF14E7BE7F13AD6AD34E228C64538D",
  "msgIndex": 0,
  "maybeTendermintLightClient": {
    "clientState": {
      "@type": "/ibc.lightclients.tendermint.v1.ClientState",
      "chainId": "devnet-2",
      "trustLevel": { "numerator": "1", "denominator": "3" },
      "trustingPeriod": 1209600000000000,
      "unbondingPeriod": 1814400000000000,
      "maxClockDrift": 5000000000,
      "frozenHeight": { "revisionNumber": "0", "revisionHeight": "0" },
      "latestHeight": { "revisionNumber": "2", "revisionHeight": "2" },
      "proofSpecs": [
        {
          "leafSpec": {
            "hash": "SHA256",
            "prehashKey": "NO_HASH",
            "prehashValue": "SHA256",
            "length": "VAR_PROTO",
            "prefix": "AA=="
          },
          "innerSpec": {
            "childOrder": [0, 1],
            "childSize": 33,
            "minPrefixLength": 4,
            "maxPrefixLength": 12,
            "emptyChild": null,
            "hash": "SHA256"
          },
          "maxDepth": 0,
          "minDepth": 0
        },
        {
          "leafSpec": {
            "hash": "SHA256",
            "prehashKey": "NO_HASH",
            "prehashValue": "SHA256",
            "length": "VAR_PROTO",
            "prefix": "AA=="
          },
          "innerSpec": {
            "childOrder": [0, 1],
            "childSize": 32,
            "minPrefixLength": 1,
            "maxPrefixLength": 1,
            "emptyChild": null,
            "hash": "SHA256"
          },
          "maxDepth": 0,
          "minDepth": 0
        }
      ],
      "upgradePath": ["upgrade", "upgradedIBCState"],
      "allowUpdateAfterExpiry": false,
      "allowUpdateAfterMisbehaviour": false
    },
    "consensusState": {
      "@type": "/ibc.lightclients.tendermint.v1.ConsensusState",
      "timestamp": "2021-05-04T18:02:36.089446Z",
      "root": { "hash": "bVjiQ29+V522NVFdx1BiVJnIBJV8Y1pYe9psvxZFAWg=" },
      "nextValidatorsHash": "E3DE0D2B3237A02E9C20C34F9EE04F69F5861FBC2E2722A011CA9037FC67A7EC"
    }
  },
  "signer": "cro1gdswrmwtzgv3kvf28lvtt7qv7q7myzmn466r3f",
  "clientId": "07-tendermint-0",
  "clientType": "07-tendermint"
}`

			txDecoder := utils.NewTxDecoder()
			block, _, _ := tendermint.ParseBlockResp(strings.NewReader(
				usecase_parser_test.TX_MSG_CREATE_CLIENT_BLOCK_RESP,
			))
			blockResults, _ := tendermint.ParseBlockResultsResp(strings.NewReader(
				usecase_parser_test.TX_MSG_CREATE_CLIENT_BLOCK_RESULTS_RESP,
			))

			accountAddressPrefix := "cro"
			stakingDenom := "basecro"
			cmds, err := parser.ParseBlockResultsTxsMsgToCommands(
				txDecoder,
				block,
				blockResults,
				accountAddressPrefix,
				stakingDenom,
			)
			Expect(err).To(BeNil())
			Expect(cmds).To(HaveLen(1))
			cmd := cmds[0]
			Expect(cmd.Name()).To(Equal("CreateMsgCreateClient"))

			untypedEvent, _ := cmd.Exec()
			createMsgCreateClientEvent := untypedEvent.(*event.MsgIBCCreateClient)

			regex, _ := regexp.Compile("\n?\r?\\s?")
			Expect(json.MustMarshalToString(createMsgCreateClientEvent)).To(Equal(
				strings.Replace(
					regex.ReplaceAllString(expected, ""),
					"{UUID}",
					createMsgCreateClientEvent.UUID(),
					-1,
				),
			))
		})
	})
})
