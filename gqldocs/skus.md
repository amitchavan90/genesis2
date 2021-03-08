# GraphQL Query and Mutation Examples

## Skus

#### SKU - Query List Sample
`
query {
  skus(
    search: {
      sortBy: DateCreated,
      sortDir: Ascending
    },
    limit: 100,
    offset: 0
  ) {
      skus {
          id,
          name,
          code,
          description,
          isBeef,
          isAppSku,
          isMiniappSku,
          isRetailSku,
          isPointSku,
          hasClones,
          masterPlan {
              id
          },
          video {
              id
          }
          loyaltyPoints,
          cloneParentID,
          archived,
          createdAt,
          urls {
              title,
              content
          },
          productInfo {
              title,
              content
          },
          photos {
              id
          },
          productCount
      }
      total
  }
}
`

#### SKU - Mutation Query Sample
`
mutation skuCreate {
    skuCreate(
        input: {
            name: "iPhone X"
            description: "Smartphone"
            isBeef: false
            isAppSku: true
            isPointSku: true
        }
    ) {
        id,
        name,
        description,
        isBeef,
        isAppSku,
        isPointSku
    }
}
`