## Mssql Load Generator

- This repository contains the code to generate the load on azure mssql database. This is done with contineously running write operations.


## Use this manifest

- Use this manifest to start a load on mssql server. Also provide the database username and password.

```yaml

apiVersion: v1
kind: Pod
metadata:
  name: sql-load-generator
  labels:
    app: sql-load-generator
spec:
  terminationGracePeriodSeconds: 60
  containers:
  - name: sql-load-generator
    image: uditgaurav/mssql:load
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
    env:
      - name: SERVER
        value: "test-chaos-server.database.windows.net"
      - name: DB_NAME
        value: "test123"
      - name: USERNAME
        value: ""
      - name: PASSWORD
        value: ""
      - name: PORT
        value: "1433"
      - name: TABLE_NAME
        value: "load"
```

## Logs of load-generator pod

```bash
udit@ubuntu ~/g/s/g/a/sql> kubectl logs -f sql-load-generator
time="2021-12-06T22:38:44Z" level=info msg="The mssql information is as follows" Database=test123 Port=1433 Server=test-chaos-server.database.windows.net Table Name=load
time="2021-12-06T22:38:44Z" level=info msg="[CONNECTION]: Trying to establish connection with the sql server"
time="2021-12-06T22:38:45Z" level=info msg="[CONNECTION]: Connection established successfully..."
time="2021-12-06T22:38:45Z" level=info msg="[Info]: Table created successfully..."
time="2021-12-06T22:38:45Z" level=info msg="[Load]: Starting load generator ..."
time="2021-12-06T22:48:27Z" level=info msg="Load has been successfully removed..."
```

