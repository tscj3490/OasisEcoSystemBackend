# Basics

## Base URL

The base URL for all API requests is: `https://oasis-api.merkins.io/api`. This means all the `/resources` used from now on in this documentation must be preceded by the base URL when performing the request.

## Authentication

API requests are made over HTTPS and authenticated with a key. This key is supplied as the {token} in an HTTP Authorization header: `Authorization: Bearer {token}`

To get your key:

1. Create an HTTP POST request to `/authentication` to perform the login.
2. Set the HTTP Content type header `Content-Type: application/json`.
3. Set a JSON body as follows:

    ```json
    {
      "email": "your_user@email.com",
      "password": "your-username-password",
    }
    ```

4. The response from this request will return the current user data, along with the JWT token:

    ```json
    {
      "currentUser": {
        // <Current User data>
      },
      "token": "<JWT Token>"
    }
    ```

5. Once the token is retrieved, it must be present in all the future request in the HTTP authorization header.

## Requests

### Basic Concepts

The API is organized around REST: resources are accessed with URL endpoints and manipulated using HTTP GET, POST, PUT and DELETE verbs. GET is used for reading resources, POST is used for creating resources, PUT for updating resources, and DELETE is used for deleting resources.

In general, resources identifiers are specified in the URL. In cases where objects are being created or updated, additional parameters are sent in the request body.

### Headers

When making all requests, the authorization header must be set and include the API key in the token portion of the authorization header: Authorization: Bearer {token}

In addition, when making requests that create or update a resource, parameters are sent as JSON in the POST request and the content type header must be set to the JSON content type: Content-Type: application/json

## Responses

### Basic Concepts

In most cases, responses are returned as JSON, and the HTTP status code indicates success or failure. In general, response codes in the 200-level range indicate success, response codes in the 400-level range indicate a failure caused by a request that can’t be fulfilled even though the service is working as expected, and response codes in the 500-level range indicate a failure caused by a problem with the service itself.

HTTP status codes have the following meanings:

| Code | Meaning               |
| ---- | --------------------- |
| 200  | OK                    |
| 400  | Bad request           |
| 401  | Unauthorized          |
| 403  | Forbidden             |
| 404  | Not found             |
| 500  | Internal server error |

# Current user

## Overview

Current user is the user releated to the JWT Token sent in the HTTP authorization header as Bearer

## GET /me

### Description

Gets the current user data

### Response

| Property | Type     | Description     |
| -------- | -------- | --------------- |
| `id`     | ObjectId | User identifier |
| `email`  | string   | User email      |
| `role`   | string   | User role       |


*NOTE: there are more properties but for simplicity we will omit them*

# Users

## Overview

Users of the application. They are differenciated in 5 roles.

| Role            | Description                                   |
| --------------- | --------------------------------------------- |
| ROLE_ADMIN      | The administrator role                        |
| ROLE_FARMER     | The client role                               |
| ROLE_ADVISER    | This role is the one that executes the visits |
| ROLE_WORKER     | This role is the one that executes the tasks  |
| ROLE_TECHNICIAN | This role is currently unused                 |

## GET /users/

### GET /users/ Description

Fetches all users of the application

### GET /users/ Response

#### Array of

| Property         | Type     | Description                                                       |
| ---------------- | -------- | ----------------------------------------------------------------- |
| `id`             | ObjectId | User identifier                                                   |
| `email`          | string   | User email                                                        |
| `role`           | string   | User role                                                         |
| `firstName`      | string   | User first name                                                   |
| `lastName`       | string   | User last name                                                    |
| `fullName`       | string   | User full name                                                    |
| `organizationId` | string   | currently unused                                                  |
| `cooperativeId`  | string   | Cooperative id the user is worker or client depending on his role |

## POST /users/

### POST /users/ Decription

Method to create a User, this one is done in the dashboard

### POST /users/ Body Parameters

| `email`          | string   | User email                                                        |
| `password`       | string   | User password                                                     |
| `role`           | string   | User role                                                         |
| `firstName`      | string   | User first name                                                   |
| `lastName`       | string   | User last name                                                    |
| `organizationId` | string   | currently unused                                                  |
| `cooperativeId`  | string   | Cooperative id the user is worker or client depending on his role |

## GET /users/{userId}

### GET /users/{userId} Description

Fetches the user whose id matches with the one passed on the call

### GET /users/{userId} Response

| Property         | Type     | Description                                                       |
| ---------------- | -------- | ----------------------------------------------------------------- |
| `id`             | ObjectId | User identifier                                                   |
| `email`          | string   | User email                                                        |
| `role`           | string   | User role                                                         |
| `firstName`      | string   | User first name                                                   |
| `lastName`       | string   | User last name                                                    |
| `fullName`       | string   | User full name                                                    |
| `organizationId` | string   | currently unused                                                  |
| `cooperativeId`  | string   | Cooperative id the user is worker or client depending on his role |

## POST /users/{userId}

### POST /users/{userId} Description

Updates the user whose id matches with the one passed on the call

### POST /users/{userId} Body Parameters

