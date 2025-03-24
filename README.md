![LINUXtips Logo](LINUXtips-logo.png)

# Girus CLI

## Sobre o Girus CLI

O Girus CLI é uma ferramenta de linha de comando que facilita a criação, gerenciamento e utilização da plataforma Girus - um ambiente de laboratórios interativos baseado em Kubernetes.

Desenvolvido como parte do projeto Girus da LinuxTips, o CLI simplifica o processo de implantação da plataforma em ambientes locais, permitindo que instrutores e estudantes configurem rapidamente um ambiente de laboratório completo para treinamentos técnicos.

## Recursos Principais

- **Criação de Cluster**: Implante automaticamente um cluster Kubernetes local usando Kind
- **Implantação da Plataforma**: Configure a plataforma Girus completa com backend e frontend
- **Port-forwarding Automático**: Acesse facilmente os serviços através de portas locais
- **Gerenciamento de Laboratórios**: Liste e exclua clusters existentes
- **Compatível com Múltiplos SO**: Funciona em Linux, macOS e Windows

## Requisitos

- **Go** (versão 1.21 ou superior)
- **Docker** (em execução)
- **kubectl**
- **kind** (Kubernetes in Docker)

## Instalação

### Instalação Automática (Linux/macOS)

O script de instalação verifica automaticamente as dependências necessárias e instala o Girus CLI:

```bash
curl -fsSL https://raw.githubusercontent.com/linuxtips/girus/main/girus-cli/install.sh | bash
```

O script verifica e instala automaticamente:
- Go (se não estiver instalado)
- Kind (se não estiver instalado)
- Kubectl (se não estiver instalado)
- Girus CLI

### Instalação Manual

1. Clone o repositório:
   ```bash
   git clone https://github.com/linuxtips/girus.git
   ```

2. Acesse o diretório do CLI:
   ```bash
   cd girus/girus-cli
   ```

3. Compile o CLI:
   ```bash
   go build -o girus
   ```

4. Mova o binário para um local no seu PATH:
   ```bash
   sudo mv girus /usr/local/bin/
   ```

## Uso

### Criar um Cluster Girus

Para criar um novo cluster com a plataforma Girus:

```bash
girus create
```

Opções disponíveis:
- `--file`: Utilize um arquivo YAML personalizado para implantação
- `--cluster-name`: Especifique um nome para o cluster (padrão: "girus")
- `--verbose`: Exiba informações detalhadas durante a implantação
- `--skip-port-forward`: Não configure port-forwarding automático
- `--skip-browser`: Não abra o navegador automaticamente após a implantação

### Listar Clusters Girus

Para verificar os clusters Girus existentes:

```bash
girus list
```

### Excluir um Cluster Girus

Para remover um cluster existente:

```bash
girus delete
```

Opções disponíveis:
- `--cluster-name`: Especifique o nome do cluster a ser excluído (padrão: "girus")

## Fluxo de Implantação

Ao executar `girus create`, o CLI realiza as seguintes ações:

1. Verifica se o Docker está em execução
2. Verifica a existência de clusters anteriores
3. Cria um novo cluster Kind
4. Implanta os componentes do Girus:
   - Namespace dedicado
   - Permissões e contas de serviço
   - Backend (API REST)
   - Frontend (interface web)
5. Configura port-forwarding para acesso local:
   - Backend: http://localhost:8080
   - Frontend: http://localhost:8000
6. Abre o navegador com a interface do Girus

## Solução de Problemas

### Verificar Status do Cluster

Para verificar o status do cluster:

```bash
kubectl cluster-info --context kind-girus
```

### Verificar Status dos Pods

Para verificar se todos os componentes estão em execução:

```bash
kubectl get pods -n girus
```

### Problemas Comuns

1. **Docker não está em execução**:
   ```bash
   sudo systemctl start docker
   ```

2. **Conflito de portas**:
   O CLI verifica e tenta liberar as portas 8000 e 8080 automaticamente.

3. **Problemas com Kind**:
   ```bash
   kind delete cluster --name girus
   ```
   Em seguida, tente criar novamente.

## Personalização

### Arquivo de Deployment Personalizado

Você pode usar um arquivo YAML personalizado para configurações específicas:

```bash
girus create --file minha-configuracao.yaml
```

### Adicionar Laboratórios Personalizados

Para adicionar laboratórios personalizados:

```bash
girus create --lab-file meu-laboratorio.yaml
```

## Desenvolvimento

Para contribuir com o desenvolvimento do Girus CLI:

1. Faça um fork do repositório no GitHub
2. Clone seu fork localmente
3. Crie uma branch para sua contribuição
4. Faça suas alterações
5. Execute testes locais
6. Envie um Pull Request

## Arquitetura

O Girus CLI é construído com:

- [Cobra](https://github.com/spf13/cobra): Framework de linha de comando para Go
- [Kind](https://kind.sigs.k8s.io/): Para criar clusters Kubernetes locais em contêineres Docker
- [kubectl](https://kubernetes.io/docs/reference/kubectl/): Para interagir com o cluster Kubernetes

## Licença

Este projeto é distribuído sob a licença Apache 2.0. Veja o arquivo LICENSE para mais detalhes.

## Suporte

Se encontrar problemas ou tiver dúvidas, abra uma issue no [repositório do GitHub](https://github.com/linuxtips/girus). 