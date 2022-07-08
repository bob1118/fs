package autoload

import "log"

func GetConfiguration(filename string) (string, error) {
	var content string
	var err error

	if content, err = readConfFromDatabase(filename); err != nil {
		if content, err = readConfFromFile(filename); err != nil {
			if content, err = buildConfiguration(filename); err != nil {
				log.Println(err)
			} else {
				writeConfToFile(filename, content) //nothing todo.
			}
		}
		content = updateConf(content)
		writeConfToDatabase(filename, content)
	}
	return content, err
}

func readConfFromDatabase(filename string) (string, error) { return "", nil }
func readConfFromFile(filename string) (string, error)     { return "", nil }
func buildConfiguration(filename string) (string, error)   { return "", nil }
func writeConfToFile(filename, content string)             {}
func writeConfToDatabase(filename string, content string)  {}
func updateConf(old string) string                         { return "" }
