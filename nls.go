package main

import (
	"errors"
	nls "github.com/aliyun/alibabacloud-nls-go-sdk"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type TtsUserParam struct {
	F io.Writer
	//Logger *nls.NlsLogger
}

func onTaskFailed(text string, param interface{}) {
	_, ok := param.(*TtsUserParam)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}
	log.Println(text, param)

	//p.Logger.Println("TaskFailed:", text)
}

func onSynthesisResult(data []byte, param interface{}) {
	p, ok := param.(*TtsUserParam)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}
	p.F.Write(data)
}

func onCompleted(text string, param interface{}) {
	_, ok := param.(*TtsUserParam)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	//p.Logger.Println("onCompleted:", text)
}

func onClose(param interface{}) {
	_, ok := param.(*TtsUserParam)
	if !ok {
		log.Default().Fatal("invalid logger")
		return
	}

	//p.Logger.Println("onClosed:")
}

func waitReady(ch chan bool) error {
	select {
	case done := <-ch:
		{
			if !done {
				log.Println("Wait failed")
				return errors.New("wait failed")
			}
			//logger.Println("Wait done")
		}
	case <-time.After(60 * time.Second):
		{
			log.Println("Wait timeout")
			return errors.New("wait timeout")
		}
	}
	return nil
}

var lk sync.Mutex
var fail = 0
var reqNum = 0

// GenerateWithDefaultParam 用默认配置生成合成语音，生成的语音文件保存在fname中
// config: 连接配置
// fname: 保存文件名
// voice: 发音人
// text: 合成文本
func GenerateWithDefaultParam(config *nls.ConnectionConfig, fname string, voice string, text string) {

	GenerateWithAllParam(config, fname, text, voice, "wav", 16000, 50, 0, 0, false)
}

// GenerateWithBasicParam 用基本参数生成合成语音
// config: 连接配置
// fname: 保存文件名
// text: 合成文本
// voice: 发音人
// volume: 音量，范围是0~100，默认值为50。
// speechRate: 语速，范围是-500~500，默认值为0。
// pitchRate: 语调 语调，范围是-500~500，默认值为0。
// enableSubtitle: 是否开启字幕，默认不开启。
func GenerateWithBasicParam(config *nls.ConnectionConfig, fname string, text string, voice string, volume int,
	speechRate int, pitchRate int) {
	GenerateWithAllParam(config, fname, text, voice, "wav", 16000, volume, speechRate, pitchRate, false)
}

// GenerateWithAllParam 用所有参数生成合成语音
// config: 连接配置
// fname: 保存文件名
// text: 合成文本
// voice: 发音人
// format: 音频格式，默认使用WAV。
// sampleRate: 音频采样率，支持16000或8000，默认值为16000。
// volume: 音量，范围是0~100，默认值为50。
// speechRate: 语速，范围是-500~500，默认值为0。
// pitchRate: 语调，范围是-500~500，默认值为0。
// enableSubtitle: 是否开启字幕，默认不开启。
func GenerateWithAllParam(config *nls.ConnectionConfig, fname string, text string, voice string, format string, sampleRate int,
	volume int, speechRate int, pitchRate int, enableSubtitle bool) {
	param := nls.DefaultSpeechSynthesisParam()
	// 发音人，默认值：“xiaoyun”。
	param.Voice = voice
	// 音频格式，默认使用WAV。
	param.Format = format
	// 音频采样率，支持16000或8000，默认值为16000。
	param.SampleRate = sampleRate
	// 音量，范围是0~100，默认值为50。
	param.Volume = volume
	// 语速，范围是-500~500，默认值为0。
	param.SpeechRate = speechRate
	// 语调，范围是-500~500，默认值为0。
	param.PitchRate = pitchRate
	// 是否开启字幕，默认不开启。
	param.EnableSubtitle = enableSubtitle

	ttsUserParam := new(TtsUserParam)
	fout, err := os.OpenFile(fname, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)

	ttsUserParam.F = fout
	//ttsUserParam.Logger = logger
	//第三个参数控制是否请求长文本语音合成，false为短文本语音合成
	tts, err := nls.NewSpeechSynthesis(config, nil, false,
		onTaskFailed, onSynthesisResult, nil,
		onCompleted, onClose, ttsUserParam)
	if err != nil {
		log.Fatalln(err)
		return
	}

	lk.Lock()
	reqNum++
	lk.Unlock()
	//logger.Println("SR start")
	ch, err := tts.Start(text, param, nil)
	if err != nil {
		lk.Lock()
		fail++
		lk.Unlock()
		tts.Shutdown()
	}

	err = waitReady(ch)
	if err != nil {
		lk.Lock()
		fail++
		lk.Unlock()
		tts.Shutdown()
	}
	//logger.Println("Synthesis done")
	tts.Shutdown()

}
