package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/gustavolimam/control-access/src/components/defaults"
)

// Config representa a configuração do serviço
var (
	Config SysConfig
)

// SysConfig define a estrutura de configuração do serviço
type SysConfig struct {
	Camera  CamCfg
	Jidosha CfgJidosha
	Path    PathConfig
}

// PathConfig define a estrutura de configuração dos diretórios
type PathConfig struct {
	FinalPackage string // Caminho para armazenar arquivos .xml e as imagens zoom e pan
	LogPath      string // Caminho para armazenar .txt de logs
}

// CfgJidosha define a estrutura de configuração do jidosha
type CfgJidosha struct {
	Timeout    int
	NumThreads int
}

// CamCfg define a estrutura de configuração de uma camera
type CamCfg struct {
	Address    string // Ip para conexão
	FrameRate  int
	ImgQuality int
}

// SetupConfig salva em memória as configurações do serviço
func SetupConfig() error {

	pathCfg := path.Join(defaults.GetPath(), "util", "config.json")

	// Ler arquivo de cfg salvo em disco
	file, err := ioutil.ReadFile(pathCfg)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, &Config); err != nil {
		return err
	}

	return setupPaths()
}

// setupPaths verifica se todos os diretórios existem, senão cria os mesmos
func setupPaths() error {

	if err := verifyPath(Config.Path.FinalPackage); err != nil {
		return err
	}

	if err := verifyPath(Config.Path.LogPath); err != nil {
		return err
	}

	return verifyPath(Config.Path.LogPath)
}

// verifyPath função que verifica a existência e executa a criação de diretórios
func verifyPath(path string) error {
	if ok, err := exists(path); err != nil {
		return err
	} else if !ok {
		err = os.MkdirAll(path, os.ModePerm)
		return err
	}

	return nil
}

// exists verifica a existência de um caminho
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
