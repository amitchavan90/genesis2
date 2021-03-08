import 'package:fieldapp/graphql/fragments.dart';

class GQLMutation {
  /// input: { email: String!, password: String! }
  static final String requestToken = '''
  mutation RequestToken(\$input: RequestToken) {
    RequestToken(input: \$input)
  }
  ''';

  /// input: { email, firstName, lastName, roleID, password, affiliateOrg, mobilePhone }
  static final String changeDetails = '''
  mutation changeDetails(\$input: UpdateUser!) {
    changeDetails(input: \$input) {
      ...UserFragment
    }
  }
  ${GQLFragment.user}
  ''';

  /// $oldPassword: String!, $password: String!
  static final String changePassword = '''
  mutation changePassword(\$oldPassword: String!, \$password: String!) {
    changePassword(oldPassword: \$oldPassword, password: \$password)
  }
  ''';

  /// $email: String!, $viaSMS: Boolean
  static final String forgotPassword = '''
  mutation forgotPassword(\$email: String!, \$viaSMS: Boolean) {
    forgotPassword(email: \$email, viaSMS: \$viaSMS)
  }
  ''';

  /// $token: String!, $password: String!, $email: NullString
  static final String resetPassword = '''
  mutation resetPassword(\$token: String!, \$password: String!, \$email: NullString) {
    resetPassword(token: \$token, password: \$password, email: \$email)
  }
  ''';

  /// $ids: [ID!]!, $action: Action!, $value: BatchActionInput
  static final String productBatchAction = '''
  mutation productBatchAction(\$ids: [ID!]!, \$action: Action!, \$value: BatchActionInput) {
    productBatchAction(ids: \$ids, action: \$action, value: \$value)
  }
  ''';

  /// $ids: [ID!]!, $action: Action!, $value: BatchActionInput
  static final String cartonBatchAction = '''
  mutation cartonBatchAction(\$ids: [ID!]!, \$action: Action!, \$value: BatchActionInput) {
    cartonBatchAction(ids: \$ids, action: \$action, value: \$value)
  }
  ''';

  /// $ids: [ID!]!, $action: Action!, $value: BatchActionInput
  static final String palletBatchAction = '''
  mutation palletBatchAction(\$ids: [ID!]!, \$action: Action!, \$value: BatchActionInput) {
    palletBatchAction(ids: \$ids, action: \$action, value: \$value)
  }
  ''';

  /// ```
  /// input {
  ///   trackActionCode: String!
  ///
  ///   productIDs: [String!]
  ///   cartonIDs: [String!]
  ///   palletIDs: [String!]
  ///   containerIDs: [String!]
  ///
  ///   productScanTimes: [Time!]
  ///   cartonScanTimes: [Time!]
  /// 	palletScanTimes: [Time!]
  /// 	containerScanTimes: [Time!]
  ///
  ///   cartonPhotoBlobIDs: [String!]
  ///   productPhotoBlobIDs: [String!]
  ///
  ///   memo: NullString
  ///   locationGeohash: NullString
  ///   locationName: NullString
  /// }
  /// ```
  static final String recordTransaction = '''
  mutation recordTransaction(\$input: RecordTransactionInput!) {
    recordTransaction(input: \$input)
  }
  ''';

  static final String fileUpload = '''
  mutation fileUpload(\$file: Upload!) {
    fileUpload(file: \$file) {
        id
        file_url
     }
  }
  ''';
}
