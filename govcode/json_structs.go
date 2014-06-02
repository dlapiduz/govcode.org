package main

type Account struct {
	Verified     bool
	Account      string
	Organization string
}

type AccountData struct {
	PageCount  int       `json:"page_count"`
	TotalItems int       `json:"total_items"`
	PageNumber int       `json:"page_number"`
	Accounts   []Account `json:"accounts"`
}
