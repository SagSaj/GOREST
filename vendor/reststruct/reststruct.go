package reststruct

import "time"

type StructForREST struct {
	ID            int       `json:"id"`
	ParryType     string    `json:"type"`
	FromAccountID int       `json:"fromAccountID"`
	ToAccountID   int       `json:"toAccountID"`
	BetValue      float32   `json:"betValue"`
	CreatedAt     time.Time `json:"createdAt"`
	Status        string    `json:"status"`
	ArenaID       int       `json:"arenaID"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     time.Time `json:"deletedAt"`
}
