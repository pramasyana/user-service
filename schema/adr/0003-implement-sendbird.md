# 3.Sendbird - Chat Engine

Date: 2021-06-21

## Status 
Proposed

## Expected Condition:
Buyer dapat melakukan komunikasi secara langsung dengan seller, sehingga diharapkan dapat mendorong lebih banyak terjadinya transaksi.

### Background

Buyer tidak dapat melakukan komunikasi langsung dengan seller


### Sequence Diagram - buyer
![Sequence Diagram Buyer](https://i.postimg.cc/QjJsN2H3/diagram-chat.png)

### Sequence Diagram - Seller
![Sequence Diagram Seller](https://i.postimg.cc/M6DPQtP3/flow-integrasi-data-seller-merchantpng.png)



### Access token vs Session Token
Ketika user akan connect ke server sendbird, kita dapat  memilih untuk mengautentikasi user hanya dengan user_id dan dapat juga dengan menambahkan access token atau session token, tetapi yang akan kita implement adalah dengan menggunakan `user_id` dan `session_token` ---> connect (user_id, session_token), kenapa menggunakan session token ? ...  karena  session token dapat diset expired time nya dan synchronize atau samakan dengan expired-time token nya user-service. berikut perbandingan antara access token dan session token dari dokumentasi sendbird : 

|           |      Access Token      |  Session Token |
|---------- |:-------------:|:--------------------------------:|
Used for    | Stateful authentication | Stateless authentication
Work as     |Permanent credential to the system|Temporary credential to the system
Valid or active until|Revoked|Timestamp set when issued (default: the next 7 days from now)
Identification for|The user account|The user's current session
Tokens per user|Up to 10 (valid)|Up to 100 (active)
If exceeded the limit|The oldest token is revoked and the new one is added to the list.|The oldest and active token is revoked and the new one is added to the list.
Auto-revocation|No|Yes (by default the system revokes the expired tokens)


### Get Session Token - Sendbird

- Url 
`GET https://api-{application_id}.sendbird.com/v3/users/{user_id}`


- Body Response -  Get Session Token      
    - ```
        {
          "user_id": "USR20030471",
          "nickname": "Muhamad Rusdi Syahre",
          "unread_message_count": 7,
          "profile_url": "https://sendbird.com/main/img/profiles/profile_11_512px.png",
          "access_token": "800284474d5d94f3e1658c6c0794872bda4f5f72",
          "session_tokens": [
            {
              "session_token": "0c49975b05f0fd90ec60210a356d7ab2d3f1b8c3",
              "expires_at": 2037946156348
            }
          ],
          "is_online": false,
          "is_active": true,
          "created_at": 1581084870,
          "last_seen_at": 1581085696197,
          "has_ever_logged_in": true,
          "metadata": {
            "isMerchant": True,
            "merchantName": "Berkah Store"
          }
        }
        ```
           
    

### Create User - Sendbird

- Url 
`POST https://api-{application_id}.sendbird.com/v3/users`

- Body Request - Create User
    -  ```   
        {
          "user_id": "USR20030471",
          "nickname": "Muhamad Rusdi Syahren",
          "profile_url": "https://sendbird.com/main/img/profiles/profile_11_512px.png",
          "issue_access_token": true,
          "issue_session_token": true,
          "session_token_expires_at": 1542945056625,
          "metadata": {
            "isMerchant": True,
            "merchantName": "Berkah Store"
          }
        }
        ```
    
- Body Response - Create User
    - ```
        {
          "user_id": "USR20030471",
          "nickname": "Muhamad Rusdi Syahre",
          "unread_message_count": 7,
          "profile_url": "https://sendbird.com/main/img/profiles/profile_11_512px.png",
          "access_token": "800284474d5d94f3e1658c6c0794872bda4f5f72",
          "session_tokens": [
            {
              "session_token": "0c49975b05f0fd90ec60210a356d7ab2d3f1b8c3",
              "expires_at": 2037946156348
            }
          ],
          "is_online": false,
          "is_active": true,
          "created_at": 1581084870,
          "last_seen_at": 1581085696197,
          "has_ever_logged_in": true,
          "metadata": {
            "isMerchant": True,
            "merchantName": "Berkah Store"
          }
        }
        ```
        
     
### Update Session Token Expired - Sendbird

- Url 
`PUT https://api-{application_id}.sendbird.com/v3/users/{user_id}`

- Body Request - Update Session Token Expired 
    -  ```   
        {
          "session_token_expires_at": 1542945056625,
        }
        ```
    
- Body Response -  Update Session Token Expired     
    - ```
        {
          "user_id": "USR20030471",
          "nickname": "Muhamad Rusdi Syahre",
          "unread_message_count": 7,
          "profile_url": "https://sendbird.com/main/img/profiles/profile_11_512px.png",
          "access_token": "800284474d5d94f3e1658c6c0794872bda4f5f72",
          "session_tokens": [
            {
              "session_token": "0c49975b05f0fd90ec60210a356d7ab2d3f1b8c3",
              "expires_at": 2037946156348
            }
          ],
          "is_online": false,
          "is_active": true,
          "created_at": 1581084870,
          "last_seen_at": 1581085696197,
          "has_ever_logged_in": true,
          "metadata": {
            "isMerchant": True,
            "merchantName": "Berkah Store"
          }
        }
        ```
    
### Create Channel  - Sendbird

- Body Request - create channel
    -  ```   
        {
            "name": "rusdi-azlan",
            "channel_url": "private_chat_room_424",
            "is_distinct": true,
            "inviter_id": "USR20030471",
            "user_ids": ["USR20030471", "USR20030472"],
            "operator_ids": ["USR20030471"]
        }
        ```
    
- Body Response -  create channel    
    - ```
        {
            "name": "rusdi-azlan",
            "channel_url": "private_chat_room_424",
            "cover_url": "https://sendbird.com/main/img/cover/cover_08.jpg",
            "custom_type": "",
            "unread_message_count": 0,
            "data": "",
            "is_distinct": true,
            "is_public": false,
            "member_count": 3,
            "joined_member_count": 1,
            "members": [
                {
                    "user_id": "USR20030471",
                    "nickname": "rusdi",
                    "profile_url": "https://sendbird.com/main/img/profiles/profile_17_512px.png",
                    "is_active": true,
                    "is_online": false,
                    
                    "last_seen_at": 1530232836311,
                    "state": "joined",
                    "metadata": {
                      "isMerchant": false,
                      "merchantName": ""
                    }
                },
                {
                    "user_id": "USR20030472",
                    "nickname": "azlan",
                    "profile_url": "https://sendbird.com/main/img/profiles/profile_08_512px.png",
                    "is_active": true,
                    "is_online": false,
                    "last_seen_at": 1530237133254,
                    "metadata": {
                      "isMerchant": True,
                      "merchantName": "Berkah Store"
                    }
                }
            ],
            "operators": [
                 "user_id": "USR20030471",
                    "nickname": "rusdi",
                    "profile_url": "https://sendbird.com/main/img/profiles/profile_17_512px.png",
                    "is_active": true,
                    "is_online": false,
                    
                    "last_seen_at": 1530232836311,
                    "state": "joined",
                    "metadata": {
                      "isMerchant": false,
                      "merchantName": ""
                    }
            ],
            "max_length_message": 500,
            "last_message": null,
            "created_at": 1543468122
        }
        ```       


### List Channel  - Sendbird

- parameter Request - list channel
    - `?limit=5&order=latest_last_message&show_member=true`
       
    
- Body Response -  list channel    
    - ```
        {
            "channels": [
                    {
                        "name": "rusdi-azlan",
                        "channel_url": "private_chat_room_424",
                        "cover_url": "https://sendbird.com/main/img/cover/cover_08.jpg",
                        "custom_type": "",
                        "unread_message_count": 0,
                        "data": "",
                        "is_distinct": true,
                        "is_public": false,
                        "member_count": 3,
                        "joined_member_count": 1,
                        "members": [
                            {
                                "user_id": "USR20030471",
                                "nickname": "rusdi",
                                "profile_url": "https://sendbird.com/main/img/profiles/profile_17_512px.png",
                                "is_active": true,
                                "is_online": false,
                                
                                "last_seen_at": 1530232836311,
                                "state": "joined",
                                "metadata": {
                                  "isMerchant": false,
                                  "merchantName": ""
                                }
                            },
                            {
                                "user_id": "USR20030472",
                                "nickname": "azlan",
                                "profile_url": "https://sendbird.com/main/img/profiles/profile_08_512px.png",
                                "is_active": true,
                                "is_online": false,
                                "last_seen_at": 1530237133254,
                                "metadata": {
                                  "isMerchant": True,
                                  "merchantName": "Berkah Store"
                                }
                            }
                        ],
                        "operators": [
                             "user_id": "USR20030471",
                                "nickname": "rusdi",
                                "profile_url": "https://sendbird.com/main/img/profiles/profile_17_512px.png",
                                "is_active": true,
                                "is_online": false,
                                
                                "last_seen_at": 1530232836311,
                                "state": "joined",
                                "metadata": {
                                  "isMerchant": false,
                                  "merchantName": ""
                                }
                        ]
                      
                    "member_count": 3,
                    "joined_member_count": 3,
                    
                    ],
                    "delivery_receipt": {
                        "Jay": 1542762344162,
                        "Jin": 1542394323413,
                        "David": 1542543464371
                    },
                    "read_receipt": {
                        "Jay": 1542762343245,
                        "Jin": 1542394301402,
                        "David": 1542543456343
                    },
                    "last_message": {
                        "message_id": 640903435,
                        "type": "MESG",
                        "custom_type": "",
                        "mention_type": "users",
                        "mentioned_users": [],
                        "created_at": 1542762343245,
                        "updated_at": 0,
                        "is_removed": false,
                        "channel_url": "sendbird_group_channel_25108471_a1bc35f5f0d237207bc1rd343562c878fc2fd426",
                        "user": {
                            "user_id": "Jay",
                            "nickname": "Rooster",
                            "profile_url": "https://sendbird.com/main/img/profiles/profile_13_512px.png"
                            "metadata": {
                                "location": "New York",
                                "marriage": "Y"
                            }
                        },
                        "message": "Can you please make the presentation for me?",
                        "translations": {},
                        "data": "",
                        "file": {}
                    },
                    "unread_message_count": 0,
                    "unread_mention_count": 0,
                    "metadata": {
                        "background_image": "https://sendbird.com/main/img/bg/theme_013.png",
                        "text_size": "large"
                    },
    
                },
                ... # More group channels
            ],
            "next": "ansYQFFRQ1AIEUBXX1RcE2d0FUZSUlkJFVQRHB86AkAgNn8eABABBBNFX11fUlsWYnMS"
        }
        ```       


     
        
 
### Server Architecture
![Server Architecture](https://i.postimg.cc/QtLGfWg9/server-architecture-chat-engine.png)

- Ada rencana untuk menggunakan sendbird untuk chat enginenya : 
	- Pros
		- sdk ready production (ios,andorid & web)
		- integrasi dan impelementasi cukup mudah didukung oleh dokumentasi dan example project
		- tidak perlu menyeiapkan server chat
	- Cons
		- untuk custome terbatas karena harus mengikuti rule langsung dari sendbird
		- banyak limitasi karena menyesuikan dengan package yang dipilih
		
    dengan menggunakan send-bird sebagai chat engine diharapkan dapat men-solve kebutuhan untuk komunikasi langsung antara seller dan buyer, dan juga di harapkan proses development akan lebih cepat karena tidak membuat engine dari awal.

- Any suggestion?

## Data-Migration
- data yang akan di migrasikan hanya data merchant yg akan diambil dari user-service, sedangkan data user  on demand, kemudian untuk data merchant akan di migrasikan melalui hit endpoint sendbird 
    - POST `https://api-{application_id}.sendbird.com/v3/users)`

## Challenge 
- Custome chat engine (UI,flow dll) sesuai dengan kebutuhan bhinneka dengan package yg sudah dipilih