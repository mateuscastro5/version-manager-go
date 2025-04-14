# version-manager (Go Version)

Com esta ferramenta você pode controlar o versionamento Git dos seus projetos.

## Como instalar:

### Usando Go Install
```
go install github.com/be-tech/version-manager@latest
```

### Compilando manualmente
```bash
# Clone o repositório
git clone https://github.com/be-tech/version-manager.git

# Entre no diretório
cd version-manager

# Compile o projeto
go build -o v-manager

# Mova o binário para um diretório no seu PATH (opcional)
sudo mv v-manager /usr/local/bin/
```

## Como usar:
Basta digitar:
```
v-manager
```

### Configuração de Tokens para Integração com GitHub/GitLab

Para criar releases no GitHub ou GitLab, você precisa configurar o token de acesso:

1. **Arquivo .env**: Crie um arquivo `.env` no diretório onde vai executar a ferramenta com o seguinte conteúdo:

   ```
   # Para GitHub
   GITHUB_TOKEN=seu_token_aqui
   
   # Para GitLab
   GITLAB_TOKEN=seu_token_aqui
   ```

2. **Variáveis de ambiente**: Alternativamente, defina as variáveis de ambiente diretamente:

   ```bash
   # Para GitHub
   export GITHUB_TOKEN=seu_token_aqui
   
   # Para GitLab
   export GITLAB_TOKEN=seu_token_aqui
   ```

3. **Permissões necessárias**:
   - Para GitHub: O token precisa ter permissão de `repo` completo para criar releases
   - Para GitLab: O token precisa ter permissão de `api` para criar releases

## Funcionalidades

- Seleção de repositório remoto
- Seleção de branches de origem e destino para merge
- Opção de push automático para o repositório remoto
- Remoção automática de branches de origem após merge
- Gerenciamento de tags de versão (major, minor, patch, pre-releases)
- Criação de releases no GitHub/GitLab
- Suporte a arquivos .env para configuração de tokens de acesso

## Requisitos
- Git instalado e configurado
- Go 1.16 ou superior (apenas para compilação)
