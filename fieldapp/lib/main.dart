import 'package:fieldapp/auth.dart';
import 'package:fieldapp/pages/menu.dart';
import 'package:fieldapp/pages/profile.dart';
import 'package:fieldapp/pages/scannerPage.dart';
import 'package:fieldapp/pages/login.dart';
import 'package:fieldapp/pages/settings.dart';
import 'package:fieldapp/pages/start.dart';
import 'package:fieldapp/types/types.dart';
import 'package:fieldapp/widgets/appScaffold.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:graphql_flutter/graphql_flutter.dart';
import 'package:page_transition/page_transition.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter_sentry/flutter_sentry.dart';

/// API Host options
final Map<String, String> hostOptions = {
  "Australia": "https://admin.gn.l28produce.com.au",
  "China": "https://admin.gn.latitude28.cn",
  "Staging": "https://admin.gn.staging.l28produce.com.au",
  "Custom": "",
};

/// API Host
String hostOption = "Australia";
String host = hostOptions[hostOption];

SharedPreferences prefs;
final GlobalKey<NavigatorState> navKey = new GlobalKey<NavigatorState>();
User me;

enum Mode {
  // binding/unbinding items
  advance,
  // binds carton to product + applies track action
  packCartonCompleteTask,
  // applies track action to scanned item
  completeTask,
}

Mode mode;

const Color COLOUR_PRIMARY = Color(0xFF252525);
const Color COLOUR_SECONDARY = Color(0xFF1A1A1A);

const Duration timeoutDuration = Duration(seconds: 30);

// Setup GraphQL
GraphQLClient getGQLClient() => GraphQLClient(
      cache: InMemoryCache(),
      link: Auth.getLink(),
    );
GraphQLClient client = getGQLClient();
ValueNotifier<GraphQLClient> clientValueNotifier = ValueNotifier(client);
void updateGQLClient() {
  client = getGQLClient();
  clientValueNotifier.value = client;
}

// Setup App
Future<void> main() => FlutterSentry.wrap(
      () async {
        WidgetsFlutterBinding.ensureInitialized();

        // change nav bar colour
        SystemChrome.setSystemUIOverlayStyle(
          SystemUiOverlayStyle(
            systemNavigationBarColor: COLOUR_PRIMARY,
            statusBarColor: Colors.transparent,
            statusBarIconBrightness: Brightness.dark,
          ),
        );

        // get prefs
        prefs = await SharedPreferences.getInstance();
        if (prefs.containsKey("host")) host = prefs.getString("host");
        if (prefs.containsKey("hostOption")) {
          var option = prefs.getString("hostOption");
          if (hostOptions.containsKey(option)) {
            hostOption = option;
            host = hostOptions[hostOption];
          }
        }

        // start app
        await SystemChrome.setPreferredOrientations(
          [DeviceOrientation.portraitUp, DeviceOrientation.portraitDown],
        );

        runApp(MyApp());
      },
      dsn:
          'https://f5204318bd264d36962e4acabdac1cd5@o370480.ingest.sentry.io/5194354',
    );

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return GraphQLProvider(
      client: clientValueNotifier,
      child: MaterialApp(
        title: 'Genesis',
        theme: ThemeData(
          primaryColor: COLOUR_PRIMARY,
          accentColor: COLOUR_SECONDARY,
          textTheme: Theme.of(context).textTheme.apply(fontFamily: 'Konnect'),
          splashColor: COLOUR_PRIMARY,
          highlightColor: Colors.white.withOpacity(0.3),
        ),
        debugShowCheckedModeBanner: false,
        navigatorKey: navKey,
        initialRoute: '/',
        onGenerateRoute: (settings) {
          String _route = settings.name;

          switch (_route) {
            case '/':
              return PageTransition(
                type: PageTransitionType.fade,
                child: StartPage(),
              );
            case '/menu':
              return PageTransition(
                type: PageTransitionType.fade,
                child: AppScaffold(
                  builder: (isLoading, setLoading) => HomeOptionsPage(),
                  title: "Genesis",
                ),
              );
            case '/login':
              return PageTransition(
                type: PageTransitionType.fade,
                child: LoginPage(),
              );
            case "/home":
              return PageTransition(
                type: PageTransitionType.rightToLeft,
                child: AppScaffold(
                  builder: (isLoading, setLoading) =>
                      ScannerPage(setLoading: setLoading),
                  titleFontSize: mode == Mode.packCartonCompleteTask ? 17 : 20,
                  title: mode == Mode.advance
                      ? "Advanced"
                      : mode == Mode.completeTask
                          ? "Complete Task Mode"
                          : mode == Mode.packCartonCompleteTask
                              ? "Pack Carton and Complete Task Mode"
                              : "",
                ),
              );
            case '/settings':
              return PageTransition(
                type: PageTransitionType.rightToLeft,
                child: AppScaffold(
                  builder: (isLoading, setLoading) => SettingsPage(),
                  title: "Settings",
                ),
              );
            case '/profile':
              return PageTransition(
                type: PageTransitionType.rightToLeft,
                child: AppScaffold(
                  builder: (isLoading, setLoading) =>
                      ProfilePage(setLoading: setLoading),
                  title: "Profile",
                ),
              );
          }

          return null;
        },
      ),
    );
  }
}
