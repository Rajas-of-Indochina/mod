package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func writeEffects(modDir string) {
	fmt.Println("Reading heritage groups...")
	heritageGroups := getKeys(filepath.Join(modDir, "common", "culture", "pillars"), "heritage_group")
	fmt.Println("Reading heritage families...")
	heritageFamilies := getKeys(filepath.Join(modDir, "common", "culture", "pillars"), "heritage_family")
	fmt.Println("Reading language groups...")
	languageGroups := getKeys(filepath.Join(modDir, "common", "culture", "pillars"), "language_group")
	fmt.Println("Reading language families...")
	languageFamilies := getKeys(filepath.Join(modDir, "common", "culture", "pillars"), "language_family")

	fmt.Println("Creating scripted effect file...")
	outFile, err := os.Create(filepath.Join(modDir, "common", "scripted_effects", "ccu_scripted_effects.txt"))
	if err != nil {
		log.Fatal(err)
	}
	writeHeader(outFile)
	_, err = outFile.WriteString("ccu_initialize_culture = {\n")
	_, err = outFile.WriteString("\tccu_initialize_heritage_group = yes\n")
	_, err = outFile.WriteString("\tccu_initialize_heritage_family = yes\n")
	_, err = outFile.WriteString("\tccu_initialize_language_group = yes\n")
	_, err = outFile.WriteString("\tccu_initialize_language_family = yes\n")
	_, err = outFile.WriteString("}\n\n")
	writeEffect(heritageGroups, "heritage_group", outFile)
	_, err = outFile.WriteString("\n\n")
	writeEffect(heritageFamilies, "heritage_family", outFile)
	_, err = outFile.WriteString("\n\n")
	writeEffect(languageGroups, "language_group", outFile)
	_, err = outFile.WriteString("\n\n")
	writeEffect(languageFamilies, "language_family", outFile)

	fmt.Println("Creating localization files...")
	writeLocalization(modDir, heritageGroups, "heritage_group", "-")
	writeLocalization(modDir, heritageFamilies, "heritage_family", " ")
	writeLocalization(modDir, languageGroups, "language_group", "-")
	writeLocalization(modDir, languageFamilies, "language_family", "-")

}

func writeLocalization(modDir string, keys []string, varName string, delimiter string) {
	outFile, err := os.Create(filepath.Join(modDir, "localization", "english", "ccu_"+varName+"_l_english.yml"))
	if err != nil {
		log.Fatal(err)
	}
	writeLocHeader(outFile)
	for _, key := range keys {

		fields := strings.Split(key, "_")[2:]
		newKey := ""
		for _, field := range fields {
			newKey += field + delimiter
		}
		newKey = newKey[0 : len(newKey)-1]
		newKey = strings.Title(newKey)
		gameConcept := "[" + varName + "|E]"
		_, err = outFile.WriteString("culture_parameter_" + key + ":0 \"#P +[EmptyScope.ScriptValue('same_" + varName + "_cultural_acceptance')|0]#! [cultural_acceptance_baseline|E] with Cultures sharing the " + newKey + " " + gameConcept + "\"\n")
	}
}

func writeEffect(keys []string, varName string, outFile *os.File) {
	_, _ = outFile.WriteString("ccu_initialize_" + varName + " = {\n")
	for i, key := range keys {
		if i == 0 {
			_, _ = outFile.WriteString("\tif = { limit = { has_cultural_parameter = " + key + " } set_variable = { name = " + varName + " value = " + strconv.Itoa(i) + " } }\n")
		} else {
			_, _ = outFile.WriteString("\telse_if = { limit = { has_cultural_parameter = " + key + " } set_variable = { name = " + varName + " value = " + strconv.Itoa(i) + " } }\n")
		}
	}
	_, _ = outFile.WriteString("}\n\n")
}

func getKeys(inDir string, searchString string) []string {

	keys := make(map[string]int)

	fileInfo, err := ioutil.ReadDir(inDir)
	if err != nil {
		fmt.Println("failed to read directory: " + inDir)
		log.Fatal(err)
	}
	for _, file := range fileInfo {
		filePath := filepath.Join(inDir, file.Name())
		thisFile, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Failed to open file: " + filePath)
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(thisFile)
		for scanner.Scan() {
			line := cleanLine(scanner.Text())

			fields := strings.Fields(line)
			for _, field := range fields {
				if strings.Contains(field, searchString) {
					if _, ok := keys[field]; ok {
						keys[field]++
					} else {
						keys[field] = 1
					}
				}
			}
		}
		_ = thisFile.Close()
	}

	keySlice := make([]string, 0)
	for key := range keys {
		keySlice = append(keySlice, key)
	}
	keySlice = qsort(keySlice)
	return keySlice
}

// removes comment blocks from a line (string)
func cleanLine(line string) string {
	// Look out for following byte signifying a comment
	commentByte := []byte("#")[0]
	// iterate through all bytes in line
	for i := 0; i < len(line); i++ {
		// if byte = comment byte, return line up until that byte
		if line[i] == commentByte {
			if i > 1 {
				return line[0:i]
			} else {
				return ""
			}
		}
	}
	// if no comments found, return whole line
	return line
}
