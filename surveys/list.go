package surveys

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// SurveyFile contains info about an available survey
type SurveyFile struct {
	Id          string            `json:"id"`
	Label       string            `json:"label"`
	Description map[string]string `json:"description"`
	Time        time.Time         `json:"time"`
	file        string
	Study       string `json:"study"`
	Metadata    map[string]string `json:"metadata",omitempty`
}

type Symlink struct {
	path string
	info os.FileInfo
}

func readSymlink(path string) (Symlink, error) {
	s := Symlink{path: "", info: nil}

	p, e := filepath.EvalSymlinks(path)
	if e != nil {
		fmt.Printf("Error reading %s : %s", path, e)
		return s, e
	}
	i, e := os.Stat(p)
	if e != nil {
		fmt.Printf("Error reading %s : %s", path, e)
		return s, e
	}
	s.path = p
	s.info = i
	return s, nil
}

func detectSurvey(file string) (SurveyFile, bool) {
	survey := SurveyFile{file: file}
	def, err := LoadSurvey(file)
	if err != nil {
		fmt.Printf("Not a survey %s : %s\n", file, err)
		return survey, false
	}
	// Default
	study := filepath.Base(filepath.Dir(filepath.Dir(file)))
	survey.Study = study
	survey.Label = def.SurveyDefinition.Key
	survey.Description = LocalisedToMap(def.Props.Name)
	if(len(def.Metadata) > 0) {
		survey.Metadata = def.Metadata
	}
	return survey, true
}

func readDirectory(root string) (map[string]SurveyFile, error) {
	surveys := make(map[string]SurveyFile, 10)
	basePath := filepath.Clean(root)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		name := info.Name()

		if info.Mode()&fs.ModeSymlink != 0 {
			i, e := readSymlink(path)
			if e != nil {
				return nil // Skip entry
			}
			if i.info.IsDir() {
				ss, e := readDirectory(i.path)
				if e != nil {
					fmt.Printf("Error reading %s : %s", path, e)
				}
				prefix := name + "/"
				for n, s := range ss {
					s.Id = prefix + s.Id
					surveys[prefix+n] = s
				}
			}
			return nil // Stop here (dont handle symlink file)
		}

		include := filepath.Ext(name) == ".json"

		if !include {
			return nil
		}

		survey, ok := detectSurvey(path)
		if !ok {
			return nil
		}

		rel, _ := filepath.Rel(basePath, path)

		survey.Id = rel
		survey.Time = info.ModTime()
		surveys[rel] = survey
		fmt.Printf("Added name: %s (%s)\n", path, rel)
		return nil
	})

	return surveys, err
}

// SurveyList manages list of surveys
type SurveyList struct {
	root    string
	surveys map[string]SurveyFile
}

// GetList returns the list of available surveys as an array
func (l *SurveyList) GetList() []SurveyFile {
	n := len(l.surveys)
	ss := make([]SurveyFile, 0, n)

	keys := make([]string, 0, n)
	for k := range l.surveys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		survey := l.surveys[key]
		ss = append(ss, survey)
	}
	return ss
}

// Get provides the survey json file location from its id
func (l *SurveyList) Get(id string) (string, error) {
	s, ok := l.surveys[id]
	if ok {
		return s.file, nil
	}
	return "", nil
}

// Update the survey list (warning not thread safe now)
func (l *SurveyList) Update() error {
	list, err := readDirectory(l.root)
	if err != nil {
		return err
	}
	l.surveys = list
	return nil
}

// NewSurveyList returns an instance of SurveyList manager
func NewSurveyList(root string) *SurveyList {
	return &SurveyList{root: root}
}
