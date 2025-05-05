# Depurando Bash Scripts no VS Code

O VS Code não possui um depurador de Bash nativo. Para depurar arquivos `.sh`, você precisa instalar uma extensão específica.

**Extensão Recomendada:**

*   **Nome:** Bash Debug
*   **ID:** `rogalmic.bash-debug`
*   **Link:** Bash Debug no Marketplace

**Instalação:**

Você pode instalar a extensão pela interface do VS Code (procurando por `rogalmic.bash-debug`) ou diretamente pelo terminal:

```bash
code --install-extension rogalmic.bash-debug
```

### Importante:

Este é um depurador simples. Ele é muito útil para tarefas básicas como definir breakpoints, avançar linha por linha (step over/into/out) e inspecionar variáveis simples.
Recursos mais avançados do Bash, como o comando eval ou manipulações complexas de subshells, podem não funcionar como esperado ou não ser totalmente suportados pelo depurador. Use-o como uma ferramenta auxiliar para entender o fluxo do script, mas esteja ciente de suas limitações.


# Depurando Código Go

A depuração de código Go no VS Code é geralmente mais direta, pois é suportada pela extensão oficial do Go.

você só precisa adicionar breakpoints e iniciar a depuração. O VS Code irá compilar o código automaticamente e iniciar o depurador.

Pré-requisitos:

* Tenha a extensão Go para VS Code instalada. Ela geralmente é recomendada automaticamente ao abrir um projeto Go.
* A extensão Go utiliza ferramentas como gopls (o language server) e dlv (o depurador Delve).
