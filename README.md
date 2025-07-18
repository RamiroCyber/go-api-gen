# Go API Generator CLI

[![Go Version](https://img.shields.io/badge/go-1.20%2B-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

Uma ferramenta de linha de comando (CLI) em Go para automatizar a geração de módulos de API seguindo uma arquitetura em camadas: model, repository, service e controller. Suporta operações CRUD básicas e métodos customizados, com integração ao framework Fiber para controllers.

Essa CLI otimiza o desenvolvimento de APIs REST em Go, reduzindo código boilerplate e permitindo customizações rápidas.

## Recursos Principais
- Geração de módulos completos com interfaces e implementações.
- Suporte a CRUD básico (Create, Read, Update, Delete, List).
- Adição de métodos customizados via flags (ex: `--methods FindByEmail`).
- Templates baseados em Fiber para controllers.
- Estrutura modular: Arquivos gerados em `internal/modules/{nome-do-modulo}`.
- Fácil extensão para campos customizados no model ou integração com DBs.

## Requisitos
- Go 1.20 ou superior.
- Para projetos gerados: Dependências como `github.com/gofiber/fiber/v2` (instaladas via `go get`).

## Instalação

### Instalação Global via Go Install
Para usar o CLI em qualquer projeto:
1. Clone o repositório:
   ```
   git clone https://github.com/RamiroCyber/go-api-gen.git
   cd go-api-gen
   ```
2. Instale globalmente:
   ```
   go install
   ```
    - Isso coloca o binário em `$GOPATH/bin`. Certifique-se de que `$GOPATH/bin` está no seu PATH (adicione `export PATH=$PATH:$(go env GOPATH)/bin` ao seu `.bashrc` ou equivalente).

### Instalação via GitHub (para usuários remotos)
```
go install github.com/RamiroCyber/go-api-gen@latest
```

### Build Local
Se preferir um binário local:
```
go build -o go-api-gen
```
- Rode com `./go-api-gen`.

## Uso
O CLI usa [Cobra](https://github.com/spf13/cobra) para comandos intuitivos.

### Comandos Disponíveis
- `go-api-gen --help`: Mostra ajuda geral.
- `go-api-gen generate --help`: Ajuda específica para geração.

### Gerar um Módulo
```
go-api-gen generate module <nome-do-modulo> [flags]
```
- `<nome-do-modulo>`: Nome do módulo (ex: `user`), usado como pasta e prefixo de structs.
- Flags:
    - `--methods <metodo1>,<metodo2>`: Adiciona métodos customizados (ex: `FindByEmail`). Os métodos serão adicionados às interfaces e implementações com placeholders para customização.

Exemplo:
```
go-api-gen generate module user --methods FindByEmail,FindByName
```
- Isso cria `internal/modules/user/` com:
    - `model.go`: Struct do model (ex: `type User struct { ... }`).
    - `repository.go`: Interface do repository.
    - `repository_impl.go`: Implementação básica (mock ou com DB placeholder).
    - `service.go`: Interface do service.
    - `service_impl.go`: Implementação com injeção de repository.
    - `controller.go`: Controller Fiber com handlers CRUD e custom.

## Exemplo de Integração em um Projeto API
1. Crie um novo projeto Go:
   ```
   mkdir minha-api && cd minha-api
   go mod init github.com/seuusuario/minha-api
   go get github.com/gofiber/fiber/v2
   ```

2. Gere um módulo:
   ```
   go-api-gen generate module user --methods FindByEmail
   ```

3. Em `main.go` (ex: `cmd/api/main.go`), integre:
   ```go
   package main

   import (
       "github.com/gofiber/fiber/v2"
       "github.com/seuusuario/minha-api/internal/modules/user"
   )

   func main() {
       app := fiber.New()

       // Injeção de dependências
       userRepo := user.NewUserRepository() // Ajuste com DB se necessário
       userService := user.NewUserService(userRepo)
       userController := user.NewUserController(userService)

       // Rotas
       app.Post("/users", userController.Create)
       app.Get("/users/:id", userController.Read)
       app.Put("/users/:id", userController.Update)
       app.Delete("/users/:id", userController.Delete)
       app.Get("/users", userController.List)
       app.Get("/users/find-by-email", userController.FindByEmail) // Custom

       app.Listen(":3000")
   }
   ```

4. Rode o servidor:
   ```
   go run cmd/api/main.go
   ```
    - Teste endpoints com ferramentas como curl ou Postman.

## Customizações
- **Adicionar Campos ao Model**: Edite `model.go` manualmente após geração. (Futuro: Suporte via flag `--fields id:int,name:string`.)
- **Integração com DB**: No `repository_impl.go`, adicione lógica com SQL puro. Exemplo:
  ```go
  import "database/sql"

  type UserRepositoryImpl struct {
      db *sql.DB
  }

  func NewUserRepository(db *sql.DB) UserRepository {
      return &UserRepositoryImpl{db: db}
  }

  func (r *UserRepositoryImpl) Create(ctx context.Context, entity *User) error {
      // Implemente query SQL aqui
      return nil
  }
  ```
- **Métodos Customizados**: Os templates adicionam placeholders. Ajuste parâmetros (ex: mude `param string` para `email string` em `FindByEmail`).

## Contribuição
Contribuições são bem-vindas! Siga estes passos:
1. Fork o repositório.
2. Crie uma branch: `git checkout -b feature/nova-funcionalidade`.
3. Commit suas mudanças: `git commit -m 'Adiciona nova funcionalidade'`.
4. Push para a branch: `git push origin feature/nova-funcionalidade`.
5. Abra um Pull Request.

Por favor, rode testes e garanta que o código siga o estilo Go (use `go fmt`).

## Licença
Este projeto está licenciado sob a [MIT License](LICENSE).

## Contato
- Autor: Seu Nome (seu@email.com)
- Issues: [Abra uma issue no GitHub](https://github.com/seuusuario/go-api-gen/issues)

Obrigado por usar o Go API Generator CLI! 🚀