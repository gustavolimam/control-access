package report

var errorReportCh = make(ErrorReport)

// ErrorReport representa um canal para enviar erros para main
type ErrorReport chan []interface{}

func GetErrorReportCh() ErrorReport {
	return errorReportCh
}

func ReportError(msg ...interface{}) {
	errorReportCh <- msg
}
