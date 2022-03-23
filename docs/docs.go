// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate_swagger = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "{   \"app_version\": \"v1.1.1\",   \"device\": \"iphone\",   \"device_id\": \"aaaaaaaa\",   \"device_version\": \"12\",   \"dpi\": \"390x844\",   \"ip\": \"127.0.0.1\",   \"lat\": \"21.1111\",   \"lon\": \"32.4444\",   \"os\": 1,   \"os_version\": \"11\",   \"referer\": \"\" }"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/base/captcha": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "防止有人恶意攻击，尝试破解密码",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "生成图片验证码",
                "parameters": [
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpresponse.Response"
                        }
                    }
                }
            }
        },
        "/base/constList": {
            "get": {
                "description": "常量列表",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "所有常量列表",
                "parameters": [
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "enum": [
                            "1",
                            "2",
                            "3",
                            "4"
                        ],
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"登陆成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/base/headerStruct": {
            "get": {
                "description": "日常header里放一诸如验证类的东西，统一公示出来，方便使用",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "header头结构体",
                "parameters": [
                    {
                        "description": "客户端基础信息",
                        "name": "X-HeaderBaseInfo",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.HeaderBaseInfo"
                        }
                    },
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/request.Header"
                        }
                    }
                }
            }
        },
        "/base/login": {
            "post": {
                "description": "用户登陆，验证，生成token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "用户登陆",
                "parameters": [
                    {
                        "description": "用户名, 密码, 验证码",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.Login"
                        }
                    },
                    {
                        "type": "string",
                        "default": "11",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"登陆成功\"}",
                        "schema": {
                            "$ref": "#/definitions/request.Login"
                        }
                    }
                }
            }
        },
        "/base/loginThird": {
            "post": {
                "description": "用户登陆，验证，生成token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "用户登陆三方",
                "parameters": [
                    {
                        "description": "用户名, 密码, 验证码",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.Login"
                        }
                    },
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"登陆成功\"}",
                        "schema": {
                            "$ref": "#/definitions/request.LoginThird"
                        }
                    }
                }
            }
        },
        "/base/parserToken": {
            "post": {
                "description": "应用接到token后，要到后端再统计认证一下，确保准确",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "解析一个TOKEN",
                "parameters": [
                    {
                        "description": "需要验证的token值",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/v1.ParserTokenReq"
                        }
                    },
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/httpresponse.Response"
                        }
                    }
                }
            }
        },
        "/base/projectList": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "每个项目的详细信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "项目列表",
                "parameters": [
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Project"
                        }
                    }
                }
            }
        },
        "/base/register": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "用户注册账号",
                "parameters": [
                    {
                        "type": "string",
                        "default": "1",
                        "description": "来源",
                        "name": "X-Source-Type",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "6",
                        "description": "项目ID",
                        "name": "X-Project-Id",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "default": "imzgoframe",
                        "description": "访问KEY",
                        "name": "X-Access",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "客户端基础信息(json格式,参考request.HeaderBaseInfo)",
                        "name": "X-Base-Info",
                        "in": "header"
                    },
                    {
                        "description": "用户信息",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.Register"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"注册成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/base/sendSMS": {
            "post": {
                "description": "登陆、注册、通知等发送短信",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Base"
                ],
                "summary": "发送验证码",
                "parameters": [
                    {
                        "description": "手机号, 规则ID",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SendSMS"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"发送成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/sys/config": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Config",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Config",
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"登陆成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/sys/quit": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Quit",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "System"
                ],
                "summary": "Quit",
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"登陆成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/changePassword": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "用户修改密码",
                "parameters": [
                    {
                        "description": "用户名, 原密码, 新密码",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.ChangePasswordStruct"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"修改成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/deleteUser": {
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "删除用户",
                "parameters": [
                    {
                        "description": "用户ID",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.GetById"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"删除成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/getUserList": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "分页获取用户列表",
                "parameters": [
                    {
                        "description": "页码, 每页大小",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.PageInfo"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"获取成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/logout": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "用户退出",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "用户退出",
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"退出成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/setUserInfo": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "设置用户信息",
                "parameters": [
                    {
                        "description": "ID, 用户名, 昵称, 头像链接",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"success\":true,\"data\":{},\"msg\":\"设置成功\"}",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "gorm.DeletedAt": {
            "type": "object",
            "properties": {
                "time": {
                    "type": "string"
                },
                "valid": {
                    "description": "Valid is true if Time is not NULL",
                    "type": "boolean"
                }
            }
        },
        "httpresponse.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "msg": {
                    "type": "string"
                }
            }
        },
        "model.Project": {
            "type": "object",
            "properties": {
                "access": {
                    "type": "string"
                },
                "created_at": {
                    "type": "integer"
                },
                "deleted_at": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "desc": {
                    "type": "string"
                },
                "git": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "secret_key": {
                    "type": "string"
                },
                "status": {
                    "type": "integer"
                },
                "type": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "integer"
                }
            }
        },
        "model.User": {
            "type": "object",
            "properties": {
                "birthday": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "integer"
                },
                "deleted_at": {
                    "$ref": "#/definitions/gorm.DeletedAt"
                },
                "email": {
                    "type": "string"
                },
                "headerImg": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "mobile": {
                    "type": "string"
                },
                "nick_name": {
                    "type": "string"
                },
                "project_id": {
                    "type": "integer"
                },
                "recommend": {
                    "type": "string"
                },
                "robot": {
                    "type": "integer"
                },
                "sex": {
                    "type": "integer"
                },
                "status": {
                    "type": "integer"
                },
                "third_id": {
                    "type": "string"
                },
                "type": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "request.ChangePasswordStruct": {
            "type": "object",
            "properties": {
                "newPassword": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "request.GetById": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "number"
                }
            }
        },
        "request.Header": {
            "type": "object",
            "properties": {
                "access": {
                    "description": "使用网关时，不允许随意访问，得有key",
                    "type": "string"
                },
                "auto_ip": {
                    "description": "获取不到请求方IP时，系统自动获取生成",
                    "type": "string"
                },
                "base_info": {
                    "description": "收集客户端的一些基础信息，json",
                    "$ref": "#/definitions/request.HeaderBaseInfo"
                },
                "project_id": {
                    "description": "项目ID，所有的服务/项目/前端/App，均要先向管理员申请一个账号，才能用于日常请求",
                    "type": "integer"
                },
                "request_id": {
                    "description": "每次请求的唯一标识，响应时也会返回，如果请求方没有，后端会默认生成一个",
                    "type": "string"
                },
                "source_type": {
                    "description": "pc h5 ios android vr spider unknow",
                    "type": "integer"
                },
                "token": {
                    "description": "JWT用户登陆令牌(HS256 对称算法，共享一个密钥)",
                    "type": "string"
                },
                "trace_id": {
                    "description": "追踪ID，主要用于链路追踪，如果请求方没有，后端会默认生成一个，跟request略像，但给后端使用",
                    "type": "string"
                }
            }
        },
        "request.HeaderBaseInfo": {
            "type": "object",
            "properties": {
                "app_version": {
                    "description": "app/前端/服务/项目 版本号",
                    "type": "string"
                },
                "device": {
                    "description": "ipad iphone huawei mi chrome firefox ie",
                    "type": "string"
                },
                "device_id": {
                    "description": "设备ID",
                    "type": "string"
                },
                "device_version": {
                    "description": "mi8 hongmi7 ios8 ios9 ie8 ie9",
                    "type": "string"
                },
                "dpi": {
                    "description": "分辨率",
                    "type": "string"
                },
                "ip": {
                    "description": "请求方的IP",
                    "type": "string"
                },
                "lat": {
                    "description": "纬度",
                    "type": "string"
                },
                "lon": {
                    "description": "经度",
                    "type": "string"
                },
                "os": {
                    "description": "win mac android ios",
                    "type": "integer"
                },
                "os_version": {
                    "description": "win7 win9 mac10 android9",
                    "type": "string"
                },
                "referer": {
                    "description": "页面来源",
                    "type": "string"
                }
            }
        },
        "request.Login": {
            "type": "object",
            "properties": {
                "captcha": {
                    "description": "验证码",
                    "type": "string"
                },
                "captchaId": {
                    "description": "验证码-ID",
                    "type": "string"
                },
                "password": {
                    "description": "密码",
                    "type": "string"
                },
                "username": {
                    "description": "用户名：普通字符串、手机号、邮箱",
                    "type": "string"
                }
            }
        },
        "request.LoginThird": {
            "type": "object",
            "properties": {
                "Code": {
                    "type": "string"
                },
                "captcha": {
                    "type": "string"
                },
                "captchaId": {
                    "type": "string"
                },
                "platform": {
                    "type": "string"
                }
            }
        },
        "request.PageInfo": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer"
                },
                "pageSize": {
                    "type": "integer"
                }
            }
        },
        "request.Register": {
            "type": "object",
            "properties": {
                "birthday": {
                    "description": "生日",
                    "type": "integer"
                },
                "channel": {
                    "description": "来源渠道",
                    "type": "integer"
                },
                "ext_diy": {
                    "description": "自定义用户属性，暂未实现",
                    "type": "string"
                },
                "guest": {
                    "description": "类型,1普通2游客",
                    "type": "integer"
                },
                "headerImg": {
                    "description": "头像地址",
                    "type": "string"
                },
                "nickName": {
                    "description": "昵称",
                    "type": "string"
                },
                "passWord": {
                    "description": "登陆密码 转md5存储",
                    "type": "string"
                },
                "project_id": {
                    "description": "项目Id",
                    "type": "integer"
                },
                "recommend": {
                    "description": "推荐人",
                    "type": "string"
                },
                "sex": {
                    "description": "性别",
                    "type": "integer"
                },
                "third_id": {
                    "description": "三方平台ID",
                    "type": "string"
                },
                "third_type": {
                    "description": "三方平台类型",
                    "type": "integer"
                },
                "userName": {
                    "description": "用户名",
                    "type": "string"
                }
            }
        },
        "request.SendSMS": {
            "type": "object",
            "properties": {
                "app_id": {
                    "type": "string"
                },
                "code": {
                    "type": "string"
                },
                "mobile": {
                    "type": "string"
                },
                "rule_id": {
                    "type": "integer"
                }
            }
        },
        "v1.ParserTokenReq": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "X-Token",
            "in": "header"
        }
    }
}`

// SwaggerInfo_swagger holds exported Swagger Info so clients can modify it
var SwaggerInfo_swagger = &swag.Spec{
	Version:          "0.1 测试版",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "z golang 框架",
	Description:      "restful api 工具，模拟客户端请求，方便调试/测试<br/>注：这只是一个工具，不是万能的，像：动态枚举类型、公共请求header、动态常量等<br/>详细的请去 <a href=\"http://127.0.0.1:6060\" target=\"_black\">godoc</a> 里去查看",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate_swagger,
}

func init() {
	swag.Register(SwaggerInfo_swagger.InstanceName(), SwaggerInfo_swagger)
}
