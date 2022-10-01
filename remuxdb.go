package remuxdb

import (
	"encoding/json"
	"os"
	"reflect"
)

type DB struct {
	Name       string
	Collection string
}

var dir string

// Create a new database with this function 🔥!
func NewDB(name string, collection string) DB {
	dir = "./" + name + "/" + collection + ".json"
	if fileExists(dir) {
		return DB{name, collection}
	}
	os.Mkdir("./"+name, os.ModePerm)
	os.Create(dir)
	var data, _ = json.Marshal([]string{})
	os.WriteFile(dir, data, os.ModePerm)
	return DB{name, collection}
}

// Create a new database with this function without creating new files 🔥!
func newDBNotInit(name string, collection string) DB {
	dir = "./" + name + "/" + collection + ".json"
	return DB{name, collection}
}

// Run this when database was made with newDBNotInit
func (db DB) init() {
	if fileExists(dir) {
		return
	}
	os.Mkdir("./"+db.Name, os.ModePerm)
	os.Create(dir)
	var data, _ = json.Marshal([]string{})
	os.WriteFile(dir, data, os.ModePerm)
}

// Read all values from the database
func (db DB) Read(v any) {
	var data, _ = os.ReadFile(dir)
	json.Unmarshal(data, v)
}

// Add a value to the database
func (db DB) Write(v any) {
	var datatype []any
	var data, _ = os.ReadFile(dir)
	json.Unmarshal(data, &datatype)
	datatype = append(datatype, v)
	var data2, _ = json.Marshal(datatype)
	os.WriteFile(dir, data2, os.ModePerm)

}

// Remove a value from the database
func (db DB) Remove(item any) {
	var fromjson []any
	var data, _ = os.ReadFile(dir)
	json.Unmarshal(data, &fromjson)

	var fromitemraw, _ = json.Marshal(item)
	var fromitem map[string]any
	json.Unmarshal(fromitemraw, &fromitem)
	// map[string]string - v
	for i, v := range fromjson {
		if reflect.DeepEqual(v, fromitem) {
			fromjson = remove(fromjson, i)
			var data, _ = json.Marshal(fromjson)
			os.WriteFile(dir, data, os.ModePerm)
		}
	}

}

// Find a value from the database
func (db DB) Find(handler func(index int, item map[string]any) bool, value any) {
	var fromjson []map[string]any
	var data, _ = os.ReadFile(dir)
	json.Unmarshal(data, &fromjson)
	for i, v := range fromjson {
		if handler(i, v) {
			var data, _ = json.Marshal(v)
			json.Unmarshal(data, value)
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func remove(s []any, i int) []any {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
