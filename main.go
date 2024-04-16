package main

import (
	"encoding/json"
	"fmt"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/pflag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

type Config struct {
	URL           string `json:"URL"`
	APP_KEY       string `json:"APP_KEY"`
	ACCESS_KEY    string `json:"ACCESS_KEY"`
	ACCESS_SECRET string `json:"ACCESS_SECRET"`
}

func loadConfig() *Config {
	// 从用户目录/btts/config.json中读取配置
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return nil
	}
	configPath := filepath.Join(home, "btts", "config.json")
	// 判断config.json是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Println("未找到配置文件：", configPath)
		config := createConfig()
		done := createConfigFile(configPath, config)
		if !done {
			log.Println("创建配置文件失败！,请检查权限，程序可以继续执行，但下次运行依然需要再次输入配置")
		}
		return config
	}
	config := &Config{}
	err = LoadConfig(configPath, config)
	if err != nil {
		log.Println(err)
		return nil
	}
	return config
}

func createConfigFile(configPath string, config *Config) bool {
	// 创建目录
	os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	// 写入配置
	file, err := os.Create(configPath)
	if err != nil {
		log.Println(err)
		return false
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func createConfig() *Config {
	// 提示用户从控制台输入三个配置
	fmt.Println("请在控制台输入APP_KEY, ACCESS_KEY, ACCESS_SECRET")
	var URL, APP_KEY, ACCESS_KEY, ACCESS_SECRET string
	fmt.Print("URL(直接回车以使用默认值): ")
	fmt.Scanln(&URL)
	fmt.Print("APP_KEY: ")
	fmt.Scanln(&APP_KEY)
	fmt.Print("ACCESS_KEY: ")
	fmt.Scanln(&ACCESS_KEY)
	fmt.Print("ACCESS_SECRET: ")
	fmt.Scanln(&ACCESS_SECRET)
	config := &Config{
		URL:           URL,
		APP_KEY:       APP_KEY,
		ACCESS_KEY:    ACCESS_KEY,
		ACCESS_SECRET: ACCESS_SECRET,
	}
	return config
}

func LoadConfig(path string, config *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var input string
	var output string
	var printHelpMessage bool
	var clean bool
	var randomCount int
	pflag.StringVarP(&input, "input", "i", "", "输入的剧本文件路径")
	pflag.StringVarP(&output, "output", "o", "", "输出的音频文件的目录")
	pflag.BoolVarP(&printHelpMessage, "help", "h", false, "打印帮助信息")
	pflag.BoolVarP(&clean, "clean", "c", false, "清空输出目录")
	pflag.IntVarP(&randomCount, "random", "r", 0, "随机选择n个，测试效果，0表示全部生成。")
	pflag.Parse()

	if printHelpMessage {
		pflag.Usage()
		return
	}

	if len(input) == 0 {
		fmt.Println("请指定输入的剧本文件路径")
		pflag.Usage()
		return
	}

	if len(output) == 0 {
		fmt.Println("未指定输出音频地址，保存在当前目录下的output文件夹中。")
		output = "output"
	}

	// 读取配置文件
	c := loadConfig()

	// 判断output是否存在，不存在则创建
	if _, err := os.Stat(output); os.IsNotExist(err) {
		fmt.Println("创建输出目录：", output)
		os.Mkdir(output, os.ModePerm)
	} else {
		if clean {
			fmt.Println("清空输出目录：", output)
			os.RemoveAll(output)
			os.Mkdir(output, os.ModePerm)
		}
	}

	if len(c.URL) == 0 {
		log.Println("URL为空，使用默认URL：", nls.DEFAULT_URL)
		c.URL = nls.DEFAULT_URL
	}

	config, err := nls.NewConnectionConfigWithAKInfoDefault(c.URL, c.APP_KEY, c.ACCESS_KEY, c.ACCESS_SECRET)
	if err != nil {
		log.Println(err)
		return
	}

	csv, err := ReadCSV(input)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("剧本解析完成，共", len(csv), "条对话。")
	if randomCount > 0 {
		csv = sampleN(csv, randomCount)
		log.Println("随机选择", randomCount, "条对话。")
	}
	log.Println("开始语音合成任务，共", len(csv), "条对话。")
	// 创建一个新的进度条
	bar := pb.StartNew(len(csv))
	for i, dialogue := range csv {
		prefix := fmt.Sprintf("%d-%s-镜头%d-%s-%d", i, dialogue.Project, dialogue.Shot, dialogue.SpeakingRole, dialogue.Sequence)
		log.Println("Processing dialogue : ", prefix)
		filename := filepath.Join(output, prefix+".wav")
		GenerateWithBasicParam(config, filename, dialogue.Text, dialogue.Speaker, dialogue.Pitch, dialogue.Speed, dialogue.Volume)
		bar.Increment()
	}
	bar.Finish()
	fmt.Println("处理完成！输出文件保存在：", output)
}

// sampleN 从输入的[]*Dialogue中随机选择n个,如果n<=0或n>=len(dialogues)则返回全部
func sampleN(dialogues []*Dialogue, n int) []*Dialogue {
	if n <= 0 || n >= len(dialogues) {
		return dialogues
	}
	// 从dialogues中随机选择n个
	// Fisher-Yates shuffle 打乱数组，取前n个即可实现
	// 没有调用Seed函数，1.20弃用了rand.Seed()，默认使用全局随机数生成器
	rand.Shuffle(len(dialogues), func(i, j int) {
		dialogues[i], dialogues[j] = dialogues[j], dialogues[i]
	})
	return dialogues[:n]
}
