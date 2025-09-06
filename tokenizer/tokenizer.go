package tokenizer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

// HuggingFaceTokenizerConfig Hugging Face tokenizer配置结构
type HuggingFaceTokenizerConfig struct {
	Version       string                 `json:"version"`
	Truncation    interface{}            `json:"truncation"`
	Padding       interface{}            `json:"padding"`
	AddedTokens   []AddedToken           `json:"added_tokens"`
	Normalizer    map[string]interface{} `json:"normalizer"`
	PreTokenizer  map[string]interface{} `json:"pre_tokenizer"`
	PostProcessor map[string]interface{} `json:"post_processor"`
	Decoder       map[string]interface{} `json:"decoder"`
	Model         HFModel                `json:"model"`
}

// AddedToken 添加的token结构
type AddedToken struct {
	ID         int    `json:"id"`
	Content    string `json:"content"`
	SingleWord bool   `json:"single_word"`
	LStrip     bool   `json:"lstrip"`
	RStrip     bool   `json:"rstrip"`
	Normalized bool   `json:"normalized"`
	Special    bool   `json:"special"`
}

// HFModel Hugging Face模型结构
type HFModel struct {
	Type                    string         `json:"type"`
	Dropout                 interface{}    `json:"dropout"`
	UnknownToken            string         `json:"unk_token"`
	ContinuingSubwordPrefix string         `json:"continuing_subword_prefix"`
	EndOfWordSuffix         string         `json:"end_of_word_suffix"`
	FuseUnk                 bool           `json:"fuse_unk"`
	ByteFallback            bool           `json:"byte_fallback"`
	IgnoreMerges            bool           `json:"ignore_merges"`
	Vocab                   map[string]int `json:"vocab"`
	Merges                  [][]string     `json:"merges"`
}

// TokenizerConfig tokenizer配置结构
type TokenizerConfig struct {
	Vocabulary    map[string]int    `json:"vocabulary"`
	ReverseVocab  map[int]string    `json:"reverse_vocab"`
	SpecialTokens map[string]string `json:"special_tokens"`
	ModelName     string            `json:"model_name"`
	MaxTokens     int               `json:"max_tokens"`
	Merges        map[string]bool   `json:"merges"` // BPE merges
	IsBPE         bool              `json:"is_bpe"` // 是否为BPE模型
}

// Tokenizer tokenizer结构
type Tokenizer struct {
	config *TokenizerConfig
}

// TokenizerResult tokenizer结果结构
type TokenizerResult struct {
	Tokens       []string `json:"tokens"`
	TokenIDs     []int    `json:"token_ids"`
	TokenCount   int      `json:"token_count"`
	CharCount    int      `json:"char_count"`
	WordCount    int      `json:"word_count"`
	LineCount    int      `json:"line_count"`
	VocabSize    int      `json:"vocab_size"`
	UnknownCount int      `json:"unknown_count"`
	ModelName    string   `json:"model_name"`
}

// NewTokenizer 创建新的tokenizer实例
func NewTokenizer(configPath string) (*Tokenizer, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load tokenizer config: %v", err)
	}

	return &Tokenizer{
		config: config,
	}, nil
}

// loadConfig 加载tokenizer配置
func loadConfig(configPath string) (*TokenizerConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// 尝试解析为Hugging Face格式
	var hfConfig HuggingFaceTokenizerConfig
	if err := json.Unmarshal(data, &hfConfig); err == nil {
		return convertHuggingFaceConfig(&hfConfig)
	}

	// 如果不是Hugging Face格式，尝试解析为旧格式
	var config TokenizerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}

	// 构建反向词汇表
	if config.ReverseVocab == nil {
		config.ReverseVocab = make(map[int]string)
		for token, id := range config.Vocabulary {
			config.ReverseVocab[id] = token
		}
	}

	// 设置默认特殊token
	if config.SpecialTokens == nil {
		config.SpecialTokens = map[string]string{
			"pad":  "<pad>",
			"unk":  "<unk>",
			"bos":  "<s>",
			"eos":  "</s>",
			"mask": "<mask>",
		}
	}

	config.Vocabulary = make(map[string]int)
	// 添加特殊token到词汇表
	for _, specialToken := range config.SpecialTokens {
		if _, exists := config.Vocabulary[specialToken]; !exists {
			config.Vocabulary[specialToken] = len(config.Vocabulary)
		}
	}

	return &config, nil
}

