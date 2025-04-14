package main

import (
	"fmt"

	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Estrutura global para manter o estado
type AppState struct {
	FilePath string
	RawJSON  string
	Data     map[string]string
}

// Função para carregar JSON
func loadJSON(filepath string) (string, map[string]string, error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return "", nil, err
	}

	jsonStr := string(file)
	data := make(map[string]string)

	// Extrair os pares chave-valor do JSON mantendo a estrutura original
	result := gjson.Parse(jsonStr)
	result.ForEach(func(key, value gjson.Result) bool {
		data[key.String()] = value.String()
		return true
	})

	return jsonStr, data, nil
}

// Função para salvar JSON, preservando a ordem
func saveJSON(filepath string, originalJSON string, data map[string]string, addedKeys map[string]string) error {
	result := originalJSON

	// Se o JSON original está vazio, começamos com um objeto vazio
	if originalJSON == "" || originalJSON == "{}" {
		result = "{}"
	}

	// Primeiro atualizamos valores existentes sem alterar a ordem
	for key, value := range data {
		// Verifica se a chave não é nova (adicionada agora)
		if _, isNew := addedKeys[key]; !isNew {
			result, _ = sjson.Set(result, key, value)
		}
	}

	// Depois adicionamos novas chaves ao final
	for key, value := range addedKeys {
		// Garantimos que a chave seja adicionada ao final
		if !gjson.Get(result, key).Exists() {
			result, _ = sjson.Set(result, key, value)
		}
	}

	// Escreve o resultado no arquivo
	return os.WriteFile(filepath, []byte(result), 0644)
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
		RawJSON:  "",
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

			rawJSON, loadedData, err := loadJSON(jsonFilePath)
			if err != nil {
				fmt.Println(red("Erro ao carregar JSON:"), err)
				continue
			}

			// Atualiza o estado
			state.FilePath = jsonFilePath
			state.RawJSON = rawJSON
			state.Data = loadedData
			fmt.Println(green("✅ Arquivo carregado com sucesso!"))

		case "🧹 Limpar caminho do JSON":
			// Limpa o estado
			state.FilePath = ""
			state.RawJSON = ""
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

			// Cria um mapa para as chaves adicionadas
			addedKeys := make(map[string]string)
			addedKeys[newKey] = newValue

			// Atualiza o Data em memória
			state.Data[newKey] = newValue

			// Salva no arquivo preservando a ordem
			err = saveJSON(state.FilePath, state.RawJSON, state.Data, addedKeys)
			if err != nil {
				fmt.Println(red("Erro ao salvar JSON:"), err)
				continue
			}

			// Atualiza o RawJSON após a adição
			state.RawJSON, _, err = loadJSON(state.FilePath)
			if err != nil {
				fmt.Println(red("Erro ao recarregar JSON:"), err)
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

			// Remove do mapa em memória
			delete(state.Data, keyToRemove)

			// Cria um JSON temporário para remoção
			tempJSON := state.RawJSON
			tempJSON, _ = sjson.Delete(tempJSON, keyToRemove)

			// Salva no arquivo
			err = os.WriteFile(state.FilePath, []byte(tempJSON), 0644)
			if err != nil {
				fmt.Println(red("Erro ao salvar JSON:"), err)
				continue
			}

			// Atualiza o RawJSON após a remoção
			state.RawJSON = tempJSON

			fmt.Println(green("✅ Chave removida com sucesso!"))

		case "🚪 Sair":
			fmt.Println(yellow("👋 Até logo!"))
			return
		}
	}
}
