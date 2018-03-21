package trustedHW

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/consensus"
)

type Config struct{
	//parameters:

}


type TrustedHW struct{
	config Config

}
//put the author (the leader of the committee)
func (thw *TrustedHW) Author(header *types.Header) (common.Address, error){
	return header.Coinbase, nil
}

func (thw *TrustedHW) VerifyHeader (chain consensus.ChainReader, header *types.Header, seal bool) error {



	//step 1: Sanity check.
	number := header.Number.Uint64()
	//already in the local chain
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	//same as ethash, check ancestor first
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	//Step 2: check author is in the committee list

	//Step 3: check the verifier's signature

	//Step 4: check the author's signature.

	return nil
	
}
