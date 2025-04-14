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

## Funcionalidades

- Seleção de repositório remoto
- Seleção de branches de origem e destino para merge
- Opção de push automático para o repositório remoto
- Remoção automática de branches de origem após merge
- Gerenciamento de tags de versão (major, minor, patch, pre-releases)

## Requisitos
- Git instalado e configurado
- Go 1.16 ou superior (apenas para compilação)
