package reader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/sebastianxyzsss/tardigrade/globals"
	"github.com/sebastianxyzsss/tardigrade/logger"
	"gopkg.in/yaml.v2"
)

var ll = logger.SetupLog()

var TardiContentDir string = ".tardigrade"

var TardiContent string = TardiContentDir + "/tardicontent.yml"

var TardiHistory string = TardiContentDir + "/tardihistory.yml"

func Unmarshall(mapStr string) *map[interface{}]interface{} {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(mapStr), &m)
	if err != nil {
		ll.Error().Msg("error:" + err.Error())
	}
	return &m
}

func Marshall(mm *map[interface{}]interface{}) (*string, error) {
	y, err := yaml.Marshal(mm)
	if err != nil {
		return nil, err
	}
	str := string(y)
	return &str, nil
}

func GetFileAsString(filePath string) *string {

	yamlContent, err := os.ReadFile(filePath)
	if err != nil {
		ll.Debug().Msg("unable to read file: " + err.Error())
		return nil
	}

	yamlContentStr := string(yamlContent)

	return &yamlContentStr
}

func GetHttpRequestAsString(urlString string) *string {

	url := strings.Replace(urlString, "/https://", "https://", -1)
	url = strings.Replace(url, "/http://", "http://", -1)

	ll.Debug().Msg("url with content: " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ll.Debug().Msg("unable to read url: " + err.Error())
		return nil
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ll.Debug().Msg("unable to get url response: " + err.Error())
		return nil
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		ll.Debug().Msg("unable to parse response: " + err.Error())
	}

	yamlContentStr := string(b)
	return &yamlContentStr
}

func getLocalContent() *map[interface{}]interface{} {
	yamlContent := GetFileAsString(TardiContent)
	if yamlContent == nil {
		ll.Debug().Msg("unable to get local content")
		return nil
	} else {
		ll.Debug().Msg("getting local content")
	}
	yamlAsMap := Unmarshall(*yamlContent)
	return yamlAsMap
}

func getUserHomeFileContent(fileName string) *map[interface{}]interface{} {
	userDirName, err := os.UserHomeDir()
	if err != nil {
		ll.Debug().Msg("unable to get user home dir: " + err.Error())
		return nil
	}
	ll.Debug().Msg("user home dir name: " + userDirName)

	yamlContent := GetFileAsString(userDirName + "/" + fileName)
	if yamlContent == nil {
		ll.Debug().Msg("unable to get content " + fileName)
		return nil
	} else {
		ll.Debug().Msg("getting content " + fileName)
	}
	yamlAsMap := Unmarshall(*yamlContent)
	return yamlAsMap
}

func getUserHomeContent() *map[interface{}]interface{} {
	return getUserHomeFileContent(TardiContent)
}

func GetHistoryPath() (*string, error) {
	userDirName, err := os.UserHomeDir()
	if err != nil {
		ll.Debug().Msg("unable to get user home dir: " + err.Error())
		return nil, err
	}
	userTardiHistory := userDirName + "/" + TardiHistory
	return &userTardiHistory, nil
}

func getHistoryContent() *map[interface{}]interface{} {
	userTardiHistory, err := GetHistoryPath()
	if err != nil {
		ll.Debug().Msg("unable to get history path: " + err.Error())
		return nil
	}
	ll.Debug().Msg("userTardiHistory: " + *userTardiHistory)
	if DoesFileNotExist(*userTardiHistory) {
		ll.Debug().Msg("history does not exist, so creating")
		WriteToFile(*userTardiHistory, firstInHistory())
	}
	return getUserHomeFileContent(TardiHistory)
}

func appendFileListContent(m *map[interface{}]interface{}, filesToRead []string) *map[interface{}]interface{} {
	if filesToRead == nil {
		return m
	}

	for _, fileToRead := range filesToRead {

		ll.Debug().Msg("reading file (or url): " + fileToRead)

		var yamlContent *string = nil

		if strings.Contains(fileToRead, "/http") {
			yamlContent = GetHttpRequestAsString(fileToRead)
		} else {
			yamlContent = GetFileAsString(fileToRead)
		}

		if yamlContent == nil {
			ll.Info().Msg("unable to get file content for: " + fileToRead)
			continue
		} else {
			ll.Debug().Msg("found content for: " + fileToRead)
		}
		yamlAsMap := Unmarshall(*yamlContent)
		m = appendMap(m, yamlAsMap)
	}
	return m
}

