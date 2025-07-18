package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var moduleName string
var customMethods []string

// Comando pai: "generate"
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Comandos para gerar componentes",
	Long:  "Comandos para gerar módulos e outros componentes da API.",
}

// Subcomando: "module [name]"
var moduleCmd = &cobra.Command{
	Use:   "module [name]",
	Short: "Gera um módulo com model, repository, service e controller",
	Args:  cobra.ExactArgs(1), // Exige exatamente 1 arg: o nome do módulo
	Run: func(cmd *cobra.Command, args []string) {
		moduleName = args[0]
		generateModule(moduleName, customMethods)
	},
}

func init() {
	// Adicione flags ao subcomando "module"
	moduleCmd.Flags().StringSliceVar(&customMethods, "methods", []string{}, "Métodos customizados, ex: FindByEmail")

	// Aninhe "module" sob "generate"
	generateCmd.AddCommand(moduleCmd)
}

// Função generateModule (mantenha como antes, com templates e lógica de geração)
func generateModule(name string, customs []string) {
	titleName := strings.ToTitle(name)
	data := struct {
		ModuleName      string
		TitleModuleName string
		CustomMethods   []string
	}{
		ModuleName:      name,
		TitleModuleName: titleName,
		CustomMethods:   customs,
	}

	templates := map[string]string{
		"model.go.tmpl":                filepath.Join("internal/modules", name, "model.go"),
		"repository_interface.go.tmpl": filepath.Join("internal/modules", name, "repository.go"),
		"repository_impl.go.tmpl":      filepath.Join("internal/modules", name, "repository_impl.go"),
		"service_interface.go.tmpl":    filepath.Join("internal/modules", name, "service.go"),
		"service_impl.go.tmpl":         filepath.Join("internal/modules", name, "service_impl.go"),
		"controller.go.tmpl":           filepath.Join("internal/modules", name, "controller.go"),
	}

	for tmplFile, outFile := range templates {
		tmplPath := filepath.Join("templates", tmplFile)
		t, err := template.ParseFiles(tmplPath)
		if err != nil {
			fmt.Printf("Erro ao parsear template %s: %v\n", tmplFile, err)
			return
		}

		dir := filepath.Dir(outFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Erro ao criar diretório: %v\n", err)
			return
		}

		f, err := os.Create(outFile)
		if err != nil {
			fmt.Printf("Erro ao criar arquivo %s: %v\n", outFile, err)
			return
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {

			}
		}(f)

		if err := t.Execute(f, data); err != nil {
			fmt.Printf("Erro ao executar template %s: %v\n", tmplFile, err)
		}
	}

	fmt.Printf("Módulo %s gerado com sucesso!\n", name)
}
