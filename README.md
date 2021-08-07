# Quick key value db evaluation

## Scope of evaluation

1. Prefix range query read performance
2. Write performance

## Steps

1. Insert 10 million random records with the following structure
    - key: address+integer+integer
    - value: transaction hash
2. Generating skewed data
    - Mimicking the total number of transactions on ethereum
    - 2^0 + 2^1 + .... + 2^29
    - For each two power value, equal number of entries will be created with the same address
3. Exposing a metrics api
    - Tracks time writes for each address
    - Tracks time for reads for each address