func firstInHistory() string {
	var historyData = `history:
  - pwd
`
	return historyData
}

func getDummyStr(prefix string) string {
	var dummyData = `group10:
  - ls -la
  - echo hi
`
	return prefix + dummyData
}

func getDummyContent() *map[interface{}]interface{} {
	var dummyData = getDummyStr("")
	yamlAsMap := Unmarshall(dummyData)
	return yamlAsMap
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func getUserDirName() string {
	userDirName, err := os.UserHomeDir()
	if err != nil {
		ll.Error().Msg("unable to get user home dir: " + err.Error())
		return ""
	}
	return userDirName
}

func WriteToFile(fileName string, content string) {

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	ll.Debug().Msg("bytes written successfully to " + fileName)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateNewUserContentFile() {
	userDirName := getUserDirName()

	if CheckContentFileExists(userDirName) {
		ll.Debug().Msg("> userhome content file - already Exists")
	} else {
		ll.Debug().Msg("> userhome content file - Does Not exist, creating new one")

		tardiContentDirName := userDirName + "/" + TardiContentDir
		tardiContentFileName := userDirName + "/" + TardiContent
		dummyGroupPrefix := "userhome "

		CreateNewContentFile(tardiContentDirName, tardiContentFileName, dummyGroupPrefix)
	}
}

func CreateNewLocalContentFile() {
	ll.Info().Msg("attempting creating local file")

	dummyGroupPrefix := "local-"

	CreateNewContentFile(TardiContentDir, TardiContent, dummyGroupPrefix)
}

func CheckContentFileExists(dir string) bool {
	tardiContentDirName := dir + "/" + TardiContentDir
	tardiContentFileName := dir + "/" + TardiContent

	if _, err := os.Stat(tardiContentDirName); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(tardiContentFileName); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateNewContentFile(tardiContentDirName string, tardiContentFileName string, dummyGroupPrefix string) {

	if _, err := os.Stat(tardiContentDirName); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(tardiContentDirName, os.ModePerm)
		if err != nil {
			ll.Error().Msg("something went wrong with making the path: " + err.Error())
			return
		}
	}

	if checkFileExists(tardiContentFileName) {
		ll.Warn().Msg(tardiContentFileName + " file already exists")
		return
	}

	newfile, err := os.Create(tardiContentFileName)
	defer newfile.Close()
	if err != nil {
		ll.Error().Msg("unable to create file: " + err.Error())
		return
	}

	var dummyData = getDummyStr(dummyGroupPrefix)

	_, err2 := newfile.WriteString(dummyData)

	if err2 != nil {
		ll.Error().Msg("unable to write to file: " + err.Error())
		return
	}

	ll.Debug().Msg("created file")
}

func appendMap(m1 *map[interface{}]interface{}, m2 *map[interface{}]interface{}) *map[interface{}]interface{} {
	if m1 == nil && m2 == nil {
		return nil
	}
	if m1 == nil {
		return m2
	}
	if m2 != nil {
		for k, v := range *m2 {
			(*m1)[k] = v
		}
	}
	return m1
}

func DumpMapAsYaml(m *map[interface{}]interface{}) {
	d, err := yaml.Marshal(m)
	if err != nil {
		ll.Error().Msg("error:" + err.Error())
	}
	fmt.Printf("= m dump:\n%s\n\n", string(d))
}

func DoesFileNotExist(filePath string) bool {
	_, error := os.Stat(filePath)
	return os.IsNotExist(error)
}

// main entry point
func GetRawMapContent(filesToRead []string) *map[interface{}]interface{} {

	yamlAsMap := getUserHomeContent()

	localYamlAsMap := getLocalContent()

	historyAsMap := getHistoryContent()

	yamlAsMap = appendMap(yamlAsMap, localYamlAsMap)

	yamlAsMap = appendMap(yamlAsMap, historyAsMap)

	if globals.FilterAction == globals.FilterFiles || globals.FilterAction == globals.FilterNone {
		yamlAsMap = appendFileListContent(yamlAsMap, filesToRead)
	}

	if yamlAsMap == nil {
		ll.Warn().Msg("unable to get any content, getting dummy content")
		yamlAsMap = getDummyContent()
	}

	return yamlAsMap
}
