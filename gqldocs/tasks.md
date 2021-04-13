# GraphQL Query and Mutation Examples

## Task

#### TASK - Query List Sample
`
query {
  userTasks(
    search: {
      sortBy: DateCreated,
      sortDir: Ascending
    },
    limit: 100,
    offset: 0
  ) {
      userTasks {
          id,
          isComplete,
          status,
          task {
              id,
              title,
              description,
          },
          user {
              id,
              firstName,
              lastName,
              email
          }
          userSubtasks {
              id,
              isComplete,
              status
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
mutation userTaskCreate {
    userTaskCreate(
        input: {
            taskID: "64f611f2-475a-4276-9b9a-1dcc5536f3ce"
        }
    ) {
        id,
        isComplete,
        status,
        task {
            id,
            title,
            description
        },
        user {
            id,
            firstName,
            lastName,
            email
        },
        subtasks {
            id,
            isComplete,
            status
        },
        createdAt
    }
}
`