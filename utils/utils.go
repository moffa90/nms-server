package utils

import (
	"github.com/moffa90/triadNMS/assets"
	"github.com/moffa90/triadNMS/constants"
	"github.com/moffa90/triadNMS/utils/func-maps"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
)

func RenderPage(w io.Writer, path string, i interface{}){
	tplBytes, _ := assets.Asset(path)
	header, _ := assets.Asset(constants.TEMPLATE_PAGE_HEADER_PATH)
	footer, _ := assets.Asset(constants.TEMPLATE_PAGE_FOOTER_PATH)
	menu, _ := assets.Asset(constants.TEMPLATE_PAGE_MENU_PATH)
	cellInfo, _ := assets.Asset(constants.TEMPLATE_PARTIAL_CELL_INFO_INFO_PATH)
	remoteCellTemplate, _ := assets.Asset(constants.TEMPLATE_PARTIAL_CELL_REMOTE_INFO_PATH)

	tpl := template.New("templates")
	tpl.New("header").Funcs(func_maps.Generic).Parse(string(header))
	tpl.New("footer").Parse(string(footer))
	tpl.New("menu").Parse(string(menu))
	tpl.New("cellInfo").Parse(string(cellInfo))
	tpl.New("remoteCellInfo").Parse(string(remoteCellTemplate))
	tpl.New("render").Funcs(func_maps.Generic).Parse(string(tplBytes))
	if ErrTpl :=  tpl.ExecuteTemplate(w, "render", i); ErrTpl != nil{
		log.Println(ErrTpl)
	}
}

func ExecScript(path string) ([]byte, error){
	cmd := exec.Command(path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed running %s with %s\n", path, err)
		return nil, err
	}
	log.Println(string(output))
	return output, nil
}

func ExistsFile(path string) bool{
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}


func CreateFile(path string) (*os.File, error){
	f, err := os.Create(path)

	if err != nil {
		log.Println(err.Error())
	}

	return f, err
}


func BackupFile(path string){
	if ExistsFile(path){
		from, err := os.Open(path)
		if err != nil {
			log.Println(err)
		}
		defer from.Close()

		to, err := os.OpenFile(path + ".bak", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Println(err)
		}
	}
}

func SetEnvVar(key string, value string){
	if err := os.Setenv(key, value); err != nil {
		log.Println(err)
	}
}

func GetEnvVar(key string) string{
	return os.Getenv(key)
}
