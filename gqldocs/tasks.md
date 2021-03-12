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
          isTimeBound,
          isPeopleBound,
          isProductRelevant,
          isFinal,
          finishDate,
          maximumPeople,
          sku {
             id,
             name,
             code
          },
          subtasks {
              id,
              title,
              description
          },
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
    sku {
        id,
        name,
        code
    },
    subtasks {
        id,
        title,
        description
    }
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
            subtasks: [
                {
                    title: "Subtask 1"
                    description: "Subtask One"
                },
                {
                    title: "Subtask 2"
                    description: "Subtask Two"
                }
            ]
        }
    ) {
        title,
        description,
        loyaltyPoints,
        isTimeBound,
        isPeopleBound,
        isProductRelevant,
        isFinal,
        finishDate,
        maximumPeople,
        sku {
            id,
            name,
            code
        }
        subtasks {
            id,
            title,
            description
        }
    }
}
`