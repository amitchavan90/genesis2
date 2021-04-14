class GQLFragment {
  static final String role = '''
  fragment RoleFragment on Role {
    id
    name
    permissions
    tier

    trackActions {
      id
      code
      archived
      name
      nameChinese
      requirePhotos
    }
  }
  ''';

  static final String user = '''
  fragment UserFragment on User {
    id
    firstName
    lastName
    email
    verified
    role {
      ...RoleFragment
    }
    affiliateOrg
    mobilePhone
    mobileVerified
  }
  $role
  ''';

  static final String product = '''
  fragment ProductFragment on Product {
    id
    code
    archived
    carton {
      id
      code
    }
    sku {
      id
      code
    }
    order {
      id
      code
    }
  }
  ''';

  static final String carton = '''
	fragment CartonFragment on Carton {
		id
		code
		archived
		createdAt
		productCount

		pallet {
			id
			code
			container {
				id
				code
			}
		}
		sku {
			id
			code
		}
		order {
			id
			code
		}
	}
''';
}