// convertHuggingFaceConfig 转换Hugging Face配置为内部配置
func convertHuggingFaceConfig(hfConfig *HuggingFaceTokenizerConfig) (*TokenizerConfig, error) {
	config := &TokenizerConfig{
		Vocabulary:    make(map[string]int),
		ReverseVocab:  make(map[int]string),
		SpecialTokens: make(map[string]string),
		Merges:        make(map[string]bool),
		ModelName:     "GLM4.5",
		MaxTokens:     512,
	}

	// 检查是否为BPE模型
	config.IsBPE = hfConfig.Model.Type == "BPE"

	// 添加基础词汇表
	if hfConfig.Model.Vocab != nil {
		for token, id := range hfConfig.Model.Vocab {
			config.Vocabulary[token] = id
			config.ReverseVocab[id] = token
		}
	}

	// 处理BPE merges
	if config.IsBPE && len(hfConfig.Model.Merges) > 0 {
		for _, merge := range hfConfig.Model.Merges {
			// merge是数组格式，如 ["Ġ", "ĠĠĠ"]
			if len(merge) == 2 {
				mergeKey := merge[0] + " " + merge[1]
				config.Merges[mergeKey] = true
			}
		}
	}

	// 添加特殊tokens
	for _, addedToken := range hfConfig.AddedTokens {
		config.Vocabulary[addedToken.Content] = addedToken.ID
		config.ReverseVocab[addedToken.ID] = addedToken.Content

		if addedToken.Special {
			// 映射特殊token类型
			switch addedToken.Content {
			case "<s>":
				config.SpecialTokens["bos"] = addedToken.Content
			case "</s>":
				config.SpecialTokens["eos"] = addedToken.Content
			case "<pad>":
				config.SpecialTokens["pad"] = addedToken.Content
			case "<unk>":
				config.SpecialTokens["unk"] = addedToken.Content
			case "<mask>":
				config.SpecialTokens["mask"] = addedToken.Content
			case "[MASK]":
				config.SpecialTokens["mask"] = addedToken.Content
			}
		}
	}

	// 确保有基本的特殊token
	if _, exists := config.SpecialTokens["unk"]; !exists {
		unkToken := hfConfig.Model.UnknownToken
		if unkToken == "" {
			unkToken = "<unk>"
		}
		config.SpecialTokens["unk"] = unkToken
		if _, exists := config.Vocabulary[unkToken]; !exists {
			config.Vocabulary[unkToken] = len(config.Vocabulary)
		}
	}

	return config, nil
}

// Encode 将文本编码为token IDs
func (t *Tokenizer) Encode(text string) ([]int, error) {
	tokens := t.tokenize(text)
	var tokenIDs []int

	for _, token := range tokens {
		if id, exists := t.config.Vocabulary[token]; exists {
			tokenIDs = append(tokenIDs, id)
		} else {
			// 使用unknown token
			if unkID, exists := t.config.Vocabulary[t.config.SpecialTokens["unk"]]; exists {
				tokenIDs = append(tokenIDs, unkID)
			} else {
				tokenIDs = append(tokenIDs, 0)
			}
		}
	}

	return tokenIDs, nil
}

// Decode 将token IDs解码为文本
func (t *Tokenizer) Decode(tokenIDs []int) (string, error) {
	var tokens []string

	for _, id := range tokenIDs {
		if token, exists := t.config.ReverseVocab[id]; exists {
			tokens = append(tokens, token)
		} else {
			tokens = append(tokens, t.config.SpecialTokens["unk"])
		}
	}

	return strings.Join(tokens, " "), nil
}

