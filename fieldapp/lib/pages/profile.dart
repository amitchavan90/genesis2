import 'package:auto_size_text/auto_size_text.dart';
import 'package:fieldapp/graphql/mutations.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/utils.dart';
import 'package:fieldapp/widgets/commonWidgets.dart';
import 'package:fieldapp/widgets/dialogs.dart';
import 'package:flutter/material.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:fieldapp/main.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:fieldapp/widgets/changePassword.dart';

class ProfilePage extends StatefulWidget {
  final void Function(bool) setLoading;
  ProfilePage({Key key, @required this.setLoading}) : super(key: key);

  @override
  _ProfilePageState createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final GlobalKey<ScaffoldState> scaffoldKey = GlobalKey<ScaffoldState>();
  final formKey = GlobalKey<FormState>();
  final Map<String, TextEditingController> controllers =
      new Map<String, TextEditingController>();
  final Map<String, FocusNode> focusNodes = new Map<String, FocusNode>();

  @override
  void dispose() {
    controllers.forEach((key, value) => value.dispose());
    focusNodes.forEach((key, value) => value.dispose());
    controllers.clear();
    focusNodes.clear();
    super.dispose();
  }

  void onSubmit() async {
    if (!formKey.currentState.validate()) return;

    FocusScope.of(context).unfocus(); // close keyboard

    widget.setLoading(true);

    User updatedMe;
    try {
      bool timeout = false;
      QueryResult result = await client
          .mutate(
        MutationOptions(
          documentNode: gql(GQLMutation.changeDetails),
          variables: {
            'input': {
              'email': controllers["Email"].text,
              'firstName': controllers["First Name"].text,
              'lastName': controllers["Last Name"].text,
              'mobilePhone': controllers["Mobile"].text,
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
        widget.setLoading(false);
        showErrorDialog(
          context,
          "Timed out",
        );
        return;
      }

      if (result.hasException || result.data == null) {
        widget.setLoading(false);
        showErrorDialog(
          context,
          result.exception.graphqlErrors.length == 0
              ? "An issue occured"
              : result.exception.graphqlErrors[0].message,
        );
        return;
      }

      // Get user
      updatedMe = User.fromJson(result.data["changeDetails"]);
    } catch (e) {
      showErrorDialog(context, "An issue occured");
    }

    widget.setLoading(false);

    if (updatedMe != null) me = updatedMe;
    showSuccessDialog(context, "Changes Saved", "");
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      key: scaffoldKey,
      body: ListView(
        padding: EdgeInsets.all(20),
        children: <Widget>[
          // Heading
          userCard(),
          SizedBox(height: 18),
          // Form
          Form(
            key: formKey,
            autovalidate: true,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: <Widget>[
                // Email
                basicTextField(
                  "Email",
                  me != null ? me.email : "",
                  controllers,
                  focusNodes,
                  keyboardType: TextInputType.emailAddress,
                  textInputAction: TextInputAction.next,
                  validator: (value) {
                    if (value.isEmpty)
                      return ' This field is required ';
                    else if (!isEmail(value))
                      return ' Please enter a valid email address ';
                    return null;
                  },
                ),
                // First Name
                basicTextField(
                  "First Name",
                  me != null ? me.firstName : "",
                  controllers,
                  focusNodes,
                  textInputAction: TextInputAction.next,
                  validator: (value) =>
                      value.isEmpty ? ' This field is required ' : null,
                ),
                // Last Name
                basicTextField(
                  "Last Name",
                  me != null ? me.lastName : "",
                  controllers,
                  focusNodes,
                  textInputAction: TextInputAction.next,
                  validator: (value) =>
                      value.isEmpty ? ' This field is required ' : null,
                ),
                // Mobile
                basicTextField(
                  "Mobile",
                  me != null ? me.mobilePhone : "",
                  controllers,
                  focusNodes,
                  textInputAction: TextInputAction.done,
                  keyboardType: TextInputType.phone,
                ),

                // Submit
                Row(
                  children: <Widget>[
                    Expanded(
                      child: RaisedButton(
                        child: Text(
                          "Save Changes",
                          style: TextStyle(
                            color: Colors.white,
                          ),
                        ),
                        color: COLOUR_PRIMARY,
                        onPressed: onSubmit,
                      ),
                    ),
                  ],
                ),

                // Change Password
                Align(
                  alignment: Alignment.topRight,
                  child: RaisedButton(
                    child: Text(
                      "Change Password",
                    ),
                    onPressed: () {
                      showDialog<void>(
                        context: context,
                        barrierDismissible: false,
                        builder: (BuildContext context) => ChangePassword(),
                      );
                    },
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget userCard() {
    return Row(
      children: <Widget>[
        Padding(
          child: FaIcon(
            FontAwesomeIcons.userCircle,
            size: 36,
          ),
          padding: EdgeInsets.only(right: 8),
        ),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: <Widget>[
            AutoSizeText(
              '${me.firstName} ${me.lastName}',
              style: TextStyle(fontSize: 20),
            ),
            AutoSizeText(
              me.email,
              style: TextStyle(fontSize: 16, color: Colors.grey.shade600),
            )
          ],
        ),
      ],
    );
  }
}
