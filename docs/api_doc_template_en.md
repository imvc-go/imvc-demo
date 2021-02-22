#XXX System API Documentation

| Date | Description | Reviser |
| ---- | ---- | ------ |
|      |      |        |
|      |      |        |

[toc]



## 1. Sketch

This is a sketch of this system.

**API Address：**

* Stage URL：`https://stage.xxx.com` or `grpc://172.17.0.1:8000`
* Production URL：`https://prod.xxx.com`

**Matters need attention：**

* replace the matters need attention here



## 2. API

### 1 Register API 

**sketch of this API：** 

- API for sign up new users 

**URL Path：** 
- ` api/user/register `
  

**Request method：**
- POST 

**Form or parameters：** 

|Parameter|Required|Type|Description|
|:----    |:---|:----- |-----   |
|username |Yes |string | user name   |
|password |Yes |string | user password    |
|name     |No  |string | user nickname    |

 **Response structure**

``` json
  {
    "error_code": 0,
    "data": {
      "uid": "1",
      "username": "12154545",
      "name": "user1",
      "groupid": 2 ,
      "reg_time": "1436864169",
      "last_login_time": "0",
    }
  }
```

 **Response field description** 

|Field|Type|Description|
|:-----  |:-----|-----                           |
|groupid |int   |user group id，1：super user；2：normal user ... |

 **Note** 

- If there have some notes of this API, add it here.


