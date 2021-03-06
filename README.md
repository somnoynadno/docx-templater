# docx-templater

Простейшая утилита для генерации шаблонных docx-документов на основе [UniDoc](https://www.unidoc.io/).

Сервис представляет собой простейшее веб-приложение, использующее порт 4040.

Логика простая: сервис проходится по вашему docx-файлу и заменяет одни вхождения строк на другие.

## Конфигурация 

Для запуска требуются следующие переменные окружения:

- **LICENSE_KEY** - лицензионный ключ, полученный на https://cloud.unidoc.io/
- **TEMPLATES_FOLDER** - папка с документами-шаблонами
- **SAVE_FOLDER** - папка, в которую будут сохраняться созданные документы

## Старт приложения

### Классическая схема

1. Клонируем: ``` $ git clone https://github.com/somnoynadno/docx-templater```

2. Билдим: ``` $ cd docx-templater && go build -o main .```

3. Задаем переменные окружения, которые я перечислил выше (или хардкодим их в ```main.go```, если удобнее)

4. Запускаем бинарник: ``` $ ./main``` 

5. Создаем папку с шаблонами и добавляем в неё шаблонный документ

6. Создаем папку, в которую будут записываться созданные документы

7. Пробуем исполнить запрос и проверяем результат

### Альтернативная схема

Можно использовать любой понравившийся кусок кода из ```main.go``` или запустить
приложение в докере, используя ```Dockerfile```

## Использование 

### Создание документа

**POST** http://localhost:4040/docgen/create

```
{
    "Filename": "your_document.docx", // лежит в папке, заданной в TEMPLATES_FOLDER
    "Templates": {
        "{some_name}": "WHATEVER YOU WANT",
        "{start_another_name}": WILL BE REPLACED BY THIS"
    }
}
```

Если в вашем docx-файле имелись строки {some_name} и {start_another_name}, то они получат другие значения, которые вы им присвоили.

В ответе сервера вы получите имя созданного документа, по которому можно обратиться в следующем запросе.

### Получение документа

**GET** http://localhost:4040/docgen/static/<document_name>

В результате сервер выдаст вам сгенерированный по вашему шаблону файл.

## Лицензия

Это небольшая обертка над UniDoc API, которая успешно используется у меня на проде. 

Я просто решил поделиться.

Можете использовать этот код в любых своих целях.
