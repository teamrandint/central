package tests

type MockLogger struct {
}

func (MockLogger) QuoteServer(server string, transNum int, reply *transactionserver.QuoteReply) {

}

func (MockLogger) AccountTransaction(server string, transNum int, action string, user interface{}, funds interface{}) {

}

func (MockLogger) SystemError(server string, transNum int, command string, user interface{}, stock interface{}, filename interface{},
	funds interface{}, errorMsg interface{}) {

}

func (MockLogger) SystemEvent(server string, transNum int, command string, username interface{}, stock interface{},
	filename interface{}, funds interface{}) {

}

func (MockLogger) DumpLog(filename string, username interface{}) {

}
