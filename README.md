# Consulta CEP
#### Versão: 1.0

#### **Descrição:** Consulta CEP é uma ferramenta de linha de comando que permite consultar duas APIs diferentes (BrasilAPI e ViaCEP) para obter informações sobre um determinado código postal brasileiro (CEP).

## Uso:

* Execute o programa e insira um CEP válido de 8 dígitos quando solicitado.
* O programa consultará as duas APIs e exibirá os resultados em formato JSON.
* Se ocorrer um erro, o programa exibirá uma mensagem de erro.
* Para sair do programa, pressione CTRL+C.

## Recursos:

* Consulta duas APIs diferentes (BrasilAPI e ViaCEP) para obter informações sobre um CEP
* Exibe os resultados em formato JSON
* Trata erros e timeouts
* Permite sair do programa com CTRL+C

## Dependências:

* Go 1.17 ou posterior
* Pacote github.com/common-nighthawk/go-figure para arte ASCII

## Construção e Execução:

* Clone o repositório e navegue até o diretório do projeto.
* Execute go build main.go para construir o programa.
* Execute ./main para executar o programa.

## Código-fonte:

O código-fonte está disponível no arquivo main.go e é escrito em Go.

## Agradecimentos:

Agradecemos à equipe do **BrasilAPI** e **ViaCEP** por fornecer as APIs utilizadas neste projeto.