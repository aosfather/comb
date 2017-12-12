package aiml

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ASTERISK Sentence = CreateSentence(" * ", []int32{0, 2}, " * ")
	regex    []string = []string{"\\.", "\\*", "\\+", "\\[", "\\^", "\\-", "\\]", "\\(", "\\)", "\\?", "\\|", "\\{", "\\}", "\\$"}
)

func EscapeRegex(source string) string {
	splitter := strings.Replace(source, "\\\\", "\\\\\\\\", -1)

	for _, special := range regex {
		splitter = strings.Replace(splitter, special, "\\"+special, -1)
	}

	return splitter

}

//构建句子对象
func CreateSentence(original string, mapping []int32, normalized string) Sentence {
	s := Sentence{}
	s.original = original
	s.mappings = mapping
	s.setNormalized(normalized)
	return s
}

type Sentence struct {
	original   string
	mappings   []int32 // The index mappings of normalized elements to original elements.
	normalized string
	splitted   []string // The normalized entry, splitted in an array of words.
}

func (this *Sentence) TrimOriginal() string {
	return this.original
}

func (this *Sentence) Original(begin, end int) string {
	if begin < 0 {
		return ""
	}
	//开始值过大
	for {
		if begin >= 0 && this.mappings[begin] == 0 {
			begin--
		} else {
			break
		}
	}

	//处理end
	n := len(this.mappings)

	for {
		if end < n && this.mappings[end] == 0 {
			end++
		}
	}

	if end >= n {
		end = n - 1
	}

	result := this.original[this.mappings[begin] : this.mappings[end]+1]
	result = strings.Replace(result, "^[^A-Za-z0-9]+|[^A-Za-z0-9]+$", " ", -1)
	return result
}

func (this *Sentence) Normalized(index int) string {
	return this.splitted[index]
}

func (this *Sentence) Length() int {
	return len(this.splitted)
}

func (this *Sentence) setNormalized(str string) {
	this.normalized = str
	if str != "" {
		this.splitted = strings.Split(strings.TrimSpace(str), " ")
	}

}

/*-------------------------------------------------------------

---------------------------------------------------------------*/
type SentenceSplitter struct {
	protection map[string]string //Map of sentence-protection substitution patterns
	splitters  []string          //List of sentence-spliting patterns
	pattern    *regexp.Regexp    //The regular expression which will split entries by sentence splitters
}

func (this *SentenceSplitter) Init(protection map[string]string, splitters []string) {
	this.protection = protection

	for index, value := range splitters {
		splitters[index] = EscapeRegex(value)
	}

	this.splitters = splitters
	splitPattern := "\\s*("

	splitPattern += strings.Join(this.splitters, "|")
	splitPattern += ")\\s*"

	this.pattern = regexp.MustCompile(splitPattern)
}

/*-------------------------------------------------------------
                            请求回复

--------------------------------------------------------------*/
type Request struct {
	Sentences []Sentence
	Original  string
}

func (this *Request) Empty() bool {
	return this.Sentences == nil || this.Original == ""
}

func (this *Request) Last(index int) Sentence {
	n := len(this.Sentences)
	return this.Sentences[n-(1+index)]
}

func (this *Request) TrimOriginal() string {
	return strings.TrimSpace(this.Original)

}

type Response struct {
	Request
}

func (this *Response) Append(output string) {
	if !strings.HasSuffix(this.Original, " ") {
		this.Original += " "
	}

	this.Original += output

}

/*-------------------------------------------------------
                     Tokenizer
--------------------------------------------------------*/
type Tokenizer struct {
	ignoreWhitespace bool
	splitters        []string
	pattern          *regexp.Regexp
}

func (this *Tokenizer) Init(splitters ...string) {

	this.afterSetProperty()

}

func (this *Tokenizer) InitByConfig(splitters ...string) {

	this.afterSetProperty()

}
func (this *Tokenizer) IsIgnoreWhitespace() bool {
	return this.ignoreWhitespace
}

func (this *Tokenizer) GetSplitters() []string {
	return this.splitters
}

func (this *Tokenizer) Tokenize(input string) (output []string) {
	output = this.pattern.FindAllString(input, -1)
	return

}

func (this *Tokenizer) afterSetProperty() {
	if this.splitters == nil {
		return
	}

	for index, value := range this.splitters {
		this.splitters[index] = EscapeRegex(value)
	}

	expression := strings.Join(this.splitters, "|")

	if this.ignoreWhitespace {
		expression = "(" + expression + ")\\s*|\\s+"
	} else {
		expression = "(" + expression + "|\\s+)"
	}
	var err error
	this.pattern, err = regexp.Compile(expression)
	if err != nil {
		fmt.Println("tokenizer init failed: " + err.Error())
		this.pattern = nil
	}

}
