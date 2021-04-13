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
          weight,
          weightUnit,
          price,
          currency,
          isBeef,
          isAppSku,
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
          categories {
              id,
              name
          },
          productCategories {
              id,
              name
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
            name: "Galaxy S21"
            description: "Smartphone"
            weight: 3
            weightUnit: "Kilogram"
            price: 999,
            currency: "AUD"
            isBeef: false
            isAppSku: true
            isPointSku: true
            categories: [
                {
                    name: "Cat 1"
                },
                {
                    name: "Cat 2"
                }
            ]
            productCategories: [
                {
                    name: "Product Cat 1"
                },
                {
                    name: "Product Cat 2"
                }
            ]
        }
    ) {
        id,
        name,
        description,
        isBeef,
        isAppSku,
        isPointSku,
        weight,
        weightUnit,
        price,
        currency,
        categories {
            id,
            name
        },
        productCategories {
            id,
            name
        }
    }
}
`