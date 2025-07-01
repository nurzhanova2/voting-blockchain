package dto

// AddBlockRequest — тело запроса при добавлении блока
type AddBlockRequest struct {
	VoteHash string `json:"vote_hash"`
}
