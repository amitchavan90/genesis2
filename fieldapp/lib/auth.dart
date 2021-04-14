import 'package:fieldapp/graphql/mutations.dart';
import 'package:fieldapp/graphql/queries.dart';
import 'package:fieldapp/types/types.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:fieldapp/main.dart';

class Auth {
  static Link getLink() {
    final HttpLink httpLink = HttpLink(uri: '$host/api/gql/query');
    final AuthLink authLink =
        AuthLink(getToken: () => prefs.getString("token"));
    return authLink.concat(httpLink);
  }

  static Future<GQLResponse> signIn(String email, String password) async {
    GQLResponse response = new GQLResponse();

    try {
      bool timeout = false;
      QueryResult result = await client
          .mutate(
        MutationOptions(
          documentNode: gql(GQLMutation.requestToken),
          variables: {
            'input': {
              'email': email,
              'password': password,
            }
          },
        ),
      )
          .timeout(
        timeoutDuration,
        onTimeout: () {
          timeout = true;
          return QueryResult();
        },
      );

      if (timeout) {
        response.message = "Timed out";
        return response;
      }

      // Error check
      if (result.hasException || result.data == null) {
        if (result.exception.graphqlErrors.length != 0)
          response.message = result.exception.graphqlErrors[0].message;
        else if (result.exception.clientException != null)
          response.message = result.exception.clientException.message;
        return response;
      }

      // Set token
      String token = result.data["RequestToken"];
      await prefs.setString("host", host);
      await prefs.setString("hostOption", hostOption);
      await prefs.setString("token", 'Bearer $token');
      updateGQLClient();

      response.success = true;
      response.message = "Success";
    } catch (e) {
      print(e.toString());
    }

    return response;
  }

  static void signOut() => prefs.remove("token");

  static Future<GQLResponse> getMe() async {
    GQLResponse response = new GQLResponse();

    try {
      bool timeout = false;
      QueryResult result = await client
          .query(
        QueryOptions(
          documentNode: gql(GQLQuery.me),
          fetchPolicy: FetchPolicy.networkOnly,
        ),
      )
          .timeout(
        timeoutDuration,
        onTimeout: () {
          timeout = true;
          return QueryResult();
        },
      );

      if (timeout) {
        response.message = "Timed out while fetching user.";
        return response;
      }

      // Error check
      if (result.hasException || result.data == null) {
        if (result.exception.graphqlErrors.length == 0) return response;
        response.message = result.exception.graphqlErrors[0].message;
        return response;
      }

      // Set me
      me = User.fromJson(result.data["me"]);

      response.success = true;
      response.message = "Success";
    } catch (e) {
      print(e.toString());
    }

    return response;
  }

  static Future<GQLResponse> forgotPassword(String email) async {
    GQLResponse response = new GQLResponse();

    try {
      bool timeout = false;
      QueryResult result = await client
          .mutate(
        MutationOptions(
          documentNode: gql(GQLMutation.forgotPassword),
          variables: {
            'email': email,
            'viaSMS': true,
          },
        ),
      )
          .timeout(
        timeoutDuration,
        onTimeout: () {
          timeout = true;
          return QueryResult();
        },
      );

      if (timeout) {
        response.message = "Timed out";
        return response;
      }

      // Error check
      if (result.hasException || result.data == null) {
        if (result.exception.graphqlErrors.length == 0) return response;
        response.message = result.exception.graphqlErrors[0].message;
        return response;
      }

      response.success = true;
      response.message = "Success";
    } catch (e) {
      print(e.toString());
    }

    return response;
  }

  static Future<GQLResponse> verifyResetToken(
    String token,
    String email,
  ) async {
    GQLResponse response = new GQLResponse();

    try {
      bool timeout = false;
      QueryResult result = await client
          .query(
        QueryOptions(
          documentNode: gql(GQLQuery.verifyResetToken),
          variables: {
            'token': token,
            'email': email,
          },
        ),
      )
          .timeout(
        timeoutDuration,
        onTimeout: () {
          timeout = true;
          return QueryResult();
        },
      );

      if (timeout) {
        response.message = "Timed out while verifying reset token.";
        return response;
      }

      // Error check
      if (result.hasException || result.data == null) {
        if (result.exception.graphqlErrors.length == 0) return response;
        response.message = result.exception.graphqlErrors[0].message;
        return response;
      }

      response.success = true;
      response.message = "Valid";
    } catch (e) {
      print(e.toString());
    }

    return response;
  }

  static Future<GQLResponse> resetPassword(
    String token,
    String password,
    String email,
  ) async {
    GQLResponse response = new GQLResponse();

    try {
      bool timeout = false;
      QueryResult result = await client
          .mutate(
        MutationOptions(
          documentNode: gql(GQLMutation.resetPassword),
          variables: {
            'token': token,
            'password': password,
            'email': email,
          },
        ),
      )
          .timeout(
        timeoutDuration,
        onTimeout: () {
          timeout = true;
          return QueryResult();
        },
      );

      if (timeout) {
        response.message = "Timed out";
        return response;
      }

      // Error check
      if (result.hasException || result.data == null) {
        if (result.exception.graphqlErrors.length == 0) return response;
        response.message = result.exception.graphqlErrors[0].message;
        return response;
      }

      response.success = true;
      response.message = "Success";
    } catch (e) {
      print(e.toString());
    }

    return response;
  }

  static Future<GQLResponse> appVersionCheck() async {
    GQLResponse response = new GQLResponse();

    try {
      bool timeout = false;
      QueryResult result = await client
          .query(
        QueryOptions(
          documentNode: gql(GQLQuery.fieldappVersion),
          fetchPolicy: FetchPolicy.networkOnly,
        ),
      )
          .timeout(
        timeoutDuration,
        onTimeout: () {
          timeout = true;
          return QueryResult();
        },
      );

      if (timeout) {
        response.message = "Timed out while doing app version check.";
        return response;
      }

      // Error check
      if (result.hasException || result.data == null) {
        if (result.exception.graphqlErrors.length == 0) return response;
        response.message = result.exception.graphqlErrors[0].message;
        return response;
      }

      response.success = true;
      response.message = result.data["settings"]["fieldappVersion"];
    } catch (e) {
      print(e.toString());
    }

    return response;
  }
}
