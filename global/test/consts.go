package test

// Note: at the moment I'm using a fixed test user but might change it to a random user in the future
const (
	UserEmail = "admin@example.com"
	UserName  = "admin"
)

// DoubleAbsErr is the absolute tolerable double error in tests
const DoubleAbsErr = float64(0.001)

const alphabet = "abcdefghijklmnopqrstuvwxyz"
