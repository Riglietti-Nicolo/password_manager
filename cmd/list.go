package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"bufio"
	"strings"
	"os"
	"encoding/json"
)


var listCmd = &cobra.Command{
	Use: "list",
	Short: "visualizza la lista delle entry salvate",
	Run: func (cmd *cobra.Command, args []string){
		fmt.Println("comando list")
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

		for i, entry := range vault.Entries {
			fmt.Printf(`
%d: Nome Sito:	    %s
   Username:	    %s
   Password: 	    [nascosta]
   URL:		    %s
   Data creazione:  %s

			`, i, entry.NomeSito, entry.Username, entry.URL, entry.DataCreazione)
		}
		
	},
}

func init(){
	rootCmd.AddCommand(listCmd)
}