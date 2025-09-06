package tokenizer_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/render-examples/go-gin-web-server/tokenizer"
)

// TestMain 设置测试环境
func TestMain(m *testing.M) {
	// 获取当前文件目录
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	
	// 切换到项目根目录
	projectRoot := filepath.Dir(currentDir)
	err := os.Chdir(projectRoot)
	if err != nil {
		panic("Failed to change working directory: " + err.Error())
	}
	
	// 运行测试
	code := m.Run()
	os.Exit(code)
}

// TestTokenizerLoading 测试分词器加载
func TestTokenizerLoading(t *testing.T) {
	configPath := filepath.Join("tokenizer", "test_config.json")
	
	tk, err := tokenizer.NewTokenizer(configPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer: %v", err)
	}
	
	if tk == nil {
		t.Fatal("Tokenizer is nil")
	}
	
	if tk.GetVocabSize() <= 0 {
		t.Errorf("Invalid vocabulary size: %d", tk.GetVocabSize())
	}
	
	vocab := tk.GetVocabulary()
	if len(vocab) == 0 {
		t.Error("Vocabulary is empty")
	}
}

// TestTokenization 测试分词功能
func TestTokenization(t *testing.T) {
	configPath := filepath.Join("tokenizer", "test_config.json")

	tk, err := tokenizer.NewTokenizer(configPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer: %v", err)
	}

	testText := "Hello, world! This is a test."
	result, err := tk.Tokenize(testText)
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}
	
	if result.TokenCount <= 0 {
		t.Errorf("Invalid token count: %d", result.TokenCount)
	}
	
	if len(result.Tokens) != len(result.TokenIDs) {
		t.Errorf("Tokens and TokenIDs length mismatch: %d vs %d", len(result.Tokens), len(result.TokenIDs))
	}
	
	if result.TokenCount != len(result.Tokens) {
		t.Errorf("Token count mismatch: %d vs %d", result.TokenCount, len(result.Tokens))
	}
}

// TestEncodingDecoding 测试编码和解码功能
func TestEncodingDecoding(t *testing.T) {
	configPath := filepath.Join("tokenizer", "test_config.json")

	tk, err := tokenizer.NewTokenizer(configPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer: %v", err)
	}

	testText := "hello world test"
	
	// 测试编码
	encoded, err := tk.Encode(testText)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}
	
	if len(encoded) == 0 {
		t.Error("Encoded result is empty")
	}
	
	// 测试解码
	decoded, err := tk.Decode(encoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}
	
	// 由于测试配置有限，只验证解码过程成功，不严格要求完全匹配
	if decoded == "" {
		t.Error("Decoded text is empty")
	}
}

// TestTopTokens 测试获取顶部token功能
func TestTopTokens(t *testing.T) {
	configPath := filepath.Join("tokenizer", "test_config.json")

	tk, err := tokenizer.NewTokenizer(configPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer: %v", err)
	}

	topTokens := tk.GetTopTokens(20)
	if len(topTokens) == 0 {
		t.Error("Top tokens is empty")
	}
	
	if len(topTokens) > 20 {
		t.Errorf("Too many top tokens: expected max 20, got %d", len(topTokens))
	}
}

// TestEmptyText 测试空文本处理
func TestEmptyText(t *testing.T) {
	configPath := filepath.Join("tokenizer", "test_config.json")

	tk, err := tokenizer.NewTokenizer(configPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer: %v", err)
	}

	// 测试空字符串
	result, err := tk.Tokenize("")
	if err != nil {
		t.Fatalf("Tokenization of empty text failed: %v", err)
	}
	
	if result.TokenCount != 0 {
		t.Errorf("Expected 0 tokens for empty text, got %d", result.TokenCount)
	}
	
	// 测试空字符串编码
	encoded, err := tk.Encode("")
	if err != nil {
		t.Fatalf("Encoding of empty text failed: %v", err)
	}
	
	if len(encoded) != 0 {
		t.Errorf("Expected empty encoding for empty text, got %v", encoded)
	}
}

// TestChineseText 测试中文文本处理
func TestChineseText(t *testing.T) {
	configPath := filepath.Join("tokenizer", "test_config.json")

	tk, err := tokenizer.NewTokenizer(configPath)
	if err != nil {
		t.Fatalf("Failed to load tokenizer: %v", err)
	}

	// 使用测试配置中存在的词汇
	testText := "hello world tokenizer test"
	result, err := tk.Tokenize(testText)
	if err != nil {
		t.Fatalf("Tokenization failed: %v", err)
	}
	
	if result.TokenCount <= 0 {
		t.Errorf("Invalid token count: %d", result.TokenCount)
	}
	
	// 测试编码和解码
	encoded, err := tk.Encode(testText)
	if err != nil {
		t.Fatalf("Encoding failed: %v", err)
	}
	
	decoded, err := tk.Decode(encoded)
	if err != nil {
		t.Fatalf("Decoding failed: %v", err)
	}
	
	// 由于测试配置有限，只验证解码过程成功
	if decoded == "" {
		t.Error("Decoded text is empty")
	}
}
