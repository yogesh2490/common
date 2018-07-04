package common

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

var config_map_lock sync.RWMutex
var json_data = map[string]string{}
var config_loaded = false
var default_loaded = false

/*
	This is used to load configuration values.
*/

type IConfigGetter interface {
	MustGetConfigVar(variableName string) string
	SafeGetConfigVar(variableName string) string
}

//Implements the IConfigGetter interface
/*
	File path for .json file to load configurations.
	There should also be a default configuration file at <<config_file_path>>.dist
	The file should be in the following json format:
	{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3"
	}
*/
type configGetter struct {
	config_file_path string
}

/*
	This method to create a configGetter with injected information.
*/
func GetConfigGetter(config_file_path string) IConfigGetter {
	return &configGetter{
		config_file_path: config_file_path,
	}
}

func (this configGetter) MustGetConfigVar(variableName string) string {
	config := os.Getenv(variableName)
	if config == "" {
		this.loadDefaults()
		this.loadConfig()
		ok := false
		config_map_lock.Lock()
		config, ok = json_data[variableName]
		config_map_lock.Unlock()
		if !ok {
			panic("FATAL ERROR COULD NOT LOAD VAR : " + variableName)
		}
	}
	return config
}

/*

*/
func (this configGetter) SafeGetConfigVar(variableName string) string {
	config := os.Getenv(variableName)
	if config == "" {
		this.loadDefaults()
		this.loadConfig()
		ok := false
		config_map_lock.Lock()
		config, ok = json_data[variableName]
		config_map_lock.Unlock()
		if !ok {
			return ""
		}
	}
	return config
}

/*
	Loads the config ovverridden values.
*/
func (this configGetter) loadConfig() {
	if !config_loaded {
		loadConfigFile(this.config_file_path)
		config_loaded = true
	}
}

/*
	Loads the default config values
*/
func (this configGetter) loadDefaults() {
	if !config_loaded {
		loadConfigFile(this.config_file_path)
		config_loaded = true
	}
}

/*
	Load json contents of config file into config_data
	config_file: Path to JSON file.
*/
func loadConfigFile(config_file string) {
	if _, err := os.Stat(config_file); err == nil {
		reader_result, read_err := ioutil.ReadFile(config_file)
		if read_err != nil {
			panic("FATAL ERROR COULD NOT LOAD CONFIG FILE " + config_file + " Err: " + read_err.Error())
		}

		config_map_lock.Lock()
		json_err := json.Unmarshal(reader_result, &json_data)
		config_map_lock.Unlock()
		if json_err != nil {
			panic("FATAL ERROR COULD NOT UNMARSHAL CONFIG FILE " + config_file + " Err: " + json_err.Error())
		}
	}
}