| `email`          | string   | User email                                                        |
| `password`       | string   | User password                                                     |
| `role`           | string   | User role                                                         |
| `firstName`      | string   | User first name                                                   |
| `lastName`       | string   | User last name                                                    |
| `organizationId` | string   | currently unused                                                  |
| `cooperativeId`  | string   | Cooperative id the user is worker or client depending on his role |

## DELETE /users/{userId}

### DELETE /users/{userId} Description

Updates the user whose id matches with the one passed on the call changing active to false so it will no longer be fetched by the GET methods

## POST /users/{userId}/termsOfService

### POST /users/{userId}/termsOfService Description

Updates the user whose id matches with the one passed on the call changing termsOfService to true

## GET /users/technicians

### GET /users/technicians Description

Fetches all users whose role is "ROLE_TECHNICIAN"

## GET /users/workers

### GET /users/workers Description

Fetches all users whose role is "ROLE_WORKER"

## GET /users/cooperative/{cooperativeId}

### GET /users/cooperative/{cooperativeId} Description

Fetches all user whose role is "ROLE_FARMER" and cooperativeId is the same as the one passed on the call

# Cooperatives

## Overview

Cooperatives are the entity that houses the plantations, farmers and workers.

## GET /cooperatives/

### GET /cooperatives/ Description

Fetches all the cooperatives and returns them in an Array.

### GET /cooperatives/ Response

#### Array of

| Property    | Type        | Description                    |
| ----------- | ----------- | ------------------------------ |
| `id`        | ObjectId    | Autogenerated cooperative id   |
| `name`      | string      | The name of the cooperative    |
| `address`   | string      | The address of the cooperative |
| `taxIdCode` | string      | The cooperative C.I.F code     |
| `location`  | Coordinates | Coordinates of the cooperative |

## POST /cooperatives/

### POST /cooperatives/ Description

Creates a cooperative and returns the newly created cooperative object with its properties.

### POST /cooperatives/ Body Parameters

| Property    | Type        | Description                    |
| ----------- | ----------- | ------------------------------ |
| `name`      | string      | The name of the cooperative    |
| `address`   | string      | The address of the cooperative |
| `taxIdCode` | string      | The cooperative C.I.F code     |
| `location`  | Coordinates | Coordinates of the cooperative |

### POST /cooperatives/ Response

| Property    | Type        | Description                    |
| ----------- | ----------- | ------------------------------ |
| `id`        | ObjectId    | Autogenerated cooperative id   |
| `name`      | string      | The name of the cooperative    |
| `address`   | string      | The address of the cooperative |
| `taxIdCode` | string      | The cooperative C.I.F code     |
| `location`  | Coordinates | Coordinates of the cooperative |

### Coordinates

| Property | Type   | Description                                  |
| -------- | ------ | -------------------------------------------- |
| `lat`    | number | The latitude of the cooperative coordinates  |
| `lng`    | number | The longitude of the cooperative coordinates |

## GET /cooperatives/{cooperativeId}

### GET /cooperatives/{cooperativeId} Description

Fetches the cooperative that matches the cooperativeId and returns it.

### GET /cooperatives/{cooperativeId} Response

| Property    | Type        | Description                    |
| ----------- | ----------- | ------------------------------ |
| `id`        | ObjectId    | Autogenerated cooperative id   |
| `name`      | string      | The name of the cooperative    |
| `address`   | string      | The address of the cooperative |
| `taxIdCode` | string      | The cooperative C.I.F code     |
| `location`  | Coordinates | Coordinates of the cooperative |

## POST /cooperatives/{cooperativeId}

### POST /cooperatives/{cooperativeId} Description

Updates the cooperative that matches the cooperativeId passed with the properties of the body and returns it.

### POST /cooperatives/{cooperativeId} Body Parameters

| Property    | Type        | Description                    |
| ----------- | ----------- | ------------------------------ |
| `name`      | string      | The name of the cooperative    |
| `address`   | string      | The address of the cooperative |
| `taxIdCode` | string      | The cooperative C.I.F code     |
| `location`  | Coordinates | Coordinates of the cooperative |

### POST /cooperatives/{cooperativeId} Response

| Property    | Type        | Description                    |
| ----------- | ----------- | ------------------------------ |
| `id`        | ObjectId    | Autogenerated cooperative id   |
| `name`      | string      | The name of the cooperative    |
| `address`   | string      | The address of the cooperative |
| `taxIdCode` | string      | The cooperative C.I.F code     |
| `location`  | Coordinates | Coordinates of the cooperative |

## DELETE /cooperatives/{cooperativeId}

### DELETE /cooperatives/{cooperativeId} Description

Updates the cooperative that matches the cooperativeId changing its active to false so the cooperative will no longer be fetched.

## GET /cooperatives/{cooperativeId}/visitCards

### GET /cooperatives/{cooperativeId}/visitCards Description

Fetches the cooperative's visitCards that matches the cooperativeId and returns them.

### GET /cooperatives/{cooperativeId}/visitCards Response

#### Array of

| Property        | Type              | Description                           |
| --------------- | ----------------- | ------------------------------------- |
| `id`            | ObjectId          | Autogenerated visitCard id            |
| `name`          | string            | The name of the visitCard             |
| `samplings`     | Array of Sampling | The samplings of the visitCard        |
| `cooperativeId` | string            | The id of the visitCard's cooperative |

