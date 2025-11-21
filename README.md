

# API 测试说明

### **结构图**

elderly-care-backend/
├── common/
│   ├── constants/
│   │   ├── account_constant.go
│   │   ├── common_constant.go
│   │   ├── kafka_constant.go
│   │   ├── redis_constant.go
│   │   └── response_constant.go
│   ├── custom/
│   │   ├── global_lock.go
│   │   ├── kafka_operator.go
│   │   └── set.go
│   ├── factories/
│   │   └── oss_factory.go
│   └── server_error/
│       └── errors.go
├── config/
│   ├── config.go
│   ├── config.yaml
│   ├── db_config.go
│   ├── kafka_config.go
│   ├── logger_config.go
│   ├── map_config.go
│   ├── redis_config.go
│   ├── redsync_config.go
│   └── snowflake_config.go
├── controllers/
│   ├── account_controller.go
│   ├── chat_controller.go
│   ├── evaluation_controller.go
│   ├── file_controller.go
│   ├── location_controller.go
│   ├── sos_controller.go
│   └── task_controller.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── dto/
│   ├── account_dto/
│   │   ├── account_update_dto.go
│   │   ├── claims.go
│   │   └── register_dto.go
│   ├── chat_dto.go
│   ├── sos_dto.go
│   └── task_dto.go
├── global/
│   └── global.go
├── middlewares/
│   └── verify_middleware.go
├── models/
│   ├── account_evaluation.go
│   ├── account.go
│   ├── base_model.go
│   ├── contact_list.go
│   ├── location.go
│   ├── message.go
│   ├── navigation.go
│   ├── sos.go
│   └── task.go
├── routes/
│   ├── account_route.go
│   ├── chat_route.go
│   ├── evaluation_route.go
│   ├── file_route.go
│   ├── location_route.go
│   ├── router.go
│   ├── sos_route.go
│   └── task_route.go
├── services/
│   ├── address_service.go
│   ├── amap_service.go
│   ├── realtime_location_service.go
│   ├── task_matching_service.go
│   └── websocket_service.go
├── utils/
│   ├── account_util.go
│   ├── common_util.go
│   ├── file_util.go
│   ├── location_utils.go
│   ├── music_util.go
│   └── parse_str_util.go
├── vo/
│   ├── music_vo/
│   │   ├── music_list_vo.go
│   │   ├── playlist_detail_vo.go
│   │   ├── playlist_vo.go
│   │   └── source_lyrics_vo.go
│   ├── account_vo.go
│   ├── account_evaluation_vo.go
│   ├── contact_list_vo.go
│   ├── response_vo.go
│   ├── sos_vo.go
│   └── task_vo.go
├── docker-compose.yaml
├── go.mod
├── go.sum
├── log/
├── main.go
└── README.md



### API接口速查表

 账户模块 (/account)

POST /account/register - 用户注册，需要：昵称、手机、密码、头像、性别、年龄

POST /account/login - 用户登录，需要：手机号、密码、登录类型

PUT /account - 更新账户信息，需要认证

PUT /account/changePassword - 修改密码，需要：原密码、新密码

GET /account/checkPhone - 检查手机号是否存在，需要：手机号

GET /account - 获取当前用户信息，需要认证

 

定位模块 (/location)

POST /location/update - 更新实时位置，需要：经纬度、地址

GET /location/user/:userId - 获取指定用户位置，需要：用户ID

GET /location/nearby - 获取附近用户，需要：经纬度、半径、角色

GET /location/history - 获取位置历史，需要认证

GET /location/reverse-geocode - 坐标转地址，需要：经纬度

 

导航模块 (/location/navigation)

GET /location/navigation - 基础路径规划，需要：起点坐标、终点坐标

GET /location/navigation/to-target - 到目标导航，需要：目标经纬度

GET /location/navigation/user - 用户间导航，需要：目标用户ID

GET /location/navigation/location - 历史位置导航，需要：位置记录ID

 

求助模块 (/sos)

POST /sos/emergency - 触发紧急求助，需要：经纬度、求助信息

POST /sos/:sosId/accept - 接受求助，需要：SOS ID

PUT /sos/:sosId/resolve - 解决求助，需要：SOS ID

GET /sos/current - 获取当前求助，需要认证

 

任务模块 (/tasks)

POST /tasks - 创建任务，需要：标题、描述、经纬度、报酬

GET /tasks/nearby - 获取附近任务，需要：经纬度、半径

POST /tasks/:taskId/accept - 接受任务，需要：任务ID

 

实时通信

GET /ws - WebSocket连接，用于实时位置推送

 

文档

GET /swagger/*any - API文档页面
