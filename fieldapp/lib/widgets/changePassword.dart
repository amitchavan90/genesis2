import 'package:fieldapp/graphql/mutations.dart';
import 'package:fieldapp/main.dart';
import 'package:flutter/material.dart';
import 'package:fieldapp/widgets/commonWidgets.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:graphql_flutter/graphql_flutter.dart';

class ChangePassword extends StatefulWidget {
  final String token;
  final String email;
  ChangePassword({Key key, this.token, this.email}) : super(key: key);

  @override
  _ChangePasswordState createState() => _ChangePasswordState();
}

class _ChangePasswordState extends State<ChangePassword> {
  final formKey = GlobalKey<FormState>();
  final Map<String, TextEditingController> controllers =
      new Map<String, TextEditingController>();
  final Map<String, FocusNode> focusNodes = new Map<String, FocusNode>();
  bool loading = false;

  bool showPassword = false;
  bool showNewPassword = false;

  void onSubmit() async {
    // Validate
    if (!formKey.currentState.validate()) return;

    bool resetPassword = widget.token != null;

    // get passwords
    String password = resetPassword ? "" : controllers["Current Password"].text;
    String newPassword = controllers["New Password"].text;

    // Changing to the same password?
    if (!resetPassword && password == newPassword) {
      controllers["Current Password"].clear();
      controllers["New Password"].clear();

      showWarningDialog(
        context,
        "",
        "Please enter a new password",
      );

      return;
    }

    setState(() => loading = true);

    // Change Password
    try {
      bool timeout = false;
      QueryResult result = await client
          .mutate(
        MutationOptions(
          documentNode: gql(resetPassword
              ? GQLMutation.resetPassword
              : GQLMutation.changePassword),
          variables: resetPassword
              ? {
                  "token": widget.token,
                  "password": newPassword,
                  "email": widget.email,
                }
              : {
                  "oldPassword": password,
                  "password": newPassword,
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

      setState(() => loading = false);

      if (timeout) {
        showErrorDialog(context, "Timed out");
        return;
      }

      if (result.hasException || result.data == null) {
        showErrorDialog(
          context,
          result.exception.graphqlErrors.length == 0
              ? "Password Change Failed"
              : result.exception.graphqlErrors[0].message,
        );
        return;
      }

      Navigator.of(context).pop(); // close pop-up
      showSuccessDialog(context, "Password Change Successful", "");

      if (!resetPassword) controllers["Current Password"].clear();
      controllers["New Password"].clear();
    } catch (e) {
      setState(() => loading = false);
      showErrorDialog(context, "Password Change Failed");
    }
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: <Widget>[
        Form(
          key: formKey,
          autovalidate: true,
          child: SimpleDialog(
            contentPadding: EdgeInsets.all(10),
            children: <Widget>[
              // Current Password
              widget.token == null
                  ? Row(
                      children: <Widget>[
                        Expanded(
                          child: basicTextField(
                            "Current Password",
                            "",
                            controllers,
                            focusNodes,
                            obscureText: !showPassword,
                            validator: (value) {
                              if (value.isEmpty)
                                return ' This field is required ';
                              return null;
                            },
                            textInputAction: TextInputAction.next,
                          ),
                        ),
                        SizedBox(
                          width: 35,
                          child: FlatButton(
                            padding: EdgeInsets.symmetric(horizontal: 5),
                            child: FaIcon(
                              showPassword
                                  ? FontAwesomeIcons.eye
                                  : FontAwesomeIcons.eyeSlash,
                              size: 20,
                            ),
                            onPressed: () =>
                                setState(() => showPassword = !showPassword),
                          ),
                        ),
                      ],
                    )
                  : Container(),

              // New Password
              Row(
                children: <Widget>[
                  Expanded(
                    child: basicTextField(
                      "New Password",
                      "",
                      controllers,
                      focusNodes,
                      obscureText: !showNewPassword,
                      validator: (value) {
                        if (value.isEmpty) return ' This field is required ';
                        if (value.length < 4) return ' Password is too short ';
                        return null;
                      },
                      onFieldSubmitted: (_) => onSubmit(),
                    ),
                  ),
                  SizedBox(
                    width: 35,
                    child: FlatButton(
                      padding: EdgeInsets.symmetric(horizontal: 5),
                      child: FaIcon(
                        showNewPassword
                            ? FontAwesomeIcons.eye
                            : FontAwesomeIcons.eyeSlash,
                        size: 20,
                      ),
                      onPressed: () =>
                          setState(() => showNewPassword = !showNewPassword),
                    ),
                  ),
                ],
              ),

              RaisedButton(
                child: Text(
                  "Change Password",
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                color: COLOUR_PRIMARY,
                onPressed: onSubmit,
              ),

              SizedBox(height: 5),

              FlatButton(
                child: Text(
                  "Cancel",
                  style: TextStyle(
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                color: Colors.black.withOpacity(0.45),
                onPressed: () => Navigator.of(context).pop(),
              ),
            ],
          ),
        ),
        loadingOverlay(loading),
      ],
    );
  }
}
