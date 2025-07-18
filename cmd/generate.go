package cmd

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

var moduleName string
var customMethods []string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Comandos para gerar componentes",
	Long:  "Comandos para gerar módulos e outros componentes da API.",
}

var moduleCmd = &cobra.Command{
	Use:   "module [name]",
	Short: "Gera um módulo com model, repository, service e controller",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moduleName = args[0]
		generateModule(moduleName, customMethods)
	},
}

func init() {
	moduleCmd.Flags().StringSliceVar(&customMethods, "methods", []string{}, "Métodos customizados, ex: FindByEmail")
	generateCmd.AddCommand(moduleCmd)
}

func generateModule(name string, customs []string) {
	moduleDir := strings.ToLower(name)
	titleName := strings.Title(moduleDir)

	data := struct {
		ModuleName      string
		TitleModuleName string
		CustomMethods   []string
	}{
		ModuleName:      moduleDir,
		TitleModuleName: titleName,
		CustomMethods:   customs,
	}

	funcMap := template.FuncMap{
		"title": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
		},
	}

	templates := map[string]string{
		"model.go.tmpl":                filepath.Join("internal/modules", moduleDir, "model.go"),
		"repository_interface.go.tmpl": filepath.Join("internal/modules", moduleDir, "repository.go"),
		"repository_impl.go.tmpl":      filepath.Join("internal/modules", moduleDir, "repository_impl.go"),
		"service_interface.go.tmpl":    filepath.Join("internal/modules", moduleDir, "service.go"),
		"service_impl.go.tmpl":         filepath.Join("internal/modules", moduleDir, "service_impl.go"),
		"controller.go.tmpl":           filepath.Join("internal/modules", moduleDir, "controller.go"),
	}

	for tmplFile, outFile := range templates {
		tmplContent, err := templatesFS.ReadFile("templates/" + tmplFile)
		if err != nil {
			fatalf("Erro ao ler template embutido %s: %v", tmplFile, err)
		}

		t := template.New(tmplFile).Funcs(funcMap)
		t, err = t.Parse(string(tmplContent))
		if err != nil {
			fatalf("Erro ao parsear template %s: %v", tmplFile, err)
		}

		dir := filepath.Dir(outFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fatalf("Erro ao criar diretório: %v", err)
		}

		f, err := os.Create(outFile)
		if err != nil {
			fatalf("Erro ao criar arquivo %s: %v", outFile, err)
		}

		if err := t.Execute(f, data); err != nil {
			f.Close()
			fatalf("Erro ao executar template %s: %v", tmplFile, err)
		}
		f.Close()
	}

	fmt.Printf("Módulo %s gerado com sucesso!\n", name)
}

func fatalf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}
