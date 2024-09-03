package account

type Account struct {
	ID      int64   `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Pass    string  `json:"password"`
	NoteIDS []int64 `json:"noteIds"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
