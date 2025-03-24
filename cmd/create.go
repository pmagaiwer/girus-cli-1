package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	deployFile     string
	clusterName    string
	verboseMode    bool
	useExternalFile bool
	labFile        string
	skipPortForward bool
	skipBrowser     bool
)

// defaultDeployment contém o YAML de deployment padrão do Girus
const defaultDeployment = `apiVersion: v1
kind: Namespace
metadata:
  name: girus
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: girus-sa
  namespace: girus
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: girus-role
  namespace: girus
rules:
  - apiGroups: [""]
    resources: ["pods", "pods/log", "pods/exec"]
    verbs: ["get", "list", "create", "delete", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: girus-cluster-role
rules:
  - apiGroups: [""]
    resources: ["namespaces"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
  - apiGroups: [""]
    resources: ["pods", "pods/log", "pods/exec"]
    verbs: ["get", "list", "watch", "create", "delete"]
  - apiGroups: ["apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: [""]
    resources: ["services", "configmaps", "secrets", "serviceaccounts"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
  - apiGroups: ["rbac.authorization.k8s.io"]
    resources: ["roles", "rolebindings"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: girus-cluster-rolebinding
subjects:
  - kind: ServiceAccount
    name: girus-sa
    namespace: girus
roleRef:
  kind: ClusterRole
  name: girus-cluster-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: girus-rolebinding
  namespace: girus
subjects:
  - kind: ServiceAccount
    name: girus-sa
    namespace: girus
roleRef:
  kind: Role
  name: girus-role
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: girus-config
  namespace: girus
data:
  config.yaml: |
    lab:
      defaultImage: "ubuntu:latest"
      podNamePrefix: "lab"
      containerName: "linux-lab"
      command:
        - "sleep"
        - "infinity"
      resources:
        cpuRequest: "100m"
        cpuLimit: "500m"
        memoryRequest: "64Mi"
        memoryLimit: "256Mi"
      envVars:
        TERM: "xterm-256color"
        SHELL: "/bin/bash"
        privileged: false
    # Outras configurações podem ser adicionadas aqui
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: basic-linux-lab
  namespace: girus
  labels:
    app: girus-lab-template
data:
  lab.yaml: |
    name: linux-basics
    title: "Introdução ao Linux"
    description: "Laboratório básico para praticar comandos Linux essenciais"
    duration: 30m
    tasks:
      - name: "Navegação básica"
        description: "Pratique comandos básicos de navegação"
        steps:
          - "Use 'pwd' para ver o diretório atual"
          - "Liste os arquivos com 'ls -la'"
          - "Crie um diretório chamado 'test' com 'mkdir test'"
        validation:
          - command: "test -d test"
            expectedOutput: ""
            errorMessage: "Diretório 'test' não foi criado"
      - name: "Manipulação de arquivos"
        description: "Aprenda a criar e manipular arquivos"
        steps:
          - "Entre no diretório 'test' com 'cd test'"
          - "Crie um arquivo chamado 'hello.txt' com 'echo hello > hello.txt'"
          - "Visualize o conteúdo com 'cat hello.txt'"
        validation:
          - command: "cat test/hello.txt"
            expectedOutput: "hello"
            errorMessage: "Arquivo 'hello.txt' não existe ou não tem o conteúdo esperado"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: girus-backend
  namespace: girus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: girus-backend
  template:
    metadata:
      labels:
        app: girus-backend
    spec:
      serviceAccountName: girus-sa
      containers:
        - name: backend
          image: linuxtips/girus-backend:0.1
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
            - name: GIN_MODE
              value: "release"
            - name: LAB_DEFAULT_IMAGE
              valueFrom:
                configMapKeyRef:
                  name: girus-config
                  key: lab.defaultImage
                  optional: true
          volumeMounts:
            - name: config-volume
              mountPath: /app/config
      volumes:
        - name: config-volume
          configMap:
            name: girus-config
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: girus-frontend
  namespace: girus
spec:
  replicas: 1
  selector:
    matchLabels:
      app: girus-frontend
  template:
    metadata:
      labels:
        app: girus-frontend
    spec:
      containers:
        - name: frontend
          image: linuxtips/girus-frontend:0.1
          imagePullPolicy: IfNotPresent
          ports:
            - containerPort: 80
          volumeMounts:
            - name: nginx-config
              mountPath: /etc/nginx/conf.d
      volumes:
        - name: nginx-config
          configMap:
            name: nginx-config
---
apiVersion: v1
kind: Service
metadata:
  name: girus-backend
  namespace: girus
spec:
  selector:
    app: girus-backend
  ports:
    - port: 8080
      targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: Service
metadata:
  name: girus-frontend
  namespace: girus
spec:
  selector:
    app: girus-frontend
  ports:
    - port: 80
      targetPort: 80
  type: ClusterIP
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-config
  namespace: girus
data:
  default.conf: |
    server {
        listen 80;
        server_name localhost;
        root /usr/share/nginx/html;
        index index.html;
        
        # Compressão
        gzip on;
        gzip_vary on;
        gzip_min_length 1000;
        gzip_proxied any;
        gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/xml+rss text/javascript;
        gzip_comp_level 6;
        
        # Cache para recursos estáticos
        location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
            expires 30d;
            add_header Cache-Control "public, no-transform";
        }
        
        # Redirecionar todas as requisições API para o backend
        location /api/ {
            proxy_pass http://girus-backend:8080/api/;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_buffering off;
            proxy_request_buffering off;
        }
        
        # Configuração para WebSockets
        location /ws/ {
            proxy_pass http://girus-backend:8080/ws/;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "upgrade";
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_read_timeout 86400;
        }
        
        # Configuração para React Router
        location / {
            try_files $uri $uri/ /index.html;
        }
    }
`

