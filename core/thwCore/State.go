package thwCore


import (
	"github.com/ethereum/go-ethereum/common"
)

//This package is created to solve the cyclic dependency



type State interface{
	Init(hc interface{}) error
	IsCommittee(addr common.Address, num uint64) (bool, error)
	IsNextCommittee(addr common.Address, num uint64) (bool, error)
	AddCandidate(candidate *Candidate) error
	FakeConsensus(addr common.Address, number uint64) (bool, error)
	NewTerm (term *Term) error
	CandidateCount() uint64
	IsValidator (addr *common.Address, num uint64 ) (bool, error)
	ValidatorCount () uint64

}