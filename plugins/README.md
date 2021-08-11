# How plugin work?

## plugin state and callback flow
```mermaid
graph TD;
    A[write your plugin]--copy files to plugins-->B[start]
    B --run app--> C[unmigrated]
    C --"triforce call migrateUp() when added or updated"--> D[ready]
    D --"call migrateDown() when removed"--> C
    D --"call execute() when start cal data"--> V0[started]
```

## collector and enricher flow
```mermaid
graph TD;
    A[plugin unregistered]--"call plugin.execute() to add collectors and enrichers"-->B[plugin started]
    B --"call each collector/enrichers.name() to register"--> C[collector and enricher register]
    C --"when start collector, run `collector.dependencies(pks)`"--> D[got collector dependencies map]
    D --"run `collector.isDataPrepared(pks) && collector.collectData(pks)` one by one"--> E[collect success]
    E --"when start enricher, run `enricher.dependencies(pks)`"--> F[got enricher dependencies map]
    F --"when start enricher, run `enricher.dependencies(pks)`"--> G{collectors ready?}
    G --"ready"--> I{"enricher has no dependency && enricher.couldLazyLoad(pks)"}
    G --"not ready"--> H[fail]
    I --"no"--> J[start cal data]
    I --"yes"--> K[finish]
    J --"run `enricher.isDataPrepared(pks) && collector.collectData(pks)`"--> K[enricher data ready]
    K --"collector.queryData()"--> M[finish]
```

## picture generate by mermaid

view on github: `https://github.com/BackMarket/github-mermaid-extension`

or view online: `https://mermaid-js.github.io/mermaid-live-editor`
