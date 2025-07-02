package models

type Choice struct {
    ID         int    `db:"id"`
    ElectionID int    `db:"election_id"`
    Text       string `db:"text"`
    Count      int    `db:"count,omitempty"`
}
