import 'package:fieldapp/graphql/fragments.dart';

class GQLQuery {
  static final String me = '''
  {
    me {
      ...UserFragment
    }
  }
  ${GQLFragment.user}
  ''';

  static final String getObject = '''
  query getObject(\$id: String!) {
    getObject(id: \$id) {
      product {
        ...ProductFragment
      }
      carton {
        id
        code
        archived
        pallet {
          id
          code
        }
        order {
          id
          code
        }
      }
      pallet {
        id
        code
        archived
        container {
          id
          code
        }
      }
      container {
        id
        code
        archived
      }
    }
  }
  ${GQLFragment.product}
  ''';

  /// $input: { productIDs, cartonIDs, palletIDs, containerIDs }
  static final String getObjects = '''
  query getObjects(\$input: GetObjectsRequest!) {
    getObjects(input: \$input) {
      products {
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
      cartons {
        id
        code
        archived
        pallet {
          id
          code
        }
        order {
          id
          code
        }
      }
      pallets {
        id
        code
        archived
        container {
          id
          code
        }
      }
      containers {
        id
        code
        archived
      }
    }
  }
  ''';

  /// $token: String!, $email: NullString
  static final String verifyResetToken = '''
  query verifyResetToken(\$token: String!, \$email: NullString) {
    verifyResetToken(token: \$token, email: \$email)
  }
  ''';

  /// $search: SearchFilter!, $limit: Int!, $offset: Int!
  static final String skus = '''
  query skus(\$search: SearchFilter!, \$limit: Int!, \$offset: Int!) {
    skus(search: \$search, limit: \$limit, offset: \$offset) {
      skus {
        id
        name
        code
        description
        archived

        masterPlan {
          id
          file_url
        }
      }
      total
    }
  }
  ''';

  static final String fieldappVersion = '''
  {
		settings {
			fieldappVersion
		}
	}
  ''';
}
