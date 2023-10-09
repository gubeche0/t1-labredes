# Chat

Para enviar uma mensagem para todos os clientes conectados, basta digitar a mensagem e pressionar enter.

Para enviar para um cliente específico, digite o comando `/sendPrivate <user> <mensagem>`, onde `<user>` é o usuario do cliente e `<mensagem>` é a mensagem a ser enviada.

Para enviar um arquivo para todos os clientes conectados, digite o comando `/sendFile <caminho>`, onde `<caminho>` é o caminho do arquivo a ser enviado.

Para enviar um arquivo para um cliente específico, digite o comando `/sendFileTo <user> <caminho>`, onde `<user>` é o usuario do cliente e `<caminho>` é o caminho do arquivo a ser enviado.

Para listar os usuarios conectados, digite o comando `/listUsers`.

Para sair do chat, digite o comando `/exit

Para ver os comandos disponíveis, digite o comando `/help`.

Para limpar a tela, digite o comando `/clear`.
## Execução

### TCP
Execute o servidor TCP com o comando:

```bash
$ go run cmd/serverTCP/main.go -listen <ip>
```

Então execute os clientes TCP com o comando:

```bash
$ go run cmd/clientTCP/main.go -destination <ip>
```

### UDP
Execute o servidor UDP com o comando:

```bash
$ go run cmd/serverUDP/main.go -listen <ip>
```

Então execute os clientes TCP com o comando:

```bash
$ go run cmd/clientUDP/main.go -destination <ip>
```

