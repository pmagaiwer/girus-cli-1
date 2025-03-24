#!/bin/bash

set -e

echo "=== Instalando o Girus CLI ==="

# Detectar sistema operacional
OS="linux"
if [[ "$OSTYPE" == "darwin"* ]]; then
    OS="darwin"
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OSTYPE" == "cygwin" ]]; then
    OS="windows"
fi

echo "Sistema operacional detectado: $OS"

# Função para instalar Go
install_go() {
    echo "Instalando Go (última versão estável)..."
    
    # Obter a última versão estável do Go
    LATEST_GO_VERSION=$(curl -s https://golang.org/VERSION?m=text | head -n 1 | sed 's/go//')
    
    if [ "$OS" == "linux" ]; then
        # Linux
        wget -q https://golang.org/dl/go${LATEST_GO_VERSION}.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go${LATEST_GO_VERSION}.linux-amd64.tar.gz
        rm go${LATEST_GO_VERSION}.linux-amd64.tar.gz
        
        # Adicionar ao PATH temporariamente
        export PATH=$PATH:/usr/local/go/bin
        
        echo "Go instalado em /usr/local/go"
        echo "Para tornar permanente, adicione ao seu arquivo ~/.bashrc ou ~/.zshrc:"
        echo 'export PATH=$PATH:/usr/local/go/bin'
    
    elif [ "$OS" == "darwin" ]; then
        # MacOS
        if command -v brew &> /dev/null; then
            brew install go
        else
            wget -q https://golang.org/dl/go${LATEST_GO_VERSION}.darwin-amd64.pkg
            sudo installer -pkg go${LATEST_GO_VERSION}.darwin-amd64.pkg -target /
            rm go${LATEST_GO_VERSION}.darwin-amd64.pkg
            
            # Adicionar ao PATH temporariamente
            export PATH=$PATH:/usr/local/go/bin
            
            echo "Go instalado em /usr/local/go"
            echo "Para tornar permanente, adicione ao seu arquivo ~/.bash_profile ou ~/.zshrc:"
            echo 'export PATH=$PATH:/usr/local/go/bin'
        fi
    
    elif [ "$OS" == "windows" ]; then
        # Windows - sugestão para instalação manual
        echo "Instalação automática do Go não suportada no Windows."
        echo "Por favor, baixe e instale Go manualmente de https://golang.org/dl/"
        echo "Após a instalação, reabra o terminal e execute este script novamente."
        exit 1
    fi
    
    # Verificar a instalação
    if ! command -v go &> /dev/null; then
        echo "Falha ao instalar o Go. Por favor, instale manualmente."
        exit 1
    fi
    
    echo "Go ${LATEST_GO_VERSION} instalado com sucesso!"
}

# Função para instalar Kind
install_kind() {
    echo "Instalando Kind..."
    
    if [ "$OS" == "linux" ] || [ "$OS" == "darwin" ]; then
        # Linux/Mac
        curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-$(uname)-amd64
        chmod +x ./kind
        sudo mv ./kind /usr/local/bin/kind
    
    elif [ "$OS" == "windows" ]; then
        # Windows
        echo "Instalação automática do Kind não suportada no Windows."
        echo "Por favor, baixe e instale Kind manualmente:"
        echo "https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
        echo "Após a instalação, reabra o terminal e execute este script novamente."
        exit 1
    fi
    
    # Verificar a instalação
    if ! command -v kind &> /dev/null; then
        echo "Falha ao instalar o Kind. Por favor, instale manualmente."
        exit 1
    fi
    
    echo "Kind instalado com sucesso!"
}

# Função para instalar Kubectl
install_kubectl() {
    echo "Instalando Kubectl..."
    
    if [ "$OS" == "linux" ]; then
        # Linux
        curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
        chmod +x kubectl
        sudo mv kubectl /usr/local/bin/
    
    elif [ "$OS" == "darwin" ]; then
        # MacOS
        if command -v brew &> /dev/null; then
            brew install kubectl
        else
            curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
            chmod +x kubectl
            sudo mv kubectl /usr/local/bin/
        fi
    
    elif [ "$OS" == "windows" ]; then
        # Windows
        echo "Instalação automática do Kubectl não suportada no Windows."
        echo "Por favor, baixe e instale Kubectl manualmente:"
        echo "https://kubernetes.io/docs/tasks/tools/install-kubectl/"
        echo "Após a instalação, reabra o terminal e execute este script novamente."
        exit 1
    fi
    
    # Verificar a instalação
    if ! command -v kubectl &> /dev/null; then
        echo "Falha ao instalar o Kubectl. Por favor, instale manualmente."
        exit 1
    fi
    
    echo "Kubectl instalado com sucesso!"
}

# Verificar se Go está instalado
if ! command -v go &> /dev/null; then
    echo "Go não está instalado."
    read -p "Deseja instalar Go automaticamente? (S/n): " INSTALL_GO
    INSTALL_GO=${INSTALL_GO:-S}
    
    if [[ "$INSTALL_GO" =~ ^[Ss]$ ]]; then
        install_go
    else
        echo "Erro: Go é necessário para compilar o Girus CLI."
        echo "Visite: https://golang.org/doc/install"
        exit 1
    fi
else
    # Verificar versão do Go
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    MAJOR_VERSION=$(echo $GO_VERSION | cut -d. -f1)
    MINOR_VERSION=$(echo $GO_VERSION | cut -d. -f2)

    if [ "$MAJOR_VERSION" -lt 1 ] || ([ "$MAJOR_VERSION" -eq 1 ] && [ "$MINOR_VERSION" -lt 22 ]); then
        echo "Aviso: Versão do Go é $GO_VERSION, mas Go 1.22 ou superior é recomendado."
        read -p "Deseja continuar mesmo assim? (S/n): " CONTINUE_OLD_GO
        CONTINUE_OLD_GO=${CONTINUE_OLD_GO:-S}
        
        if [[ ! "$CONTINUE_OLD_GO" =~ ^[Ss]$ ]]; then
            echo "Instalação cancelada. Por favor, atualize o Go."
            exit 1
        fi
    else
        echo "✅ Go $GO_VERSION encontrado."
    fi
fi

# Verificar se o GOROOT está configurado
if [ -z "$GOROOT" ]; then
    # Tentar detectar o GOROOT
    if [ -d "/usr/local/go" ]; then
        echo "GOROOT não configurado. Configurando para /usr/local/go"
        export GOROOT=/usr/local/go
        export PATH=$GOROOT/bin:$PATH
    else
        echo "Aviso: GOROOT não está configurado. Se ocorrer erro de compilação, execute:"
        echo "export GOROOT=/caminho/para/go"
        echo "export PATH=\$GOROOT/bin:\$PATH"
    fi
fi

# Verificar se Kind está instalado
if ! command -v kind &> /dev/null; then
    echo "Kind não está instalado."
    read -p "Deseja instalar Kind automaticamente? (S/n): " INSTALL_KIND
    INSTALL_KIND=${INSTALL_KIND:-S}
    
    if [[ "$INSTALL_KIND" =~ ^[Ss]$ ]]; then
        install_kind
    else
        echo "Aviso: Kind não está instalado. É necessário para criar clusters."
        echo "Visite: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
    fi
else
    echo "✅ Kind encontrado."
fi

# Verificar se Kubectl está instalado
if ! command -v kubectl &> /dev/null; then
    echo "Kubectl não está instalado."
    read -p "Deseja instalar Kubectl automaticamente? (S/n): " INSTALL_KUBECTL
    INSTALL_KUBECTL=${INSTALL_KUBECTL:-S}
    
    if [[ "$INSTALL_KUBECTL" =~ ^[Ss]$ ]]; then
        install_kubectl
    else
        echo "Aviso: Kubectl não está instalado. É necessário para gerenciar clusters."
        echo "Visite: https://kubernetes.io/docs/tasks/tools/install-kubectl/"
    fi
else
    echo "✅ Kubectl encontrado."
fi

# Baixar dependências
echo "Baixando dependências..."
go mod tidy

# Compilar o projeto
echo "Compilando o Girus CLI..."
go build -o girus .

# Verificar se a compilação foi bem-sucedida
if [ ! -f "./girus" ]; then
    echo "Erro: A compilação falhou."
    exit 1
fi

echo "Girus CLI compilado com sucesso."

# Perguntar se o usuário deseja instalar o girus no PATH
if [ "$OS" != "windows" ]; then
    read -p "Deseja instalar o Girus CLI em /usr/local/bin? (S/n): " INSTALL_CHOICE
    INSTALL_CHOICE=${INSTALL_CHOICE:-S}

    if [[ "$INSTALL_CHOICE" =~ ^[Ss]$ ]]; then
        echo "Instalando o Girus CLI em /usr/local/bin..."
        sudo mv ./girus /usr/local/bin/
        echo "Instalação concluída! Você pode agora usar o comando 'girus' no terminal."
    else
        echo "O binário está disponível em $(pwd)/girus"
        echo "Você pode movê-lo manualmente para um diretório no PATH quando desejar."
    fi
else
    echo "O binário está disponível em $(pwd)/girus"
    echo "No Windows, recomendamos mover o arquivo para uma pasta no seu PATH manualmente."
fi

# Verificar se Docker está instalado e em execução
if ! command -v docker &> /dev/null; then
    echo "⚠️  Aviso: Docker não encontrado."
    echo "O Docker é necessário para criar clusters Kind e executar o Girus."
    echo "Por favor, instale o Docker adequado para seu sistema operacional:"
    echo "  - Linux: https://docs.docker.com/engine/install/"
    echo "  - macOS: https://docs.docker.com/desktop/install/mac-install/"
    echo "  - Windows: https://docs.docker.com/desktop/install/windows-install/"
else
    # Verificar se o Docker está em execução
    if ! docker info &> /dev/null; then
        echo "⚠️  Aviso: Docker não está em execução."
        echo "Inicie o Docker antes de criar clusters com o Girus."
    else
        echo "✅ Docker está instalado e em execução."
    fi
fi

echo "=== Instalação do Girus CLI concluída ==="
echo "Para começar, execute:"
echo "  girus create cluster    # Cria um cluster Kubernetes com Girus"
echo "  girus list labs         # Lista os laboratórios disponíveis"
echo "  girus help              # Mostra ajuda e comandos disponíveis" 