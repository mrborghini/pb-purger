# pb-purger

This application purges old entries inside of [PocketBase](https://pocketbase.io/) older than a certain time given in the following JSON object:

```json
[
  {
    "pbUrl": "http://localhost:8090",
    "accountCollection": "users",
    "pbUsername": "",
    "pbPassword": "",
    "collections": [
      {
        "name": "your_collection_name",
        "deletionTimeSeconds": 2678400
      },
      {
        "name": "your_collection_name2",
        "deletionTimeSeconds": 600
      }
    ]
  },
  {
    "pbUrl": "http://localhost:8080",
    "accountCollection": "clients",
    "pbUsername": "myaccount",
    "pbPassword": "mine",
    "collections": [
      {
        "name": "your_collection_name",
        "deletionTimeSeconds": 2678400
      },
      {
        "name": "your_collection_name2",
        "deletionTimeSeconds": 600
      }
    ]
  }
]
```

You can customize it easily in this json format.

## Setup

```bash
cp config_example.json config.json
```

This step is always required.

### Docker

```bash
docker compose up -d
```

### Native

```bash
./start.sh
```