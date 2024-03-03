# Reproducer for Dapr JSON Serialization issue

Issue: [https://github.com/dapr/dapr/issues/7580](https://github.com/dapr/dapr/issues/7580)


To run locally using the Dapr CLI:

```
dapr run --app-id reproducer --dapr-grpc-port 60001 go run main.go
```

I am using the Redis Stack for Query API, so you need to run: 

```
docker run redis/redis-stack 6379:6379
```


My StateStore YAML looks like this: 

```
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: statestore
spec:
  type: state.redis
  version: v1
  metadata:
  - name: keyPrefix
    value: name 
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""
  - name: actorStateStore
    value: "true"
  - name: queryIndexes
    value: |
      [
        {
          "name": "voteIndex",
          "indexes": [
           {
            "key": "type",
            "type": "TEXT" 
           }
          ]
        }
      ]     
```

After running the code you should see the 3 keys in Redis where the Save and the Save Bulk values are correctly serialized into JSON. While ExecuteStateTransaction is encoded using Base64 encoding, breaking the Query API.