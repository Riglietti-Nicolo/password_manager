package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"syscall"
	"encoding/base64"

)

var editCmd = &cobra.Command{
	Use: "edit",
	Short: "modifica una password del vault",
	Run: func(cmd *cobra.Command, args []string) {

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
		
		fmt.Println("Nuova password: ")
		newPasswordBytes, _ := term.ReadPassword(int(syscall.Stdin))
		newPassword := string(newPasswordBytes)		

		
		key := derivateKey(passPhrase, vault.Salt)

		encryptedPasswordBytes, err := encrypt(newPassword, key)

		if err != nil {
    		fmt.Println("Errore nella cifratura:", err)
    		return
		}

		encryptedPassword := base64.StdEncoding.EncodeToString(encryptedPasswordBytes)
		found := false
		for i := range vault.Entries {
			if vault.Entries[i].NomeSito == args[0]{
				vault.Entries[i].Password = encryptedPassword
				found = true
				break
			}
		}

		if !found {
			fmt.Println("Entry non trovata")
			return
		}
		

		data, err = json.MarshalIndent(vault, "", "  ")
		if err != nil {
    		fmt.Println("Errore nella serializzazione del vault:", err)
    		return
		}

		err = os.WriteFile(vault_name, data, 0600)
		if err != nil {
    		fmt.Println("Errore nel salvataggio del vault:", err)
    		return
		}

		fmt.Println("Nuova password aggiornata con successo:")

	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}