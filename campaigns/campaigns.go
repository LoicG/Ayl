package data

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Content struct {
	Title       string
	Description string
	Landing     string
}

type MyList struct {
	Elements map[string]struct{}
}

func (l *MyList) UnmarshalJSON(b []byte) error {
	tmp := []string{}
	if err := json.Unmarshal(b, &tmp); err != nil {
		return err
	}
	l.Elements = make(map[string]struct{})
	for _, v := range tmp {
		l.Elements[v] = struct{}{}
	}
	return nil
}

func (l *MyList) MarshalJSON() ([]byte, error) {
	data := []string{}
	for k := range l.Elements {
		data = append(data, k)
	}
	return json.Marshal(data)
}

func MakeList(input []string) MyList {
	elements := make(map[string]struct{})
	for _, v := range input {
		elements[v] = struct{}{}
	}
	return MyList{
		Elements: elements,
	}
}

type Campaign struct {
	Price      float32
	Content    *Content
	Countries  MyList `json:"countries,omitempty"`
	Devices    MyList `json:"devices,omitempty"`
	Placements MyList `json:"placements,omitempty"`
}

type Campaigns struct {
	Elements map[string]*Campaign `json:"campaigns"`
}

func LoadJSONFile(path string) (*Campaigns, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	result := &Campaigns{}
	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func SaveJSONFile(path string, campaigns *Campaigns) error {
	data, err := json.MarshalIndent(campaigns, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(path, "campaigns.json"), data, 0644)
}
