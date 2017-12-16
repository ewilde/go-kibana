# Search for index patterns
/kibana/elasticsearch/logzioCustomerKibanaIndex/index-pattern/_search?stored_fields=

**Request body**
`{"query":{"match_all":{}},"size":100}`

**Response body**
```json
{
  "took": 0,
  "timed_out": false,
  "_shards": {
    "total": 1,
    "successful": 1,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 1,
    "hits": [
      {
        "_index": "logzioCustomerKibanaIndex",
        "_type": "index-pattern",
        "_id": "[logzioCustomerIndex]YYMMDD",
        "_score": 1
      }
    ]
  }
}
```

# Create search
/kibana/elasticsearch/logzioCustomerKibanaIndex/search/9c2f2320-e252-11e7-96f8-397bd34fab6c

**Request body**

```json
{
  "title": "test s",
  "description": "",
  "hits": 0,
  "columns": [
    "message"
  ],
  "sort": [
    "@timestamp",
    "desc"
  ],
  "version": 1,
  "kibanaSavedObjectMeta": {
    "searchSourceJSON": "{\"index\":\"[logzioCustomerIndex]YYMMDD\",\"highlightAll\":true,\"version\":true,\"query\":{\"query_string\":{\"query\":\"message:\\\"GET\\\"\",\"analyze_wildcard\":true}},\"filter\":[]}"
  },
  "_createdBy": {
    "userId": 19429,
    "fullName": "Edward Wilde",
    "username": "edward.wilde@form3.tech"
  },
  "_createdAt": 1513423009383,
  "_updatedBy": {
    "userId": 19429,
    "fullName": "Edward Wilde",
    "username": "edward.wilde@form3.tech"
  },
  "_updatedAt": 1513423009383
}
```

**Response**
```json
{
  "_index": "logzioCustomerKibanaIndex",
  "_type": "search",
  "_id": "9c2f2320-e252-11e7-96f8-397bd34fab6c",
  "_version": 1,
  "result": "created",
  "_shards": {
    "total": 2,
    "successful": 2,
    "failed": 0
  },
  "created": true
}
```
