# GraphQL Query and Mutation Examples

## Referral

#### REFERRAL - Query List Sample
`
query {
  referrals(
    search: {
      sortBy: DateCreated,
      sortDir: Ascending
    },
    limit: 100,
    offset: 0
  ) {
      referrals {
          id,
          userID,
          referredByID,
          isRedemmed,
          createdAt
      }
      total
  }
}
`

#### REFERRAL - Query Object Sample
`
query referral {
  referral(userID: "bd7c3b92-c127-47ab-bb44-f22ec4cfe448") {
    id,
    userID,
    referredByID,
    isRedemmed,
    createdAt
  }
}
`