### Sampling

| Property            | Type            | Description                                                         |
| ------------------- | --------------- | ------------------------------------------------------------------- |
| `name`              | string          | The name of the sampling                                            |
| `seasonDescription` | string          | The description of the sampling season                              |
| `items`             | Array of string | The headers of the inputs that will be filled when creating a visit |

## GET /cooperatives/{cooperativeId}/workers

### GET /cooperatives/{cooperativeId}/workers Description

Fetches the cooperative's Users with role worker that matches the cooperativeId and returns them.

### GET /cooperatives/{cooperativeId}/workers Response

#### Array of User

# DeclaredCultivationCrops

## Overview

DeclaredCultivationCrops are the entity that holds the type of crops a plot can have. The data in this is loaded from the PAC so normally the Delete, Update and Create won't be used

## GET /declaredCultivationsCrops/

Fetches all the declared cultivation crops.

### GET /declaredCultivationsCrops/ Response

#### Array of

| Property                  | Type     | Description                                                                                   |
| ------------------------- | -------- | --------------------------------------------------------------------------------------------- |
| `id`                      | ObjectId | Autogenerated declaredCultivationsCrop id                                                     |
| `name`                    | string   | The name of the declaredCultivationsCrop                                                      |
| `codeDeclaredCultivation` | number   | The code of the declaredCultivationsCrop, this is unique                                      |
| `season`                  | string   | The season of the declaredCultivationsCrop                                                    |
| `colorCode`               | string   | The color code of the declaredCultivationCrop, this will color the plots which have this crop |

## POST /declaredCultivationsCrops/

### POST /declaredCultivationsCrops/ Description

Creates a new DeclaredCultivationCrop.

### POST /declaredCultivationsCrops/ Body

| Property                  | Type     | Description                                                                                   |
| ------------------------- | -------- | --------------------------------------------------------------------------------------------- |
| `name`                    | string   | The name of the declaredCultivationsCrop                                                      |
| `codeDeclaredCultivation` | number   | The code of the declaredCultivationsCrop, this is unique                                      |
| `season`                  | string   | The season of the declaredCultivationsCrop                                                    |
| `colorCode`               | string   | The color code of the declaredCultivationCrop, this will color the plots which have this crop |

# VariatyCrops

## VariatyCrops overview

VariatyCrops are the entity that holds the varieties a crop can have. The data in this is loaded from the PAC so normally the Delete, Update and Create won't be used

## GET /variatyCrops/

### GET /variatyCrops/ Description

Fetches all the crops' varieties.

### GET /variatyCrops/ Response

| Property              | Type     | Description                                                  |
| --------------------- | -------- | ------------------------------------------------------------ |
| `id`                  | ObjectId | Autogenerated variatyCrop id                                 |
| `name`                | string   | The name of the variatyCrop                                  |
| `codeVariaty`         | number   | The code of the variatyCrop, this is unique                  |
| `declaredCultivation` | string   | The name of the declaredCultivationsCrop the variety is from |

## POST /variatyCrops/

### POST /variatyCrops/ Description

Creates a new variaty with the body sent

### POST /variatyCrops/ Body

| Property              | Type     | Description                                                  |
| --------------------- | -------- | ------------------------------------------------------------ |
| `name`                | string   | The name of the variatyCrop                                  |
| `codeVariaty`         | number   | The code of the variatyCrop, this is unique                  |
| `declaredCultivation` | string   | The name of the declaredCultivationsCrop the variety is from |

## POST /variatyCrops/{variatyId}

### POST /variatyCrops/{variatyId} Description

Updates the variaty whose id matches the one in the call with the body sent

### POST /variatyCrops/{variatyId} Body

| Property              | Type     | Description                                                  |
| --------------------- | -------- | ------------------------------------------------------------ |
| `name`                | string   | The name of the variatyCrop                                  |
| `codeVariaty`         | number   | The code of the variatyCrop, this is unique                  |
| `declaredCultivation` | string   | The name of the declaredCultivationsCrop the variety is from |

## DELETE /variatyCrops/{variatyId}

### DELETE /variatyCrops/{variatyId} Description

Updates the variaty whose id matches the one in the call changing their active property to false so it will no longer be fetched

## GET /variatyCrops/cultivation/

### GET /variatyCrops/cultivation/ Description

Fetches all the crops' varieties whose declaredCultivation matches the one sent in the call.

# Files

## Overview

Files are uploaded and downloaded using these methods

## POST /files/

### POST /files/ Description

Method to upload file into the designated folder. It returns the uploaded file

### POST /files/ Response

| Property     | Type     | Description                    |
| ------------ | -------- | ------------------------------ |
| `id`         | ObjectId | Autogenerated file id          |
| `fileName`   | string   | The name of the file           |
| `uploadedAt` | date     | the date the file was uploaded |

## GET /files/{fileId}

### GET /files/{fileId} Description

This is the route to download the file which fileId matches the one passed in the call

# Plantations

## Overview

Plantations are entities that hold the plantation information

## POST /plantations/

### POST /plantations/ Description

