#!/usr/bin/env bash
set -e

# ASCII Art Banner para o Girus
echo ""
cat << "EOF"
   ██████╗ ██╗██████╗ ██╗   ██╗███████╗
  ██╔════╝ ██║██╔══██╗██║   ██║██╔════╝
  ██║  ███╗██║██████╔╝██║   ██║███████╗
  ██║   ██║██║██╔══██╗██║   ██║╚════██║
  ╚██████╔╝██║██║  ██║╚██████╔╝███████║
   ╚═════╝ ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚══════╝
EOF

# Configurações e variáveis
GIRUS_CODENAME="Maracatu"
GIRUS_VERSION="0.1.0"
KIND_VERSION="0.27.0"
DOWNLOAD_TOOL="none"
ORIGINAL_DIR=$(pwd)
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

# Verificar se o script está sendo executado como root (sudo)
if [ "$(id -u)" -eq 0 ]; then
    echo "❌ ERRO: Este script não deve ser executado como root ou com sudo."
    echo "   Por favor, execute sem sudo. O script solicitará elevação quando necessário."
    exit 1
fi

# Verificar se o terminal é interativo
IS_INTERACTIVE=0
if [ -t 0 ]; then
    IS_INTERACTIVE=1
fi

# Forçar modo interativo para o script completo
IS_INTERACTIVE=1

# Função para verificar se o comando curl ou wget está disponível
check_download_tool() {
    echo -e "\nVerificando se curl ou wget estão disponíveis..."
    if command -v curl &> /dev/null; then
        DOWNLOAD_TOOL="curl"
        echo "O comando curl está disponível e será utilizado."
    elif command -v wget &> /dev/null; then
        DOWNLOAD_TOOL="wget"
        echo "O comando wget está disponível e será utilizado."
    else
        echo "curl ou wget não disponíveis. Por favor, instale um deles e tente novamente."
        exit 1
    fi
}

# Função para pedir confirmação ao usuário (interativo) ou mostrar ação padrão (não-interativo)
ask_user() {
    local prompt="$1"
    local default="$2"
    local variable_name="$3"
    
    # Modo sempre interativo - perguntar ao usuário
    read -rp "$prompt" response
    # Se resposta for vazia, usar o padrão
    response=${response:-$default}
    
    # Exportar a resposta para a variável solicitada
    eval "$variable_name=\"$response\""
}

# Função para verificar a arquitetura do sistema
check_arch() {
    echo -e "\nVerificando a arquitetura do sistema..."
    ARCH="$(uname -m)"
    if [ "$ARCH" == "x86_64" ] || [ "$ARCH" == "amd64" ]; then
        ARCH="amd64"
    elif [ "$ARCH" == "aarch64" ] || [ "$ARCH" == "arm64" ]; then
        ARCH="arm64"
    else
        echo "❌ Arquitetura não suportada: $(uname -m)"
        exit 1
    fi
    echo "Arquitetura detectada: $ARCH"
}

# Função para verificar o sistema operacional
check_os() {
    echo -e "\nVerificando tipo de sistema operacional..."
    OS="$(uname -s)"
    if [ "$OS" == "Linux" ]; then
        OS="linux"
        BINARY_URL="https://github.com/badtuxx/girus-cli/releases/download/v$GIRUS_VERSION/girus-$OS-$ARCH"
    elif [ "$OS" == "Darwin" ]; then
        OS="darwin"
        BINARY_URL="https://github.com/badtuxx/girus-cli/releases/download/v$GIRUS_VERSION/girus-$OS-$ARCH"
    elif [ "$OS" == "CYGWIN" ] || [ "$OS" == "MINGW" ] || [ "$OS" == "MSYS" ]; then
        OS="windows"
        BINARY_URL="https://github.com/badtuxx/girus-cli/releases/download/v$GIRUS_VERSION/girus-$OS-$ARCH.exe"
    else
        echo "❌ Sistema operacional não suportado: $(uname -s)"
        exit 1
    fi
    echo "Sistema operacional detectado: $OS"
    echo "URL de download: $BINARY_URL"
}

