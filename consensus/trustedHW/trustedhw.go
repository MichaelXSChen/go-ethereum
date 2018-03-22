package trustedHW


import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/naoina/toml/ast"
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	errNoCommittee = errors.New("not a committee member")
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
	//What to do about the parent.


	return thw.verifySeal(chain, header)
}

func (thw *TrustedHW) verifySeal(chain consensus.ChainReader, header *types.Header) error {
	//Step 2: check author is in the committee list

	//Step 3: check the verifier's signature

	//Step 4: check the author's signature.

	return  nil
}

func (thw *TrustedHW) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	//does not support uncles
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}


func (thw *TrustedHW) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return thw.verifySeal(chain, header)
}

//Read through the Chain and Determine whether addr is in the committee.
func (thw *TrustedHW) isCommittee (chain consensus.ChainReader, addr common.Address, number uint64) bool{

	return true
}

func (thw *TrustedHW) Prepare(chain consensus.ChainReader, header *types.Header) error {
	//A header is prepared only when a consensus has been made.
	number := header.Number.Uint64()
	if ! thw.isCommittee(chain, number){
		return errNoCommittee
	}
	return nil
}


//ensuring no uncles are set. No
func (thw *TrustedHW) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error){


}

func (thw *TrustedHW) Seal (chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error){
	//attempt to achieve consensus.
}

func (thw *TrustedHW) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {

}


func (thw *TrustedHW) APIs(chain consensus.ChainReader) []rpc.API{

}

