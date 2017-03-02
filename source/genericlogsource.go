package source

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"

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

var genericFileNamePatterns = []*regexp.Regexp{
	regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?\\.[0-9]+\\.gz$"),
	regexp.MustCompile("^([a-zA-Z0-9_-]+)(\\.log)?\\.[0-9]+\\.bz2$"),
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
					newPrefix := fmt.Sprintf("%s%s/", prefix, fileName)

					for _, subCategory := range s.collectCategories(subDir, newPrefix, depth-1) {
						categories[fmt.Sprintf("%s%s", prefix, subCategory.Name)] = subCategory
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
				fullName := fmt.Sprintf("%s%s", prefix, name)
				category := genericLogCategory{Name: fullName}

				for _, file := range files {
					category.Files = append(category.Files, path.Join(dir, file.Name()))
				}

				sort.Sort(sort.Reverse(sort.StringSlice(category.Files)))

				categories[fullName] = category
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
	s.loadCategories()
	_, found := s.categories[category]

	return found
}

func (s *GenericLogSource) GetChain(category string) LogChain {
	if !s.ContainsCategory(category) {
		return nil
	}

	return NewGenericLogChain(s.categories[category].Files)
}
