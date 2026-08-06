# `products/list` — Guia da estrutura

Este README explica o papel de cada arquivo do pacote `products/list` e mostra
um endpoint fictício (`GET /products`, com filtros e paginação) atravessando
todas as camadas, do roteamento até a query SQL.

A ideia central: **cada arquivo tem uma única responsabilidade**. Isso deixa
claro onde mexer quando algo muda (ex: trocar de SQLite para Postgres só
afeta `repository.go`; mudar o formato da resposta só afeta `dto.go` e
`mapper.go`).

```
productsHandler/list/
├── routes.go     → registra as rotas do grupo
├── handler.go    → camada HTTP (bind, chamada ao service, resposta)
├── dto.go        → contratos de entrada/saída da API
├── models.go     → representação interna dos dados (espelha o banco)
├── mapper.go     → converte entre models e dtos
├── service.go    → regras de negócio (validação, paginação, orquestração)
├── repository.go → acesso ao banco (SQL puro)
└── swagger.go    → anotações Swagger isoladas do handler
```

---

## 1. `routes.go`

Só faz o "cabeamento" das rotas com o Echo. Não tem lógica nenhuma — se
alguém quiser saber quais endpoints esse pacote expõe, basta abrir este
arquivo.

```go
package list

import "github.com/labstack/echo/v5"

func RegisterRoutes(g *echo.Group) {
	g.GET("", Handler)
}
```

No `internal/server/routes.go` (ou onde você monta os grupos), isso vira:

```go
productsGroup := e.Group("/productsHandler", middleware.RequireAuth)
list.RegisterRoutes(productsGroup)
```

---

## 2. `dto.go`

Define os **contratos da API**: o que o cliente manda (query params, body) e
o que ele recebe de volta. Nada de tipos de banco aqui — isso é o que fica
"público", documentado pelo Swagger.

```go
package list

type ListProductsRequest struct {
	Category string  `query:"category" example:"electronics"`
	MinPrice float64 `query:"min_price" example:"10.00"`
	MaxPrice float64 `query:"max_price" example:"500.00"`
	Page     int     `query:"page" example:"1"`
	Limit    int     `query:"limit" example:"20"`
}

type ProductResponse struct {
	ID          int     `json:"id" example:"42"`
	Name        string  `json:"name" example:"Teclado Mecânico RGB"`
	Description string  `json:"description" example:"Switches azuis, ABNT2"`
	Price       float64 `json:"price" example:"259.90"`
	Category    string  `json:"category" example:"electronics"`
	Stock       int     `json:"stock" example:"15"`
}

type ListProductsResponse struct {
	Products []ProductResponse `json:"productsHandler"`
	Total    int               `json:"total" example:"134"`
	Page     int               `json:"page" example:"1"`
	Limit    int               `json:"limit" example:"20"`
}
```

---

## 3. `dto.go`

Representa os dados **como eles existem no banco**. É a "verdade interna" —
o campo `CategoryID` (chave estrangeira) aparece aqui, mas nunca vaza pro
DTO, que expõe `Category` (nome, não ID).

```go
package list

import "time"

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	CategoryID  int
	Stock       int
	CreatedAt   time.Time
}
```

---

## 4. `mapper.go`

A ponte entre `dto.go` (banco) e `dto.go` (API). Sem essa camada, ou o
handler faz essa conversão manualmente (repetitivo) ou o model acaba
vazando direto pra resposta JSON (acoplamento perigoso — qualquer mudança de
schema quebra o contrato da API).

```go
package list

func ToProductResponse(p Product, categoryName string) ProductResponse {
	return ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    categoryName,
		Stock:       p.Stock,
	}
}

func ToProductResponseList(productsHandler []Product, categories map[int]string) []ProductResponse {
	result := make([]ProductResponse, 0, len(productsHandler))
	for _, p := range productsHandler {
		result = append(result, ToProductResponse(p, categories[p.CategoryID]))
	}
	return result
}
```

---

## 5. `repository.go`

Isola **todo** o SQL. Nada de query solta em handler ou service. Recebe
parâmetros já "prontos" (filtros) e devolve dados crus (`[]Product`), sem
saber nada sobre HTTP.

```go
package list

import "g0/internal/database"

type Filter struct {
	Category string
	MinPrice float64
	MaxPrice float64
	Offset   int
	Limit    int
}

const listProductsQuery = `
	SELECT p.id, p.name, p.description, p.price, p.category_id, p.stock, p.created_at
	FROM productsHandler p
	JOIN categories c ON c.id = p.category_id
	WHERE (? = '' OR c.name = ?)
	  AND p.price >= ?
	  AND (? = 0 OR p.price <= ?)
	ORDER BY p.id
	LIMIT ? OFFSET ?;
`

const countProductsQuery = `
	SELECT COUNT(*)
	FROM productsHandler p
	JOIN categories c ON c.id = p.category_id
	WHERE (? = '' OR c.name = ?)
	  AND p.price >= ?
	  AND (? = 0 OR p.price <= ?);
