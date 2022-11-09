package surveys

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	"errors"
)

var errorInvalidFormat = errors.New("Invalid survey format")

type SurveyJSON struct {
	Props struct {
		Name        []LocalisedObject `json:"name"`
		Description []LocalisedObject `json:"description"`
	} `json:"props"`
	SurveyDefinition struct {
		Key string `json:"key"`
	} `json:"surveyDefinition"`
	Metadata map[string]string `json:"metadata",omitempty`
}

type LocalisedObject struct {
	Code string `json:"code"`
	// For texts
	Parts []ExpressionArg `json:"parts"`
}

type ExpressionArg struct {
	Str string `json:"str,omitempty"`
}

func LoadSurvey(file string) (*SurveyJSON, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	p := SurveyJSON{}
	err = json.Unmarshal(data, &p)
	if err != nil {
		//fmt.Println("error:", err)
		return nil, err
	}
	if(p.SurveyDefinition.Key == "") {
		return nil, errorInvalidFormat
	}
	return &p, nil
}

func LocalisedToMap(o []LocalisedObject) map[string]string {
	desc := make(map[string]string, 0)
	for _, loc := range o {
		code := loc.Code
		pp := make([]string, 0)
		for _, p := range loc.Parts {
			if p.Str != "" {
				pp = append(pp, p.Str)
			}
		}
		if len(pp) > 0 {
			desc[code] = strings.Join(pp, " ")
		}
	}
	return desc
}
