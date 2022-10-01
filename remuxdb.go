package remuxdb

import (
	"encoding/json"
	"os"
	"reflect"
)

type DB struct {
	Name       string
	Collection string
	Dir        string
}

var dir string

// Create a new database with this function 🔥!
func NewDB(name string, collection string) DB {
	dir = "./" + name + "/" + collection + ".json"
	if fileExists(dir) {
		return DB{name, collection, dir}
	}
	os.Mkdir("./"+name, os.ModePerm)
	os.Create(dir)
	var data, _ = json.Marshal([]string{})
	os.WriteFile(dir, data, os.ModePerm)
	return DB{name, collection, dir}
}

// Create a new database with this function without creating new files 🔥!
func NewDBNotInit(name string, collection string) DB {
	dir = "./" + name + "/" + collection + ".json"
	return DB{name, collection, dir}
}

// Run this when database was made with newDBNotInit
func (db DB) Init() {
	if fileExists(db.Dir) {
		return
	}
	os.Mkdir("./"+db.Name, os.ModePerm)
	os.Create(db.Dir)
	var data, _ = json.Marshal([]string{})
	os.WriteFile(db.Dir, data, os.ModePerm)
}

// Read all values from the database
func (db DB) Read(v any) {
	var data, _ = os.ReadFile(db.Dir)
	json.Unmarshal(data, v)
}

// Add a value to the database
func (db DB) Write(v any) {
	var datatype []any
	var data, _ = os.ReadFile(db.Dir)
	json.Unmarshal(data, &datatype)
	datatype = append(datatype, v)
	var data2, _ = json.Marshal(datatype)
	os.WriteFile(dir, data2, os.ModePerm)

}

// Remove a value from the database
func (db DB) Remove(item any) {
	var fromjson []map[string]any
	var data, _ = os.ReadFile(db.Dir)
	json.Unmarshal(data, &fromjson)

	var fromitemraw, _ = json.Marshal(item)
	var fromitem map[string]any
	json.Unmarshal(fromitemraw, &fromitem)

	for i, v := range fromjson {
		if reflect.DeepEqual(v, fromitem) {
			fromjson = remove(fromjson, i)
			var data, _ = json.Marshal(fromjson)
			os.WriteFile(db.Dir, data, os.ModePerm)
		}
	}

}

// Find a value from the database
func (db DB) Find(handler func(index int, item map[string]any) bool, value any) {
	var fromjson []map[string]any
	var data, _ = os.ReadFile(db.Dir)
	json.Unmarshal(data, &fromjson)
	for i, v := range fromjson {
		for k, v2 := range v {
			switch reflect.ValueOf(v2).Type().Name() {
			case "float64":
				v[k] = int(v2.(float64))
			}
		}
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

func remove(s []map[string]interface{}, i int) []map[string]interface{} {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
