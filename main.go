package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
)

// Estrutura global para manter o estado
type AppState struct {
	FilePath string
	Data     map[string]string
}

// Função para carregar JSON
func loadJSON(filepath string) (map[string]string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var data map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Função para salvar JSON
func saveJSON(filepath string, data map[string]string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, jsonData, 0644)
}

// Função de busca no JSON (apenas nos valores)
func searchInValues(data map[string]string, query string) map[string]string {
	results := make(map[string]string)
	query = strings.ToLower(query)

	for key, value := range data {
		if strings.Contains(strings.ToLower(value), query) {
			results[key] = value
		}
	}

	return results
}

// Função para encontrar valores duplicados
func findDuplicateValues(data map[string]string) map[string][]string {
	valueToKeys := make(map[string][]string)
	
	// Agrupa chaves por valor
	for key, value := range data {
		valueToKeys[value] = append(valueToKeys[value], key)
	}
	
	// Filtra apenas valores com mais de uma chave
	duplicates := make(map[string][]string)
	for value, keys := range valueToKeys {
		if len(keys) > 1 {
			duplicates[value] = keys
		}
	}
	
	return duplicates
}

func main() {
	// Configuração de cores
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()

	// Estado da aplicação
	state := AppState{
		FilePath: "",
		Data:     nil,
	}

	for {
		// Menu principal com emojis
		prompt := promptui.Select{
			Label: "🌐 Gerenciador de Dados JSON 🌐",
			Items: []string{
				"📂 Carregar arquivo JSON",
				"🧹 Limpar caminho do JSON",
				"🔍 Buscar no JSON",
				"📊 Encontrar valores duplicados",
				"➕ Adicionar nova chave",
				"➖ Remover chave",
				"🚪 Sair",
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Println(red("Erro no menu:"), err)
			return
		}

		switch result {
		case "📂 Carregar arquivo JSON":
			// Solicitar caminho do arquivo
			filePrompt := promptui.Prompt{
				Label: "Digite o caminho completo do arquivo JSON",
				Validate: func(input string) error {
					if _, err := os.Stat(input); os.IsNotExist(err) {
						return fmt.Errorf("arquivo não encontrado")
					}
					if filepath.Ext(input) != ".json" {
						return fmt.Errorf("o arquivo deve ser um JSON")
					}
					return nil
				},
			}

			jsonFilePath, err := filePrompt.Run()
			if err != nil {
				fmt.Println(red("Erro ao selecionar arquivo:"), err)
				continue
			}

			loadedData, err := loadJSON(jsonFilePath)
			if err != nil {
				fmt.Println(red("Erro ao carregar JSON:"), err)
				continue
			}

			// Atualiza o estado
			state.FilePath = jsonFilePath
			state.Data = loadedData
			fmt.Println(green("✅ Arquivo carregado com sucesso!"))

		case "🧹 Limpar caminho do JSON":
			// Limpa o estado
			state.FilePath = ""
			state.Data = nil
			fmt.Println(yellow("🧹 Caminho do JSON limpo!"))

		case "🔍 Buscar no JSON":
			if state.FilePath == "" || state.Data == nil {
				fmt.Println(yellow("⚠️ Carregue um arquivo JSON primeiro!"))
				continue
			}

			searchPrompt := promptui.Prompt{
				Label: "Digite sua busca",
			}

			query, err := searchPrompt.Run()
			if err != nil {
				fmt.Println(red("Erro na busca:"), err)
				continue
			}

			results := searchInValues(state.Data, query)
			if len(results) == 0 {
				fmt.Println(yellow("🔎 Nenhum resultado encontrado."))
			} else {
				fmt.Println(cyan("🔍 Resultados da busca:"))
				for key, value := range results {
					fmt.Printf("🔑 %s: %s\n", green(key), yellow(value))
				}
			}

		case "📊 Encontrar valores duplicados":
			if state.FilePath == "" || state.Data == nil {
				fmt.Println(yellow("⚠️ Carregue um arquivo JSON primeiro!"))
				continue
			}

			duplicates := findDuplicateValues(state.Data)
			if len(duplicates) == 0 {
				fmt.Println(yellow("🔍 Nenhum valor duplicado encontrado."))
			} else {
				fmt.Println(cyan("📊 Valores Duplicados:"))
				for value, keys := range duplicates {
					fmt.Printf("%s: %s\n", 
						red(fmt.Sprintf("Valor duplicado: %s", value)), 
						green(fmt.Sprintf("Chaves: %v", keys)))
				}
			}

		case "➕ Adicionar nova chave":
			if state.FilePath == "" || state.Data == nil {
				fmt.Println(yellow("⚠️ Carregue um arquivo JSON primeiro!"))
				continue
			}

			keyPrompt := promptui.Prompt{
				Label: "Digite a nova chave",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("chave não pode ser vazia")
					}
					if existingValue, exists := state.Data[input]; exists {
						// Mostra a chave e valor existente em vermelho
						fmt.Printf("%s\n", red(fmt.Sprintf("❌ Chave já existe: %s = %s", input, existingValue)))
						return fmt.Errorf("chave já existe")
					}
					return nil
				},
			}

			newKey, err := keyPrompt.Run()
			if err != nil {
				fmt.Println(red("Erro ao adicionar chave:"), err)
				continue
			}

			valuePrompt := promptui.Prompt{
				Label: "Digite o valor para a chave",
			}

			newValue, err := valuePrompt.Run()
			if err != nil {
				fmt.Println(red("Erro ao adicionar valor:"), err)
				continue
			}

			state.Data[newKey] = newValue
			err = saveJSON(state.FilePath, state.Data)
			if err != nil {
				fmt.Println(red("Erro ao salvar JSON:"), err)
				continue
			}
			fmt.Println(green("✅ Chave adicionada com sucesso!"))

		case "➖ Remover chave":
			if state.FilePath == "" || state.Data == nil {
				fmt.Println(yellow("⚠️ Carregue um arquivo JSON primeiro!"))
				continue
			}

			keyPrompt := promptui.Prompt{
				Label: "Digite a chave a ser removida",
				Validate: func(input string) error {
					if input == "" {
						return fmt.Errorf("chave não pode ser vazia")
					}
					if _, exists := state.Data[input]; !exists {
						return fmt.Errorf("chave não encontrada")
					}
					return nil
				},
			}

			keyToRemove, err := keyPrompt.Run()
			if err != nil {
				fmt.Println(red("Erro ao remover chave:"), err)
				continue
			}

			delete(state.Data, keyToRemove)
			err = saveJSON(state.FilePath, state.Data)
			if err != nil {
				fmt.Println(red("Erro ao salvar JSON:"), err)
				continue
			}
			fmt.Println(green("✅ Chave removida com sucesso!"))

		case "🚪 Sair":
			fmt.Println(yellow("👋 Até logo!"))
			return
		}
	}
}