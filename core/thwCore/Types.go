package thwCore

import "github.com/ethereum/go-ethereum/common"

type Candidate struct{
	Referee common.Address
	Addr common.Address
	JoinRound uint64    //at which round the candidate joined as candidate
	Term int64         //the ``total'' term for the candidate
}

type Term struct{
	Start  uint64
	Len    uint64
	Seed   int64
}
