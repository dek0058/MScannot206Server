# MScannot206Server&nbsp;![Go](https://img.shields.io/badge/Language-Go-00ADD8?style=flat&logo=go&logoColor=white) ![Go Version](https://img.shields.io/badge/Version-1.25.4-00ADD8?style=flat&logo=go&logoColor=white)

이 프로젝트는 [MScannot206](https://github.com/dek0058/MScannot206) 클라이언트를 보조하기 위한 콘솔 서버 입니다.

[메이플스토리 월드 크리에이터 이용약관](https://github.com/dek0058/MScannot206)을 준수하며, 해당 프로젝트는 비공식 프로젝트임을 알립니다.

## 📋 요구사항

 - [Go](https://go.dev/doc/install)
 - [MongoDB](https://www.mongodb.com/try/download/community)

## 📚 API Documentation

상세한 API 명세는 아래 문서들을 참고해주세요.

- [🔐 로그인/인증 API (Login)](document/login.md)
- [👤 유저/캐릭터 API (User)](document/user.md)

## 🏗️ 메인 아키텍처

```mermaid
graph TD
    classDef user fill:#ffffff,stroke:#333,stroke-width:2px,color:#000000,font-weight:bold;
    classDef client fill:#E3F2FD,stroke:#1565C0,stroke-width:2px,color:#000000,font-weight:bold;
    classDef server fill:#E8F5E9,stroke:#2E7D32,stroke-width:2px,color:#000000,font-weight:bold;
    classDef db fill:#FFF3E0,stroke:#EF6C00,stroke-width:2px,color:#000000,font-weight:bold;

    User((User)):::user
    Client[Client]:::client

    subgraph Server_Area [Server Side]
        direction TB
        Services[Services]:::server
        Repositories[Repositories]:::server
    end

    subgraph Data_Area [Persistence Layer]
        DB[("MongoDB")]:::db
    end

    User--->|1.Connect|Client
    Client -->|2.API Request| Services
    Services -->|3.Input Data| Repositories
    Repositories -->|4.Query| DB
    DB -.->|5.Result| Repositories
    Repositories -.->|6.Output Data| Services
    Services -.->|7.API Response| Client
```
