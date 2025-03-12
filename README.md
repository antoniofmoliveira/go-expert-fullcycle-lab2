# Lab2

## Solução

* implementação de projeto contendo o docker-compose.yaml para executar todo o projeto
  * otel/opentelemetry-collector:0.63.1
  * zipkin/zipkin:latest
  * jaegertracing/all-in-one:latest
  * servicoa
    * subprojeto `servicoa`
      * POST  <http://localhost:8080/cep>
      * serviço simples validador de cep
      * repassa cep validado ao `servicob`
      * retorna resultado da consulta ou erro de validação
  * servicob
    * subprojeto `servicob`
      * GET <http://servicob:8081/?cep={{cep}}>
      * serviço que consulta cep e temperatura
      * consulta cep
      * extrai cidade da reposta
      * consulta temperatura da cidade
      * retorna clima da cidade em °C, °F e °K e o nome da cidade ou erro
  * execução
    * docker-compose up no folder raiz
    * POST no endereço <http://localhost:8080>:
      * usando (Rest Client)<https://marketplace.visualstudio.com/items?itemName=humao.rest-client> (arquivo pronto em `servicoa/api/post.http`)

        ```http
        POST http://localhost:8080/cep HTTP/1.1
        Host: localhost:8080
        Content-Type: application/json

        {
            "cep":"39408078"
        }
        ```

    * respostas

        ```http
        HTTP/1.1 200 OK
        Content-Type: application/json
        Date: Wed, 12 Mar 2025 11:33:50 GMT
        Content-Length: 68
        Connection: close

        {
        "city": "Montes Claros",
        "temp_C": 25.1,
        "temp_F": 77.18,
        "temp_K": 298.1
        }
        ```

        ```http
        HTTP/1.1 422 Unprocessable Entity
        Content-Type: text/plain; charset=utf-8
        X-Content-Type-Options: nosniff
        Date: Wed, 12 Mar 2025 11:35:48 GMT
        Content-Length: 16
        Connection: close

        invalid zipcode
        ```

  * traces
    * jaeger <http://localhost:16686/> (selecionar Service 'servicoa')
    * zipkin <http://localhost:9411/> (query `serviceName=servicoa`)
  * observações:
    * quando usando docker o endereço do zipkin deve ser o do container (não pode usar localhost) no código
    * o serviço externo de clima precisa escapar o nome da cidade
    * para tudo funcionar como previsto é preciso encadear os contextos
    * o código para otel-collector está defasado e só funciona com a versão 0.63.1
    * para as versões mais recentes do otel-collector o código retorna erro com conexão recusada
    * o serviço externo de cep não retorna 404 quando não encontra o cep. Retorna 200 com página html de erro

## Objetivo

Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

## Requisitos - Serviço A (responsável pelo input)

* O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }
* O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING
  * Caso seja válido, será encaminhado para o Serviço B via HTTP
  * Caso não seja válido, deve retornar:
    * Código HTTP: 422
    * Mensagem: invalid zipcode

## Requisitos - Serviço B (responsável pela orquestração)

* O sistema deve receber um CEP válido de 8 digitos
* O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.
* O sistema deve responder adequadamente nos seguintes cenários:
  * Em caso de sucesso:
    * Código HTTP: 200
    * Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }
  * Em caso de falha, caso o CEP não seja válido (com formato correto):
    * Código HTTP: 422
    * Mensagem: invalid zipcode
​​​  * Em caso de falha, caso o CEP não seja encontrado:
    * Código HTTP: 404
    * Mensagem: can not find zipcode

Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

* Implementar tracing distribuído entre Serviço A - Serviço B
* Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura

## Dicas

* Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: <https://viacep.com.br/>
* Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: <https://www.weatherapi.com/>
* Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32
* Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273
  * Sendo F = Fahrenheit
  * Sendo C = Celsius
  * Sendo K = Kelvin
* Para dúvidas da implementação do OTEL, você pode clicar aqui
* Para implementação de spans, você pode clicar aqui
* Você precisará utilizar um serviço de collector do OTEL
* Para mais informações sobre Zipkin, você pode clicar aqui

## Entrega

* O código-fonte completo da implementação.
* Documentação explicando como rodar o projeto em ambiente dev.
* Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.
