package surveys

import (
	"fmt"
	"log"
	"os"
	"io/fs"
	"path/filepath"
)

// SurveyFile contains info about an available survey
type SurveyFile struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	file        string
}

type Symlink struct {
	path string
	info os.FileInfo
}

func readSymlink(path string) (Symlink, error) {
	s := Symlink{path: "", info: nil}
	
	p, e := filepath.EvalSymlinks(path)
	if(e != nil) {
		fmt.Printf("Error reading %s : %s", path, e)
		return s, e
	}
	i, e := os.Stat(p)
	if(e != nil) {
		fmt.Printf("Error reading %s : %s", path, e)
		return s, e
	}
	s.path = p
	s.info = i
	return s, nil
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
		
		if(info.Mode() & fs.ModeSymlink != 0) {
			i, e := readSymlink(path)
			if(e != nil) {
				return nil // Skip entry
			}
			if(i.info.IsDir()) {
				ss, e := readDirectory(i.path)
				if(e != nil) {
					fmt.Printf("Error reading %s : %s", path, e)
				}
				prefix := name + "/"
				for n, s := range ss {
					s.Id = prefix + s.Id
					s.Description = prefix + s.Description
					surveys[ prefix + n ] = s
				}
			}
			return nil // Stop here (dont handle symlink file)
		}

		include := filepath.Ext(name) == ".json"

		if !include {
			return nil
		}

		rel, _ := filepath.Rel(basePath, path)

		survey := SurveyFile{
			Id:          rel,
			Label:       name,
			Description: rel,
			file:        path,
		}

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
	ss := make([]SurveyFile, 0, len(l.surveys))
	for _, survey := range l.surveys {
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