// Tokenize 对文本进行完整的tokenize处理
func (t *Tokenizer) Tokenize(text string) (*TokenizerResult, error) {
	tokens := t.tokenize(text)
	tokenIDs, err := t.Encode(text)
	if err != nil {
		return nil, err
	}

	// 统计信息
	charCount := len([]rune(text))
	wordCount := len(strings.Fields(text))
	lineCount := len(strings.Split(text, "\n"))

	// 统计unknown tokens
	unknownCount := 0
	for _, token := range tokens {
		if _, exists := t.config.Vocabulary[token]; !exists {
			unknownCount++
		}
	}

	return &TokenizerResult{
		Tokens:       tokens,
		TokenIDs:     tokenIDs,
		TokenCount:   len(tokens),
		CharCount:    charCount,
		WordCount:    wordCount,
		LineCount:    lineCount,
		VocabSize:    len(t.config.Vocabulary),
		UnknownCount: unknownCount,
		ModelName:    t.config.ModelName,
	}, nil
}

// tokenize 基础tokenization逻辑
func (t *Tokenizer) tokenize(text string) []string {
	// 如果是BPE模型，使用BPE算法
	if t.config.IsBPE {
		return t.bpeTokenize(text)
	}

	// 首先尝试直接分词（适用于预训练tokenizer）
	if tokens := t.directTokenize(text); len(tokens) > 0 {
		return tokens
	}

	// 回退到基础分词
	return t.basicTokenize(text)
}

// bpeTokenize BPE分词算法
func (t *Tokenizer) bpeTokenize(text string) []string {
	// 预处理：转换为字符级别
	words := t.preTokenize(text)
	var tokens []string

	for _, word := range words {
		if word == "" {
			continue
		}

		// 如果整个词在词汇表中，直接使用
		if _, exists := t.config.Vocabulary[word]; exists {
			tokens = append(tokens, word)
			continue
		}

		// 否则使用BPE算法
		wordTokens := t.applyBPE(word)
		tokens = append(tokens, wordTokens...)
	}

	return tokens
}

