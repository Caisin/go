package html2pdf

import (
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
