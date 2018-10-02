# Implementation

## CODE
**VERSION "1.0.0"**

### Installation

go get -u github.com/dink10/app_test

go build -o app_test

### Usage

```
Usage: converter [options...]

Options:
  -f      	Input csv file.
  -e    	Allow empty values in the input csv file (default: false)
  -d    	Allow any header type as is (default: false)
  -cpus         Number of used cpu cores (default for current machine is %d cores)
```

#### Examples

./app_test -f resources/test.csv

./app_test -f resources/test2.csv -e

./app_test -f resources/test2.csv -d

### Benchmarks

Сomparison was made with a library with big amount of stars [github.com/olekukonko/tablewriter](github.com/olekukonko/tablewriter) -> it's first result as BenchmarkThirdPartyTable-8 

```
BenchmarkThirdPartyTable-8        100000             66489 ns/op            7586 B/op        389 allocs/op
BenchmarkTable-8                  500000              8393 ns/op            1760 B/op         65 allocs/op
```

### Run tests

cd table/

go test

go test -run none -bench . -benchtime 3s -benchmem

## SHELL

cd bash/

cat anagram.input | ./anagram.sh > result.txt

## DB

```
2. alter table Table1 add key (ID2);
9. alter table Table2 add key (ID1,ID3);
```

Вариант индекса #2 по ID 2 имеет смысл использовать, т.к. по нему происходит выборка значения 
в условном выражении `T1.ID2 BETWEEN 600 AND 700`, результат которой будет в дальнейшем использоваться 
в следующем условии `T1.ID1 & 3 = 0`. Т.к. в этом условии побитовая операция, но не вижу необходимости 
создавать индекс на `T1.ID1`.
Для выполнения операции `INNER JOIN` механизмом базы данных, даже простейшим Nested Loop (или его производными)
для записей из `T1` выбираются записи `T2`. Как раз тут и придется к месту индекс #9 по (ID1,ID3).
А уже в дальнейшем для фильтрации будет использоваться условие `T2.ID3 BETWEEN 600 AND 700`

------------

# Tasks

## CODE

Разработать конвертор из CSV файла в таблицу из ASCII символов.
Первая строка файла задает типы столбцов.
Следующие строки - сами данные ( разделитель - точка с запятой ).

Типы:
int - целое число ( выравнивание вправо )
string - строка, строковые данные бьются на слова и выводятся в столбик.
money - денежная единица, форматирование 2 занака после запятой и
разделитель разрядов - пробел.

Исходные данные ( как пример ):
```
int;string;money
1;aaa bbb ccc;1000.33
5;aaaa bbb;0.01
13;aa bbbb;10000.00
На выходе скрипта:
```

На выходе скрипта:

```
+-----------------+
| 1|aaa | 1 000,33|
|  |bbb |         |
|  |ccc |         |
+--+----+---------+
| 5|aaaa|     0,01|
|  |bbb |         |
+--+----+---------+
|13|aa  |10 000,00|
|  |bbbb|         |
+--+----+---------+
```

Критерии оценки:
1. Читабельность кода
2. Готовность к production использованию

## SHELL

Решить задачу c подсчетом анаграмм в текстовом файле с помощью утилит командной строки.
Условия: Дан файл где каждая строка это одно слово, на высоте надо получить другой файл где будет анаграмма и кол-во раз которое она встречается в исходном файле.

## DB

Допустим, в MySQL есть две таблицы Table1 и Table2 с одинаковым списком полей:

```sql
ID1 int(11) NOT NULL default 0,
ID2 int(11) NOT NULL default 0,
ID3 int(11) NOT NULL default 0,
Value varchar(255) NOT NULL default ''
```

В обеих таблицах по миллиону записей;
значения полей ID1, ID2, ID3 равномерно распределены между 0 и 1000;
ключей в таблицах нет.

Эти таблицы участвуют только в запросе вида:

```sql
SELECT
        T1.ID3,
        COUNT(*),
        COUNT(DISTINCT T1.ID1,T1.ID2),
        SUM(T1.ID1+T2.ID2),
        CONCAT(T1.Value,T2.Value)
FROM
        Table1 T1
INNER JOIN
        Table2 T2 USING(ID1,ID3)
WHERE
        T1.ID2 BETWEEN 600 AND 700
AND
        T1.ID1 & 3 = 0
AND
        T2.ID3 BETWEEN 600 AND 700
GROUP BY
        T1.ID3
```


Какие индексы оптимальнее создать для таблиц Table1 и Table2?

1. alter table Table1 add key (ID1);
2. alter table Table1 add key (ID2);
3. alter table Table1 add key (ID3);
4. alter table Table1 add key (ID1,ID3);
5. alter table Table1 add key (ID1,ID2);
6. alter table Table2 add key (ID1);
7. alter table Table2 add key (ID2);
8. alter table Table2 add key (ID3);
9. alter table Table2 add key (ID1,ID3);

Отметье один или несколько индексов и прокомментируйте свой ответ.