Method to create a plantation

### POST /plantations/ Body Parameters

| Property        | Type              | Description                                                |
| --------------- | ----------------- | ---------------------------------------------------------- |
| `name`          | string            | name of the plantation                                     |
| `cooperativeId` | ObjectId          | id of the cooperative that works on the plantation         |
| `ownerId`       | Array of ObjectId | ids of users that can manage the plantation ands its plots |
| `ownerPlot`     | string            | name of the owner of the plantation, can be empty          |

## GET /plantations/

### GET /plantations/ Description

Fetches the plantations that the currentUser has access to

## GET /plantations/plots

### GET /plantations/plots Description

Same as GET `/plantations/` but every plantation has a new Array property with their plots

## POST /plantations/{plantationId}

### POST /plantations/{plantationId} Description

Updates a plantation whose id matches the one given

### POST /plantations/{plantationId} Body Parameters

| Property        | Type              | Description                                                |
| --------------- | ----------------- | ---------------------------------------------------------- |
| `name`          | string            | name of the plantation                                     |
| `cooperativeId` | ObjectId          | id of the cooperative that works on the plantation         |
| `ownerId`       | Array of ObjectId | ids of users that can manage the plantation ands its plots |
| `ownerPlot`     | string            | name of the owner of the plantation, can be empty          |

## DELETE /plantations/{plantationId}

### DELETE /plantations/{plantationId} Description

Updates the plantation that matches the plantationId changing its active to false so the plantation will no longer be fetched.

## GET /plantations/{plantationId}/issues

### GET /plantations/{plantationId}/issues Description

Fetches the sent issues (excluding the finished ones) located in the plantation whose id matches the one given.

## GET /plantations/{plantationId}/issues/finished

### GET /plantations/{plantationId}/issues/finished Description

Fetches the finished issues located in the plantation whose id matches the one given.

## GET /plantations/{plantationId}/sigpacPlots

### GET /plantations/{plantationId}/sigpacPlots Description

Fetches the sigpac plots that have the plot codes of the plots of the plantation located in the plantation whose id matches the one given.

## GET /plantations/{plantationId}/sigpacPlotsWithNoUser

### GET /plantations/{plantationId}/sigpacPlotsWithNoUser Description

Fetches the sigpac plots that have the plot codes of the plots of the plantation located in the plantation whose id matches the one given. This method is used in the webview for the maps on the mobile app

## GET /plantations/{plantationId}/myIssues

### GET /plantations/{plantationId}/myIssues Description

Fetches the not yet sent issues which plantationId matches the one in the call

# Plots

## Overview

Plots are entities that hold the plots information.

## POST /plots/

### POST /plots/ Description

Create a new plot

### POST /plots/ Body Parameters

| Property         | Type     | Description                                                     |
| ---------------- | -------- | --------------------------------------------------------------- |
| `name`           | string   | name of the plot                                                |
| `plantationId`   | ObjectId | id of the plantation the plot is located                        |
| `plotCode`       | PlotCode | identifier that is unique for every plot                        |
| `cropsCode`      | number   | code of the crop that is going to be in the plot                |
| `varietyCode`    | number   | code of the variety of the crop that is going to be in the plot |
| `supDecl`        | number   | area of the plot, this data comes from the sigpacPlots          |
| `presAlegSigpac` | string   | sigpac code of the plot, this data comes from the sigpacPlots   |

### PlotCode

| Property      | Type   | Description      |
| ------------- | ------ | ---------------- |
| `provinceRec` | number | province code    |
| `townRec`     | number | town code        |
| `polygonRec`  | number | polyygon code    |
| `plotNum`     | number | plot number      |
| `enclosure`   | number | enclosure number |

## GET /plots/{plotId}

### GET /plots/{plotId} Description

Fetches the plot whose id matches the one given in the call

## GET /plots/{plotId} Response

| Property         | Type     | Description                                                     |
| ---------------- | -------- | --------------------------------------------------------------- |
| `id`             | ObjectId | Autogenerated plot id                                           |
| `name`           | string   | Name of the plot                                                |
| `plantationId`   | ObjectId | Id of the plantation the plot is located                        |
| `plotCode`       | PlotCode | identifier that is unique for every plot                        |
| `cropsCode`      | number   | Code of the crop that is going to be in the plot                |
| `varietyCode`    | number   | Code of the variety of the crop that is going to be in the plot |
| `supDecl`        | number   | Area of the plot, this data comes from the sigpacPlots          |
| `presAlegSigpac` | string   | Sigpac code of the plot, this data comes from the sigpacPlots   |

## POST /plots/{plotId}

### POST /plots/{plotId} Description

Updates the plot whose id matches the one given in the call

### POST /plots/{plotId} Body Parameters

| Property      | Type   | Description                                                     |
| ------------- | ------ | --------------------------------------------------------------- |
| `name`        | string | name of the plot                                                |
| `cropsCode`   | number | code of the crop that is going to be in the plot                |
| `varietyCode` | number | code of the variety of the crop that is going to be in the plot |

## DELETE /plots/{plotId}

### DELETE /plots/{plotId} Description

Updates the plot whose id matches the one given in the call changing active to false so it will no longer be fetched by the GET methods

