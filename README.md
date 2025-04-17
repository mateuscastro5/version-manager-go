# Git-Manager

Com esta ferramenta você pode controlar o versionamento Git dos seus projetos.

## Como instalar:

### Usando NPM
```
npm install @be-tech/git-manager
```
### Usando yarn
```
yarn add @be-tech/git-manager
```
### Usando PNPM
```
pnpm install @be-tech/git-manager
```

## Como usar:
Basta digitar e seguir as instruções:
```
git-manager
```

### Configuração de Tokens para Integração com GitHub/GitLab PARA RELEASES

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
- Node v20 ou superior
