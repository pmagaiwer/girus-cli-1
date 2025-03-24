# Girus CLI

![Girus Logo](https://raw.githubusercontent.com/linuxtips/girus/main/web/public/girus-logo.svg)

[![Build Status](https://github.com/linuxtips/girus/actions/workflows/build.yml/badge.svg)](https://github.com/linuxtips/girus/actions/workflows/build.yml)
[![Docker Status](https://github.com/linuxtips/girus/actions/workflows/docker.yml/badge.svg)](https://github.com/linuxtips/girus/actions/workflows/docker.yml)
[![Test Status](https://github.com/linuxtips/girus/actions/workflows/test.yml/badge.svg)](https://github.com/linuxtips/girus/actions/workflows/test.yml)

## Sobre o Girus CLI

O Girus CLI é uma ferramenta de linha de comando que facilita a criação, gerenciamento e utilização da plataforma Girus - um ambiente de laboratórios interativos baseado em Kubernetes.

Desenvolvido como parte do projeto Girus da LINUXtips, o CLI simplifica o processo de implantação da plataforma em ambientes locais, permitindo que instrutores e estudantes configurem rapidamente um ambiente de laboratório completo para treinamentos técnicos.

## Recursos Principais

- **Criação de Cluster**: Implante automaticamente um cluster Kubernetes local usando Kind
- **Implantação da Plataforma**: Configure a plataforma Girus completa com backend e frontend
- **Port-forwarding Automático**: Acesse facilmente os serviços através de portas locais
- **Gerenciamento de Laboratórios**: Liste e exclua clusters existentes
- **Compatível com Múltiplos SO**: Funciona em Linux, macOS e Windows
- **Integração com Docker**: Suporte completo para contêineres e ambientes isolados
- **Atualizações Automáticas**: Sistema de verificação e atualização de dependências

## Requisitos

- **Go** (versão 1.22 ou superior)
- **Docker** (em execução)
- **kubectl**
- **kind** (Kubernetes in Docker)

## Instalação

### Instalação Automática (Linux/macOS)

O script de instalação verifica automaticamente as dependências necessárias e instala o Girus CLI:

```bash
curl -fsSL https://girus.linuxtips.io | bash
```

Ou usando o repositório diretamente:

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

## Comandos

### Criar Recursos (`create`)

```bash
# Criar um novo cluster Girus
girus create cluster

# Opções disponíveis:
Cria um cluster Kind com o nome "girus" e implanta todos os componentes necessários.
Por padrão, o deployment embutido no binário é utilizado.

Usage:
  girus create cluster

Flags:
   -f, --file string         Arquivo YAML para deployment do Girus (opcional)
   -h, --help                help for cluster
   --skip-browser        Não abrir o navegador automaticamente
   --skip-port-forward   Não perguntar sobre configurar port-forwarding
   -v, --verbose             Modo detalhado com output completo em vez da barra de progresso
```

### Listar Recursos (`list`)

```bash
# Listar todos os clusters
girus list clusters

# Saída do comando list clusters:
Obtendo lista de clusters Kind...

Clusters Kind disponíveis:
==========================
✅ girus (cluster Girus)
   Pods:
   └─ girus-backend-5dc9b6679f-255z5    Running   true
   └─ girus-frontend-5b8668554d-t552m   Running   true
```

### Excluir Recursos (`delete`)

```bash
# Excluir um cluster
girus delete cluster

# Opções disponíveis:
  -f, --force    Força a exclusão sem confirmação
  -v, --verbose  Modo detalhado com output completo
```

## Fluxo de Trabalho Típico

1. **Criar um novo ambiente**:
   ```bash
   girus create cluster
   ```
   Isso irá:
   - Criar um cluster Kind
   - Configurar o namespace Girus
   - Implantar o backend e frontend
   - Configurar port-forwarding (8080 para backend, 8000 para frontend)
   - Abrir o navegador com a interface

2. **Verificar laboratórios disponíveis**:
   ```bash
   girus list labs
   ```

3. **Monitorar o ambiente**:
   ```bash
   girus list clusters
   ```

4. **Limpar o ambiente**:
   ```bash
   girus delete cluster
   ```

## Desenvolvimento

### Configuração do Ambiente
1. Fork o repositório
2. Clone localmente
3. Instale as dependências:
   ```bash
   go mod download
   ```

### Executando Testes
```bash
go test -v ./...
```

### Linting
```bash
golangci-lint run
```

### Build Local
```bash
go build -v -o girus -ldflags="-X 'github.com/linuxtips/girus/girus-cli/cmd.Version=dev'" ./main.go
```

## Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto é distribuído sob a licença GPLv3. Veja o arquivo `LICENSE` para mais detalhes.

## Suporte

- **Issues**: Use o [GitHub Issues](https://github.com/badtuxx/girus-cli/issues)
- **Discussões**: Participe das [Discussões no GitHub](https://github.com/badtuxx/girus-cli/discussions)
- **Documentação**: Visite nossa [Wiki](https://github.com/badtuxx/girus-cli/wiki)

## Mantenedores
- Jeferson Fernando ([@badtuxx](https://github.com/badtuxx))
- LINUXtips ([@linuxtips](https://github.com/linuxtips))

---

Feito com 💚 pela LINUXtips 