# SetsBuilder

Программа строит множества **СЛЕД** и **ПЕРВ** для символов грамматики:

![image](https://github.com/DanArmor/taEkzUtils/assets/39347109/993590a9-8505-497d-a0bd-2e319dd4c302)


Также она строит множества **ВЫБОР** для правил и находит пересечения в них:

![image](https://github.com/DanArmor/taEkzUtils/assets/39347109/e27a644c-7315-4eaf-b8cd-fcfc55aece76)

### Требования к вводу:
* В папке, где происходит запуск, должен находиться файл `rules.json` с описаниями грамматики.
* Первый символ, с которого начинается вывод - `S`
* Все нетерминалы и терминалы должны состоять из одного символа - недопустимы нетерминалы/терминалы вида `S'` / `a'` и т. п.

Пример содержимого `rules.json`:
```json
{
    "rules": [
        {
            "left": "S",
            "right": "S;O"
        },
        {
            "left": "S",
            "right": "Z;"
        },
        {
            "left": "O",
            "right": "Y[S]"
        },
        {
            "left": "O",
            "right": "Y[S][S]"
        },
        {
            "left": "O",
            "right": "{[S]Y}"
        },
        {
            "left": "O",
            "right": "{Y[S]}"
        },
        {
            "left": "O",
            "right": "a=Y"
        },
        {
            "left": "Y",
            "right": "(Y|Y)"
        },
        {
            "left": "Y",
            "right": "(Y&Y)"
        },
        {
            "left": "Y",
            "right": "!(Y)"
        },
        {
            "left": "Y",
            "right": "a"
        },
        {
            "left": "Z",
            "right": "O"
        }
    ]
}
```
