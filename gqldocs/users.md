# GraphQL Query and Mutation Examples

## Users

#### USER - Query List Sample
`
query {
  users(
    search: {
      sortBy: DateCreated,
      sortDir: Ascending
    },
    limit: 100,
    offset: 0
  ) {
      users {
          id,
          email,
          firstName,
          lastName,
          mobilePhone,
          organisation {
              id
          },
          role {
              id
          }
      }
      total
  }
}
`

#### USER - Query Object Sample
`
query user {
  user(email: "username@example.com") {
    id,
    email,
    firstName,
    lastName,
    mobilePhone,
    organisation {
        id
    },
    role {
        id
    }
  }
}
`

#### USER - Mutation Query Sample
`
mutation userCreate {
    userCreate(
        input: {
            email: "username@example.com"
            firstName: "User"
            lastName: "Name"
            roleID: "5d23a6d4-6ea4-462c-9e13-b7def16af3df"
            password: "password"
            affiliateOrg: "e39e4131-fd0e-4560-b129-59a4c7fe1f3f"
            mobilePhone: "1234567890"
        }
    ) {
        id,
        email,
        firstName,
        lastName
    }
}
`