package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads configs from a file to an interface
func LoadConfig(fileName string, iFace interface{}) (err error) {
	if err = createIfDoesntExist(fileName); err != nil {
		return err
	}

	if strings.HasSuffix(fileName, ".json") { // if json file
		if err = readJSONFromFile(fileName, iFace); err != nil {
			return
		}
	} else if strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml") { // if yaml file
		if err = readYamlFromFile(fileName, iFace); err != nil {
			return
		}
		// log.Printf("interface %+v", iFace)
	} else {
		return errors.New("no supported file type located")
	}

	return nil
}

// SaveConfig saves interfaces to a file
func SaveConfig(file string, iFace interface{}) (err error) {
	if strings.HasSuffix(file, ".json") { // if json file
		if err := writeJSONToFile(file, iFace); err != nil {
			return err
		}
	} else if strings.HasSuffix(file, ".yml") || strings.HasSuffix(file, ".yaml") { // if yaml file
		if err = writeYamlToFile(file, iFace); err != nil {
			return
		}
		// log.Printf("interface %+v", iFace)
	} else {
		return errors.New("no supported file type located")
	}

	return nil
}

// File management
func writeJSONToFile(file string, iFace interface{}) (err error) {
	jData, err := json.MarshalIndent(iFace, "", "  ")
	if err != nil {
		return
	}

	// create a file with a supplied name
	if jsonFile, err := os.Create(file); err != nil {
		return err
	} else if _, err = jsonFile.Write(jData); err != nil {
		return err
	}

	return
}

func readJSONFromFile(file string, iFace interface{}) (err error) {
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		if err := jsonFile.Close(); err != nil {
			return
		}
	}()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err = json.Unmarshal(byteValue, iFace); err != nil {
		return err
	}

	// return the json byte value.
	return nil
}

func writeYamlToFile(file string, iFace interface{}) (err error) {
	ydata, err := yaml.Marshal(iFace)
	if err != nil {
		return
	}

	// create a file with a supplied name
	yamlFile, err := os.Create(file)
	if err != nil {
		return
	}

	if _, err = yamlFile.Write(ydata); err != nil {
		return
	}

	return
}

func readYamlFromFile(file string, iFace interface{}) (err error) {
	yamlFile, err := os.Open(file)
	if err != nil {
		return
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer func() {
		if err := yamlFile.Close(); err != nil {
			return
		}
	}()

	byteValue, _ := ioutil.ReadAll(yamlFile)
	if err = yaml.Unmarshal(byteValue, iFace); err != nil {
		return
	}

	return
}

// Exists reports whether the named file or directory exists.
func createIfDoesntExist(name string) (err error) {
	p, file := path.Split(name)

	// if confDir exists carry on
	if _, err := os.Stat(name); err != nil {
		// if file doesn't exist
		if os.IsNotExist(err) {
			// stat
			if _, err = os.Stat(name); err != nil {
				if file == "" {
					if err = os.Mkdir(p, 0755); err != nil {
						return err
					}
				} else {
					if fileCheck, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644); err != nil {
						return err
					} else {
						if err := fileCheck.Close(); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return
}
