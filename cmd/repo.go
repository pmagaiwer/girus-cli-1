package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/badtuxx/girus-cli/internal/repo"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Gerencia repositórios de laboratórios",
	Long:  `Gerencia repositórios de laboratórios, permitindo adicionar, remover, listar e atualizar repositórios.`,
}

var repoAddCmd = &cobra.Command{
	Use:   "add [nome] [url]",
	Short: "Adiciona um novo repositório",
	Long:  `Adiciona um novo repositório de laboratórios com o nome e URL especificados.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]
		description, _ := cmd.Flags().GetString("description")

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		if err := rm.AddRepository(name, url, description); err != nil {
			return err
		}

		fmt.Printf("Repositório '%s' adicionado com sucesso.\n", name)
		return nil
	},
}

var repoRemoveCmd = &cobra.Command{
	Use:   "remove [nome]",
	Short: "Remove um repositório",
	Long:  `Remove um repositório de laboratórios pelo nome.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		if err := rm.RemoveRepository(name); err != nil {
			return err
		}

		fmt.Printf("Repositório '%s' removido com sucesso.\n", name)
		return nil
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos os repositórios",
	Long:  `Lista todos os repositórios de laboratórios configurados.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		repos := rm.ListRepositories()
		if len(repos) == 0 {
			fmt.Println("Nenhum repositório configurado.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NOME\tURL\tDESCRIÇÃO")
		for _, r := range repos {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, r.URL, r.Description)
		}
		w.Flush()

		return nil
	},
}

var repoUpdateCmd = &cobra.Command{
	Use:   "update [nome] [url]",
	Short: "Atualiza um repositório",
	Long:  `Atualiza um repositório de laboratórios existente com novos dados.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]
		description, _ := cmd.Flags().GetString("description")

		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		if err := rm.UpdateRepository(name, url, description); err != nil {
			return err
		}

		fmt.Printf("Repositório '%s' atualizado com sucesso.\n", name)
		return nil
	},
}

func init() {
	repoCmd.AddCommand(repoAddCmd, repoRemoveCmd, repoListCmd, repoUpdateCmd)

	// Flags para os comandos
	repoAddCmd.Flags().String("description", "", "Descrição do repositório")
	repoUpdateCmd.Flags().String("description", "", "Nova descrição do repositório")
}
