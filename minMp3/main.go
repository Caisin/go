package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	files, _ := ioutil.ReadDir(".")
	dir, _ := os.Getwd()

	path := dir + string(os.PathSeparator)
	outPath := path + "out"
	_ = os.RemoveAll(outPath + "/")
	_ = os.Mkdir(outPath, os.ModePerm)
	log.Printf(dir)
	for _, file := range files {
		name := file.Name()
		suffix := strings.HasSuffix(name, ".mp3")
		if suffix {
			outFile := outPath + string(os.PathSeparator) + name
			log.Printf("开始压缩文件 %s", name)
			log.Printf("输出目录 %s", outFile)
			ffmpegPath := path + "ffmpeg.exe"
			cmd := exec.Command(ffmpegPath,
				"-i", name,
				"-b:a", "16k",
				"-acodec", "mp3",
				"-ar", "11025",
				"-ac", "1", outFile)
			// 获取输出对象，可以从该对象中读取输出结果
			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}
			// 保证关闭输出流
			defer stdout.Close()
			// 运行命令
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
			// 读取输出结果
			opBytes, err := ioutil.ReadAll(stdout)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(opBytes))
		}
	}
	log.Printf("压缩完成!")
}
