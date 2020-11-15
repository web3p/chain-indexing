package tmcosmosutils_test

import (
	"encoding/base64"

	"github.com/crypto-com/chainindex/internal/tmcosmosutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const tendermintPubKey = "tWY6qzpOg/6HFj2X3a8+tzIAehW7k2MWOgrjotcWCuI="
const tendermintAddress = "18253C74D541A2B5BF492E8B5910A268AFED4D48"
const consensusAddress = "tcrocnclcons1rqjncax4gx3tt06f9694jy9zdzh76n2gg5yqfd"
const consensusPubKey = "tcrocnclconspub1zcjduepqk4nr42e6f6plapck8ktamte7kueqq7s4hwfkx936pt3694ckpt3q50kp4w"

var _ = Describe("tmcosmosutils", func() {
	Describe("TmAddressFromTmPubKey", func() {
		It("should work", func() {
			pubKey, _ := base64.StdEncoding.DecodeString(tendermintPubKey)
			Expect(tmcosmosutils.TmAddressFromTmPubKey(
				pubKey,
			)).To(Equal(tendermintAddress))
		})
	})

	Describe("ConsensusAddressFromTmPubKey", func() {
		It("should work", func() {
			pubKey, _ := base64.StdEncoding.DecodeString(tendermintPubKey)
			Expect(tmcosmosutils.ConsensusAddressFromTmPubKey(
				"tcrocnclcons", pubKey,
			)).To(Equal(consensusAddress))
		})
	})

	Describe("ConsensusPubKeyFromTmPubKey", func() {
		It("should work", func() {
			pubKey, _ := base64.StdEncoding.DecodeString(tendermintPubKey)
			Expect(tmcosmosutils.ConsensusPubKeyFromTmPubKey(
				"tcrocnclconspub", pubKey,
			)).To(Equal(consensusPubKey))
		})
	})
})
