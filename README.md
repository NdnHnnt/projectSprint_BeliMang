# projectSprint_BeliMang API Documentation

This project provides a set of APIs for managing users, admins, and merchants.

## Admin Endpoints

- `POST /admin/register`: Register a new admin.
- `POST /admin/login`: Login as an admin.
- `POST /admin/merchant`: Register a new merchant. Requires admin bearer token.
- `GET /admin/merchant`: Get merchant details. Requires admin bearer token.
- `POST /admin/merchant/:merchantId/items`: Register a new item for a merchant. Requires admin bearer token.

## User Endpoints

- `POST /user/register`: Register a new user.
- `POST /user/login`: Login as a user.
- `GET /user/merchant/nearby/:lat,:long`: Get nearby merchants. Requires user bearer token.

Please refer to the individual handler files for more details on the request and response formats:

- [Admin Handler](handler/admin.handler.go)
- [User Handler](handler/user.handler.go)
- [Merchant Handler](handler/merchant.handler.go)