package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

//go:embed templates/*.tmpl
var templatesFS embed.FS

type Field struct {
	Name string
	Type string
	Tags string
}

var generateCmd = &cobra.Command{
	Use:   "generate module <name> [flags]",
	Short: "Generate a new API module",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		moduleName := args[1]
		customMethodsStr, _ := cmd.Flags().GetString("methods")
		fieldsStr, _ := cmd.Flags().GetString("fields")
		db, _ := cmd.Flags().GetString("db")
		validator, _ := cmd.Flags().GetBool("validator")
		tests, _ := cmd.Flags().GetBool("tests")
		auth, _ := cmd.Flags().GetString("auth")
		swagger, _ := cmd.Flags().GetBool("swagger")
		docker, _ := cmd.Flags().GetBool("docker")
		ci, _ := cmd.Flags().GetBool("ci")
		rootPackage, _ := cmd.Flags().GetString("root-package")

		var fields []Field
		if fieldsStr != "" {
			for _, f := range strings.Split(fieldsStr, ",") {
				parts := strings.Split(f, ":")
				if len(parts) < 2 {
					return fmt.Errorf("invalid field format: %s", f)
				}
				name := strings.Title(parts[0])
				typeAndTags := parts[1]
				typeParts := strings.Split(typeAndTags, " ")
				typ := typeParts[0]
				tags := ""
				if len(typeParts) > 1 {
					tags = strings.Join(typeParts[1:], " ")
				}
				fields = append(fields, Field{Name: name, Type: typ, Tags: tags})
			}
		}

		var customMethods []string
		if customMethodsStr != "" {
			customMethods = strings.Split(customMethodsStr, ",")
			var filtered []string
			for _, m := range customMethods {
				if m != "" {
					filtered = append(filtered, m)
				}
			}
			customMethods = filtered
		}

		data := map[string]interface{}{
			"ModuleName":      moduleName,
			"TitleModuleName": strings.Title(moduleName),
			"CustomMethods":   customMethods,
			"Fields":          fields,
			"DB":              db,
			"Validator":       validator,
			"Tests":           tests,
			"Auth":            auth,
			"Swagger":         swagger,
			"RootPackage":     rootPackage,
		}

		moduleDir := filepath.Join("internal", "modules", moduleName)
		if err := os.MkdirAll(moduleDir, 0755); err != nil {
			return fmt.Errorf("failed to create module directory %s: %w", moduleDir, err)
		}

		tmplFiles := []string{
			"model.go.tmpl",
			"repository.go.tmpl",
			"repository_impl.go.tmpl",
			"service.go.tmpl",
			"service_impl.go.tmpl",
			"controller.go.tmpl",
		}

		if tests {
			tmplFiles = append(tmplFiles, "service_test.go.tmpl")
		}
		if auth == "jwt" {
			middlewareDir := filepath.Join("pkg", "middleware")
			if err := os.MkdirAll(middlewareDir, 0755); err != nil {
				return fmt.Errorf("failed to create middleware directory %s: %w", middlewareDir, err)
			}
			if err := renderTemplate("auth.go.tmpl", filepath.Join(middlewareDir, "auth.go"), data); err != nil {
				return err
			}
		}
		if swagger {
			swaggerDir := "swagger"
			if err := os.MkdirAll(swaggerDir, 0755); err != nil {
				return fmt.Errorf("failed to create swagger directory %s: %w", swaggerDir, err)
			}
			fmt.Println("Swagger enabled. Run `swag init` to generate swagger.json.")
		}
		if docker {
			if err := renderTemplate("Dockerfile.tmpl", "Dockerfile", data); err != nil {
				return err
			}
		}
		if ci {
			ciDir := filepath.Join(".github", "workflows")
			if err := os.MkdirAll(ciDir, 0755); err != nil {
				return fmt.Errorf("failed to create CI directory %s: %w", ciDir, err)
			}
			if err := renderTemplate("ci.yml.tmpl", filepath.Join(ciDir, "ci.yml"), data); err != nil {
				return err
			}
		}

		errorsDir := filepath.Join("pkg", "errors")
		if err := os.MkdirAll(errorsDir, 0755); err != nil {
			return fmt.Errorf("failed to create errors directory %s: %w", errorsDir, err)
		}
		if err := renderTemplate("errors.go.tmpl", filepath.Join(errorsDir, "errors.go"), data); err != nil {
			return err
		}

		for _, tmpl := range tmplFiles {
			outputFile := strings.Replace(tmpl, ".tmpl", "", 1)
			if err := renderTemplate(tmpl, filepath.Join(moduleDir, outputFile), data); err != nil {
				return err
			}
		}

		fmt.Printf("Módulo %s gerado com sucesso em %s!\n", moduleName, moduleDir)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of go-api-gen",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("go-api-gen version %s\n", Version)
	},
}

func renderTemplate(tmplFile, output string, data map[string]interface{}) error {
	tmplPath := "templates/" + tmplFile
	tmplContent, err := templatesFS.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", tmplPath, err)
	}

	tmpl, err := template.New(tmplFile).Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"add":   func(a, b int) int { return a + b },
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", tmplFile, err)
	}

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", output, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to render template %s to %s: %w", tmplFile, output, err)
	}
	return nil
}

func listEmbeddedFiles() {
	fmt.Println("Debug: List of embedded templates:")
	err := fs.WalkDir(templatesFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Error walking path %s: %v\n", path, err)
			return err
		}
		if !d.IsDir() {
			fmt.Println(path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error listing embedded files: %v\n", err)
	}
}

func init() {
	generateCmd.Flags().String("methods", "", "Comma-separated custom methods (ex: FindByEmail)")
	generateCmd.Flags().String("fields", "", "Comma-separated fields: name:type [tags] (ex: name:string required,email:string email)")
	generateCmd.Flags().String("db", "postgres", "Database type: postgres, mysql, gorm, mongo")
	generateCmd.Flags().Bool("validator", true, "Enable validator/v10 for validations")
	generateCmd.Flags().Bool("tests", false, "Generate unit tests with mocks")
	generateCmd.Flags().String("auth", "", "Authentication type: jwt")
	generateCmd.Flags().Bool("swagger", false, "Generate Swagger annotations")
	generateCmd.Flags().Bool("docker", false, "Generate Dockerfile")
	generateCmd.Flags().Bool("ci", false, "Generate GitHub Actions CI workflow")
	generateCmd.Flags().String("root-package", "github.com/user/project", "Root package for the generated project")
}

func main() {
	listEmbeddedFiles()
	rootCmd := &cobra.Command{Use: "go-api-gen"}
	rootCmd.AddCommand(generateCmd, versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
