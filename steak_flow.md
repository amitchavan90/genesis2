# Product Registration Flow (rev g) STEAK API

A registered product a single steak.

There will be 4 steps.

0. View: Launch Mini Program, view product, get distributor code in URL
1. Detail: Customer info QR code below the meat in WeChat, get meat info, close chain, points available
2. Close chain: Customer login to WeChat, mini program will communicate with L28 API backend, get points and details
3. Final confirmation: Customer sees thank you page inside WeChat

## Notes

- URLs below are all for test/staging
- `steakID` is unique for each product, it is used to make sure QR1 and QR2 are scanned in order
- `product_id` is now `productID`

## Step 0 (Launch WeChat Mini Program, view product) (QR1)

Execute a GET request using the QR code:

- URL: https://consumer.gn.staging.latitude28.cn/api/steak/view?productID=295928d1-1a08-4bda-a966-a2f815ce9b5f

On Success:

- Server response code: 302 Found
- Redirect to /view? with `productID` and `distributorCode` parameters

On fail:

- Server response code: 4XX or 5XX

MP will now GET with new redirected URL

- URL: https://consumer.gn.staging.latitude28.cn/view?productID=50ea8ee0-9725-4603-b559-c4ab6628ee9f&distributorCode=ABCD

On Success:

- Server response code: 200 OK
- Static HTML Product Info page

On fail:

- Server response code: 4XX or 5XX

## Step 1 (Retrive product points and details and close) (QR2)

Execute a GET request using the QR code:

- URL: https://consumer.gn.staging.latitude28.cn/api/steak/detail?steakID=295928d1-1a08-4bda-a966-a2f815ce9b5f&productID=50ea8ee0-9725-4603-b559-c4ab6628ee9f

Note:

- `closeID` maybe present depending if the product is already registered or not
- `productID` is appended to the URL by the MP
- Server will verify the `productID` ties to `steakID`, if it fails, MP need to advice Customer that it is wrong QR2 and start from beginning and try scan QR1 again
- On success, MP will render page using server JSON data

On success:

- Server response code: 200 OK
  - product not yet registered
    ```json
    {
    	"success": true,
    	"bonusPointsExpiry": "2020-05-31T08:22:11.78Z",
    	"bonusPoints": 10,
    	"basePoints": 5,
    	"distributorCode": "ABCD", // could be empty if not exist
    	"closeID": "02e7adc0-923f-4c18-88a6-03121cf8d728"
    }
    ```
  - product already registered
    ```json
    {
    	"success": true,
    	"bonusPointsExpiry": "2020-05-31T08:22:11.78Z",
    	"bonusPoints": 10,
    	"basePoints": 5,
    	"distributorCode": "ABCD"
    }
    ```

On fail:

- Server response code: 4XX or 5XX
- success will show as false, `message` may be present
- Response Body:
  ```json
  {
  	"success": false,
  	"message": "invalid steak uuid"
  }
  ```

## Step 2 (Close chain from inside WeChat)

Call Genesis API with the steakID token received from the previous step in order to:

- Close the chain
- Receive the number of points that steak is worth

Note: MP and the server will make sure the closeID is the correct one for the steakID, or it will stop/fail.

Execute a POST request:

- URL: https://consumer.gn.staging.latitude28.cn/api/steak/close
- Request Body:
  ```json
  {
  	"steakID": "295928d1-1a08-4bda-a966-a2f815ce9b5f",
  	"closeID": "02e7adc0-923f-4c18-88a6-03121cf8d728",
  	"weChatID": "beef-fan-888"
  }
  ```

On success:

- Server response code: 200 OK
- success will show as true, `message` may be present.
- Response Body:
  ```json
  {
  	"success": true,
  	"steakID": "295928d1-1a08-4bda-a966-a2f815ce9b5f",
  	"customerID": "426fd0ab-4115-4f82-bba9-268948ad2073",
  	"weChatID": "beef-fan-888",
  	"pointsGiven": 15 // base points + bonus points
  }
  ```

On fail:

- Server response code: 4XX or 5XX
- success will show as false, `message` may be present.
- Response Body:
  ```json
  {
  	"success": false,
  	"message": "Product already registered"
  }
  ```

## Step 3 (Final confirmation)

If desired, call Genesis with a GET request to see the success page and present it to the user.

Execute a GET request:

- URL: https://consumer.gn.staging.latitude28.cn/api/steak/final?steakID=295928d1-1a08-4bda-a966-a2f815ce9b5f
- Note: if there is `&demo=true` in url parameter, it will bypass the product register check

Server response:

On success:

- Server response code: 200 OK
- Static HTML confirmation page

On fail:

- Server response code: 4XX or 5XX
- Response Body:
  ```json
  {
  	"success": false,
  	"message": "cannot find product"
  }
  ```