`

func FindProducts(f Filter) ([]Product, error) {
	rows, err := database.DB.Query(
		listProductsQuery,
		f.Category, f.Category,
		f.MinPrice,
		f.MaxPrice, f.MaxPrice,
		f.Limit, f.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsHandler []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Description,
			&p.Price, &p.CategoryID, &p.Stock, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		productsHandler = append(productsHandler, p)
	}
	return productsHandler, nil
}

func CountProducts(f Filter) (int, error) {
	var total int
	err := database.DB.QueryRow(
		countProductsQuery,
		f.Category, f.Category,
		f.MinPrice,
		f.MaxPrice, f.MaxPrice,
	).Scan(&total)
	return total, err
}
```

---

## 6. `service.go`

As **regras de negócio**: validação de paginação, orquestração entre
repositórios, decisão de valores padrão. É a camada que o `handler.go`
chama — ele não sabe nada de SQL nem de bcrypt/regras, só delega.

```go
package list

func ListProducts(req ListProductsRequest) (ListProductsResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}

	limit := req.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := Filter{
		Category: req.Category,
		MinPrice: req.MinPrice,
		MaxPrice: req.MaxPrice,
		Offset:   (page - 1) * limit,
		Limit:    limit,
	}

	productsHandler, err := FindProducts(filter)
	if err != nil {
		return ListProductsResponse{}, err
	}

	total, err := CountProducts(filter)
	if err != nil {
		return ListProductsResponse{}, err
	}

	// Em um cenário real, isso viria de um repository de categorias,
	// ou de um JOIN que já traga o nome direto no Product.
	categories := map[int]string{}
	for _, p := range productsHandler {
		categories[p.CategoryID] = req.Category
	}

	return ListProductsResponse{
		Products: ToProductResponseList(productsHandler, categories),
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}
```

---

## 7. `handler.go`

A camada HTTP. Faz **só três coisas**: parseia a requisição, chama o
service, traduz o resultado (ou erro) para uma resposta HTTP. Não tem SQL,
não tem regra de negócio.

```go
package list

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func Handler(c *echo.Context) error {
	req := new(ListProductsRequest)
	if err := c.Bind(req); err != nil {
		return c.String(http.StatusBadRequest, "Parâmetros inválidos")
	}

	resp, err := ListProducts(*req)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
```

---

## 8. `swagger.go`

Mantém as anotações Swagger **fora** do handler, para ele não ficar poluído
com comentários gigantes. Como o `swag` associa a anotação ao próximo
`func` declarado logo abaixo dela (em qualquer arquivo do pacote), a
convenção é criar uma função "fantasma" só para pendurar a documentação:

```go
package list

// ProductsDocs existe apenas para hospedar a anotação Swagger deste
// endpoint, mantendo handler.go limpo. Não é chamada em nenhum lugar.
//
// @Summary Lista produtos
// @Description Retorna produtos filtrados por categoria e faixa de preço, com paginação
// @Tags Produtos
// @Accept json
// @Produce json
// @Param category query string false "Categoria do produto"
// @Param min_price query number false "Preço mínimo"
// @Param max_price query number false "Preço máximo"
// @Param page query int false "Página" default(1)
// @Param limit query int false "Itens por página" default(20)
// @Success 200 {object} ListProductsResponse
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /productsHandler [get]
func ProductsDocs() {}
```

> ⚠️ Atenção: o `swag` (swaggo) lê o comentário **imediatamente acima** da
> declaração — por isso a função vazia é necessária. Se preferir evitar essa
> "gambiarra", a alternativa mais comum é manter a anotação direto acima do
> `Handler` em `handler.go`; separar em `swagger.go` é só uma questão de
> organização/gosto quando os endpoints têm documentação muito extensa.

---

## Fluxo completo da requisição

```
GET /productsHandler?category=electronics&min_price=50&page=1&limit=10
        │
        ▼
routes.go        → direciona para Handler
handler.go       → Bind(ListProductsRequest) → chama ListProducts(req)
service.go       → normaliza paginação, monta Filter, chama repository
repository.go    → roda SQL, devolve []Product + total
mapper.go        → converte []Product em []ProductResponse
handler.go       → c.JSON(200, ListProductsResponse{...})
```

Resposta final:

```json
{
  "productsHandler": [
    {
      "id": 42,
      "name": "Teclado Mecânico RGB",
      "description": "Switches azuis, ABNT2",
      "price": 259.90,
      "category": "electronics",
      "stock": 15
    }
  ],
  "total": 134,
  "page": 1,
  "limit": 10
}
```

## Por que separar assim?

- **Testabilidade**: `service.go` e `repository.go` não dependem de Echo —
  dá pra testar com `go test` puro, sem subir servidor HTTP.
- **Troca de banco**: mudar de SQLite pra Postgres afeta só `repository.go`.
- **Contrato estável**: `dto.go` pode ficar igual mesmo que `dto.go` mude
  (ex: adicionar uma coluna no banco não quebra a resposta da API).
- **Handler enxuto**: fica fácil ler o que o endpoint faz sem se perder em
  SQL ou regra de negócio.