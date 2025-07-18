package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var moduleName string
var customMethods []string // Para métodos extras como FindByEmail

var generateCmd = &cobra.Command{
	Use:   "module [name]",
	Short: "Gera um módulo com model, repository, service e controller",
	Args:  cobra.ExactArgs(1), // Exige exatamente 1 argumento: o nome do módulo
	Run: func(cmd *cobra.Command, args []string) {
		moduleName = args[0]
		// Aqui vamos gerar os arquivos
		generateModule(moduleName, customMethods)
	},
}

func init() {
	generateCmd.Flags().StringSliceVar(&customMethods, "methods", []string{}, "Métodos customizados, ex: FindByEmail")
	// Torna o flag persistente se precisar
}

func generateModule(name string, customs []string) {
	// Capitalize o nome
	titleName := strings.ToTitle(name) // Agora usada no struct

	data := struct {
		ModuleName      string
		TitleModuleName string // Novo campo para o nome capitalizado
		CustomMethods   []string
	}{
		ModuleName:      name,      // Minúsculo, ex: "user"
		TitleModuleName: titleName, // Capitalizado, ex: "User"
		CustomMethods:   customs,
	}

	// Lista de templates e arquivos de saída (igual ao anterior)
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

		// Crie diretórios se não existirem
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
		defer f.Close()

		if err := t.Execute(f, data); err != nil {
			fmt.Printf("Erro ao executar template %s: %v\n", tmplFile, err)
		}
	}

	fmt.Printf("Módulo %s gerado com sucesso!\n", name)
}
