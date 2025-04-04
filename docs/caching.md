## Принцип работы кэша

Кэш реализован с помощью стратегии LRU, а в качестве паттерна я использовал фасад, в котором находится хранилище и кэш. 
В данный момент реализовано кэширование пользователей и заказов, в данной диаграмме описана работа механизма кэширования заказов. 

```mermaid
zenuml
    title Кэширование заказов
    @Actor User
    web as Router
    serv as Service
    fac as Facade
    cache as Cache
    repo as OrderRepo
    
    User->web.run {
        while(true) {
            User->web: request to path
            
            web->serv.executeCommand {
                try {
                    checkArgs
                    try {
                        serv->fac: executeCommand
                        fac->cache: checkOrders(params)
                        
                        if (!cache.noSuchOrders) {
                            try {
                                fac->repo: processOrders(params)
                                repo->fac: return orders
                            } catch {
                                @return
                                repo->fac: err          
                            }
                            
                            repo->cache: put(orders)
                        }
                        
                    } catch {
                        @return
                        fac->serv: err                         
                    }
                } catch {
                    @return
                    serv->web: err 
                }
                @return
                web->User: http.Error
            }
        }    
    }
```