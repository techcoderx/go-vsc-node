package stateEngine

type TxPacket struct {
	TxId string
	Ops  []VSCTransaction
}

type TxOutput struct {
	Ok        bool
	Logs      []string
	RcUsed    int64
	LedgerIds []string
}

type TxResult struct {
	Success bool
	Ret     string
	RcUsed  int64
}

type TxResultWithId struct {
	Success bool
	Ret     string
	TxId    string
	RcUsed  int64
}

// More information about the TX
type TxSelf struct {
	TxId                 string
	BlockId              string
	BlockHeight          uint64
	Index                int
	OpIndex              int
	Timestamp            string
	RequiredAuths        []string
	RequiredPostingAuths []string
}

type OplogOutputEntry struct {
	Id        string `json:"id" bson:"id"`
	Ok        bool   `json:"ok" bson:"ok"`
	LedgerIdx []int  `json:"lidx" bson:"lidx"`
}
