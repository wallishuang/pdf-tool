package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/spf13/viper"
)

type Config struct {
	Operation string
	Merge     MergeInfo
	Split     SplitInfo
}

type MergeInfo struct {
	InputFolder  string
	OutputFolder string
	OutputFile   string
}

type SplitInfo struct {
	InputFile    string
	Pages        []string
	OutputFolder string
}

func getConf() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("read config toml error: %v", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}

func main() {

}

func GoWork() {
	config := getConf()
	if config.Operation == "merge" {
		MergeCreateFile(config)
	} else if config.Operation == "split" {
		SplitFile(config)
	} else if config.Operation == "both" {
		MergeCreateFile(config)
		SplitFile(config)
	} else {
		fmt.Println("no such operation. ", config.Operation)
	}
	time.Sleep(3000)
}

func SplitFile(config *Config) {
	var conf *pdfcpu.Configuration
	pages := config.Split.Pages
	outFolder := config.Split.OutputFolder
	if _, err := os.Stat(outFolder); os.IsNotExist(err) {
		err := os.MkdirAll(outFolder, 755)
		if err != nil {
			fmt.Println("create new folder error. ", err.Error())
		}
	}

	err := api.ExtractPagesFile(config.Split.InputFile, outFolder, pages, conf)
	if err != nil {
		fmt.Println("split pdf error.", err.Error())
	}

	fmt.Println("Split file ok.")
}

func MergeCreateFile(config *Config) {
	var conf *pdfcpu.Configuration
	inFiles := []string{}
	files, _ := ioutil.ReadDir(config.Merge.InputFolder)

	for _, file := range files {
		inFiles = append(inFiles, config.Merge.InputFolder+"/"+file.Name())
	}

	outputFolder := config.Merge.OutputFolder
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.MkdirAll(outputFolder, 755)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	err := api.MergeCreateFile(inFiles, outputFolder+"/"+config.Merge.OutputFile, conf)
	if err != nil {
		fmt.Println("merge pdf error.", err.Error())
	}

	fmt.Println("Merge file ok.")
}
