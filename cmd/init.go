package cmd

import (
	"bufio"
	"fmt"
	"log"
	"syscall"
	"os"
	"strings"
	"crypto/rand"
	"crypto/sha256"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/term"
	"encoding/json"
	"encoding/base64"
)

type VaultEntry struct {
    NomeSito      string `json:"nome_sito"`
    Username      string `json:"username"`
    Password      string `json:"password"`
    URL           string `json:"url"`
    DataCreazione string `json:"data_creazione"`
}

type Vault struct {
	Salt []byte          `json:"salt"`
	Check string		 `json:"check"`
    Entries []VaultEntry `json:"entries"`
}

var initCmd = &cobra.Command{
    Use:   "init",
    Short: "Inizializza un nuovo vault protetto da passphrase",
    Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Crea un nuovo vault: ")
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("nome del vault: ")
		vault_name, _ := reader.ReadString('\n') 
		vault_name = strings.TrimSpace(vault_name)
		vault_name = strings.ReplaceAll(vault_name, " ", "_")
		fmt.Println("Scegli una passphrase: ")
		passPhrase, err := term.ReadPassword(int(syscall.Stdin))

		if err != nil {
			fmt.Println("\nErrore nella lettura della passPhrase: ", err)
			return
		}

		fmt.Println("\nPassPhrase letta con successo!")
		_ = passPhrase

		salt, err := generateSalt(16)
		if err != nil {
			log.Fatal("Errore nella generazione del salt: ", err)
		}

		key := derivateKey(passPhrase, salt)
		// fmt.Println("Passphrase derivata con successo!", key)

		checkString := "vault_check"
		encryptedCheck, err := encrypt(checkString, key) 
		if err != nil {
		    fmt.Println("Errore nella cifratura del check: ", err)
		}
		

		vault := Vault {
			Salt: salt,
			Check: base64.StdEncoding.EncodeToString(encryptedCheck),
			Entries: []VaultEntry{},
		}

		data, err := json.MarshalIndent(vault, "", " ")
		
		if err != nil {
			fmt.Println("Errore nella serializzazione del vault!")
		}

		err = os.WriteFile(vault_name+".json", data, 0600)
		if err != nil {
		    fmt.Println("Errore nella scrittura del file vault.json!", err)
		}

		fmt.Println("Vault inizializzato con successo!")

	},
}

func init() {
    rootCmd.AddCommand(initCmd)
}

func generateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	return salt, nil
}

func derivateKey(passPhrase []byte, salt []byte) []byte {
	iterations := 100_000
	keyLength := 32
	return pbkdf2.Key(passPhrase, salt, iterations, keyLength, sha256.New)
}
