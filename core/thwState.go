package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/"
)

type THWState struct {
	//members
	candidateList []common.Address
	candidateCount int64
	//rand seed from block
	latestRand int64
	//parameters


}


func (thwState *THWState) IsCommittee(addr common.Address) (bool, error){

}


func (thwState *THWState) IsNextCommittee(addr common.Address) (bool, error){

}

func (thwState *THWState) AddCandidate(addr common.Address) error{
	return nil
}