// preTokenize 预分词处理
func (t *Tokenizer) preTokenize(text string) []string {
	// 简单的预分词：按空格和标点符号分割
	var words []string
	var currentWord strings.Builder

	for _, r := range text {
		if unicode.IsSpace(r) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else if unicode.IsPunct(r) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			words = append(words, string(r))
		} else {
			currentWord.WriteRune(r)
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}

// applyBPE 应用BPE算法
func (t *Tokenizer) applyBPE(word string) []string {
	// 转换为字符列表
	characters := strings.Split(word, "")

	// 如果没有merges，回退到字符级别
	if len(t.config.Merges) == 0 {
		return characters
	}

	// 应用BPE合并规则
	for {
		// 寻找可以合并的相邻字符对
		bestMergeIndex := -1

		for i := 0; i < len(characters)-1; i++ {
			pair := characters[i] + " " + characters[i+1]
			if t.config.Merges[pair] {
				bestMergeIndex = i
				break // 找到第一个可以合并的就处理
			}
		}

		if bestMergeIndex == -1 {
			break // 没有更多可以合并的
		}

		// 合并字符
		merged := characters[bestMergeIndex] + characters[bestMergeIndex+1]
		characters = append(characters[:bestMergeIndex], append([]string{merged}, characters[bestMergeIndex+2:]...)...)
	}

	return characters
}

// directTokenize 直接使用词汇表进行分词
func (t *Tokenizer) directTokenize(text string) []string {
	var tokens []string
	words := strings.Fields(text)

	for _, word := range words {
		// 尝试完整匹配
		if _, exists := t.config.Vocabulary[word]; exists {
			tokens = append(tokens, word)
			continue
		}

		// 尝试小写匹配
		lowerWord := strings.ToLower(word)
		if _, exists := t.config.Vocabulary[lowerWord]; exists {
			tokens = append(tokens, lowerWord)
			continue
		}

		// 如果没有找到，尝试子词匹配
		subTokens := t.subwordTokenize(word)
		tokens = append(tokens, subTokens...)
	}

	return tokens
}

// subwordTokenize 子词分词
func (t *Tokenizer) subwordTokenize(word string) []string {
	var tokens []string
	lowerWord := strings.ToLower(word)

	// 简单的子词匹配：尝试从长到短匹配
	for length := len(lowerWord); length > 0; length-- {
		for i := 0; i <= len(lowerWord)-length; i++ {
			subword := lowerWord[i : i+length]
			if _, exists := t.config.Vocabulary[subword]; exists {
				tokens = append(tokens, subword)
				// 跳过已匹配的部分
				i += length - 1
				break
			}
		}
	}

	// 如果没有匹配到任何子词，回退到字符级别
	if len(tokens) == 0 {
		for _, r := range lowerWord {
			tokens = append(tokens, string(r))
		}
	}

	return tokens
}

// basicTokenize 基础分词逻辑
func (t *Tokenizer) basicTokenize(text string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			currentToken.WriteRune(r)
		} else {
			if currentToken.Len() > 0 {
				token := strings.ToLower(currentToken.String())
				if _, exists := t.config.Vocabulary[token]; exists {
					tokens = append(tokens, token)
				}
				currentToken.Reset()
			}

			// 处理标点符号和空格
			if unicode.IsSpace(r) {
				continue
			}

			// 处理常见标点
			punct := string(r)
			if _, exists := t.config.Vocabulary[punct]; exists {
				tokens = append(tokens, punct)
			}
		}
	}

	// 添加最后一个token
	if currentToken.Len() > 0 {
		token := strings.ToLower(currentToken.String())
		if _, exists := t.config.Vocabulary[token]; exists {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

// GetVocabSize 获取词汇表大小
func (t *Tokenizer) GetVocabSize() int {
	return len(t.config.Vocabulary)
}

// GetVocabulary 获取词汇表
func (t *Tokenizer) GetVocabulary() map[string]int {
	return t.config.Vocabulary
}

// GetTopTokens 获取前N个最常用的token
func (t *Tokenizer) GetTopTokens(n int) []string {
	type tokenCount struct {
		token string
		id    int
	}

	var tokens []tokenCount
	for token, id := range t.config.Vocabulary {
		tokens = append(tokens, tokenCount{token, id})
	}

	// 按ID排序（通常ID越小越常用）
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].id < tokens[j].id
	})

	var result []string
	for i := 0; i < n && i < len(tokens); i++ {
		result = append(result, tokens[i].token)
	}

	return result
}

// FindTokenByID 根据ID查找token
func (t *Tokenizer) FindTokenByID(id int) (string, bool) {
	token, exists := t.config.ReverseVocab[id]
	return token, exists
}

// FindTokensByPrefix 根据前缀查找tokens
func (t *Tokenizer) FindTokensByPrefix(prefix string) []string {
	var result []string
	for token := range t.config.Vocabulary {
		if strings.HasPrefix(token, prefix) {
			result = append(result, token)
		}
	}
	return result
}

// PrintStats 打印tokenizer统计信息
func (t *Tokenizer) PrintStats() {
	fmt.Printf("Tokenizer Statistics:\n")
	fmt.Printf("Model Name: %s\n", t.config.ModelName)
	fmt.Printf("Model Type: %s\n", map[bool]string{true: "BPE", false: "Word"}[t.config.IsBPE])
	fmt.Printf("Vocabulary Size: %d\n", len(t.config.Vocabulary))
	fmt.Printf("Special Tokens: %v\n", t.config.SpecialTokens)
	fmt.Printf("Max Tokens: %d\n", t.config.MaxTokens)
	if t.config.IsBPE {
		fmt.Printf("BPE Merges: %d\n", len(t.config.Merges))
	}

	// 显示前20个最常用的token
	topTokens := t.GetTopTokens(20)
	fmt.Printf("Top 20 tokens: %v\n", topTokens)
}
