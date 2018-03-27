package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/lists/arraylist"
	"sync"
	"fmt"
	"errors"
	"encoding/binary"
)

var(
	ErrNoCandidate = errors.New("Not a Candidate")
)


type Candidate struct{
	addr common.Address
	joinRound int64    //at which round the candidate joined as candidate
	term int64         //the ``total'' term for the candidate
}

type Term struct{
	start  uint64
	len    uint64
	seed   int64
}



type THWState struct {
	bc *BlockChain
	mu sync.Mutex
	//members
	candidateList *hashmap.Map
	candidateCount int64
	//rand seed from block
	//latestRand int64
	//latestBlock int64
	//Committee term
	CommitteeTerms *arraylist.List  //the start block number of each committee.

	//parameters
	committeeRatio int64 //On average, there is 1 committee every ``x'' candidates
	committeeMaxTerm int64 //One committee can serve ``x'' terms

}

func (thws *THWState) Init(bc *BlockChain,){ //TODO: set parameters
	thws.mu.Lock()
	thws.candidateList = hashmap.New()
	thws.CommitteeTerms = arraylist.New()
	thws.bc = bc
	thws.candidateCount = 0
	thws.mu.Unlock()
}


func (thws *THWState) findTerm (num uint64) (*Term, error){
	thws.mu.Lock()
	defer thws.mu.Unlock()

	it := thws.CommitteeTerms.Iterator()

	for it.End(); it.Prev(); {
		t, _ := it.Value().(*Term)
		if num > t.start {
			if num > t.start + t.len{
				return nil, errors.New("[Find term], out of bound")
			}else{
				return t, nil
			}
		}
	}
	return nil, errors.New("[Find term], not found")
}


func (thws *THWState) IsCommittee(addr common.Address, num uint64) (bool, error){
	t, err := thws.findTerm(num)
	if err != nil {
		return false, err
	}
	block := thws.bc.GetBlockByNumber(t.start)
	seed := block.Header().TrustRand

	return thws.checkCommittee(addr, seed), nil
}


func (thws *THWState) IsNextCommittee(addr common.Address, num uint64) (bool, error){
	seed := thws.bc.GetBlockByNumber(num).Header().TrustRand
	return thws.checkCommittee(addr, seed), nil
}




func (thws *THWState) AddCandidate(candidate *Candidate) error{
	thws.mu.Lock()
	defer thws.mu.Unlock()

	ret, found := thws.candidateList.Get(candidate.addr)
	if found == true{
		//found an candidate
		c, ok := ret.(*Candidate)
		if !ok {
			fmt.Println("Wrong type in the hash map")
			panic("Wrong type in the hash map")
		}
		//renew the cointract
		c.term = c.term + candidate.term
	}else{
		//not found in the list
		thws.candidateList.Put(candidate.addr, candidate)
		thws.candidateCount++
	}
	return nil
}


//a simple/fake checkCommittee function
func (thws *THWState) checkCommittee(addr common.Address, rand int64) bool{
	x := int64(addrToInt(addr))
	m := thws.committeeRatio
	if (x - rand) % m == 0 {
		return true
	}else{
		return false
	}

}

func addrToInt (address common.Address) uint64{
	return binary.BigEndian.Uint64(address[0:8]) +
		binary.BigEndian.Uint64(address[8:16]) + uint64(binary.BigEndian.Uint32(address[16:20]))
}

