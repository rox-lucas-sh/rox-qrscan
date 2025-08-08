package vertex

import "google.golang.org/genai"

type Item struct {
	Descricao      string  `json:"descricao"`
	Quantidade     float64 `json:"quantidade"`
	PrecoUnitario  float64 `json:"preco_unitario"`
	PrecoTotalItem float64 `json:"preco_total_item"`
}

type NotaFiscal struct {
	EmitenteNome   string  `json:"emitente_nome"`
	EmitenteCNPJ   string  `json:"emitente_cnpj"`
	DataEmissao    string  `json:"data_emissao"`
	HoraEmissao    string  `json:"hora_emissao"`
	ValorTotal     float64 `json:"valor_total"`
	FormaPagamento string  `json:"forma_pagamento"`
	ChaveAcesso    string  `json:"chave_acesso"`
	Itens          []Item  `json:"itens"`
}

func getSchema() *genai.Schema {
	// size := int64(44)
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"emitente_nome": {
				Type: genai.TypeString,
			},
			"emitente_cnpj": {
				Type: genai.TypeString,
			},
			"data_emissao": {
				Type: genai.TypeString,
			},
			"hora_emissao": {
				Type: genai.TypeString,
			},
			"valor_total": {
				Type: genai.TypeNumber,
			},
			"forma_pagamento": {
				Type: genai.TypeString,
			},
			"chave_acesso": {
				Type: genai.TypeString,
				// MaxLength: &size,
				// MinLength: &size,
				// Pattern:   "^[0-9]{44}$",
			},
			"itens": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"descricao": {
							Type: genai.TypeString,
						},
						"quantidade": {
							Type: genai.TypeNumber,
						},
						"preco_unitario": {
							Type: genai.TypeNumber,
						},
						"preco_total_item": {
							Type: genai.TypeNumber,
						},
					},
					Required: []string{"descricao", "quantidade", "preco_unitario", "preco_total_item"},
				},
			},
		},
		// Required: []string{
		// 	"emitente_nome",
		// 	"emitente_cnpj",
		// 	"data_emissao",
		// 	"hora_emissao",
		// 	"valor_total",
		// 	"forma_pagamento",
		// 	"chave_acesso",
		// 	"itens",
		// },
	}
}
