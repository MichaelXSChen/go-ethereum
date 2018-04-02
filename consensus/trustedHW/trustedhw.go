package trustedHW


import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"
	"github.com/ethereum/go-ethereum/rpc"
	//"github.com/naoina/toml/ast"
	//"time"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/log"
)

var (
	// errUnknownBlock is returned when the list of signers is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	errNoCommittee = errors.New("not a committee member")
	errNoLeader = errors.New("not leader, cannot generate a block")
)





type TrustedHW struct{
	config *params.THWConfig


}

func New (config *params.THWConfig) *TrustedHW{
	//set missing configs
	conf := *config

	return &TrustedHW{
		config: &conf,
	}
}






//put the author (the leader of the committee)
func (thw *TrustedHW) Author(header *types.Header) (common.Address, error){
	return header.Coinbase, nil
}


//used to verify header downloaded from other peers.
func (thw *TrustedHW) VerifyHeader (chain consensus.ChainReader, header *types.Header, seal bool) error {

	return thw.verifyHeader(chain, header, nil);
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
//
// XS: its an async function.
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


	//Step 2: check author is in the committee list

	//Step 3: check the verifier's signature

	//Step 4: check the author's signature.



	return nil
}



func (thw *TrustedHW) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	//does not support uncles
	if len(block.Uncles()) > 0 {
		return errors.New("uncles not allowed")
	}
	return nil
}

//double check the seal of an outgoing message.
func (thw *TrustedHW) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	//It is currently a double check.
	return nil
}

//Read through the Chain and Determine whether addr is in the committee.
func (thw *TrustedHW) isCommittee (chain consensus.ChainReader, addr common.Address, number uint64, fake bool) (bool, error) {
	state := chain.GetThwState()
	return state.IsCommittee(addr, number, fake)
}


func (thw *TrustedHW) Prepare(chain consensus.ChainReader, header *types.Header) error {
	//A header is prepared only when a consensus has been made.
	number := header.Number.Uint64()
	log.THW("Preparing block", "number", number )

	ret, err := thw.isCommittee(chain, header.Coinbase, number, thw.config.FakeConsensus)
	if err != nil {
		return err
	}
	if !ret{
		return errNoCommittee
	}
	header.Nonce = types.BlockNonce{} //empty

	header.Difficulty = big.NewInt(1)
	header.MixDigest = common.Hash{} //empty

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	return nil
}


//ensuring no uncles are set. No
func (thw *TrustedHW) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error){


	header.Root = state.IntermediateRoot(true)
	header.UncleHash = types.CalcUncleHash(nil)
	//TODO: Whether the rewards should come from here.


	return types.NewBlock(header, txs, nil, receipts), nil

}

func (thw *TrustedHW) Seal (chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error){
	//attempt to achieve consensus.
	state := chain.GetThwState()

	ret, err := state.FakeConsensus(block.Coinbase(), block.NumberU64())
	if err != nil{
		return nil, err
	}
	if ret == false{
		return nil, errNoLeader
	}
	//it is elected as the leader and can append block.

	//Step 1. DO traditional Paxos consensus
	//elected as the leader.
	//TODO: step 2, use verifier to avoid network partition. Next round.
	//Step 2. Ask for verfication from the verifier groups.
	return block, nil

}
//Main function to achieve consensus.
func (thw *TrustedHW) invokeConsensus (chain consensus.ChainReader, number *big.Int) (elected bool, seed uint64){
	return true, 0
}


func (thw *TrustedHW) CalcDifficulty(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	//Can use this function to change the protocol parameters.

	return big.NewInt(1)
}


func (thw *TrustedHW) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{{
		Namespace: "thw",
		Version:   "1.0",
		Service:   &API{chain: chain, thw:thw},
		Public:    false,
	}}
}




func invokeConsensus() uint64{
	return uint64(0)
}

//as
func consensus_thread (asCommittee <-chan uint64, abort <-chan bool, results chan<- uint64) error {
	for{//forever
		//get a term of length asCommittee
		termLen := <- asCommittee
		for i:= uint64(0); i<termLen; i++ {
			rand := invokeConsensus()
			select {
			case _ = <- abort:
				break
			case results <- rand:
			}

		}

	}
}
