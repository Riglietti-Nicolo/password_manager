package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"bufio"
	"os"
	"strings"
	"encoding/json"
	"syscall"

)

var deleteCmd = &cobra.Command {
	Use: "delete",
	Short: "elimina un' entry",
	Run: func (cmd *cobra.Command, args []string)  {
		if args == nil {
			fmt.Println("Parametro mancante: <nome sito>")
			return
		}
		
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("nome del vault a cui vuoi accedere: ")
		vault_name, _ := reader.ReadString('\n') 
		vault_name = strings.TrimSpace(vault_name)
		vault_name = strings.ReplaceAll(vault_name, " ", "_")
		vault_name = vault_name+".json"

		data, err := os.ReadFile(vault_name)

		if err != nil {
			fmt.Println("Errore nella lettura del vault: ", err)
			return
		}
		
		var vault Vault
		err = json.Unmarshal(data, &vault)
		if err != nil {
			fmt.Println("Errore nella decodifica del vault: ", err)
			return
		}

		fmt.Println("inserisci la Passphrase: ")
		passPhrase, err := term.ReadPassword(int(syscall.Stdin))

		if err != nil {
			fmt.Println("Errore nella lettura del passphrase: ", err)
			return
		}

		ok, err := checkPassphrase(&vault, passPhrase)

		if err != nil {
    		fmt.Println("Errore nella verifica della passphrase:", err)
    		return
		}
		
		if !ok {
		    fmt.Println("Passphrase errata, accesso negato.")
		    return
		}

		fmt.Println("Vault caricato correttamente!")

		fmt.Printf("sei sicuro di voler eliminare le informazioni relative a %s (S/n): ", args[0])
		var scelta string
		fmt.Scanf("%s", &scelta)
		if scelta == "n" || scelta == "N" {
			fmt.Println("Informazioni non eliminate")
			return
		}else if scelta == "S" || scelta == "s"{
			for i, entry := range vault.Entries {
				if entry.NomeSito == args[0] {
					vault.Entries = append(vault.Entries[:i], vault.Entries[i+1:]... )
					break
				}
			}

			newData, err := json.MarshalIndent(vault, "", "  ")
			if err != nil {
				fmt.Println("Errore nella serializzazione del vault:", err)
				return
			}

			err = os.WriteFile(vault_name, newData, 0600)
			if err != nil {
				fmt.Println("Errore nel salvataggio del vault:", err)
				return
			}

			fmt.Println("Informazioni eliminate con successo!")
		} else {
			fmt.Println("input non valido")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}