# Função para instalar Docker
install_docker() {
    echo "Instalando Docker..."
    
    if [ "$OS" == "linux" ]; then
        # Linux (script de conveniência do Docker)
        echo "Baixando o script de instalação do Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        echo "Executando o script de instalação (será solicitada senha de administrador)..."
        sudo sh get-docker.sh
        
        # Adicionar usuário atual ao grupo docker
        echo "Adicionando usuário atual ao grupo docker..."
        sudo usermod -aG docker "$USER"
        
        # Iniciar o serviço
        echo "Iniciando o serviço Docker..."
        sudo systemctl enable --now docker
        
        # Limpar arquivo de instalação
        rm get-docker.sh
    
    elif [ "$OS" == "darwin" ]; then
        # MacOS
        echo "No macOS, o Docker Desktop precisa ser instalado manualmente."
        echo "Por favor, baixe e instale o Docker Desktop para Mac:"
        echo "https://docs.docker.com/desktop/mac/install/"
        echo "Após a instalação, reinicie seu terminal e execute este script novamente."
        exit 1
    
    elif [ "$OS" == "windows" ]; then
        # Windows
        echo "No Windows, o Docker Desktop precisa ser instalado manualmente."
        echo "Por favor, baixe e instale o Docker Desktop para Windows:"
        echo "https://docs.docker.com/desktop/windows/install/"
        echo "Após a instalação, reabra o terminal e execute este script novamente."
        exit 1
    fi
    
    # Verificar a instalação
    if ! command -v docker &> /dev/null; then
        echo "❌ Falha ao instalar o Docker."
        echo "Por favor, instale manualmente seguindo as instruções em https://docs.docker.com/engine/install/"
        exit 1
    fi
    
    echo "Docker instalado com sucesso!"
    echo "NOTA: Pode ser necessário reiniciar seu sistema ou fazer logout/login para que as permissões de grupo sejam aplicadas."
}