// waitForPodsReady espera até que os pods do Girus (backend e frontend) estejam prontos
func waitForPodsReady(namespace string, timeout time.Duration) error {
	fmt.Println("\nAguardando os pods do Girus inicializarem...")
	
	start := time.Now()
	bar := progressbar.NewOptions(100,
		progressbar.OptionSetDescription("Inicializando Girus..."),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)

	backendReady := false
	frontendReady := false
	backendMessage := ""
	frontendMessage := ""
	
	for {
		if time.Since(start) > timeout {
			bar.Finish()
			fmt.Println("\nStatus atual dos componentes:")
			if backendReady {
				fmt.Printf("✅ Backend: Pronto\n")
			} else {
				fmt.Printf("❌ Backend: %s\n", backendMessage)
			}
			if frontendReady {
				fmt.Printf("✅ Frontend: Pronto\n")
			} else {
				fmt.Printf("❌ Frontend: %s\n", frontendMessage)
			}
			return fmt.Errorf("timeout ao esperar pelos pods do Girus (5 minutos)")
		}

		// Verificar o backend
		if !backendReady {
			br, msg, err := getPodStatus(namespace, "app=girus-backend")
			if err == nil {
				backendReady = br
				backendMessage = msg
			}
		}

		// Verificar o frontend
		if !frontendReady {
			fr, msg, err := getPodStatus(namespace, "app=girus-frontend")
			if err == nil {
				frontendReady = fr
				frontendMessage = msg
			}
		}

		// Se ambos estiverem prontos, vamos verificar a conectividade
		if backendReady && frontendReady {
			// Verificar se conseguimos acessar a API
			healthy, err := checkHealthEndpoint()
			if err != nil || !healthy {
				// Ainda não está respondendo, vamos continuar tentando
				bar.Add(1)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			
			bar.Finish()
			fmt.Println("\n✅ Backend: Pronto")
			fmt.Println("✅ Frontend: Pronto")
			fmt.Println("✅ Aplicação: Respondendo")
			return nil
		}

		bar.Add(1)
		time.Sleep(500 * time.Millisecond)
	}
}

// getPodStatus verifica o status de um pod e retorna uma mensagem descritiva
func getPodStatus(namespace, selector string) (bool, string, error) {
	// Verificar se o pod existe
	cmd := exec.Command("kubectl", "get", "pods", "-n", namespace, "-l", selector, "-o", "jsonpath={.items[0].metadata.name}")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	err := cmd.Run()
	if err != nil {
		return false, "Pod não encontrado", err
	}
	
	podName := strings.TrimSpace(out.String())
	if podName == "" {
		return false, "Pod ainda não criado", nil
	}
	
	// Verificar a fase atual do pod
	phaseCmd := exec.Command("kubectl", "get", "pod", podName, "-n", namespace, "-o", "jsonpath={.status.phase}")
	var phaseOut bytes.Buffer
	phaseCmd.Stdout = &phaseOut
	
	err = phaseCmd.Run()
	if err != nil {
		return false, "Erro ao verificar status", err
	}
	
	phase := strings.TrimSpace(phaseOut.String())
	if phase != "Running" {
		return false, fmt.Sprintf("Status: %s", phase), nil
	}
	
	// Verificar se todos os containers estão prontos
	readyCmd := exec.Command("kubectl", "get", "pod", podName, "-n", namespace, "-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
	var readyOut bytes.Buffer
	readyCmd.Stdout = &readyOut
	
	err = readyCmd.Run()
	if err != nil {
		return false, "Erro ao verificar prontidão", err
	}
	
	readyStatus := strings.TrimSpace(readyOut.String())
	if readyStatus != "True" {
		return false, "Containers inicializando", nil
	}
	
	return true, "Pronto", nil
}

// checkHealthEndpoint verifica se a aplicação está respondendo ao endpoint de saúde
func checkHealthEndpoint() (bool, error) {
	// Verificar o mapeamento de porta do serviço
	cmd := exec.Command("kubectl", "get", "svc", "-n", "girus", "girus-backend", "-o", "jsonpath={.spec.ports[0].nodePort}")
	var out bytes.Buffer
	cmd.Stdout = &out
	
	err := cmd.Run()
	if err != nil {
		// Tentar verificar diretamente o endpoint interno se não encontrarmos o NodePort
		healthCmd := exec.Command("kubectl", "exec", "-n", "girus", "deploy/girus-backend", "--", "wget", "-q", "-O-", "-T", "2", "http://localhost:8080/api/v1/health")
		return healthCmd.Run() == nil, nil
	}
	
	nodePort := strings.TrimSpace(out.String())
	if nodePort == "" {
		// Porta não encontrada, tentar verificar o serviço internamente
		healthCmd := exec.Command("kubectl", "exec", "-n", "girus", "deploy/girus-backend", "--", "wget", "-q", "-O-", "-T", "2", "http://localhost:8080/api/v1/health")
		return healthCmd.Run() == nil, nil
	}
	
	// Tentar acessar via NodePort
	healthCmd := exec.Command("curl", "-s", "-o", "/dev/null", "-w", "%{http_code}", fmt.Sprintf("http://localhost:%s/api/v1/health", nodePort))
	var healthOut bytes.Buffer
	healthCmd.Stdout = &healthOut
	
	err = healthCmd.Run()
	if err != nil {
		return false, err
	}
	
	statusCode := strings.TrimSpace(healthOut.String())
	return statusCode == "200", nil
}

// setupPortForward configura port-forward para os serviços do Girus
func setupPortForward(namespace string) error {
	// Verificar se as portas já estão em uso
	if portInUse(8000) {
		// Se porta 8000 está em uso, verificar se é nosso port-forward
		killProcess := exec.Command("sh", "-c", "lsof -ti:8000 | xargs -r kill -9")
		killProcess.Run() // Ignorar erros aqui
	}
	
	if portInUse(8080) {
		// Se porta 8080 está em uso, verificar se é nosso port-forward
		killProcess := exec.Command("sh", "-c", "lsof -ti:8080 | xargs -r kill -9")
		killProcess.Run() // Ignorar erros aqui
	}
	
	// Criar um arquivo para descartar saídas de erro
	devNull, _ := os.Open(os.DevNull)
	defer devNull.Close()
	
	// Criar um diretório para os arquivos PID se não existir
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}
	pidDir := filepath.Join(homeDir, ".girus")
	os.MkdirAll(pidDir, 0755) // Criar diretório se não existir
	
	// Arquivos para armazenar PIDs
	backendPidFile := filepath.Join(pidDir, "backend.pid")
	frontendPidFile := filepath.Join(pidDir, "frontend.pid")
	
	// Port-forward para o backend (8080) em background
	backendCmd := exec.Command("kubectl", "port-forward", "-n", namespace, "svc/girus-backend", "8080:8080")
	backendCmd.Stderr = devNull
	backendCmd.Stdout = devNull
	if err := backendCmd.Start(); err != nil {
		return fmt.Errorf("erro ao configurar port-forward para o backend: %v", err)
	}
	
	// Salvar PID do processo de backend
	if err := os.WriteFile(backendPidFile, []byte(fmt.Sprintf("%d", backendCmd.Process.Pid)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Aviso: não foi possível salvar o PID do backend: %v\n", err)
	}
	
	// Verificar se o backend está acessível
	time.Sleep(2 * time.Second) // Dar tempo para o port-forward inicializar
	backendOk := false
	for i := 0; i < 5; i++ { // Tentar algumas vezes
		healthCmd := exec.Command("curl", "-s", "--max-time", "1", "http://localhost:8080/api/v1/health")
		if healthCmd.Run() == nil {
			backendOk = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	
	if !backendOk {
		return fmt.Errorf("não foi possível conectar ao backend após configurar port-forward")
	}
	
	// Port-forward para o frontend (8000) em background
	frontendCmd := exec.Command("kubectl", "port-forward", "-n", namespace, "svc/girus-frontend", "8000:80")
	frontendCmd.Stderr = devNull
	frontendCmd.Stdout = devNull
	if err := frontendCmd.Start(); err != nil {
		// Se falhar, tentar matar o processo do backend antes de retornar
		backendCmd.Process.Kill()
		return fmt.Errorf("erro ao configurar port-forward para o frontend: %v", err)
	}
	
	// Salvar PID do processo de frontend
	if err := os.WriteFile(frontendPidFile, []byte(fmt.Sprintf("%d", frontendCmd.Process.Pid)), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Aviso: não foi possível salvar o PID do frontend: %v\n", err)
	}
	
	// Verificar se o frontend está acessível
	time.Sleep(2 * time.Second) // Dar tempo para o port-forward inicializar
	frontendOk := false
	for i := 0; i < 5; i++ { // Tentar algumas vezes
		checkCmd := exec.Command("curl", "-s", "--max-time", "1", "-o", "/dev/null", "-w", "%{http_code}", "http://localhost:8000")
		var out bytes.Buffer
		checkCmd.Stdout = &out
		if err := checkCmd.Run(); err == nil {
			statusCode := strings.TrimSpace(out.String())
			if statusCode == "200" || statusCode == "301" || statusCode == "302" {
				frontendOk = true
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
	
	if !frontendOk {
		// Se frontend falhar, matar o backend também
		backendCmd.Process.Kill()
		frontendCmd.Process.Kill()
		return fmt.Errorf("não foi possível conectar ao frontend após configurar port-forward")
	}
	
	return nil
}

// portInUse verifica se uma porta está em uso
func portInUse(port int) bool {
	checkCmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port))
	return checkCmd.Run() == nil
}

// openBrowser abre o navegador com a URL especificada
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("não foi possível abrir o navegador (sistema operacional não suportado)")
	}

	return cmd.Start()
}

var createCmd = &cobra.Command{
	Use:   "create [subcommand]",
	Short: "Comandos para criar recursos",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cria o cluster Girus",
	Long: `Cria um cluster Kind com o nome "girus" e implanta todos os componentes necessários.
Por padrão, o deployment embutido no binário é utilizado.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Verificar silenciosamente se o cluster já existe
		checkCmd := exec.Command("kind", "get", "clusters")
		output, err := checkCmd.Output()
		
		// Ignorar erros na checagem, apenas assumimos que não há clusters
		if err == nil {
			clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
			
			// Verificar se o cluster "girus" já existe
			clusterExists := false
			for _, cluster := range clusters {
				if cluster == clusterName {
					clusterExists = true
					break
				}
			}
			
			if clusterExists {
				fmt.Printf("⚠️  Cluster Girus já existe.\n")
				fmt.Print("Deseja substituí-lo? [s/N]: ")
				
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.ToLower(strings.TrimSpace(response))
				
				if response != "s" && response != "sim" && response != "y" && response != "yes" {
					fmt.Println("Operação cancelada.")
					return
				}
				
				// Excluir o cluster existente
				fmt.Printf("Excluindo cluster Girus existente...\n")
				
				deleteCmd := exec.Command("kind", "delete", "cluster", "--name", clusterName)
				if verboseMode {
					deleteCmd.Stdout = os.Stdout
					deleteCmd.Stderr = os.Stderr
					if err := deleteCmd.Run(); err != nil {
						fmt.Fprintf(os.Stderr, "❌ Erro ao excluir o cluster existente: %v\n", err)
						fmt.Println("   Por favor, exclua manualmente com 'kind delete cluster --name girus' e tente novamente.")
						os.Exit(1)
					}
				} else {
					// Usar barra de progresso
					bar := progressbar.NewOptions(100,
						progressbar.OptionSetDescription("Excluindo cluster existente..."),
						progressbar.OptionSetWidth(50),
						progressbar.OptionShowBytes(false),
						progressbar.OptionSetPredictTime(false),
						progressbar.OptionThrottle(65*time.Millisecond),
						progressbar.OptionShowCount(),
						progressbar.OptionSpinnerType(14),
						progressbar.OptionFullWidth(),
					)
					
					var stderr bytes.Buffer
					deleteCmd.Stderr = &stderr
					
					// Iniciar o comando
					err := deleteCmd.Start()
					if err != nil {
						fmt.Fprintf(os.Stderr, "❌ Erro ao iniciar exclusão: %v\n", err)
						os.Exit(1)
					}
					
					// Atualizar a barra de progresso
					done := make(chan struct{})
					go func() {
						for {
							select {
							case <-done:
								return
							default:
								bar.Add(1)
								time.Sleep(100 * time.Millisecond)
							}
						}
					}()
					
					// Aguardar o final do comando
					err = deleteCmd.Wait()
					close(done)
					bar.Finish()
					
					if err != nil {
						fmt.Fprintf(os.Stderr, "❌ Erro ao excluir o cluster existente: %v\n", err)
						fmt.Println("   Detalhes técnicos:", stderr.String())
						fmt.Println("   Por favor, exclua manualmente com 'kind delete cluster --name girus' e tente novamente.")
						os.Exit(1)
					}
				}
				
				fmt.Println("✅ Cluster existente excluído com sucesso.")
			}
		}
		
		// Criar o cluster Kind
		fmt.Println("🔄 Criando cluster Girus...")

		if verboseMode {
			// Executar normalmente mostrando o output
			createClusterCmd := exec.Command("kind", "create", "cluster", "--name", clusterName)
			createClusterCmd.Stdout = os.Stdout
			createClusterCmd.Stderr = os.Stderr

			if err := createClusterCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao criar o cluster Girus: %v\n", err)
				fmt.Println("   Possíveis causas:")
				fmt.Println("   • Docker não está em execução")
				fmt.Println("   • Permissões insuficientes")
				fmt.Println("   • Conflito com cluster existente")
				os.Exit(1)
			}
		} else {
			// Usando barra de progresso (padrão)
			bar := progressbar.NewOptions(100,
				progressbar.OptionSetDescription("Criando cluster..."),
				progressbar.OptionSetWidth(50),
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetPredictTime(false),
				progressbar.OptionThrottle(65*time.Millisecond),
				progressbar.OptionShowCount(),
				progressbar.OptionSpinnerType(14),
				progressbar.OptionFullWidth(),
			)

			// Executar comando sem mostrar saída
			createClusterCmd := exec.Command("kind", "create", "cluster", "--name", clusterName)
			var stderr bytes.Buffer
			createClusterCmd.Stderr = &stderr
			
			// Iniciar o comando
			err := createClusterCmd.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao iniciar o comando: %v\n", err)
				os.Exit(1)
			}

			// Atualizar a barra de progresso enquanto o comando está em execução
			done := make(chan struct{})
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						bar.Add(1)
						time.Sleep(200 * time.Millisecond)
					}
				}
			}()

			// Aguardar o final do comando
			err = createClusterCmd.Wait()
			close(done)
			bar.Finish()

			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao criar o cluster Girus: %v\n", err)
				
				// Traduzir mensagens de erro comuns
				errMsg := stderr.String()
				
				if strings.Contains(errMsg, "node(s) already exist for a cluster with the name") {
					fmt.Println("   Erro: Já existe um cluster com o nome 'girus' no sistema.")
					fmt.Println("   Por favor, exclua-o primeiro com 'kind delete cluster --name girus'")
				} else if strings.Contains(errMsg, "permission denied") {
					fmt.Println("   Erro: Permissão negada. Verifique as permissões do Docker.")
				} else if strings.Contains(errMsg, "Cannot connect to the Docker daemon") {
					fmt.Println("   Erro: Não foi possível conectar ao serviço Docker.")
					fmt.Println("   Verifique se o Docker está em execução com 'systemctl status docker'")
				} else {
					fmt.Println("   Detalhes técnicos:", errMsg)
				}
				
				os.Exit(1)
			}
		}

		fmt.Println("✅ Cluster Girus criado com sucesso!")

		// Aplicar o manifesto de deployment do Girus
		fmt.Println("\n📦 Implantando o Girus no cluster...")

		// Determinar se vamos usar o arquivo externo ou o deployment embutido
		deployContent := defaultDeployment
		if useExternalFile {
			// Verificar se o arquivo de deploy existe
			if _, err := os.Stat(deployFile); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "❌ Erro: arquivo de deployment '%s' não encontrado\n", deployFile)
				os.Exit(1)
			}
			
			// Ler o conteúdo do arquivo
			content, err := os.ReadFile(deployFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao ler arquivo de deployment: %v\n", err)
				os.Exit(1)
			}
			deployContent = string(content)
		}

		// Criar um arquivo temporário para o deployment
		tempFile, err := os.CreateTemp("", "girus-deploy-*.yaml")
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Erro ao criar arquivo temporário: %v\n", err)
			os.Exit(1)
		}
		defer os.Remove(tempFile.Name()) // Limpar o arquivo temporário ao finalizar

		// Escrever o conteúdo no arquivo temporário
		if _, err := tempFile.WriteString(deployContent); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Erro ao escrever no arquivo temporário: %v\n", err)
			os.Exit(1)
		}
		tempFile.Close()

		if verboseMode {
			// Executar normalmente mostrando o output
			applyCmd := exec.Command("kubectl", "apply", "-f", tempFile.Name())
			applyCmd.Stdout = os.Stdout
			applyCmd.Stderr = os.Stderr

			if err := applyCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao aplicar o manifesto do Girus: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Usar barra de progresso para o deploy (padrão)
			bar := progressbar.NewOptions(100,
				progressbar.OptionSetDescription("Implantando Girus..."),
				progressbar.OptionSetWidth(50),
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetPredictTime(false),
				progressbar.OptionThrottle(65*time.Millisecond),
				progressbar.OptionShowCount(),
				progressbar.OptionSpinnerType(14),
				progressbar.OptionFullWidth(),
			)

			// Executar comando sem mostrar saída
			applyCmd := exec.Command("kubectl", "apply", "-f", tempFile.Name())
			var stderr bytes.Buffer
			applyCmd.Stderr = &stderr
			
			// Iniciar o comando
			err := applyCmd.Start()
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao iniciar o comando: %v\n", err)
				os.Exit(1)
			}

			// Atualizar a barra de progresso enquanto o comando está em execução
			done := make(chan struct{})
			go func() {
				for {
					select {
					case <-done:
						return
					default:
						bar.Add(1)
						time.Sleep(100 * time.Millisecond)
					}
				}
			}()

			// Aguardar o final do comando
			err = applyCmd.Wait()
			close(done)
			bar.Finish()

			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ Erro ao aplicar o manifesto do Girus: %v\n", err)
				fmt.Println("   Detalhes técnicos:", stderr.String())
				os.Exit(1)
			}
		}

		// Aguardar os pods do Girus ficarem prontos
		if err := waitForPodsReady("girus", 5*time.Minute); err != nil {
			fmt.Fprintf(os.Stderr, "Aviso: %v\n", err)
			fmt.Println("Recomenda-se verificar o estado dos pods com 'kubectl get pods -n girus'")
		} else {
			fmt.Println("Todos os componentes do Girus estão prontos e em execução!")
		}

		fmt.Println("Girus implantado com sucesso no cluster!")

		// Configurar port-forward automaticamente (a menos que --skip-port-forward tenha sido especificado)
		if !skipPortForward {
			fmt.Print("\n🔌 Configurando acesso aos serviços do Girus... ")
			
			if err := setupPortForward("girus"); err != nil {
				fmt.Println("⚠️")
				fmt.Printf("Não foi possível configurar o acesso automático: %v\n", err)
				fmt.Println("\nVocê pode tentar configurar manualmente com os comandos:")
				fmt.Println("kubectl port-forward -n girus svc/girus-backend 8080:8080")
				fmt.Println("kubectl port-forward -n girus svc/girus-frontend 8000:80")
			} else {
				fmt.Println("✅")
				fmt.Println("Acesso configurado com sucesso!")
				fmt.Println("📊 Backend: http://localhost:8080")
				fmt.Println("🖥️  Frontend: http://localhost:8000")
				
				// Abrir o navegador se não foi especificado para pular
				if !skipBrowser {
					fmt.Println("\n🌐 Abrindo navegador com o Girus...")
					if err := openBrowser("http://localhost:8000"); err != nil {
						fmt.Printf("⚠️  Não foi possível abrir o navegador: %v\n", err)
						fmt.Println("   Acesse manualmente: http://localhost:8000")
					}
				}
			}
		} else {
			fmt.Println("\n⏩ Port-forward ignorado conforme solicitado")
			fmt.Println("\nPara acessar o Girus posteriormente, execute:")
			fmt.Println("kubectl port-forward -n girus svc/girus-backend 8080:8080")
			fmt.Println("kubectl port-forward -n girus svc/girus-frontend 8000:80")
		}
		
		// Exibir mensagem de conclusão
		fmt.Println("\n" + strings.Repeat("─", 60))
		fmt.Println("✅ GIRUS PRONTO PARA USO!")
		fmt.Println(strings.Repeat("─", 60))
		
		// Sugerir a adição de templates de laboratório
		fmt.Println("📋 PRÓXIMOS PASSOS:")
		fmt.Println("  • Para aplicar os templates de laboratório, execute:")
		fmt.Println("    ./update-labs.sh")
		fmt.Println("\n  • Ou aplique manualmente:")
		fmt.Println("    kubectl apply -f ./labs/")
		fmt.Println(strings.Repeat("─", 60))
	},
}

var createLabCmd = &cobra.Command{
	Use:   "lab [lab-id] ou -f [arquivo]",
	Short: "Cria um novo laboratório no Girus",
	Long:  "Adiciona um novo laboratório ao Girus a partir de um arquivo de manifesto ConfigMap, ou cria um ambiente de laboratório a partir de um ID de template existente.\nOs templates de laboratório são armazenados no diretório /labs na raiz do projeto.",
	Run: func(cmd *cobra.Command, args []string) {
		// Verificar qual modo estamos
		if labFile != "" {
			// Modo de adicionar template a partir de arquivo
			addLabFromFile(labFile, verboseMode)
		} else {
			fmt.Fprintf(os.Stderr, "Erro: Você deve especificar um arquivo de laboratório com a flag -f\n")
			fmt.Println("\nExemplo:")
			fmt.Println("  girus create lab -f meulaboratorio.yaml      # Adiciona um novo template a partir do arquivo")
			fmt.Println("  girus create lab -f /home/user/REPOS/strigus/labs/basic-linux.yaml      # Adiciona um template do diretório /labs")
			os.Exit(1)
		}
	},
}

// addLabFromFile adiciona um novo template de laboratório a partir de um arquivo
func addLabFromFile(labFile string, verboseMode bool) {
	// Verificar se o arquivo existe
	if _, err := os.Stat(labFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "❌ Erro: arquivo '%s' não encontrado\n", labFile)
		os.Exit(1)
	}

	fmt.Println("🔍 Verificando ambiente Girus...")
	
	// Verificar se há um cluster Girus ativo
	checkCmd := exec.Command("kubectl", "get", "namespace", "girus", "--no-headers", "--ignore-not-found")
	checkOutput, err := checkCmd.Output()
	if err != nil || !strings.Contains(string(checkOutput), "girus") {
		fmt.Fprintf(os.Stderr, "❌ Nenhum cluster Girus ativo encontrado\n")
		fmt.Println("   Use 'girus create cluster' para criar um cluster ou 'girus list clusters' para ver os disponíveis.")
		os.Exit(1)
	}

	// Verificar o pod do backend (silenciosamente, só mostra mensagem em caso de erro)
	backendCmd := exec.Command("kubectl", "get", "pods", "-n", "girus", "-l", "app=girus-backend", "-o", "jsonpath={.items[0].status.phase}")
	backendOutput, err := backendCmd.Output()
	if err != nil || string(backendOutput) != "Running" {
		fmt.Fprintf(os.Stderr, "❌ O backend do Girus não está em execução\n")
		fmt.Println("   Verifique o status dos pods com 'kubectl get pods -n girus'")
		os.Exit(1)
	}

	// Ler o arquivo para verificar se é um ConfigMap válido
	content, err := os.ReadFile(labFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Erro ao ler o arquivo '%s': %v\n", labFile, err)
		os.Exit(1)
	}

	// Verificação simples se o arquivo parece ser um ConfigMap válido
	fileContent := string(content)
	if !strings.Contains(fileContent, "kind: ConfigMap") ||
		!strings.Contains(fileContent, "app: girus-lab-template") {
		fmt.Fprintf(os.Stderr, "❌ O arquivo não é um manifesto de laboratório válido\n")
		fmt.Println("   O arquivo deve ser um ConfigMap com a label 'app: girus-lab-template'")
		os.Exit(1)
	}

	fmt.Printf("📦 Processando laboratório: %s\n", labFile)

	// Aplicar o ConfigMap no cluster usando kubectl apply
	if verboseMode {
		fmt.Println("   Aplicando ConfigMap no cluster...")
	}
	
	// Aplicar o ConfigMap no cluster
	if verboseMode {
		// Executar normalmente mostrando o output
		applyCmd := exec.Command("kubectl", "apply", "-f", labFile)
		applyCmd.Stdout = os.Stdout
		applyCmd.Stderr = os.Stderr
		if err := applyCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Erro ao aplicar o laboratório: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Usar barra de progresso
		bar := progressbar.NewOptions(100,
			progressbar.OptionSetDescription("   Aplicando laboratório"),
			progressbar.OptionSetWidth(50),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
		)

		// Executar comando sem mostrar saída
		applyCmd := exec.Command("kubectl", "apply", "-f", labFile)
		var stderr bytes.Buffer
		applyCmd.Stderr = &stderr
		
		// Iniciar o comando
		err := applyCmd.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Erro ao iniciar o comando: %v\n", err)
			os.Exit(1)
		}

		// Atualizar a barra de progresso enquanto o comando está em execução
		done := make(chan struct{})
		go func() {
			for {
				select {
				case <-done:
					return
				default:
					bar.Add(1)
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()

		// Aguardar o final do comando
		err = applyCmd.Wait()
		close(done)
		bar.Finish()

		if err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Erro ao aplicar o laboratório: %v\n", err)
			if verboseMode {
				fmt.Fprintf(os.Stderr, "   Detalhes: %s\n", stderr.String())
			}
			os.Exit(1)
		}
	}

	// Extrair o ID do lab (name) do arquivo YAML para mostrar na mensagem
	var labID string
	// Procurar pela linha 'name:' dentro do bloco lab.yaml:
	labNameCmd := exec.Command("sh", "-c", fmt.Sprintf("grep -A10 'lab.yaml:' %s | grep 'name:' | head -1", labFile))
	labNameOutput, err := labNameCmd.Output()
	if err == nil {
		nameLine := strings.TrimSpace(string(labNameOutput))
		parts := strings.SplitN(nameLine, "name:", 2)
		if len(parts) >= 2 {
			labID = strings.TrimSpace(parts[1])
		}
	}
	
	// Extrair também o título para exibição
	var labTitle string
	labTitleCmd := exec.Command("sh", "-c", fmt.Sprintf("grep -A10 'lab.yaml:' %s | grep 'title:' | head -1", labFile))
	labTitleOutput, err := labTitleCmd.Output()
	if err == nil {
		titleLine := strings.TrimSpace(string(labTitleOutput))
		parts := strings.SplitN(titleLine, "title:", 2)
		if len(parts) >= 2 {
			labTitle = strings.TrimSpace(parts[1])
			labTitle = strings.Trim(labTitle, "\"'")
		}
	}
	
	fmt.Println("\n🔄 Reiniciando backend para carregar o template...")
	
	// O backend apenas carrega os templates na inicialização
	if verboseMode {
		// Mostrar o output da reinicialização
		fmt.Println("   (O backend do Girus carrega os templates apenas na inicialização)")
		restartCmd := exec.Command("kubectl", "rollout", "restart", "deployment/girus-backend", "-n", "girus")
		restartCmd.Stdout = os.Stdout
		restartCmd.Stderr = os.Stderr
		if err := restartCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Erro ao reiniciar o backend: %v\n", err)
			fmt.Println("   O template foi aplicado, mas pode ser necessário reiniciar o backend manualmente:")
			fmt.Println("   kubectl rollout restart deployment/girus-backend -n girus")
		}
		
		// Aguardar o reinício completar
		fmt.Println("   Aguardando o reinício do backend completar...")
		waitCmd := exec.Command("kubectl", "rollout", "status", "deployment/girus-backend", "-n", "girus", "--timeout=60s")
		waitCmd.Stdout = os.Stdout
		waitCmd.Stderr = os.Stderr
		waitCmd.Run()
	} else {
		// Usar barra de progresso
		bar := progressbar.NewOptions(100,
			progressbar.OptionSetDescription("   Reiniciando backend"),
			progressbar.OptionSetWidth(50),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetPredictTime(false),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionSpinnerType(14),
			progressbar.OptionFullWidth(),
		)
		
		// Reiniciar o deployment do backend
		restartCmd := exec.Command("kubectl", "rollout", "restart", "deployment/girus-backend", "-n", "girus")
		var stderr bytes.Buffer
		restartCmd.Stderr = &stderr
		
		err := restartCmd.Run()
		if err != nil {
			bar.Finish()
			fmt.Fprintf(os.Stderr, "\n⚠️  Erro ao reiniciar o backend: %v\n", err)
			if verboseMode {
				fmt.Fprintf(os.Stderr, "   Detalhes: %s\n", stderr.String())
			}
			fmt.Println("   O template foi aplicado, mas pode ser necessário reiniciar o backend manualmente:")
			fmt.Println("   kubectl rollout restart deployment/girus-backend -n girus")
		} else {
			// Aguardar o reinício completar
			waitCmd := exec.Command("kubectl", "rollout", "status", "deployment/girus-backend", "-n", "girus", "--timeout=60s")
			
			// Iniciar o comando
			err := waitCmd.Start()
			if err != nil {
				bar.Finish()
				fmt.Fprintf(os.Stderr, "\n⚠️  Erro ao verificar status do reinício: %v\n", err)
			} else {
				// Atualizar a barra de progresso enquanto o comando está em execução
				done := make(chan struct{})
				go func() {
					for {
						select {
						case <-done:
							return
						default:
							bar.Add(1)
							time.Sleep(100 * time.Millisecond)
						}
					}
				}()
				
				// Aguardar o final do comando
				waitCmd.Wait()
				close(done)
			}
			bar.Finish()
		}
	}
	
	// Aguardar mais alguns segundos para que o backend reinicie completamente
	fmt.Println("   Aguardando inicialização completa...")
	time.Sleep(3 * time.Second)
	
	// Após reiniciar o backend, verificar se precisamos recriar o port-forward
	portForwardStatus := checkPortForwardNeeded()
	
	// Se port-forward é necessário, configurá-lo silenciosamente sem perguntar ao usuário
	if portForwardStatus {
		fmt.Print("🔌 Reconfigurando port-forward... ")
		
		// Iniciar o port-forward em uma goroutine
		errChan := make(chan error, 1)
		go func() {
			// Redirecionar saída do setupPortForward para evitar mensagens de erro no terminal
			oldStdout := os.Stdout
			oldStderr := os.Stderr
			os.Stdout, _ = os.Open(os.DevNull)
			os.Stderr, _ = os.Open(os.DevNull)
			defer func() {
				os.Stdout = oldStdout
				os.Stderr = oldStderr
			}()
			
			errChan <- setupPortForward("girus")
		}()
		
		// Aguardar um pouco para o port-forward ser estabelecido
		time.Sleep(2 * time.Second)
		
		// Verificar se houve erro (com timeout para não bloquear)
		var pfSuccess bool = true
		select {
		case err := <-errChan:
			if err != nil {
				fmt.Println("⚠️  Não foi possível configurar automaticamente")
				pfSuccess = false
			} else {
				fmt.Println("✅ Configurado!")
			}
		default:
			// Timeout - assumimos que está ok
			fmt.Println("✅ Configurado!")
		}
		
		// Verificar se a conexão com o backend está funcionando
		if pfSuccess {
			healthCheck := exec.Command("curl", "-s", "--max-time", "2", "http://localhost:8080/api/v1/health")
			if healthCheck.Run() != nil {
				fmt.Println("⚠️  O backend parece não estar respondendo na porta 8080")
				fmt.Println("   Se necessário, configure manualmente com:")
				fmt.Println("   kubectl port-forward -n girus svc/girus-backend 8080:8080")
				fmt.Println("   kubectl port-forward -n girus svc/girus-frontend 8000:80")
			}
		}
	}
	
	// Desenhar uma linha separadora
	fmt.Println("\n" + strings.Repeat("─", 60))
	
	// Exibir informações sobre o laboratório adicionado
	fmt.Println("✅ LABORATÓRIO ADICIONADO COM SUCESSO!")
	
	if labTitle != "" && labID != "" {
		fmt.Printf("\n📚 Título: %s\n", labTitle)
		fmt.Printf("🏷️  ID: %s\n", labID)
	} else if labID != "" {
		fmt.Printf("\n🏷️  ID do Laboratório: %s\n", labID)
	}

	fmt.Println("\n📋 PRÓXIMOS PASSOS:")
	fmt.Println("  • Para ver todos os laboratórios disponíveis:")
	fmt.Println("    girus list labs")
	
	fmt.Println("\n  • Para verificar detalhes do template adicionado:")
	if labID != "" {
		fmt.Printf("    kubectl describe configmap -n girus | grep -A20 %s\n", labID)
	} else {
		fmt.Println("    kubectl get configmaps -n girus -l app=girus-lab-template")
		fmt.Println("    kubectl describe configmap <nome-do-configmap> -n girus")
	}
	
	// Linha final
	fmt.Println(strings.Repeat("─", 60))
}

// checkPortForwardNeeded verifica se o port-forward para o backend precisa ser reconfigurado
func checkPortForwardNeeded() bool {
	// Verificar se a porta 8080 está em uso
	portCheckCmd := exec.Command("lsof", "-i", ":8080")
	if portCheckCmd.Run() != nil {
		// Porta 8080 não está em uso, precisamos de port-forward
		return true
	}
	
	// Porta está em uso, mas precisamos verificar se é o kubectl port-forward e se está funcional
	// Verificar se o processo é kubectl port-forward
	processCmd := exec.Command("sh", "-c", "ps -eo pid,cmd | grep 'kubectl port-forward' | grep '8080' | grep -v grep")
	if processCmd.Run() != nil {
		// Não encontrou processo de port-forward ativo ou válido
		return true
	}
	
	// Verificar se a conexão com o backend está funcionando
	healthCmd := exec.Command("curl", "-s", "--head", "--max-time", "2", "http://localhost:8080/api/v1/health")
	return healthCmd.Run() != nil // Retorna true (precisa de port-forward) se o comando falhar
}

func init() {
	createCmd.AddCommand(createClusterCmd)
	createCmd.AddCommand(createLabCmd)

	// Flags para createClusterCmd
	createClusterCmd.Flags().StringVarP(&deployFile, "file", "f", "", "Arquivo YAML para deployment do Girus (opcional)")
	createClusterCmd.Flags().BoolVarP(&verboseMode, "verbose", "v", false, "Modo detalhado com output completo em vez da barra de progresso")
	createClusterCmd.Flags().BoolVarP(&skipPortForward, "skip-port-forward", "", false, "Não perguntar sobre configurar port-forwarding")
	createClusterCmd.Flags().BoolVarP(&skipBrowser, "skip-browser", "", false, "Não abrir o navegador automaticamente")

	// Flags para createLabCmd
	createLabCmd.Flags().StringVarP(&labFile, "file", "f", "", "Arquivo de manifesto do laboratório (ConfigMap)")
	createLabCmd.Flags().BoolVarP(&verboseMode, "verbose", "v", false, "Modo detalhado com output completo em vez da barra de progresso")

	// definir o nome do cluster como "girus" sempre
	clusterName = "girus"
} 