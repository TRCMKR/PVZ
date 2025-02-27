## Диаграмма классов

```mermaid
---
title: Диаграмма классов
---
classDiagram
    class App {
        -orderStorage
        -stringBuilder
        
        %% Главный функционал
        +NewApp()
        +Run()
        -executeCommand()
        
        %% CLI функционал
        -clearScr()
        -help()
        -exit()
        -acceptOrder()
        -acceptOrders()
        -returnOrders()
        -formMessage()
        -processOrders()
        -parseOptionalArg()
        -userOrders()
        -returns()
        -orderHistory()
        
        %% Функционал отрисовщика
        -draw()
        -rawDrawer()
        -clearNLinesUp()
        -makePages()
        -printCurrentPage()
        -input()
        -getArrowKeys()
        -setupTerminal()
        -restoreTerminal()
        -changePage()
        -pagedDrawer()
        -printWindow()
        -changeWindow()
        -scrolledDrawer()
    }
    <<stuct>> App
    
    class storage {
        +AcceptOrder()
        +AcceptOrders()
        +ReturnOrder()
        +ProcessOrders()
        +UserOrders()
        +Returns()
        +OrderHistory()
        +Save()
    }
    <<interface>> storage
    
    class JsonStorage {
        -data
        -path
        
        +New()
        +Save()
        +AcceptOrder()
        +AcceptOrders()
        +ReturnOrder()
        +ProcessOrders()
        +UserOrders()
        +Returns()
        +OrderHistory()
        -parseInputDate()
        -processOrder()
    }
    <<struct>> JsonStorage


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

    App ..> storage
    App ..> Fabric
    App ..> Packaging
    Fabric --> Packaging : creates
    storage <|-- JsonStorage : implements
    Packaging <|-- Box : implements
    Packaging <|-- Bag : implements
    Packaging <|-- Wrap : implements
```

## Диаграмма последовательности
На самом деле класс App состоит из 3 разных слоёв: commands, cli, drawer.
Хоть они и не выделены в отдельные объекты, буду описывать взаимодействие между ними, как между разными объектами 
```mermaid
zenuml
    title Диаграмма последовательности для acceptOrder
    @Actor User
    cli as AppCLI
    coms as AppCommands
   
    drawer as AppDrawer
    Fabric
    Storage
    
    User->cli.run {
        while(true) {
            User->cli: text input
            
            cli->coms.executeCommand {
                try {
                    checkArgs
                    coms->Fabric.GetPackaging {
                        return packaging
                    }
                    coms->Storage.AcceptOrder {
                        if (err) {
                            return err
                        }
                        
                        return nil
                    }
                } catch {
                    @return
                    coms->cli: nil, raw, err 
                }
                @return
                coms->cli: strings, raw, nil
            }
            if(err) {
                printError
            } else {
                cli->drawer: draw
            }
        }    
    }

```