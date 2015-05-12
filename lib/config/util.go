import (
	"encoding/json"
	"os"
	"io/ioutil"
)

func FileExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) { return err }
	return nil
}

func SaveJson(filename string, c interface{}) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	
	err = ioutil.WriteFile(filename, data, os.ModePerm)
	if err != nil {
		return err
	}
	
	return nil
}

func Load(filename string, c interface{}) error {
	err := FileExists(filename)
	if err != nil {
		return nil
	}
	
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	
	err = json.Unmarshal(data, c)
	if err != nil {
		return err
	}
	
	return nil
}