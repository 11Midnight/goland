package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)
func convertYaml(i interface{}) (interface{}, error){
	var err error
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)],err = convertYaml(v)
			for key,value:= range m2 {
				switch v := value.(type){
				case string:
					switch key{
					case "endpoint":{}
					case "domain":
						m2["endpoint"] = fmt.Sprintf("https://%s", v)
					case "ip":
						m2["endpoint"] = fmt.Sprintf("http://%s", v)
					default:
						m2[key] = strings.ReplaceAll(v, "http:", "https:")
					}
				}

			}
		}
		return m2,err
	case []interface{}:
		for i, v := range x {
			x[i],err = convertYaml(v)
		}
	case interface{}:
		return i,err
	default:
		return nil,errors.New("Invalid type")
	}
	return i,err
}
func main() {
	writer := bufio.NewWriter(os.Stdout)
	scanner:=bufio.NewScanner(os.Stdin)
	scanner.Scan()
	fileNameYaml:=scanner.Text()
	yamlD, err := ioutil.ReadFile(fileNameYaml)
	if err != nil {
		fmt.Println(fmt.Errorf("Error reading YAML file: %s\n", err))
		return
	}
	var body interface{}
	err = yaml.Unmarshal([]byte(yamlD), &body)
	if  err != nil {
		fmt.Println(fmt.Errorf("Error unmarshal YAML file: %s\n", err))
		return
	}
	body,err = convertYaml(body)
	if  err != nil {
		fmt.Println(fmt.Errorf("Error convert body: %s\n", err))
		return
	}
	b, err := json.Marshal(body)
	if  err != nil {
		fmt.Println(fmt.Errorf("Error marshal to JSON file: %s\n", err))
		return
	}
	writer.WriteString(fmt.Sprintf("Output: %s\n", b))
	fileNameJson:=strings.Replace(fileNameYaml,".yaml",".json", -1)
	jsonFile, err := os.Create("./"+fileNameJson)
	defer jsonFile.Close()
	if err != nil {
		fmt.Println(fmt.Errorf("Error create JSON  file: %s\n", err))
		return
	}
	jsonFile.Write(b)
	writer.WriteString("JSON data written to "+ jsonFile.Name())
	writer.Flush()
}