package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"
)

const (
	outputDirectory = "argocd-doctor-output"
)

type Rule struct {
	ID             string
	Name           string
	Description    string
	FilesTriggered []string
}

// Final template that contains all information. Rendered with web/index.html.tpl
type templateData struct {
	Now           time.Time //when the page was generated
	RulesFailed   []Rule
	RulesPassed   []Rule
	RulesInactive []Rule
}

func main() {

	// Discover GitHub releases from all related repositories
	// versionsFound := discoverVersions()
	// log.Printf("Found %d versions of the GitOps Runtime\n", len(versionsFound))

	log.Println("Creating output directory at ./" + outputDirectory)
	err := os.MkdirAll(outputDirectory, 0755)
	if err != nil {
		log.Fatal(err)
	}

	copy("report-template/html-mockup/style.css", outputDirectory+"/style.css")
	copy("report-template/html-mockup/favicon.png", outputDirectory+"/favicon.png")

	log.Println("Rendering HTML report at ./" + outputDirectory + "/index.html")

	tmpl, err := template.ParseFiles("report-template/html-mockup/index.html")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(outputDirectory + "/index.html")
	if err != nil {
		log.Fatal("create file: ", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			log.Fatal("close file: ", cerr)
		}
	}()

	tData := templateData{
		Now: time.Now(),
		// VersionsFound: versionsFound,
	}

	err = tmpl.Execute(f, tData)

	if err != nil {
		log.Fatal("Could not render HTML template: ", err)
	}

}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string, dst string) {
	// Open the source file
	srcFile, err := os.Open(src)
	checkErr(err)
	defer func() {
		if cerr := srcFile.Close(); cerr != nil {
			log.Fatal("close src file: ", cerr)
		}
	}()

	// Create the destination file
	dstFile, err := os.Create(dst)
	checkErr(err)
	defer func() {
		if cerr := dstFile.Close(); cerr != nil {
			log.Fatal("close dst file: ", cerr)
		}
	}()

	// Copy the contents from src to dst
	_, err = io.Copy(dstFile, srcFile)
	checkErr(err)
}