## GET /plots/{plotId}/sigpac

### GET /plots/{plotId}/sigpac Description

Fetches the sigpacPlot with the same plotCode as the plot whose id matches the one given in the call

## GET /plots/{plotId}/sigpacPlotWithNoUser

### GET /plots/{plotId}/sigpacPlotWithNoUser Description

Fetches the sigpacPlot with the same plotCode as the plot whose id matches the one given in the call. This call is used to print the plots on the webview map of the mobile app

## GET /plot/plantation/{plantationId}

### GET /plot/plantation/{plantationId} Description

Fetches the plot whose plantationId matches the one sent in the call.

## GET /plot/plantation/{plantationId}/geo

### GET /plot/plantation/{plantationId}/geo Description

Fetches the plot whose plantationId matches the one sent in the call with a new property `geoJSON` that has the data necesary to print the plot on a map

# Issues

## Overview

Issues are entities that hold the issues that are encountered on the field. These entities increase the number of parameters when they are updated to a new state until the issue is finished

## POST /issues/

### POST /issues/ Description

Method to generate an issue. An issue can be created at three different stages: from scratch, with a visit, or with a task.
  Needed data that is not filled will autocomplete.
  
#### Initial stage issue Body

| Property            | Type          | Description                                                               |
| ------------------- | ------------- | ------------------------------------------------------------------------- |
| `id`                | ObjectId      | Autogenerated Issue id                                                    |
| `pictures`          | Array of File | Photo files of the issue                                                  |
| `cropsType`         | string        | the name of the crop that has the issue                                   |
| `cropsVariety`      | string        | the variety of the crop that has the issue (empty if there's none)        |
| `location`          | IssueLocation | has the info of where issue is located                                    |
| `issuer`            | User          | the user that generated the issue                                         |
| `issueNumber`       | number        | identifying isssue number                                                 |
| `issueComment`      | string        | a brief description of the issue                                          |
| `status`            | string        | the current status of the issue                                           |
| `sendToCooperative` | bool          | specifies if the issue is stored by the issuer or sent to the cooperative |

#### Issue possible statuses

| Statuses                    |
| --------------------------- |
| STATUS_CREATED              |
| STATUS_PENDING_VISIT        |
| STATUS_PENDING_WORKORDER    |
| STATUS_WORK_ORDER           |
| STATUS_ASSIGNED             |
| STATUS_WORKING              |
| STATUS_PAUSED               |
| STATUS_PENDING_ASSESSMENT   |
| STATUS_FINISHED             |
| STATUS_CANCELLED            |
| STATUS_DELAYED              |

#### IssueLocation

| Property         | Type        | Description                                       |
| ---------------- | ----------- | ------------------------------------------------- |
| `plantationId`   | ObjectId    | id of the plantation where the issue is located   |
| `plotId`         | ObjectId    | id of the plot where the issue is located         |
| `plantationName` | string      | name of the plantation where the issue is located |
| `coordinates`    | Coordinates | coordinates where the issue was generated         |
| `plotCode`       | PlotCode    | code of the plot where the issue is located       |

## GET /issues/

### GET /issues/ Description

Fetches the issues the currentUser has access to.

### GET /issues/ Response

The response's body depends on the status of the issue, the properties on every status are specified on the document `Issue status evolution.odt`

| Property            | Type          | Description                                                               |
| ------------------- | ------------- | ------------------------------------------------------------------------- |
| `id`                | ObjectId      | Autogenerated Issue id                                                    |
| `pictures`          | Array of File | Photo files of the issue                                                  |
| `cropsType`         | string        | the name of the crop that has the issue                                   |
| `cropsVariety`      | string        | the variety of the crop that has the issue (empty if there's none)        |
| `location`          | IssueLocation | has the info of where issue is located                                    |
| `issuer`            | User          | the user that generated the issue                                         |
| `issueNumber`       | number        | identifying isssue number                                                 |
| `issueComment`      | string        | a brief description of the issue                                          |
| `status`            | string        | the current status of the issue                                           |
| `sendToCooperative` | bool          | specifies if the issue is stored by the issuer or sent to the cooperative |
| `visitInfo`         | VisitInfo     | information on the visit to do or done                                    |
| `workOrderInfo`     | WorkOrderInfo | information on the task to do or done                                     |

#### VisitInfo

| Property         | Type                                 | Description                                                         |
| ---------------- | ------------------------------------ | ------------------------------------------------------------------- |
| `visitDate`      | date                                 | Date the visit takes place                                          |
| `technicianId`   | ObjectId                             | Id of the technician that visits                                    |
| `technician`     | User                                 | The technician that visits                                          |
| `needsWorkOrder` | bool                                 | Determines if the visit needs to generate a WorkOrder               |
| `templateId`     | ObjectId                             | id of the visitCard that was filled                                 |
| `data`           | Map of string, array of SamplingItem | the data that was filled in the visitCard                           |
| `visitComment`   | string                               | a brief description of the visit or any comment made by the visitor |

#### WorkOrderInfo

| Property               | Type                      | Description                                                          |
| ---------------------- | ------------------------- | -------------------------------------------------------------------- |
| `workDate`             | date                      | Date the task takes place                                            |
| `workDateId`           | ObjectId                  | Id of the worker that carries the task out                           |
| `worker`               | User                      | The worker that carries the task out                                 |
| `treatment`            | string                    | The treatmentType the the task has                                   |
| `treatmentDescription` | Map of string, string     | the treatment parameters data filled                                 |
| `comments`             | string                    | comments made by the worker or the person issuing the worker to work |
| `workOrderTimer`       | WorKOrderTimer            | timer that tracks down the minutes the task being worked on          |
| `totalCost`            | number                    | The total cost of the task                                           |
| `accepted`             | bool                      | Determines if the task is accepted by the client or not              |
| `products`             | Array of WorkOrderProduct | If needed, holds the products needed to carry the task out           |
| `completed`            | bool                      | Determines if the task is completed or not                           |

#### WorkOrderTimer

| Property        | Type     | Description                        |
| --------------- | -------- | ---------------------------------- |
| `startTime`     | time     | Time the task started being worked |
| `timesPaused`   | number   | times the task has been paused     |
| `lastPauseTime` | time     | last time the timer was stopped    |
| `timeWorking`   | duration | the total duration of the task     |
| `restartTime`   | time     | the time the task was restarted    |

#### WorkOrderProduct

| Property               | Type                  | Description                             |
| ---------------------- | --------------------- | --------------------------------------- |
| `productId`            | ObjectId              | the id of the product that will be used |
| `productToUse`         | Product               | the product that will be used           |
| `treatmentDescription` | Map of string, string | the treatment parameters data filled    |

#### Product

## POST /issues/{issueId}

### POST /issues/{issueId} Description

Updates the issue whose id matches the one given in the call, the status will change in consequence of what's given and will autofill not passed data

#### POST /issues/{issueId} Body

| Property            | Type          | Description                                                               |
| ------------------- | ------------- | ------------------------------------------------------------------------- |
| `id`                | ObjectId      | Autogenerated Issue id                                                    |
| `pictures`          | Array of File | Photo files of the issue                                                  |
| `cropsType`         | string        | the name of the crop that has the issue                                   |
| `cropsVariety`      | string        | the variety of the crop that has the issue (empty if there's none)        |
| `location`          | IssueLocation | has the info of where issue is located                                    |
| `issuer`            | User          | the user that generated the issue                                         |
| `issueNumber`       | number        | identifying isssue number                                                 |
| `issueComment`      | string        | a brief description of the issue                                          |
| `status`            | string        | the current status of the issue                                           |
| `sendToCooperative` | bool          | specifies if the issue is stored by the issuer or sent to the cooperative |
| `visitInfo`         | VisitInfo     | information on the visit to do or done                                    |
| `workOrderInfo`     | WorkOrderInfo | information on the task to do or done                                     |

## GET /issues/{issueId}

### GET /issues/{issuesId} Description

Fetches the issue whose id matches the one sent in the call.

### GET /issues/{issuesId} Body

The response's body depends on the status of the issue, the properties on every status are specified on the document `Issue status evolution.odt`

## POST /issues/{issueId}/work

### POST /issues/{issueId}/work Description

Updates the WorkOrderTimer to start or restart working on the issue's task

## POST /issues/{issueId}/pause

### POST /issues/{issueId}/pause Description

Updates the WorkOrderTimer to pause the task

## POST /issues/{issueId}/finish

### POST /issues/{issueId}/finish Description

Updates the issue to mark it as finished.

## POST /issues/{issueId}/accept

### POST /issues/{issueId}/accept Description

Updates the WorkOrder to mark the task as accepted

## GET /issues/finished

### GET /issues/finished Description

Fetches the finished issues the currentUser has access to.

## DELETE /issues/{issuesId}

### DELETE /issues/{issuesId} Description

Updates the issue whose id matches the one given in the call to turn active to false so the issue will no longer be fetched

## POST /issues/task

### POST /issues/task Description

Creates an issue in task state autofilling needed fields. I'm not quite sure if this method is being used since I didn't create it.

## GET /issues/states

### GET /issues/states Description

Returns an array of string containing all the different issues state ids.

## GET /issues/treatments

### GET /issues/treatments Description

Returns an array of string containing all the different treatment ids.

# Employees

## Employees overview

Employees holds User entity that aren't ROLE_FARMER

## GET /employees/

### GET /employees/ Description

Fetches the users whose role isn't role farmer

## GET /employees/cooperative/{cooperativeId}

### GET /employees/cooperative/{cooperativeId} Description

Fetches the users whose role isn't role farmer and cooperativeId matches the one given in the call

# Params

## Params overview

This Entity associates the treatments with the data that will be filled when selecting the treatment

## GET /params/

### GET /params/ Description

Fetches every Entity and returns them

### GET /params/ Response

| Property                 | Type     | Description                                                                 |
| ------------------------ | -------- | --------------------------------------------------------------------------- |
| `id`                     | ObjectId | Autogenerated Param id                                                      |
| `name`                   | string   | name of the field that will be filled                                       |
| `treatmentCultivation`   | bool     | wether or not this field will appear when selecting treatment Cultivation   |
| `treatmentSowing`        | bool     | wether or not this field will appear when selecting treatment Sowing        |
| `treatmentPlant`         | bool     | wether or not this field will appear when selecting treatment Plant         |
| `treatmentFertilization` | bool     | wether or not this field will appear when selecting treatment Fertilization |
| `treatmentPruning`       | bool     | wether or not this field will appear when selecting treatment Pruning       |
| `treatmentPhytosanitary` | bool     | wether or not this field will appear when selecting treatment Phytosanitary |
| `treatmentIrrigation`    | bool     | wether or not this field will appear when selecting treatment Irrigation    |
| `treatmentTwigsRemoval`  | bool     | wether or not this field will appear when selecting treatment TwigsRemoval  |
| `treatmentHarvesting`    | bool     | wether or not this field will appear when selecting treatment Harvesting    |
| `treatmentAnalysis`      | bool     | wether or not this field will appear when selecting treatment Analysis      |

## GET /params/{treatment}

### GET /params/{treatment} Description

Fetches every param that has true in the treatment passed in the call

### GET /params/{treatment} Body

| Property                 | Type     | Description                                                                 |
| ------------------------ | -------- | --------------------------------------------------------------------------- |
| `id`                     | ObjectId | Autogenerated Param id                                                      |
| `name`                   | string   | name of the field that will be filled                                       |
| `treatmentCultivation`   | bool     | wether or not this field will appear when selecting treatment Cultivation   |
| `treatmentSowing`        | bool     | wether or not this field will appear when selecting treatment Sowing        |
| `treatmentPlant`         | bool     | wether or not this field will appear when selecting treatment Plant         |
| `treatmentFertilization` | bool     | wether or not this field will appear when selecting treatment Fertilization |
| `treatmentPruning`       | bool     | wether or not this field will appear when selecting treatment Pruning       |
| `treatmentPhytosanitary` | bool     | wether or not this field will appear when selecting treatment Phytosanitary |
| `treatmentIrrigation`    | bool     | wether or not this field will appear when selecting treatment Irrigation    |
| `treatmentTwigsRemoval`  | bool     | wether or not this field will appear when selecting treatment TwigsRemoval  |
| `treatmentHarvesting`    | bool     | wether or not this field will appear when selecting treatment Harvesting    |
| `treatmentAnalysis`      | bool     | wether or not this field will appear when selecting treatment Analysis      |

## POST /params/

### POST /params/ Decription

Creates a new param

### POST /params/ Body

| Property                 | Type     | Description                                                                 |
| ------------------------ | -------- | --------------------------------------------------------------------------- |
| `id`                     | ObjectId | Autogenerated Param id                                                      |
| `name`                   | string   | name of the field that will be filled                                       |
| `treatmentCultivation`   | bool     | wether or not this field will appear when selecting treatment Cultivation   |
| `treatmentSowing`        | bool     | wether or not this field will appear when selecting treatment Sowing        |
| `treatmentPlant`         | bool     | wether or not this field will appear when selecting treatment Plant         |
| `treatmentFertilization` | bool     | wether or not this field will appear when selecting treatment Fertilization |
| `treatmentPruning`       | bool     | wether or not this field will appear when selecting treatment Pruning       |
| `treatmentPhytosanitary` | bool     | wether or not this field will appear when selecting treatment Phytosanitary |
| `treatmentIrrigation`    | bool     | wether or not this field will appear when selecting treatment Irrigation    |
| `treatmentTwigsRemoval`  | bool     | wether or not this field will appear when selecting treatment TwigsRemoval  |
| `treatmentHarvesting`    | bool     | wether or not this field will appear when selecting treatment Harvesting    |
| `treatmentAnalysis`      | bool     | wether or not this field will appear when selecting treatment Analysis      |

## DELETE /params/{paramId}

### DELETE /params/{paramId} Description

Updates the param that matches in id with the id passed in the call changing their active to false so it will no longer be fetched

## POST /params/{paramId}

### POST /params/{paramId} Description

Updates the param that matches in id with the id passed in the call

# Products

## Products overview

Entity that holds the product information. These data was imported from a db passed to us.

## GET /products/

### GET /products/ Description

Fetches every product stored

### GET /products/ Response

| Property         | Type     | Description                     |
| ---------------- | -------- | ------------------------------- |
| `id`             | ObjectId | Autogenerated Product id        |
| `name`           | string   | name of the product             |
| `formula`        | string   | formula of the product          |
| `holder`         | string   | holder of the product's formula |
| `registryNumber` | number   |                                 |

## POST /products/

### POST /products/ Description

Creates a new product and inserts it into the database

### POST /products/ Body

| Property         | Type   | Description                     |
| ---------------- | ------ | ------------------------------- |
| `name`           | string | name of the product             |
| `formula`        | string | formula of the product          |
| `holder`         | string | holder of the product's formula |
| `registryNumber` | number |                                 |

## DELETE /products/{productId}

### DELETE /products/{productId} Description

Updates the product that matches in id with the id passed in the call changing their active to false so it will no longer be fetched

# Provinces

## Provinces overview

Entity that holds the provinces information

## GET /provinces/

### GET /provinces/ Description

Fetches every province from the database

### GET /provinces/ Response

| Property | Type     | Description               |
| -------- | -------- | ------------------------- |
| `id`     | ObjectId | Autogenerated Province id |
| `name`   | string   | name of the province      |
| `code`   | number   | code of the product       |

# Town

## Town overview

Entity that holds the towns information

## GET /towns/

### GET /towns/ Description

Fetches the towns that matches the province number passed in the call

### GET /towns/ Response

| Property       | Type     | Description                 |
| -------------- | -------- | --------------------------- |
| `id`           | ObjectId | Autogenerated Town id       |
| `name`         | string   | name of the town            |
| `code`         | number   | code of the town            |
| `province`     | number   | code of the town's province |
| `provinceName` | string   | name of the town            |

# SigpacCrops

## SigpacCrops overview

Entity that holds the sigpacCrops information. This data was imported from sigpac.

## GET /sigpacCrops/

### GET /sigpacCrops/ Description

Fetches every sigpac crops from the database

### GET /sigpacCrops/ Response

| Property            | Type     | Description                 |
| ------------------- | -------- | --------------------------- |
| `id`                | ObjectId | Autogenerated sigpacCrop id |
| `code`              | string   | code of the sigpac crop     |
| `sigpacDescription` | number   | description given by sigpac |

## GET /sigpacCrops/{sigpacCropId}

### GET /sigpacCrops/{sigpacCropId} Description

Fetches the sigpac crop whose id matches the one passed in the call

# SigpacPlots

## SigpacPlots overview

Entity that holds the sigpacPlots information. This data was imported from sigpac.

## GET /sigpacPlots/

### GET /sigpacPlots/ Description

Fetches the sigpac plot that matches the provinceRec, townRec, polygonNum, plotNum and enclosure passed in the call

### GET /sigpacPlots/ Response

| Property      | Type     | Description                                            |
| ------------- | -------- | ------------------------------------------------------ |
| `id`          | ObjectId | Autogenerated sigpacPlot id                            |
| `gid`         | string   |                                                        |
| `sigpacId`    | number   | id given by sigpac                                     |
| `provinceRec` | number   | province code                                          |
| `townRec`     | number   | town code                                              |
| `polygonNum`  | number   | polygon code                                           |
| `plotNum`     | number   | plot code                                              |
| `enclosure`   | number   | enclosure code                                         |
| `usageCode`   | string   |                                                        |
| `area`        | number   | area of the plot                                       |
| `region`      | string   |                                                        |
| `gc`          | string   |                                                        |
| `version`     | number   |                                                        |
| `geoJSON`     | geoJSON  | this geoJSON holds the plot's polygon shown on the map |

# Statistics

## GET /statistics/{plantationId}

### GET /statistics/{plantationId} Description

Fetches the statistics of the plantation whose id matches the one passed on the call

### GET /statistics/{plantationId} Response

| Property         | Type   | Description                                                         |
| ---------------- | ------ | ------------------------------------------------------------------- |
| `totalCost`      | number | the total sum of the cost of every finished issue of the plantation |
| `finishedIssues` | number | the total number of finished issues of the plantation               |

# VisitCards

## VisitCards overview

Entity that holds the visitCards template information.

## GET /visitCards/

### GET /visitCards/ Description

Fetches all the visitCards in the db.

### GET /visitCards/ Response

| Property        | Type              | Description                            |
| --------------- | ----------------- | -------------------------------------- |
| `id`            | ObjectId          | Autogenerated visitCard id             |
| `samplings`     | Array of Sampling | The samplings of the visitCard         |
| `cooperativeId` | ObjectId          | id of the cooperative of the visitCard |

## GET /visitCards/{visitCardId}

### GET /visitCards/{visitCardId} Description

Fetches the visitCards whose id matches the one in the call.

### GET /visitCards/ Response

| Property        | Type              | Description                            |
| --------------- | ----------------- | -------------------------------------- |
| `id`            | ObjectId          | Autogenerated visitCard id             |
| `samplings`     | Array of Sampling | The samplings of the visitCard         |
| `cooperativeId` | ObjectId          | id of the cooperative of the visitCard |

## POST /visitCards/

### POST /visitCards/ Description

Creates the visitCard with the body sent.

### POST /visitCards/ Body

| Property        | Type              | Description                            |
| --------------- | ----------------- | -------------------------------------- |
| `samplings`     | Array of Sampling | The samplings of the visitCard         |
| `cooperativeId` | ObjectId          | id of the cooperative of the visitCard |

## POST /visitCards/{visitId}

### POST /visitCards/{visitId} Description

Updates the visitCard with whose id matches the one in the call the body sent.

### POST /visitCards/{visitId} Body

| Property        | Type              | Description                            |
| --------------- | ----------------- | -------------------------------------- |
| `samplings`     | Array of Sampling | The samplings of the visitCard         |
| `cooperativeId` | ObjectId          | id of the cooperative of the visitCard |

## DELETE /visitCards/{visitId}

### DELETE /visitCards/{visitId} Description

Updates the visitCard with whose id matches the one in the call changing their active prop to false so it will no longer be fetched.

