package trustedHW

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/consensus"
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
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

	return thw.verifyHeader(chain, header, nil);
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
//XS: its an async function.
func (thw *TrustedHW) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		for i, header := range headers {
			err := thw.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}
// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (thw *TrustedHW) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	//step 1: Sanity check.

	if header.Number == nil {
		return errUnknownBlock
	}
	number := header.Number.Uint64()
	//already in the local chain
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	//same as ethash, check ancestor first
	var parent *types.Header

	if parents == nil || len(parents) == 0 {
		parent = chain.GetHeader(header.ParentHash, number-1)
		if parent == nil {
			return consensus.ErrUnknownAncestor
		}
	}else{
		parent = parents[0]
	}
	//Step 2: check author is in the committee list

	//Step 3: check the verifier's signature

	//Step 4: check the author's signature.

	return nil
}


