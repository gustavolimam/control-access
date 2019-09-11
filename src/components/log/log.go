package log

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/gustavolimam/control-access/src/components/config"
)

// Service representa um serviço para identificar a origem do log
type Service string

const (
	logFormat  = "[02-01-2006 15:04:05.00000]"
	fileFormat = "02-01-2006-15-04-05"
	fatalFile  = "fatal.log"
)

var (
	logMutex    sync.Mutex
	changeMutex sync.Mutex
	fileName    string
	logFile     *os.File
	LogDay      int
)

// Log salva e printa informações sobre o sistema. Esse package
// funciona da seguinte forma:
// 1. Toda iniciação de sistema cria um arquivo txt de log com o timestamp
// 2. Concatena logs no arquivo txt;
func Log(logService Service, data ...interface{}) {
	logtime := time.Now()
	go func() {
		changeMutex.Lock()
		if logtime.Day() != LogDay {
			if err := closeLogPackage(); err != nil {
				Fatal("LOG", "Erro fechando arquivo de log do dia anterior. ", err,
					". Erro na inserção do log. Log service: ", logService, " Mensagem: ", data)
			}
			_, err := CreateLogFile()
			if err != nil {
				Fatal("LOG", "Erro criando novo arquivo de log. ", err,
					". Erro na inserção do log. Log service: ", logService, " Mensagem: ", data)
			}
		}
		changeMutex.Unlock()

		text := ""
		for _, value := range data {
			text += fmt.Sprint(value)
		}

		final := fmt.Sprintf("%s %s %s",
			logtime.Format(logFormat),
			getServiceText(logService),
			text)

		fmt.Println(final)
		saveLog(final)
	}()
}

// Fatal salva, printa informações sobre o sistema
// e encerra o processo
func Fatal(logService Service, data ...interface{}) {
	if logFile != nil {
		Log(logService, "FATAL ERROR ---", data)
	} else {
		createFatalLog("LOG", "Arquivo de log não existente. Erro na inserção do log. Log service: ",
			logService, " Mensagem: ", data)
	}

	changeMutex.Lock()
	// Caso não tenha arquivo de log, cria um para registrar erro fatal
	// do serviço
	createFatalLog(logService, data)

	if logFile != nil {
		if err := closeLogPackage(); err != nil {
			createFatalLog("LOG", "Falha na tentativa de fechar o arquivo de log: ", err)
		}
		if config.Config != (config.SysConfig{}) {
			fmt.Println("Não entendi")
		} else {
			createFatalLog("LOG", "Nao foi possivel enviar log via FTP (sem config.json)")
		}
	} else {
		createFatalLog("LOG", "Nao foi possivel enviar log via FTP (sem arquivo de log)")
	}

	os.Exit(1)
}

func getServiceText(s Service) string {
	return fmt.Sprintf("[%s]", s)
}

func saveLog(text string) {
	logMutex.Lock()
	defer logMutex.Unlock()
	if _, err := logFile.Write([]byte(text + "\r\n")); err != nil {
		createFatalLog("LOG", err)
	}
}

func getNewFileName() string {
	LogDay = time.Now().Day()
	return path.Join(config.Config.Path.LogPath, "ControleAcesso-LOG-"+
		time.Now().Format(fileFormat)+
		".log")
}

func CreateLogFile() (newLogFile string, err error) {
	logMutex.Lock()
	defer logMutex.Unlock()
	fileName = getNewFileName()
	newLogFile = fileName
	logFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		Log("LOG", "Arquivo de log criado com sucesso : ", fileName)
	}
	return newLogFile, err
}

func createFatalLog(logService Service, data ...interface{}) {
	timelog := time.Now()
	text := ""
	for _, value := range data {
		text += fmt.Sprint(value)
	}

	fileData := []byte(fmt.Sprintf("%s %s %s\r\n",
		timelog.Format(logFormat),
		getServiceText(logService),
		text))

	fileFatal, err := os.OpenFile(fatalFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer fileFatal.Close()

	if err != nil {
		return
	}

	fileFatal.Write(fileData)
}

func closeLogPackage() (err error) {
	logMutex.Lock()
	defer logMutex.Unlock()
	err = logFile.Close()
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	return
}
