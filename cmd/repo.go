package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/badtuxx/girus-cli/internal/common"
	"github.com/badtuxx/girus-cli/internal/repo"
	"github.com/spf13/cobra"
)

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: common.T("Gerencia repositórios de laboratórios", "Gestiona repositorios de laboratorios"),
	Long:  common.T(`Gerencia repositórios de laboratórios, permitindo adicionar, remover, listar e atualizar repositórios.`, `Gestiona repositorios de laboratorios, permitiendo agregar, eliminar, listar y actualizar repositorios.`),
}

var repoAddCmd = &cobra.Command{
	Use:   "add [nome] [url]",
	Short: common.T("Adiciona um novo repositório", "Agrega un nuevo repositorio"),
	Long:  common.T(`Adiciona um novo repositório de laboratórios com o nome e URL especificados.`, `Agrega un nuevo repositorio de laboratorios con el nombre y URL especificados.`),
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

		fmt.Printf(common.T("Repositório '%s' adicionado com sucesso.\n", "Repositorio '%s' agregado con éxito.\n"), name)
		return nil
	},
}

var repoRemoveCmd = &cobra.Command{
	Use:   "remove [nome]",
	Short: common.T("Remove um repositório", "Elimina un repositorio"),
	Long:  common.T(`Remove um repositório de laboratórios pelo nome.`, `Elimina un repositorio de laboratorios por su nombre.`),
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

		fmt.Printf(common.T("Repositório '%s' removido com sucesso.\n", "Repositorio '%s' eliminado con éxito.\n"), name)
		return nil
	},
}

var repoListCmd = &cobra.Command{
	Use:   "list",
	Short: common.T("Lista todos os repositórios", "Lista todos los repositorios"),
	Long:  common.T(`Lista todos os repositórios de laboratórios configurados.`, `Lista todos los repositorios de laboratorios configurados.`),
	RunE: func(cmd *cobra.Command, args []string) error {
		rm, err := repo.NewRepositoryManager()
		if err != nil {
			return err
		}

		repos := rm.ListRepositories()
		if len(repos) == 0 {
			fmt.Println(common.T("Nenhum repositório configurado.", "Ningún repositorio configurado."))
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, common.T("NOME\tURL\tDESCRIÇÃO", "NOMBRE\tURL\tDESCRIPCIÓN"))
		for _, r := range repos {
			fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, r.URL, r.Description)
		}
		w.Flush()

		return nil
	},
}

var repoUpdateCmd = &cobra.Command{
	Use:   "update [nome] [url]",
	Short: common.T("Atualiza um repositório", "Actualiza un repositorio"),
	Long:  common.T(`Atualiza um repositório de laboratórios existente com novos dados.`, `Actualiza un repositorio de laboratorios existente con nuevos datos.`),
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

		fmt.Printf(common.T("Repositório '%s' atualizado com sucesso.\n", "Repositorio '%s' actualizado con éxito.\n"), name)
		return nil
	},
}

func init() {
	repoCmd.AddCommand(repoAddCmd, repoRemoveCmd, repoListCmd, repoUpdateCmd)

	// Flags para os comandos
	repoAddCmd.Flags().String("description", "", common.T("Descrição do repositório", "Descripción del repositorio"))
	repoUpdateCmd.Flags().String("description", "", common.T("Nova descrição do repositório", "Nueva descripción del repositorio"))
}
