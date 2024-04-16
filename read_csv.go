package main

import "os"
import "github.com/gocarina/gocsv"

type Dialogue struct {
	Project      string `csv:"项目"`
	Shot         int    `csv:"镜头"`
	SpeakingRole string `csv:"发言角色"`
	Sequence     int    `csv:"序号"`
	Text         string `csv:"文本"`
	Speaker      string `csv:"发言人"`
	Pitch        int    `csv:"语调"`
	Speed        int    `csv:"语速"`
	Volume       int    `csv:"音量"`
	AudioFormat  string `csv:"音频格式"`
}

func ReadCSV(filename string) ([]*Dialogue, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	dialogues := []*Dialogue{}

	if err := gocsv.UnmarshalFile(file, &dialogues); err != nil { // Load Employees from file
		return nil, err
	}
	return dialogues, nil
}
