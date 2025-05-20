# Repositórios de Laboratórios GIRUS

Este documento explica como criar e gerenciar repositórios de laboratórios para a plataforma GIRUS.

## Estrutura de um Repositório

Um repositório de laboratórios GIRUS deve seguir a seguinte estrutura:

```
repositorio/
├── index.yaml           # Índice do repositório
└── labs/               # Diretório contendo os laboratórios
    ├── lab1/
    │   ├── lab.yaml    # Definição do laboratório
    │   └── assets/     # Recursos do laboratório (opcional)
    └── lab2/
        ├── lab.yaml
        └── assets/
```

### Arquivo index.yaml

O arquivo `index.yaml` é o ponto de entrada do repositório e deve seguir este formato:

```yaml
apiVersion: v1
generated: "2024-03-20T10:00:00Z"
entries:
  lab-name:
    - name: lab-name
      version: "1.0.0"
      description: "Descrição do laboratório"
      keywords:
        - keyword1
        - keyword2
      maintainers:
        - "Nome <email@exemplo.com>"
      url: "https://github.com/seu-repo/raw/main/labs/lab-name/lab.yaml"
      created: "2024-03-20T10:00:00Z"
      digest: "sha256:hash-do-arquivo"
```

### Arquivo lab.yaml

Cada laboratório deve ter um arquivo `lab.yaml` que define sua estrutura e conteúdo:

```yaml
apiVersion: girus.linuxtips.io/v1
kind: Lab
metadata:
  name: lab-name
  version: "1.0.0"
  description: "Descrição do laboratório"
  author: "Nome do Autor"
  created: "2024-03-20T10:00:00Z"
spec:
  environment:
    image: ubuntu:22.04
    resources:
      cpu: "1"
      memory: "1Gi"
    volumes:
      - name: workspace
        mountPath: /workspace
        size: "1Gi"

  tasks:
    - name: "Nome da Tarefa"
      description: "Descrição da tarefa"
      steps:
        - description: "Descrição do passo"
          command: "comando"
          expectedOutput: "saída esperada"
          hint: "Dica para o usuário"

  validation:
    - name: "Nome da Validação"
      description: "Descrição da validação"
      checks:
        - command: "comando"
          expectedOutput: "saída esperada"
          errorMessage: "Mensagem de erro"
```

## Criando um Repositório

1. Crie um novo repositório Git:
   ```bash
   mkdir girus-labs
   cd girus-labs
   git init
   ```

2. Crie a estrutura básica:
   ```bash
   mkdir -p labs
   touch index.yaml
   ```

3. Adicione seu primeiro laboratório:
   ```bash
   mkdir -p labs/meu-lab
   touch labs/meu-lab/lab.yaml
   ```

4. Edite o arquivo `index.yaml` com as informações do seu repositório.

5. Edite o arquivo `lab.yaml` com a definição do seu laboratório.

6. Faça commit das alterações:
   ```bash
   git add .
   git commit -m "Primeiro laboratório"
   ```

7. Publique o repositório em um serviço como GitHub, GitLab ou Bitbucket.

## Hospedando o Repositório

O repositório pode ser hospedado em qualquer serviço que permita acesso aos arquivos via HTTP/HTTPS. Algumas opções:

1. **GitHub**: Use o recurso "raw" para acessar os arquivos:
   ```
   https://github.com/seu-usuario/seu-repo/raw/main/labs/lab-name/lab.yaml
   ```

2. **GitLab**: Use o recurso "raw" para acessar os arquivos:
   ```
   https://gitlab.com/seu-usuario/seu-repo/raw/main/labs/lab-name/lab.yaml
   ```

3. **Servidor Web**: Hospede os arquivos em um servidor web próprio.

## Adicionando o Repositório ao GIRUS

Para adicionar seu repositório ao GIRUS, use o comando:

```bash
girus repo add meu-repo https://github.com/seu-usuario/seu-repo/raw/main
```

## Boas Práticas

1. **Versionamento**: Mantenha um histórico de versões dos laboratórios no `index.yaml`.

2. **Documentação**: Inclua documentação detalhada em cada laboratório.

3. **Validação**: Implemente validações robustas para garantir que os usuários completaram as tarefas corretamente.

4. **Recursos**: Mantenha os recursos (CPU, memória) em níveis razoáveis.

5. **Imagens**: Use imagens Docker oficiais e bem mantidas.

6. **Segurança**: Não inclua credenciais ou informações sensíveis nos laboratórios.

## Exemplo Completo

Veja o diretório `example/` neste repositório para um exemplo completo de um repositório de laboratórios. 