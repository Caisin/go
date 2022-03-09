package html2pdf

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestGitBookCrawler_Run(t *testing.T) {
	htmlTemplate := `
	<!DOCTYPE HTML>
	<html lang="zh-hans" >
	<head>
	<meta charset="UTF-8">
	<meta content="text/html; charset=utf-8" http-equiv="Content-Type">
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<meta name="description" content="">
	<meta charset="UTF-8">
	<link rel="stylesheet" href="{start_url}/gitbook/style.css">
	<link rel="stylesheet" href="{start_url}/gitbook/gitbook-plugin-katex/katex.min.css">
	<link rel="stylesheet" href="{start_url}/gitbook/gitbook-plugin-highlight/website.css">
	<link rel="stylesheet" href="{start_url}/gitbook/gitbook-plugin-fontsettings/website.css">
	<meta name="HandheldFriendly" content="true"/>
	<meta name="viewport" content="width=device-width, initial-scale=1, user-scalable=no">
	<meta name="apple-mobile-web-app-capable" content="yes">
	<meta name="apple-mobile-web-app-status-bar-style" content="black">
	<link rel="apple-touch-icon-precomposed" sizes="152x152" href="{start_url}/gitbook/images/apple-touch-icon-precomposed-152.png">
	<link rel="shortcut icon" href="{start_url}/gitbook/images/favicon.ico" type="image/x-icon">
	<style>
	.markdown-section pre>code {
		white-space: pre-wrap;
		word-wrap: break-word;
	}
	</style>
	</head>
	<body style="font-size: xx-large;padding:50px 10px">
	{content}
	</body>
	</html>
		`
	crawler, err := NewGitBookCrawler("http://localhost/gitbook",
		"Go语言圣经",
		"E:/code/golang/goLearn/out",
		htmlTemplate, "section.markdown-section", "ul.summary a")
	if err != nil {
		t.Error(err.Error())
		return
	}
	crawler.Run()
}

//爬取文档内容
func TestGraceUICrawler_Run(t *testing.T) {
	htmlTemplate := `
	<!DOCTYPE HTML>
	<html><head>
<meta charset="utf-8">
<title>GraceUI 底部对话框组件</title>
<meta name="keywords" content="GraceUI 底部对话框组件">
<meta name="description" content="GraceUI 底部对话框组件">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
<link rel="stylesheet" type="text/css" href="https://at.alicdn.com/t/font_611166_xrqs4mc71re.css">
<link rel="stylesheet" href="http://grace.hcoder.net/statics/css/graceV1.3.css">
<style type="text/css">
.gui-mini-map{width:200px; position:fixed; z-index:2; right:0; top:0; display:none;}
.gui-mini-map-title{line-height:18px; margin-top:18px; color:#888888 !important; font-size:13px; border-left:1px solid #F1F1F1; display:block; padding:0 15px;}
.gui-mini-map-title:hover{border-color:#166BD8;}
.gui-mlink{}
.gui-manual-body {
		margin-left: 0;
	}
</style>
</head>
	<body style="font-size: xx-large;padding:50px 10px">
<div class="gui-manual-body">
	{content}
</div>
	</body>
	</html>
		`
	crawler, err := NewGitBookCrawler("http://grace.hcoder.net/manual",
		"graceUI",
		"E:/code/golang/goLearn/out",
		htmlTemplate, "div.gui-manual-body", "#grace-accordion a")
	if err != nil {
		t.Error(err.Error())
		return
	}
	crawler.Cookies = map[string]string{
		"Hm_lpvt_0a265394dba26d422dd4ea9530ff4edd": "1588507332",
		"Hm_lvt_0a265394dba26d422dd4ea9530ff4edd":  "1586270505,1586456126,1586968552,1588451059",
		"PHPSESSID":      "6557888f19677d69b95501571672aecf",
		"UM_distinctid":  "16fed3e639e7d1-09019096f664cc-6701b35-144000-16fed3e639f783",
		"graceLoaclUser": "%5B%22150600%22%2C%225eadd7058da54%22%5D",
	}
	crawler.Run()
}

