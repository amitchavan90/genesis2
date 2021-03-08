### This is the REST API implementation
- Database settings are defined inside `genesis/server/restpkg/settings/database.go`.
- Every models are defined in the package `genesis/server/restpkg/models` (This is a temporary arrangement).
- Business logic is defined inside the package `genesis/server/restpkg/services`.
- API handlers are defined inside the package `genesis/server/restpkg/handlers`. Handlers also verifies the token authentication and permissions
- Routes are defined in the package `genesis/server/restpkg/routes`.
- Each routes are mapped into a single function MapUrls() defined in `genesis/server/restpkg/routes/routes.go`.
- This MapUrls() is called inside the file `genesis/api/api.go`. This single function handles all the REST APIs defined in the `genesis/server/restpkg`.