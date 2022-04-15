/*
Copyright © 2022 natakorn

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CreateOption struct {
	Name     string
	Path     string
	Type     string
	Override bool
}
type Resource struct {
	Template *template.Template
	Path     string
	Name     string
	FileName string
	FullPath string
}

var options CreateOption

// var resourceOptions = []string{"resource", "model", "form", "response", "route", "handler", "service"}
var resourceOptions = []string{"resource", "model", "form", "response"}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource file",
	Long:  `Todo`,
	Run: func(cmd *cobra.Command, args []string) {
		t := cmd.Flag("type").Value.String()
		n := cmd.Flag("name").Value.String()
		p := cmd.Flag("path").Value.String()
		options.Type = t
		options.Name = n
		options.Path = p
		options.Override = true
		if t == "" {
			options.Type = askForType()
		}
		if n == "" {
			options.Name = askForName()
		}
		if p == "" && options.Type != "resource" {
			options.Path = askForPath()
		}
		switch options.Type {
		case "resource":
			createResource()
		default:
			createFile()
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringP("name", "n", "", "name of the project")
	createCmd.Flags().StringP("type", "t", "", "resource type")
	createCmd.Flags().StringP("path", "p", "", "resource path")
}

func askForName() string {
	var r string
	prompt := &survey.Input{
		Message: "name of the resource",
	}
	if err := survey.AskOne(prompt, &r, survey.WithValidator(survey.Required)); err != nil {
		os.Exit(0)
	}
	return r
}

func askForType() string {
	var r string
	prompt := &survey.Select{
		Message: "what resource do you want to create",
		Options: resourceOptions,
	}
	if err := survey.AskOne(prompt, &r, survey.WithValidator(survey.Required)); err != nil {
		os.Exit(0)
	}
	return r
}

func askForPath() string {
	var r string
	prompt := &survey.Select{
		Message: "what path do you want to create",
		Options: getFolderList(),
		Default: options.Type,
	}
	if err := survey.AskOne(prompt, &r, survey.WithValidator(survey.Required)); err != nil {
		os.Exit(0)
	}
	return r
}

func getFolderList() []string {
	var names []string
	names = append(names, ".")
	cDir, err := os.Getwd()
	if err != nil {
		fmt.Println("⛔ Oops, sorry for the error but we cannot get Current Directory")
		log.Fatal(err)
	}
	files, err := os.ReadDir(cDir)
	if err != nil {
		fmt.Println("⛔ Oops, sorry for the error but we cannot Read your Current Directory")
		log.Fatal(err)
	}
	for _, f := range files {
		file, _ := os.Stat(f.Name())
		if file.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			names = append(names, f.Name())
		}
	}
	//TODO: Create new folder features
	// names = append(names, "📁 new folder")
	return names
}

func createResource() {
	dir := getCurrentDirectory()
	data := struct {
		Name string
	}{
		Name: strcase.ToCamel(options.Name),
	}
	m, _ := template.ParseFS(TemplateFS, "templates/resources/model.tmpl")
	f, _ := template.ParseFS(TemplateFS, "templates/resources/form.tmpl")
	r, _ := template.ParseFS(TemplateFS, "templates/resources/response.tmpl")
	createDirectoryIfNotExists(dir + "/model")
	createDirectoryIfNotExists(dir + "/form")
	createDirectoryIfNotExists(dir + "/resp")
	mFileName := dir + "/model/" + strcase.ToSnake(options.Name) + "_model.go"
	if _, err := os.Stat(mFileName); err == nil {
		p := &survey.Confirm{
			Message: "file already exists, do you want to overwrite it?",
			Default: false,
		}
		survey.AskOne(p, &options.Override)
	}
	fFileName := dir + "/form/" + strcase.ToSnake(options.Name) + "_form.go"
	rFileName := dir + "/resp/" + strcase.ToSnake(options.Name) + "_resp.go"
	if options.Override {
		mFile, _ := os.Create(mFileName)
		fFile, _ := os.Create(fFileName)
		rFile, _ := os.Create(rFileName)
		m.Execute(mFile, data)
		f.Execute(fFile, data)
		r.Execute(rFile, data)
		fmt.Printf("🎉 resource \u001b[32m%s\u001b[0m created \n", options.Name)
	}

}
func createFile() {
	resource := prepareResource(options)
	data := struct {
		Name string
	}{
		Name: resource.Name,
	}
	createDirectoryIfNotExists(resource.Path)
	if _, err := os.Stat(resource.FullPath); err == nil {
		msg := fmt.Sprintf("file \u001b[31m%s.go\u001b[0m already exists, overwrite it?", resource.FileName)
		p := &survey.Confirm{
			Message: msg,
			Default: false,
		}
		survey.AskOne(p, &options.Override)
	}

	if options.Override {
		file, _ := os.Create(resource.FullPath)
		resource.Template.Execute(file, data)
		fmt.Printf("🎉 file \u001b[32m%s\u001b[0m created in \u001b[32m%s\u001b[0m \n", resource.FileName, resource.FullPath)
	}
}

func getCurrentDirectory() string {
	dir, _ := os.Getwd()
	if os.Getenv("ENV") == "dev" || os.Getenv("APP_ENV") == "dev" {
		dir = dir + "/tmp"
	}
	return dir
}

func prepareResource(opt CreateOption) *Resource {
	rPath := "templates/resources/" + opt.Type + ".tmpl"
	tem, _ := template.ParseFS(TemplateFS, rPath)
	dir, _ := os.Getwd()
	//? Dev mode
	if os.Getenv("ENV") == "dev" || os.Getenv("APP_ENV") == "dev" {
		dir = dir + "/tmp"
	}
	name := transformResourceName(opt.Name, opt.Type)
	file := transformFileName(opt.Name, opt.Type) + ".go"
	path := fmt.Sprintf("%s/%s", dir, opt.Type)
	r := &Resource{
		Name:     name,
		Path:     path,
		Template: tem,
		FileName: file,
		FullPath: path + "/" + file,
	}
	return r
}
func createDirectoryIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
}

func transformFileName(name string, rType string) string {
	name = findAndReplaceName(name, rType)
	return strcase.ToSnake(name)
}

func transformResourceName(name string, rType string) string {
	name = findAndReplaceName(name, rType)
	return strcase.ToCamel(name)
}

func findAndReplaceName(name string, rType string) string {
	if rType == "response" {
		rType = "resp"
	}
	c := cases.Title(language.English)
	rType = c.String(rType)
	name = strings.TrimSuffix(name, rType)
	name = strings.TrimSuffix(name, strings.ToLower(rType))
	return name + rType
}
