# User Service

[![codecov](https://codecov.io/gh/Bhinneka/user-service/branch/development/graph/badge.svg?token=f8UB0JWdDp)](https://codecov.io/gh/Bhinneka/user-service)

#### Status
[![Coverage Status](https://codecov.io/gh/Bhinneka/user-service/branch/development/graphs/sunburst.svg?token=f8UB0JWdDp)](https://codecov.io/gh/Bhinneka/user-service)

This is user service for handling authentication using OAuth and serving about member's data. This service has two handlers in serving data – HTTP and GRPC.

built with :heart: and Go.

## Requirements

 - Golang version 1.15+


## Building and running tests

The software is designed to be able to run both on the machine on top of docker. You can see `Makefile` to see what
kind of targets supported.

### Integration Test

This test should be run with environment variable set.

```
make test
```

### Integration Test With Coverage

This test is to produce covarage report. We use `gocovmerge` to merge coverage reports from the packages. To run, use

```
make cover
```

## How to Use

This service is using _form url encoded_ on its request parameter and [jsonapi](http://jsonapi.org/) on its response body data.
You can follow this postman collection and environment to see and test the endpoints [Postman Collection](https://github.com/Bhinneka/user-service/blob/development/Micro%20Service%20-%20User.postman_collection) and [Postman Environment](https://github.com/Bhinneka/user-service/blob/development/Micro%20Service%20-%20User.postman_environment)

### Authentication

There are 4 methods — anonymous, Facebook, Google, and Azure (Active Directory) — for getting authenticated on `/api/auth` which will be determined by `grantType` parameter's value. The authentication is jwt formatted with `RS256` algorithm.

This public key can be used for another service which will handle token from user service:
```
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCoqzL5JrMzed4tb8uEoLKd42EO
sYmb0HpbicGt/OUeJxaHtt59Ew0BbpreBeiuugXweEa5xctQOxGYr27h4ZOnR0hW
Si+h5Y35CKzMEmZnzQwzQphgqww0U+e9/OAvVfCW1xWvVFr0WbhIRn+w/9DUvp+6
jKz3fIj3yQaHWVMMNQIDAQAB
-----END PUBLIC KEY-----
```

Headers:
- `Authorization: Basic Ymhpbm5la2EtbWljcm9zZXJ2aWNlcy1iMTM3MTQtNTMxMjExNTo2MjY4NjktNmU2ZTY1LTZiNjEyMC02ZDY1NmUtNzQ2MTcyLTY5MjA2NC02OTZkNjU2ZS03MzY5MDA=`
- `Content-Type: application/x-www-form-urlencoded`

#### Anonymous 

For getting anonymous authentication the service needs only `deviceId` parameters and `grantType` parameter's value is `anonymous`.

#### Facebook

For getting facebook authentication the service needs `code` - `token` from facebook login, `deviceId` parameters, and `grantType` parameter's value is `facebook`. If the email which is used does not exist, the email will be registered to the service with minimum required data.

#### Google

For getting google authentication the service needs `code` - `id_token` from google oauth login, `deviceId` parameters, and `grantType` parameter's value is `facebook`. If the email which is used does not exist, the email will be registered to the service with minimum required data.

#### Azure (Active Directory)

For getting azure authentication the service needs `code` - `code` from azure oauth login, `deviceId` parameters, and `grantType` parameter's value is `facebook`. If the email which is used does not exist, the email will be registered to the service with minimum required data.


### Membership

For accessing membership endpoints client must be authenticated to the `user-service` through `/api/auth` in order to get the token.

Headers:
- `Authorization: Bearer <TOKEN>`
- `Content-Type: application/x-www-form-urlencoded`

#### Registration

For registration client can access from `/api/register` using `POST` method to post some form data, e. g.: `firstName`, `lastName`, `email`, `password`, `rePassword`, `gender`, `dob`, and `mobile`. For `dob` must use format `DD/MM/YYYY` and `gender` must be `M` or `F`. This process returns some data which are needed by client to send email for activating membership.

#### Activation Member

After registration user must activate the membership status through `/api/activation` using `POST` method to post the `token`. This endpoint is for member who has set the password before. This process returns some data which are needed by client to send email for sending activation notification.

#### Forgot Password

For getting member's password back, client can access `/api/forgot-password` using `POST` method to post `email`. This process returns some data which are needed by client to send email for validating member's email which request forgot password.

#### Validate Token

Endpoint `/api/validate-token` is for validating token from `forgot-password` that is sent to member's email to display password form from customer facing.

#### Change Password from Forgot Password

Endpoint `/api/change-password` is for changing password after client requested `/api/forgot-password` and validated the token. To process changing password, client needs to `POST` `token`, `newPassword`, and `rePassword`.

#### Activate New Password

Endpoint `/api/activate-new-password` is for activating member who is registered from `dolphin` service and add password while activating the membership. This endpoint is only for user who does not have password and `signUpFrom` value is `dolphin`.

*Note:*
- format of birth date (*dob*): 31/12/2006 (DD/MM/YYYY)
- value of gender: `M` or `F`
- `street1` is mandatory and `street2` is optional parameter

##