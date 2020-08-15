package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Add(data map[string]interface{}, keys []string, value string) error {
	node, err := findNodeByKeys(data, keys)
	if err != nil {
		return err
	}

	if _, ok := node[keys[len(keys)-1]]; ok {
		return errors.New("key already exists")
	}

	node[keys[len(keys)-1]] = value

	return nil
}

func Set(data map[string]interface{}, keys []string, value string) error {
	node, err := findNodeByKeys(data, keys)
	if err != nil {
		return err
	}
	node[keys[len(keys)-1]] = value

	return nil
}

func findNodeByKeys(node map[string]interface{}, keys []string) (map[string]interface{}, error) {
	for _, k := range keys[:len(keys)-1] {
		if v, ok := node[k]; ok {
			switch v := v.(type) {
			case map[string]interface{}:
				node = v
			default:
				return nil, errors.New("invalid key")
			}
		} else {
			node[k] = make(map[string]interface{})
			node = node[k].(map[string]interface{})
		}
	}
	return node, nil
}

func Remove(data map[string]interface{}, keys []string) error {
	node := data
	for _, k := range keys[:len(keys)-1] {
		if _, ok := node[k]; !ok {
			return errors.New("invalid key")
		}
		node = node[k].(map[string]interface{})
	}

	if _, ok := node[keys[len(keys)-1]]; !ok {
		return errors.New("invalid key")
	}

	delete(node, keys[len(keys)-1])
	return nil
}

func Get(data map[string]interface{}, keys []string) (map[string]string, error) {
	if len(keys) == 0 {
		return getRoot(data)
	}

	node := data
	for _, k := range keys[:len(keys)-1] {
		if _, ok := node[k]; !ok {
			return nil, fmt.Errorf("key(%v) does not exists", k)
		}
		node = node[k].(map[string]interface{})
	}

	if v, ok := node[keys[len(keys)-1]]; ok {
		switch v := v.(type) {
		case map[string]interface{}:
			return getKeyValueSet(v, "", map[string]string{})
		case string:
			return map[string]string{".": v}, nil
		default:
			return nil, errors.New("unexpected data in store file")
		}
	}

	return nil, fmt.Errorf("key(%v) does not exists", keys[len(keys)-1])
}

func getRoot(data map[string]interface{}) (map[string]string, error) {
	return getKeyValueSet(data, "", map[string]string{})
}

func getKeyValueSet(node map[string]interface{}, key string, values map[string]string) (map[string]string, error) {
	var err error
	for k, v := range node {
		switch v := v.(type) {
		case map[string]interface{}:
			values, err = getKeyValueSet(v, fmt.Sprintf("%v.%v", key, k), values)
			if err != nil {
				return nil, err
			}
		case string:
			values[fmt.Sprintf("%v.%v", key, k)] = v
		default:
			return nil, errors.New("unexpected data in store file")
		}
	}
	return values, nil
}

func IsLeaf(data map[string]interface{}, keys []string) (bool, error) {
	node := data
	for _, k := range keys[:len(keys)-1] {
		if _, ok := node[k]; !ok {
			return false, fmt.Errorf("key(%v) does not exists", k)
		}
		node = node[k].(map[string]interface{})
	}

	if v, ok := node[keys[len(keys)-1]]; ok {
		switch v.(type) {
		case map[string]interface{}:
			return false, nil
		case string:
			return true, nil
		default:
			return false, errors.New("unexpected data in store file")
		}
	}

	return false, errors.New("invalid key")
}

func SplitKeys(s string) ([]string, error) {
	if !strings.HasPrefix(s, ".") || (len(s) > 1 && strings.HasSuffix(s, ".")) {
		return nil, errors.New("invalid key")
	}

	if len(s) == 1 {
		return []string{}, nil
	}

	return strings.Split(s, ".")[1:], nil
}

func Save(data map[string]interface{}, file string) error {
	bytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if os.IsNotExist(err) {
		err = os.Mkdir(filepath.Dir(file), 0755)
		if err != nil {
			return err
		}
		f, err = os.Create(file)
	}
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(bytes)
	return err
}

func Load(file string) (map[string]interface{}, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			err := Save(map[string]interface{}{}, file)
			if err != nil {
				return nil, err
			}
			return map[string]interface{}{}, nil
		}
		return nil, err
	}

	data := make(map[string]interface{})
	err = json.Unmarshal(bytes, &data)
	return data, err
}

func DataFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".store", "data.json"), nil
}
