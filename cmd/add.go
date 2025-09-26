package cmd

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	mathrand "math/rand"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var addCmd = &cobra.Command{
	Use: "add",
	Short: "aggiungi una nuova voce al vault",
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

		nomeSito, username, password, url, dataCreazione := inputNewEntry()
		
		if err != nil {
    		fmt.Println("Errore nella lettura della passphrase:", err)
    		return
		}

		key := derivateKey(passPhrase, vault.Salt)

		encryptedPasswordBytes, err := encrypt(password, key)

		if err != nil {
    		fmt.Println("Errore nella cifratura:", err)
    		return
		}

		encryptedPassword := base64.StdEncoding.EncodeToString(encryptedPasswordBytes)
		

		newEntry := VaultEntry{
			NomeSito:		nomeSito,
			Username:		username,
			Password:		encryptedPassword,
			URL:			url,
			DataCreazione: 	dataCreazione,
		}

		vault.Entries = append(vault.Entries, newEntry)

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

		fmt.Println("Voce aggiunta e vault salvato con successo!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}


func inputNewEntry() (string, string, string, string, string){

	reader := bufio.NewReader(os.Stdin)
	
	fmt.Println("Nome del sito: ")
	nomeSito, _ := reader.ReadString('\n')
	nomeSito = strings.TrimSpace(nomeSito)

	fmt.Println("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	
	fmt.Println("password (invio per generare automaticamente una password): ")
	var password string
	passwordByte, _ := term.ReadPassword(int(syscall.Stdin))
	password = string(passwordByte)
	fmt.Println()

	if password == "" {
		fmt.Println("C: lettere maiuscole\nc: lettere minuscole\nn: numeri\ns: caratteri speciali\nlunghezza password")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		passChar := strings.Split(input, ",")


		n, _ := strconv.Atoi(passChar[len(passChar)-1])
		choose := passChar[:len(passChar)-1]

		caratteri := func(start, end int) []rune {
			r := make([]rune, 0, end-start+1)
			for i := start; i<= end; i++ {
				r = append(r, rune(i))
			} 
			return r
		}

		options := map[rune][]rune{
			'C': caratteri(65, 90),
			'c': caratteri(97, 122),
			'n': caratteri(48, 57),
			's': caratteri(33, 47),
		}

		chars := []rune{}
		for _, s := range choose {

			if len(s) == 0 {
				continue
			}
			r := rune(s[0])
			if v, ok := options[r]; ok {
				chars = append(chars, v...)
			}
		}

		mathrand.Seed(time.Now().UnixNano()) //modifica la sequenza random ogni volta che si avvia il programma
		passRunes := make([]rune, n)
		for i := 0; i < n; i++ {
			passRunes[i] = chars[mathrand.Intn(len(chars))]
		}
		password = string(passRunes)
	}

	fmt.Println("URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	dataCreazione := time.Now().Format("2006-01-02 15:04:05")

	return nomeSito, username, password, url, dataCreazione
}

func encrypt(plaintext string , key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	_, err = io.ReadFull(cryptorand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	return append(nonce, ciphertext...), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("Ciphertext troppo corto")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func checkPassphrase(vault *Vault, passPhrase []byte) (bool, error) {
	key := derivateKey(passPhrase, vault.Salt)

	encryptedCheckBytes, err := base64.StdEncoding.DecodeString(vault.Check)
	if err != nil {
		return false, fmt.Errorf("Errore nella codifica check: %s", err)
	}

	decryptedCheck, err := decrypt(encryptedCheckBytes, key)
	if err != nil {
		return false, nil //passphrase errata
	}

	if string(decryptedCheck) == "vault_check" {
		return true, nil
	}

	return false, nil
}