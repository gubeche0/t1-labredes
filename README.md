# Chat

## Descrição

Chat desenvolvido para a disciplina de Laboratório de Redes de Computadores.


### Definição das mensagens

Foram criadas 8 mensagens para o chat, o primeiro byte de cada mensagem é o tipo da mensagem, os tipos são definidos da seguinte forma:

- MESSAGE_TYPE_JOIN_REQUEST: Mensagem enviada pelo cliente para o servidor através da porta de controle para solicitar a disponibilidade de um nome de usuário e requisitar o acesso ao chat.
- MESSAGE_TYPE_JOIN_REQUEST_RESPONSE: Mensagem enviada pelo servidor para o cliente através da porta de controle para informar se o nome de usuário solicitado está disponível e se o cliente pode acessar o chat.
- MESSAGE_TYPE_JOIN: Mensagem enviada para o servidor através da porta de mensagem para se conectar ao chat com o nome de usuário escolhido anteriormente.
- MESSAGE_TYPE_JOIN_RESPONSE: Mensagem enviada para o cliente confirmando a conexão ao chat.
- MESSAGE_TYPE_TEXT: Mensagem enviada para o servidor através da porta de mensagem para enviar uma mensagem de texto para um usuario ou para todos os usuarios conectados.
- MESSAGE_TYPE_FILE: Mensagem enviada para o servidor através da porta de mensagem para enviar um arquivo para um usuario ou para todos os usuarios conectados.
- MESSAGE_TYPE_LIST_USERS: Mensagem enviada para o servidor através da porta de mensagem para solicitar a lista de usuarios conectados.
- MESSAGE_TYPE_LIST_USERS_RESPONSE: Mensagem enviada para o cliente através da porta de mensagem com a lista de usuarios conectados.

Os proximos 8 bytes de cada mensagem são o tamanho da mensagem e o restante é o conteúdo da mensagem.

### Estrutura das mensagens

- MESSAGE_TYPE_JOIN_REQUEST:
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nome da origem da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
- MESSAGE_TYPE_JOIN_REQUEST_RESPONSE:
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
  - 1 byte: 1 para nome de usuário disponível, 0 para nome de usuário indisponível
- MESSAGE_TYPE_JOIN:
  - byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nome da origem da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
- MESSAGE_TYPE_JOIN_RESPONSE:
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nome da origem da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
  - 1 byte: 1 para nome de usuário disponível, 0 para nome de usuário indisponível
- MESSAGE_TYPE_TEXT
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nome da origem da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
  - x bytes: Nome do destino da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
  - x bytes: Texto da mensagem
- MESSAGE_TYPE_FILE
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nome da origem da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
  - x bytes: Nome do destino da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
  - x bytes: Nome do arquivo
  - 1 byte: `\n` para indicar o fim do nome do arquivo
  - 8 bytes: Tamanho do arquivo
  - x bytes: Conteúdo do arquivo
- MESSAGE_TYPE_LIST_USERS
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nome da origem da mensagem
  - 1 byte: `\n` para indicar o fim do nome de usuário
- MESSAGE_TYPE_LIST_USERS_RESPONSE
  - 1 byte: Tipo da mensagem
  - 8 bytes: Tamanho da mensagem
  - x bytes: Nomes dos usuarios conectados separados por `\n`

## Uso

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

