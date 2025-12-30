# WeKnora 지식 그래프

## 빠른 시작

- `.env`에 환경 변수를 설정하세요.
    - Neo4j 활성화: `NEO4J_ENABLE=true`
    - Neo4j URI: `NEO4J_URI=bolt://neo4j:7687`
    - Neo4j 사용자명: `NEO4J_USERNAME=neo4j`
    - Neo4j 비밀번호: `NEO4J_PASSWORD=password`

- Neo4j 실행
```bash
docker-compose --profile neo4j up -d
```

- 지식베이스 설정 페이지에서 엔티티/관계 추출을 활성화하고 안내에 따라 설정하세요.

## 그래프 생성

문서를 업로드하면 시스템이 자동으로 엔티티와 관계를 추출해 지식 그래프를 생성합니다.

![지식 그래프 예시](./images/graph3.png)

## 그래프 조회

`http://localhost:7474`에 접속해 `match (n) return (n)`을 실행하면 생성된 지식 그래프를 확인할 수 있습니다.  
대화 시 시스템이 그래프를 자동으로 조회해 관련 지식을 불러옵니다.
