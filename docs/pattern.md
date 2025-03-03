## Почему я выбрал Фабричный метод

Про Фабричный метод я прочитал тут <sup>[1]</sup>.
Я сразу понял, что именно он мне и нужен, потому что:

1. Позволяет освободить моё приложение от зависимости от конкретной реализации упаковок
2. Выносит код производства упаковок в отдельное место
3. Упаковки не зависят друг от друга
4. Упрощает добавление новых упаковок, достаточно определить методы из интерфейса Packaging и добавить новый тип в фабрику

```mermaid
---
title: Работа с Фабричным методом
---
classDiagram
class App
<<subsystem>> App

class Fabric
<<service>> Fabric
Fabric : +GetPackaging(pkgType string)

class Packaging
<<interface>> Packaging
Packaging : +String()
Packaging : +Pack(ord *Order)

class Bag
<<stuct>> Bag

class Box
<<stuct>> Box

class Wrap
<<stuct>> Wrap

App ..> Fabric
App ..> Packaging
Fabric --> Packaging : creates
Packaging <|-- Box : implements
Packaging <|-- Bag : implements
Packaging <|-- Wrap : implements
```

[1] https://refactoring.guru/ru/design-patterns/factory-method.