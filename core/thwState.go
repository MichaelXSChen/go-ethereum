package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/emirpasic/gods/lists/arraylist"
	"sync"
	"fmt"
	"errors"
	"encoding/binary"
	"github.com/ethereum/go-ethereum/core/thwCore"
	"github.com/ethereum/go-ethereum/log"
)

var(
	ErrNoCandidate = errors.New("Not a Candidate")
	ErrInitFailed  = errors.New("THW State Init Failed")
)


type THWState struct {
	hc *HeaderChain
	mu sync.Mutex
	//members
	candidateList *hashmap.Map    //key: addr, value: *Candidate
	candidateCount int64
	//rand seed from block
	//latestRand int64
	//latestBlock int64
	//Committee term
	CommitteeTerms *arraylist.List  //the start block number of each committee.

	//parameters
	committeeRatio uint64 //On average, there is 1 committee every ``x'' candidates
	committeeMaxTerm uint64 //One committee can serve ``x'' terms

}

func (thws *THWState) Init(headerchain interface{}) error { //TODO: set parameters
	thws.mu.Lock()
	thws.candidateList = hashmap.New()
	thws.CommitteeTerms = arraylist.New()
	hc, ok := headerchain.(*HeaderChain)
	if !ok {
		return ErrInitFailed
	}
	thws.hc = hc
	thws.candidateCount = 0
	thws.mu.Unlock()
	return nil
}


func (thws *THWState) findTerm (num uint64) (*thwCore.Term, error){
	thws.mu.Lock()
	defer thws.mu.Unlock()

	it := thws.CommitteeTerms.Iterator()

	for it.End(); it.Prev(); {
		t, _ := it.Value().(*thwCore.Term)
		if num > t.Start {
			if num > t.Start + t.Len{
				return nil, errors.New("[Find term], out of bound")
			}else{
				return t, nil
			}
		}
	}
	return nil, errors.New("[Find term], not found")
}


//a simple/fake checkCommittee function
func (thws *THWState) checkCommittee(addr common.Address, rand uint64) bool {
	x := addrToInt(addr)
	m := thws.committeeRatio
	if (x-rand)%m == 0 {
		return true
	} else {
		return false
	}
}


func (thws *THWState) IsCommittee(addr common.Address, num uint64) (bool, error){
	t, err := thws.findTerm(num)
	if err != nil {
		return false, err
	}
	seed := thws.hc.GetHeaderByNumber(t.Start).TrustRand

	return thws.checkCommittee(addr, seed), nil
}


func (thws *THWState) IsNextCommittee(addr common.Address, num uint64) (bool, error){
	seed := thws.hc.GetHeaderByNumber(num).TrustRand
	return thws.checkCommittee(addr, seed), nil
}



func (thws *THWState) AddCandidate(candidate *thwCore.Candidate) error{
	thws.mu.Lock()
	defer thws.mu.Unlock()

	ret, found := thws.candidateList.Get(candidate.Addr)
	if found == true{
		//found an candidate
		c, ok := ret.(*thwCore.Candidate)
		if !ok {
			fmt.Println("Wrong type in the hash map")
			panic("Wrong type in the hash map")
		}
		//renew the cointract
		c.Term = c.Term + candidate.Term
	}else{
		//not found in the list
		thws.candidateList.Put(candidate.Addr, candidate)
		thws.candidateCount++
	}
	fmt.Println("added candidate with addr", candidate.Addr)

	return nil
}


func addrToInt (address common.Address) uint64{
	return binary.BigEndian.Uint64(address[0:8]) +
		binary.BigEndian.Uint64(address[8:16]) + uint64(binary.BigEndian.Uint32(address[16:20]))
}

func (thws *THWState) FakeConsensus(addr common.Address, number uint64) (bool, error) {
	if _, ret := thws.candidateList.Get(addr); !ret {
		return false, ErrNoCandidate
	}

	me := addrToInt(addr)

	candidates := thws.candidateList.Keys()
	for _, c := range candidates{
		x, ok := c.(common.Address)
		if !(ok){
			log.Error("Wrong type from candidate list")
		}else{
			if me < addrToInt(x){ //not the biggest
				return false, nil
			}
		}
	}
	return true, nil
}