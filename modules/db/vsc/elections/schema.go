package elections

type ElectionResult struct {
	Epoch       uint64           `json:"epoch" bson:"epoch"`
	NetId       string           `json:"net_id" bson:"net_id"`
	BlockHeight uint64           `json:"block_height" bson:"block_height"`
	Data        string           `json:"data" bson:"data"`
	Members     []ElectionMember `json:"members" bson:"members"`
	Proposer    string           `json:"proposer" bson:"proposer"`
	Weights     []uint64         `json:"weights" bson:"weights"`
	WeightTotal uint64           `json:"weight_total" bson:"weight_total"`
}

type ElectionMember struct {
	Key     string `json:"key"`
	Account string `json:"account"`
}