//爬取文档内容
func TestGraceUICrawlerHtml(t *testing.T) {
	htmlTemplate := `
	<!DOCTYPE HTML>
	<html><head>
<meta charset="utf-8">
<title>GraceUI 随机数及字符串生成</title>
<meta name="keywords" content="graceUI,随机数及字符串生成工具">
<meta name="description" content="garceUI 随机数及字符串生成工具 [ graceUI/jsTools/random.js ]">
<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
<link rel="stylesheet" type="text/css" href="https://at.alicdn.com/t/font_611166_xrqs4mc71re.css">
<script type="text/javascript" src="{Domain}/statics/js/jquery.js"></script>
<script type="text/javascript" src="{Domain}/statics/js/grace_v1.1.js"></script>
<link rel="stylesheet" href="{Domain}/statics/css/graceV1.3.css">
<style type="text/css">
.gui-mini-map{width:200px; position:fixed; z-index:2; right:0; top:0; display:none;}
.gui-mini-map-title{line-height:18px; margin-top:18px; color:#888888 !important; font-size:13px; border-left:1px solid #F1F1F1; display:block; padding:0 15px;}
.gui-mini-map-title:hover{border-color:#166BD8;}
.gui-mlink{}
</style>
</head>
	<body style="font-size: xx-large;padding:50px 10px">
{content}
	</body>
	</html>
		`
	crawler, err := NewGitBookCrawler("http://grace.hcoder.net/manual",
		"graceUI",
		"E:/code/golang/goLearn/out/doc",
		htmlTemplate, "body", "#grace-accordion a")
	if err != nil {
		t.Error(err.Error())
		return
	}
	crawler.Cookies = map[string]string{
		"Hm_lpvt_0a265394dba26d422dd4ea9530ff4edd": "1588507332",
		"Hm_lvt_0a265394dba26d422dd4ea9530ff4edd":  "1586270505,1586456126,1586968552,1588451059",
		"PHPSESSID":      "6557888f19677d69b95501571672aecf",
		"UM_distinctid":  "16fed3e639e7d1-09019096f664cc-6701b35-144000-16fed3e639f783",
		"graceLoaclUser": "%5B%22150600%22%2C%225eadd7058da54%22%5D",
	}
	crawler.GenPdf = false
	crawler.Run()
}

func TestFlutter_Run(t *testing.T) {
	os.Setenv("WKHTMLTOPDF_PATH", "D:/work/software/wkhtmltopdf/bin")
	file, _ := ioutil.ReadFile("D:/work/code/go/goLearn/html2pdf/Flutter.html")
	htmlTemplate := string(file)
	crawler, err := NewGitBookCrawler("https://book.flutterchina.club",
		"Flutter实战·第二版",
		"D:/work/code/go/goLearn/out",
		htmlTemplate, "main.page div.theme-default-content", "aside.sidebar .sidebar-links a")
	if err != nil {
		t.Error(err.Error())
		return
	}
	crawler.GenPdf = true
	crawler.MenuSorter = func(menus []string) []string {
		sort.Slice(menus, func(i, j int) bool {
			/*c := "chapter"
			a := menus[i]
			b := menus[j]
			if strings.Contains(a, c) && !strings.Contains(b, c) {
				return false
			}
			dira := path.Dir(a)
			dirb := path.Dir(b)
			if dira == dirb {
				if strings.Contains(a, "index") && !strings.Contains(b, "index") {
					return false
				}
			}
			aN := cast.ToInt(strings.ReplaceAll(dira, c, ""))
			bN := cast.ToInt(strings.ReplaceAll(dirb, c, ""))
			if aN > bN {
				return false
			}*/
			return false
		})
		for i, menu := range menus {
			menus[i] = strings.ReplaceAll(strings.ReplaceAll(menu, " ", ""), ".md", "")
		}
		return menus
	}
	crawler.Run()
}
