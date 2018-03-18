package tests

type MockServer struct {
}

func (MockServer) TransactionNum() int {
	return 0
}
