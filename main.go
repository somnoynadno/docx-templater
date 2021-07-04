package main

import (
	"encoding/json"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var licenseKey = os.Getenv("LICENSE_KEY")
var templatesFolder = os.Getenv("TEMPLATES_FOLDER")
var saveFolder = os.Getenv("SAVE_FOLDER")


func init() {
	err := license.SetMeteredKey(licenseKey)
	if err != nil {
		panic(err)
	}
}

type RequestBody struct {
	Filename  string
	Templates map[string]string
}

var CreateDocument = func(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

	if r.Method == "OPTIONS" {
		return
	}

	body := &RequestBody{}
	err := json.NewDecoder(r.Body).Decode(body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	doc, err := document.Open(templatesFolder + body.Filename)
	if err != nil {
		log.Printf("error opening document: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer doc.Close()

	paragraphs := []document.Paragraph{}
	for _, p := range doc.Paragraphs() {
		paragraphs = append(paragraphs, p)
	}

	for _, sdt := range doc.StructuredDocumentTags() {
		for _, p := range sdt.Paragraphs() {
			paragraphs = append(paragraphs, p)
		}
	}

	for _, p := range paragraphs {
		for _, r := range p.Runs() {
			text := r.Text()

			newText := text
			for key, value := range body.Templates {
				newText = strings.Replace(newText, key, value, -1)
			}

			if text != newText {
				r.ClearContent()
				r.AddText(newText)
			}
		}
	}

	randomName := strconv.Itoa(rand.Intn(10000000)) + ".docx"
	err = doc.SaveToFile(saveFolder + randomName)
	if err != nil {
		log.Printf("error writing document: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(map[string]string{"filename": randomName})
	_, _ = w.Write(res)
}

func main() {
	http.Handle("/docgen/static/", http.StripPrefix("/docgen/static/", http.FileServer(http.Dir(saveFolder))))
	http.HandleFunc("/docgen/create", CreateDocument)
	http.ListenAndServe(":4040", nil)
}
