package source

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/kbence/logan/config"
)

func init() {
	logSourceFactories["generic"] = func(cfg *config.Configuration) LogSource {
		return &GenericLogSource{
			directories: cfg.Generic.Dirs,
			maxDepth:    cfg.Generic.MaxDepth}
	}
}

type genericLogCategory struct {
	Name  string
	Files []string
}

type genericLogMap map[string]genericLogCategory

type GenericLogSource struct {
	directories []string
	maxDepth    int
	categories  genericLogMap
}

var reGzippedLog = regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?\\.[0-9]+\\.gz$")
var reNumberedLog = regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?\\.[0-9]+$")
var reCurrentLog = regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?$")

var genericFileNamePatterns = []*regexp.Regexp{
	regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?\\.[0-9]+\\.gz$"),
	regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?\\.[0-9]+$"),
	regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?$")}

func (s *GenericLogSource) collectCategories(dir, prefix string, depth int) genericLogMap {
	categories := genericLogMap{}

	if depth >= 0 {
		fileInfos, err := ioutil.ReadDir(dir)

		if err == nil {
			categoryMap := map[string][]os.FileInfo{}

			for _, fileInfo := range fileInfos {
				fileName := fileInfo.Name()

				if fileInfo.IsDir() {
					subDir := path.Join(dir, fileName)
					prefix := fmt.Sprintf("%s%s/", prefix, fileName)

					for _, subCategory := range s.collectCategories(subDir, prefix, depth-1) {
						categories[subCategory.Name] = subCategory
					}
				} else {
					for _, pattern := range genericFileNamePatterns {
						if match := pattern.FindStringSubmatch(fileName); match != nil {
							categoryMap[match[1]] = append(categoryMap[match[1]], fileInfo)
						}
					}
				}
			}

			for name, files := range categoryMap {
				category := genericLogCategory{Name: name}

				for _, file := range files {
					category.Files = append(category.Files, path.Join(dir, file.Name()))
				}

				categories[name] = category
			}
		}
	}

	return categories
}

func (s *GenericLogSource) loadCategories() {
	if s.categories != nil {
		return
	}

	s.categories = genericLogMap{}

	for _, dir := range s.directories {
		categories := s.collectCategories(dir, "", s.maxDepth)

		for name, category := range categories {
			s.categories[name] = category
		}
	}
}

func (s *GenericLogSource) GetCategories() []string {
	s.loadCategories()
	categoryNames := []string{}

	for name, _ := range s.categories {
		categoryNames = append(categoryNames, name)
	}

	return categoryNames
}

func (s *GenericLogSource) ContainsCategory(category string) bool {
	return false
}

func (s *GenericLogSource) GetChain(category string) LogChain {
	return nil
}
