package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)
func convertYaml(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertYaml(v)
			for key,value:= range m2{
				if strings.Contains(fmt.Sprintf("%s",value),"http:"){
					result := strings.Replace(fmt.Sprintf("%s",value), "http:", "https:", -1)
					m2[key]=fmt.Sprintf("%s",result)
				}
				if key =="domain"{
					m2["endpoint"]="https://"+fmt.Sprintf("%s",value)
				}
				if key =="ip"{
					m2["endpoint"]="http://"+fmt.Sprintf("%s",value)
				}
			}
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertYaml(v)
		}
	}
	return i
}

func main() {
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
		fmt.Println(fmt.Errorf("Error reading YAML file: %s\n", err))
		return
	}
	body = convertYaml(body)
	b, err := json.Marshal(body)
	if  err != nil {
		fmt.Println(fmt.Errorf("Error reading YAML file: %s\n", err))
		return
	}
	fmt.Printf("Output: %s\n", b)
	fileNameJson:=strings.Replace(fileNameYaml,".yaml",".json", -1)
	jsonFile, err := os.Create("./"+fileNameJson)
	if err != nil {
		fmt.Println(fmt.Errorf("Error create JSON  file: %s\n", err))
		return
	}
	jsonFile.Write(b)
	defer jsonFile.Close()
	fmt.Println("JSON data written to ", jsonFile.Name())

}
