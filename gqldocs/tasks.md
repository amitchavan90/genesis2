# GraphQL Query and Mutation Examples

## Task

#### TASK - Query List Sample
`
query {
  tasks(
    search: {
      sortBy: DateCreated,
      sortDir: Ascending
    },
    limit: 100,
    offset: 0
  ) {
      tasks {
        id,
        title,
        description,
        loyaltyPoints,
        isTimeBound,
        isPeopleBound,
        isProductRelevant,
        finishDate,
        maximumPeople,
        skuID,
        createdAt
      }
      total
  }
}
`

#### TASK - Query Object Sample
`
query task {
  task(id: "bd7c3b92-c127-47ab-bb44-f22ec4cfe448") {
    id,
    title,
    description,
    loyaltyPoints,
    isTimeBound,
    isPeopleBound,
    isProductRelevant,
    finishDate,
    maximumPeople,
    skuID,
    createdAt
  }
}
`

#### TASK - Mutation Query Example
`
mutation taskCreate {
    taskCreate(
        input: {
            title: "Task Title"
            description: "Task Description"
            loyaltyPoints: 50
            isTimeBound: true
            isPeopleBound: true
            isProductRelevant: false
            finishDate: "2021-03-06T15:04:05Z"
            maximumPeople: 10
        }
    ) {
        title,
        description,
        loyaltyPoints,
        isTimeBound,
        isPeopleBound,
        isProductRelevant,
        finishDate,
        maximumPeople,
        skuID
    }
}
`