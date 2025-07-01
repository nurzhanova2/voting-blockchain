
package models

import "time"

// Block — элемент цепочки блоков голосования
type Block struct {
	Index       int       `db:"id"`            // PRIMARY KEY
	Timestamp   time.Time `db:"created_at"`    // автоматически now()
	VoteHash    string    `db:"vote_hash"`     // хеш голосования
	PrevHash    string    `db:"previous_hash"` // хеш предыдущего блока
	Hash        string    `db:"current_hash"`  // хеш текущего блока
	ElectionID  int       `db:"election_id"`   // к какому голосованию
}
