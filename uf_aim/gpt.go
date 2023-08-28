package uf_aim

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/tailscale/hujson"
	"gitlab.com/textfridayy/uno/uf"
)

type GptRequest struct {
	ApiKey string
}

var MakoModel string = "davinci:ft-text-friday-2022-10-02-02-49-17"

type MakoResponse struct {
	Original             string
	OriginalCleaned      string
	OriginalCleanedBytes []byte

	Function     AimTreeName `json:"function"`
	Noun         string      `json:"noun"`
	Adjectives   []string    `json:"adjectives"`
	Category     string      `json:"category"`
	ClothingType string      `json:"clothing_type"`
	Quantity     int         `json:"quantity"`
	ArrivalDate  string      `json:"arrival_date"`
	Store        string      `json:"store"`
	MinPrice     float64     `json:"min_price"`
	MaxPrice     float64     `json:"max_price"`
	Multiple     bool        `json:"multiple"`

	//Function     AimTreeName
	//Noun         string
	//Adjectives   []string
	//Category     string
	//ClothingType string
	//Quantity     int
	//ArrivalDate  string
	//Store        string
	//MinPrice     float64
	//MaxPrice     float64
	//Multiple     bool
}

func (mr *MakoResponse) AsQueryString() string {
	return strings.Join(mr.Adjectives[:], ", ") + " " + mr.Noun
}

func (mr *MakoResponse) AdjectivesSorted() []string {
	adjectivesSorted := mr.Adjectives
	sort.Strings(adjectivesSorted)
	return adjectivesSorted
}

func InterpretAim(
	gibs *uf.Gibs,
	content string,
) (*MakoResponse, error) {
	client := gogpt.NewClient(gibs.Conf.OpenAiKey)

	//prepare string - uppercase first letter
	firstRune, size := utf8.DecodeRuneInString(content)
	content = string(unicode.ToUpper(firstRune)) + content[size:]

	// add punctuation
	lastChar := content[len(content)-1]
	if !(lastChar == '.' || lastChar == '?') {
		content = content + "."
	}
	prompt := "Respond in JSON. " + content + " ->"
	uf.Trace(prompt)

	req := gogpt.CompletionRequest{
		Model:       "davinci:ft-text-friday-2022-10-02-03-58-57",
		MaxTokens:   100, // 16 is default
		Prompt:      prompt,
		Temperature: 0,
		//TopP:             1,
		//FrequencyPenalty: 0,
		//PresencePenalty:  0,
		Stop: []string{" ->", ">>>"},
	}
	uf.Trace("GPT	5")

	makoResponse := MakoResponse{}

	resp, err := client.CreateCompletion(gibs.Ctx, req)
	uf.Trace("GPT	6")
	if len(resp.Choices) > 0 {
		text := resp.Choices[0].Text
		makoResponse.Original = text

		// convert all double quotes to single quotes for consistency
		text = strings.ReplaceAll(text, "\"", "'")

		// remove common error from GPT. order matters here
		text = strings.ReplaceAll(text, "{\"error\":[],\"success\":{}}", "")
		text = strings.ReplaceAll(text, "{\"error\":[],\"success\":{}", "")

		// fix JSON formatting errors
		text = strings.ReplaceAll(text, "False", "false")
		text = strings.ReplaceAll(text, "True", "true")
		text = strings.ReplaceAll(text, "{ ", "{")
		text = strings.ReplaceAll(text, " }", "}")

		// find and fix end of real JSON data
		text = strings.Split(text, "}")[0]
		if text[len(text)-1] != '}' {
			text = text + "}"
		}
		uf.Trace("GPT	6.a")

		// find and fix fix start of real JSON data
		if len(text) > 1 {
			getStart := strings.Split(text, "{'Function'")
			if len(getStart) > 1 {
				text = getStart[1]
			}

			text = "{'Function'" + text
		} else {
			//do nothing
		}

		// convert all double quotes to single quotes for JSON serialization
		text = strings.ReplaceAll(text, "'", "\"")

		makoResponse.OriginalCleaned = text
		uf.Trace("GPT	6.b")

	} else {
		makoResponse.Original = ""
		makoResponse.OriginalCleaned = "{}"
	}

	uf.Trace("GPT 7")
	uf.Trace("================================ response from GPT before cleaning")
	uf.Trace(makoResponse.Original)
	uf.Trace("================================ after cleaning")
	uf.Trace(makoResponse.OriginalCleaned)

	uf.Trace("GPT 8")
	makoResponse.OriginalCleanedBytes, err = hujson.Standardize([]byte(
		makoResponse.OriginalCleaned))
	if err != nil {
		uf.Error(err)
	}

	uf.Trace("GPT 9")
	if err := json.Unmarshal(
		makoResponse.OriginalCleanedBytes,
		&makoResponse); err != nil {
		uf.Error(err)
		makoResponse.Function = AimNone
		uf.Trace("GPT 10")
		return &makoResponse, err
	}

	uf.Trace("GPT 11")
	if !makoResponse.Function.IsTreeNameValid() {
		uf.Trace("GPT 11.a")
		err = errors.New("Tree name is not valid: " + string(makoResponse.Function))
		uf.Trace("GPT 11.b")
		makoResponse.Function = AimNone
	}

	uf.Trace(makoResponse.Function)
	uf.Trace(makoResponse.Noun)
	uf.Trace("================================ done")

	return &makoResponse, err
}

//func AnswerQuestion(
//ctx *fasthttp.RequestCtx,
//pile uf.Pile,
//in AimIn,
//) (resp gogpt.CompletionResponse, err error) {
//client := gogpt.NewClient(pile.Conf.OpenAiKey)
//req := gogpt.CompletionRequest{
//Model:            "curie",
//MaxTokens:        20,
//Prompt:           in.Content,
//Temperature:      0,
//TopP:             1,
//FrequencyPenalty: 0,
//PresencePenalty:  0,
//}

//resp, err = client.CreateCompletion(ctx, req)
//return resp, err
//}