# Função para instalar Kind
install_kind() {
    echo "Instalando Kind..."
    
    if [ "$OS" == "linux" ] || [ "$OS" == "darwin" ]; then
        # Linux/Mac
        curl --progress-bar -Lo ./kind "https://kind.sigs.k8s.io/dl/v$KIND_VERSION/kind-$(uname)-amd64"
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
        curl --progress-bar -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
        chmod +x kubectl
        sudo mv kubectl /usr/local/bin/
    
    elif [ "$OS" == "darwin" ]; then
        # MacOS
        if command -v brew &> /dev/null; then
            brew install kubectl
        else
            curl --progress-bar -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl"
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

# Função para verificar se o Docker está em execução
check_docker_running() {
    if docker info &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# Verificar se o Girus CLI está no PATH
check_girus_in_path() {
    if command -v girus &> /dev/null; then
        return 0
    else
        return 1
    fi
}

# Função para verificar instalações anteriores do Girus CLI
check_previous_install() {
    echo -e "\nVerificando e limpando instalações anteriores..."
    local previous_install_found=false
    local install_locations=(
        "/usr/local/bin/girus"
        "/usr/bin/girus"
        "$HOME/.local/bin/girus"
        "./girus"
    )
    
    # Verificar instalações anteriores
    for location in "${install_locations[@]}"; do
        if [ -f "$location" ]; then
            echo "⚠️ Instalação anterior encontrada em: $location"
            previous_install_found=true
        fi
    done
    
    # Se uma instalação anterior foi encontrada, perguntar sobre limpeza
    if [ "$previous_install_found" = true ]; then
        ask_user "Deseja remover a(s) instalação(ões) anterior(es)? (S/n): " "S" "CLEAN_INSTALL"
        if [[ "$CLEAN_INSTALL" =~ ^[Ss]$ ]]; then
            echo "🧹 Removendo instalações anteriores..."
            for location in "${install_locations[@]}"; do
                if [ -f "$location" ]; then
                    echo "Removendo $location"
                    if [[ "$location" == "/usr/local/bin/girus" || "$location" == "/usr/bin/girus" ]]; then
                        sudo rm -f "$location"
                    else
                        rm -f "$location"
                    fi
                fi
            done
            echo "✅ Limpeza concluída."
        else
            echo "Continuando com a instalação sem remover versões anteriores."
        fi
    else
        echo "✅ Nenhuma instalação anterior do Girus CLI encontrada."
    fi
}

# Função para baixar e instalar o binário
download_and_install() {
    echo "📥 Baixando o Girus CLI versão $GIRUS_VERSION para $OS-$ARCH..."
    cd "$TEMP_DIR"
    
    if [ "$DOWNLOAD_TOOL" == "curl" ]; then
        echo "Usando curl para download de: $BINARY_URL"
        echo "Executando: curl -L --progress-bar \"$BINARY_URL\" -o girus"
        if ! curl -L --progress-bar "$BINARY_URL" -o girus; then
            echo "❌ Erro no curl. Tentando com opções de debug..."
            curl -L -v "$BINARY_URL" -o girus
        fi
    elif [ "$DOWNLOAD_TOOL" == "wget" ]; then
        echo "Usando wget para download de: $BINARY_URL"
        echo "Executando: wget --show-progress -q \"$BINARY_URL\" -O girus"
        if ! wget --show-progress -q "$BINARY_URL" -O girus; then
            echo "❌ Erro no wget. Tentando com opções de debug..."
            wget -v "$BINARY_URL" -O girus
        fi
    else
        echo "❌ Erro: curl ou wget não encontrados. Por favor, instale um deles e tente novamente."
        exit 1
    fi
    
    # Verificar se o download foi bem-sucedido
    if [ ! -f girus ] || [ ! -s girus ]; then
        echo "❌ Erro: Falha ao baixar o Girus CLI."
        echo "URL: $BINARY_URL"
        echo "Verifique sua conexão com a internet e se a versão $GIRUS_VERSION está disponível."
            exit 1
        fi
    
    # Tornar o binário executável
    chmod +x girus
    
    # Perguntar se o usuário deseja instalar no PATH
    echo "🔧 Girus CLI baixado com sucesso."
    ask_user "Deseja instalar o Girus CLI em /usr/local/bin? (S/n): " "S" "INSTALL_GLOBALLY"
    
    if [[ "$INSTALL_GLOBALLY" =~ ^[Ss]$ ]]; then
        echo "📋 Instalando o Girus CLI em /usr/local/bin/girus..."
        sudo mv girus /usr/local/bin/
        echo "✅ Girus CLI instalado com sucesso em /usr/local/bin/girus"
        echo "   Você pode executá-lo de qualquer lugar com o comando 'girus'"
    else
        # Copiar para o diretório original
        cp girus "$ORIGINAL_DIR/"
        echo "✅ Girus CLI copiado para o diretório atual: $(realpath "$ORIGINAL_DIR/girus")"
        echo "   Você pode executá-lo com: './girus'"
    fi
}

# Verificar se todas as dependências estão instaladas
verify_all_dependencies() {
    local all_deps_ok=true
    
    # Verificar Docker
    if command -v docker &> /dev/null && check_docker_running; then
        echo "✅ Docker está instalado e em execução."
    else
        echo "❌ Docker não está instalado ou não está em execução."
        all_deps_ok=false
    fi
    
    # Verificar Kind
    if command -v kind &> /dev/null; then
        echo "✅ Kind está instalado."
    else
        echo "❌ Kind não está instalado."
        all_deps_ok=false
    fi
    
    # Verificar Kubectl
    if command -v kubectl &> /dev/null; then
        echo "✅ Kubectl está instalado."
    else
        echo "❌ Kubectl não está instalado."
        all_deps_ok=false
    fi
    
    # Verificar Girus CLI
    if check_girus_in_path; then
        echo "✅ Girus CLI está instalado e disponível no PATH."
    else
        echo "⚠️ Girus CLI não está disponível no PATH."
        all_deps_ok=false
    fi
    
    return "$( [ "$all_deps_ok" = true ] && echo 0 || echo 1 )"
}

#---------------------------------------------------------------------------------
# Iniciar mensagem principal
#---------------------------------------------------------------------------------
echo -e "\nScript de Instalação - Versão v$GIRUS_VERSION - Codename: $GIRUS_CODENAME"
echo -e "\n=== Iniciando instalação do Girus CLI ==="
# Verificar se as ferramentas para download estão disponíveis
check_download_tool

# Verificar tipo de sistema operacional e arquitetura
check_arch
check_os

# Verificar e limpar instalações anteriores
check_previous_install

# ETAPA 1: Verificar pré-requisitos - Docker
echo -e "\n=== ETAPA 1: Verificando Docker ==="
if ! command -v docker &> /dev/null; then
    echo "Docker não está instalado."
    ask_user "Deseja instalar Docker automaticamente? (Linux apenas) (S/n): " "S" "INSTALL_DOCKER"
    
    if [[ "$INSTALL_DOCKER" =~ ^[Ss]$ ]]; then
        install_docker
    else
        echo "⚠️ Aviso: Docker é necessário para criar clusters Kind e executar o Girus."
        echo "Por favor, instale o Docker adequado para seu sistema operacional:"
        echo " - Linux: https://docs.docker.com/engine/install/"
        echo " - macOS: https://docs.docker.com/desktop/install/mac-install/"
        echo " - Windows: https://docs.docker.com/desktop/install/windows-install/"
        exit 1
    fi
else
    # Verificar se o Docker está em execução
    if ! docker info &> /dev/null; then
        echo "⚠️ Aviso: Docker está instalado, mas não está em execução."
        ask_user "Deseja tentar iniciar o Docker? (S/n): " "S" "START_DOCKER"
        
        if [[ "$START_DOCKER" =~ ^[Ss]$ ]]; then
            echo "Tentando iniciar o Docker..."
            if [ "$OS" == "linux" ]; then
                sudo systemctl start docker
                # Verificar novamente
                if ! docker info &> /dev/null; then
                    echo "❌ Falha ao iniciar o Docker. Por favor, inicie manualmente com 'sudo systemctl start docker'"
                    exit 1
                fi
            else
                echo "No macOS/Windows, inicie o Docker Desktop manualmente e execute este script novamente."
                exit 1
            fi
        else
            echo "❌ Erro: Docker precisa estar em execução para usar o Girus. Por favor, inicie-o e tente novamente."
            exit 1
        fi
    fi
    echo "✅ Docker está instalado e em execução."
fi

# ETAPA 2: Verificar pré-requisitos - Kind
echo -e "\n=== ETAPA 2: Verificando Kind ==="
if ! command -v kind &> /dev/null; then
    echo "Kind não está instalado."
    ask_user "Deseja instalar Kind automaticamente? (S/n): " "S" "INSTALL_KIND"
    
    if [[ "$INSTALL_KIND" =~ ^[Ss]$ ]]; then
        install_kind
    else
        echo "⚠️ Aviso: Kind é necessário para criar clusters Kubernetes e executar o Girus."
        echo "Você pode instalá-lo manualmente seguindo as instruções em: https://kind.sigs.k8s.io/docs/user/quick-start/#installation"
        exit 1
    fi
else
    echo "✅ Kind já está instalado."
fi

# ETAPA 3: Verificar pré-requisitos - Kubectl
echo -e "\n=== ETAPA 3: Verificando Kubectl ==="
if ! command -v kubectl &> /dev/null; then
    echo "Kubectl não está instalado."
    ask_user "Deseja instalar Kubectl automaticamente? (S/n): " "S" "INSTALL_KUBECTL"
    
    if [[ "$INSTALL_KUBECTL" =~ ^[Ss]$ ]]; then
        install_kubectl
    else
        echo "⚠️ Aviso: Kubectl é necessário para interagir com o cluster Kubernetes."
        echo "Você pode instalá-lo manualmente seguindo as instruções em: https://kubernetes.io/docs/tasks/tools/install-kubectl/"
    exit 1
    fi
else
    echo "✅ Kubectl já está instalado."
fi

# ETAPA 4: Baixar e instalar o Girus CLI
echo -e "\n=== ETAPA 4: Instalando Girus CLI ==="
echo "URL de download: $BINARY_URL"
download_and_install

# Voltar para o diretório original
cd "$ORIGINAL_DIR"

# Mensagem final de conclusão
echo -e "\n===== INSTALAÇÃO CONCLUÍDA =====\n"

# Verificar todas as dependências
verify_all_dependencies
echo ""

# Exibir instruções para próximos passos
cat << EOF
📝 PRÓXIMOS PASSOS:

1. Para criar um novo cluster Kubernetes e instalar o Girus:
   $ girus create cluster

2. Após a criação do cluster, acesse o Girus no navegador:
   http://localhost:8000

3. No navegador, inicie o laboratório Linux de boas-vindas para conhecer 
   a plataforma e começar sua experiência com o Girus!

Obrigado por instalar o Girus CLI!
EOF

exit 0 