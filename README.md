# HydroBytes-BaseStation
The Base Station is a part of a collection of systems called
"[HydroBytes](https://github.com/deezone/HydroBytes)" that uses micro
controllers to manage and monitor plant health in an automated garden.

## Introduction

The "garden" is simply a backyard patio in Brooklyn, New York. Typically
there are only two seasons - cold and hot in New York City. By
automating an urban garden ideally the space will thrive with minimum
supervision. The amount of effort to automate is besides the point, everyone needs their vices.

- **[Water Station](https://github.com/deezone/HydroBytes-WaterStation)**
- **Base Station**
- **[Plant Station](https://github.com/deezone/HydroBytes-PlantStation)**

![brooklyn-20201115 garden layout](https://raw.githubusercontent.com/deezone/HydroBytes/master/resources/gardenBrooklynDiagram-20201115.jpg)

![Garden](https://github.com/deezone/HydroBytes-WaterManagement/blob/master/resources/garden-01.png)

### YouTube Channel

[![YouTube Channel](https://github.com/deezone/HydroBytes-WaterStation/blob/master/resources/youTube-TN.png?raw=true)](https://www.youtube.com/channel/UC00A_lEJD2Hcy9bw6UuoUBA "All of the HydroBytes videos")

### Notes

Development of a Go based API is based on instruction in the amazing
courses at **[Ardan Labs](https://education.ardanlabs.com/collections?category=courses)**.

#### Starting Web Server
```
> go run ./cmd/api

2021/01/09 18:55:47 main : Started
2021/01/09 18:55:47 main : Config :
--web-address=localhost:8000
--web-read-timeout=5s
--web-write-timeout=5s
--web-shutdown-timeout=5s
--db-user=postgres
--db-host=localhost
--db-name=postgres
--db-disable-tls=true
2021/01/09 18:55:47 main : API listening on localhost:8000


^C
2021/01/09 18:57:37 main : Start shutdown
2021/01/09 18:57:37 main : Completed
```

- supported requests to `localhost:8000`:
  - `GET  /v1/station-types`
  - `GET  /v1/station-type/{id}`
  - `POST /v1/station-type`
  - `DELETE /v1/station-type/{id}`
  - `GET  /v1/station-type/{station-type-id}/stations`
  - `POST /v1/station-type/{station-type-id}/station`
  - `DELETE /v1/station/{id}`

#### Admin tools

```
> go run ./cmd/admin -h migrate
Usage: admin [options] [arguments]

OPTIONS
  --db-user/$STATIONS_DB_USER                <string>  (default: postgres)
  --db-password/$STATIONS_DB_PASSWORD        <string>  (noprint,default: postgres)
  --db-host/$STATIONS_DB_HOST                <string>  (default: localhost)
  --db-name/$STATIONS_DB_NAME                <string>  (default: postgres)
  --db-disable-tls/$STATIONS_DB_DISABLE_TLS  <bool>    (default: false)
  --help/-h
  display this help message
```

- `migrate` to update database with schema defined in code.
```
> go run ./cmd/admin migrate
Migrations complete
```

- `seed` populate the database tables with seed data for testing and development.
```
> go run ./cmd/admin seed
Seeding complete
```

#### Tests

- **Unit Tests**

**station_type** and **station**
```
> go test ./internal/station_type
ok  	github.com/deezone/HydroBytes-BaseStation/internal/station_type	13.515s
```

NOTE: test coverage reports:
```
alais gotwc='go test -coverprofile=coverage.out && go tool cover -html=coverage.out && rm coverage.out'
```

- **Functional tests**
```
# bust cache
> go clean -testcache

> go test ./cmd/api/tests/station_tests
ok  	github.com/deezone/HydroBytes-BaseStation/cmd/api/tests/station_tests	3.248s

> go test ./cmd/api/tests/station_type_tests
ok  	github.com/deezone/HydroBytes-BaseStation/cmd/api/tests/station_type_tests	2.875s
```
