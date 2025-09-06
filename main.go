package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Tool 表示一个工具的结构
type Tool struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"`
}

// 模拟的工具数据
var tools = []Tool{
	{
		ID:          1,
		Name:        "HTML转图片",
		Description: "将HTML代码转换为图片格式，方便分享和保存",
		Icon:        "bi-file-earmark-image",
		URL:         "/html2img",
	},
	{
		ID:          2,
		Name:        "JSON格式化",
		Description: "格式化和验证JSON数据，使其更易读",
		Icon:        "bi-code-square",
		URL:         "/json-formatter",
	},
	{
		ID:          3,
		Name:        "Base64编码/解码",
		Description: "将文本进行Base64编码或解码",
		Icon:        "bi-arrow-left-right",
		URL:         "/base64-encoder",
	},
	{
		ID:          4,
		Name:        "正则表达式测试",
		Description: "测试和调试正则表达式",
		Icon:        "bi-code-slash",
		URL:         "/regex-tester",
	},
	{
		ID:          5,
		Name:        "URL编码/解码",
		Description: "对URL进行编码或解码操作",
		Icon:        "bi-link-45deg",
		URL:         "/url-encoder",
	},
	{
		ID:          6,
		Name:        "哈希计算",
		Description: "计算文本或文件的MD5、SHA-1、SHA-256等哈希值",
		Icon:        "bi-shield-lock",
		URL:         "/hash-calculator",
	},
	{
		ID:          7,
		Name:        "时间戳转换",
		Description: "时间戳与日期时间相互转换，支持多种格式和时区",
		Icon:        "bi-clock-history",
		URL:         "/timestamp-converter",
	},
	{
		ID:          8,
		Name:        "UUID生成器",
		Description: "生成各种版本的UUID，支持批量生成和格式化输出",
		Icon:        "bi-key",
		URL:         "/uuid-generator",
	},
	{
		ID:          9,
		Name:        "颜色选择器",
		Description: "选择和转换不同格式的颜色值",
		Icon:        "bi-palette",
		URL:         "/color-picker",
	},
}

func main() {
	ConfigRuntime()
	StartGin()
}

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
}

// StartWorkers start starsWorker by goroutine.
func StartWorkers() {
	// go statsWorker()
}

// getTools 处理获取工具列表的请求
func getTools(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tools})
}

// indexHandler 处理主页请求
func indexHandler(c *gin.Context) {
	c.File("resources/static/index.html")
}

// html2imgHandler 处理HTML转图片工具请求
func html2imgHandler(c *gin.Context) {
	c.File("resources/static/html2img/index.html")
}

// jsonFormatterHandler 处理JSON格式化工具请求
func jsonFormatterHandler(c *gin.Context) {
	c.File("resources/static/json-formatter/index.html")
}

// regexTesterHandler 处理正则表达式测试工具请求
func regexTesterHandler(c *gin.Context) {
	c.File("resources/static/regex-tester/index.html")
}

// urlEncoderHandler 处理URL编码解码工具请求
func urlEncoderHandler(c *gin.Context) {
	c.File("resources/static/url-encoder/index.html")
}

// base64EncoderHandler 处理Base64编码解码工具请求
func base64EncoderHandler(c *gin.Context) {
	c.File("resources/static/base64-encoder/index.html")
}

// hashCalculatorHandler 处理哈希计算工具请求
func hashCalculatorHandler(c *gin.Context) {
	c.File("resources/static/hash-calculator/index.html")
}

// timestampConverterHandler 处理时间戳转换工具请求
func timestampConverterHandler(c *gin.Context) {
	c.File("resources/static/timestamp-converter/index.html")
}

// uuidGeneratorHandler 处理UUID生成器工具请求
func uuidGeneratorHandler(c *gin.Context) {
	c.File("resources/static/uuid-generator/index.html")
}

// colorPickerHandler 处理颜色选择器工具请求
func colorPickerHandler(c *gin.Context) {
	c.File("resources/static/color-picker/index.html")
}

// StartGin starts gin web server with setting router.
func StartGin() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	// router.Use(rateLimit, gin.Recovery())
	router.Use(gin.Recovery())
	// router.LoadHTMLGlob("resources/*.templ.html")
	router.Static("/static", "resources/static")
	router.GET("/", indexHandler)
	router.GET("/api/tools", getTools)
	router.GET("/html2img", html2imgHandler)
	router.GET("/json-formatter", jsonFormatterHandler)
	router.GET("/regex-tester", regexTesterHandler)
	router.GET("/url-encoder", urlEncoderHandler)
	router.GET("/base64-encoder", base64EncoderHandler)
	router.GET("/hash-calculator", hashCalculatorHandler)
	router.GET("/timestamp-converter", timestampConverterHandler)
	router.GET("/uuid-generator", uuidGeneratorHandler)
	router.GET("/color-picker", colorPickerHandler)
	// router.GET("/room/:roomid", roomGET)
	// router.POST("/room-post/:roomid", roomPOST)
	// router.GET("/stream/:roomid", streamRoom)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
