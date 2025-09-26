package cmd

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)


var getCmd = &cobra.Command{
	Use: "get",
	Short: "mostra la password relativa ad un sito",
	Run: func(cmd *cobra.Command, args []string){

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

		for _, entry := range vault.Entries {
			if entry.NomeSito == args[0]{
				key := derivateKey(passPhrase, vault.Salt)
				encryptedPass, _ := base64.StdEncoding.DecodeString(entry.Password)
				decryptPass, _ := decrypt(encryptedPass, key)
				fmt.Printf("password %s: %s", entry.NomeSito, decryptPass)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}