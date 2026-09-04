# Анализатор размера диска (Go)

[![hexlet-check](https://github.com/kimchi-lover/go-from-scratch-project-242/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/kimchi-lover/go-from-scratch-project-242/actions)

`hexlet-path-size` — консольная утилита, которая показывает размер файла или директории. Аналог `du`, но с предсказуемым выводом: одна строка вида `<размер>\t<путь>`.

Учебный проект Хекслета: https://ru.hexlet.io/programs/go-from-scratch

## Возможности

- Размер одного файла или суммарный размер файлов в директории.
- Рекурсивный обход вложенных директорий (`-r`).
- Человекочитаемый формат: `KB`, `MB`, `GB` и далее, основание 1024 (`-H`).
- Учёт скрытых файлов и директорий, имя которых начинается с точки (`-a`).
- Экспортируемая функция `GetPathSize`, которую можно использовать как библиотеку.

## Установка

```bash
git clone https://github.com/kimchi-lover/go-from-scratch-project-242.git
cd go-from-scratch-project-242
make build
```

Бинарник появится в `bin/hexlet-path-size`.

## Использование

```
hexlet-path-size [global options] <path>
```

| Флаг | Описание |
|---|---|
| `--recursive`, `-r` | считать вложенные директории |
| `--human`, `-H` | выводить размер в удобных единицах |
| `--all`, `-a` | учитывать скрытые файлы и директории |
| `--help`, `-h` | показать справку |

Демонстрация:

[![asciicast](https://asciinema.org/a/IcPUM5jT5AhGjgQu.svg)](https://asciinema.org/a/IcPUM5jT5AhGjgQu)

Примеры на тестовых данных из репозитория:

```bash
# размер файла
./bin/hexlet-path-size testdata/test.txt
6B	testdata/test.txt

# директория: только файлы первого уровня
./bin/hexlet-path-size testdata/dir
14B	testdata/dir

# директория с вложенными
./bin/hexlet-path-size -r testdata/dir
25B	testdata/dir

# скрытые файлы не учитываются без -a
./bin/hexlet-path-size testdata/hidden
4B  testdata/hidden

./bin/hexlet-path-size -r -a testdata/hidden
13B	testdata/hidden

# человекочитаемый формат
./bin/hexlet-path-size -H bin/hexlet-path-size
5.2MB	bin/hexlet-path-size
```


## Использование как библиотеки

Модуль называется `code`. Функция принимает путь и три флага и возвращает готовую строку с размером:

```go
import "code"

size, err := code.GetPathSize("testdata/dir", true, false, false) // recursive, human, all
if err != nil {
	// обработать ошибку
}
fmt.Println(size) // 25B
```

## Разработка

```bash
make build     # собрать бинарник в bin/
make test      # go test -v -race -cover ./...
make lint      # golangci-lint run
make lint-fix  # golangci-lint run --fix
```

Запись для README лежит в `docs/demo.cast`, сценарий в `docs/demo.sh`. Перезаписать и загрузить заново:

```bash
asciinema rec -f asciicast-v2 -c "sh docs/demo.sh" --overwrite docs/demo.cast && asciinema upload docs/demo.cast
```

Тестовые данные лежат в `testdata` и имеют фиксированные размеры. Файл `.gitattributes` запрещает git менять переводы строк в них, иначе размеры разошлись бы между системами.

---

<details>
<summary>Автоматические тесты Хекслета</summary>

Тесты запускаются на каждый коммит. За запуск отвечает файл `.github/workflows/hexlet-check.yml` — не удаляйте и не переименовывайте ни его, ни репозиторий.

</details>

## О Хекслете

[Хекслет](https://ru.hexlet.io/) — школа программирования: авторские программы обучения с практикой, поддержкой наставников и реальными проектами, которые остаются в резюме. Этот репозиторий — один из таких проектов.
