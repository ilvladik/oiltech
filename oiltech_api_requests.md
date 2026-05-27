# Примеры запросов к сервису OilTech

## Healthcheck

```bash
curl http://localhost:8080/healthz
```

## Создать датасет `meter_readings_raw`

```bash
curl -X POST http://localhost:8080/datasets \
  -H "Content-Type: application/json" \
  -d '{
    "code": "meter_readings_raw",
    "name": "Meter Readings Raw",
    "mqtt_topic": "oiltech/datasets/meter_readings_raw/rows",
    "columns": [
      {"code": "event_id", "name": "Event ID", "data_type": "string", "nullable": false},
      {"code": "source_system", "name": "Source System", "data_type": "string", "nullable": false},
      {"code": "source_tag", "name": "Source Tag", "data_type": "string", "nullable": false},
      {"code": "meter_id", "name": "Meter ID", "data_type": "string", "nullable": false},
      {"code": "facility_id", "name": "Facility ID", "data_type": "string", "nullable": false},
      {"code": "equipment_id", "name": "Equipment ID", "data_type": "string", "nullable": true},
      {"code": "equipment_type", "name": "Equipment Type", "data_type": "string", "nullable": true},
      {"code": "resource_type", "name": "Resource Type", "data_type": "string", "nullable": false},
      {"code": "measurement_kind", "name": "Measurement Kind", "data_type": "string", "nullable": false},
      {"code": "moment", "name": "Moment", "data_type": "datetime", "nullable": true},
      {"code": "period_start", "name": "Period Start", "data_type": "datetime", "nullable": true},
      {"code": "period_end", "name": "Period End", "data_type": "datetime", "nullable": true},
      {"code": "value", "name": "Value", "data_type": "decimal", "nullable": false},
      {"code": "unit", "name": "Unit", "data_type": "string", "nullable": false},
      {"code": "quality", "name": "Quality", "data_type": "string", "nullable": false},
      {"code": "status", "name": "Status", "data_type": "string", "nullable": false}
    ]
  }'
```

## Создать датасет `process_measurements_raw`

```bash
curl -X POST http://localhost:8080/datasets \
  -H "Content-Type: application/json" \
  -d '{
    "code": "process_measurements_raw",
    "name": "Process Measurements Raw",
    "mqtt_topic": "oiltech/datasets/process_measurements_raw/rows",
    "columns": [
      {"code": "event_id", "name": "Event ID", "data_type": "string", "nullable": false},
      {"code": "source_system", "name": "Source System", "data_type": "string", "nullable": false},
      {"code": "source_tag", "name": "Source Tag", "data_type": "string", "nullable": false},
      {"code": "sensor_id", "name": "Sensor ID", "data_type": "string", "nullable": false},
      {"code": "facility_id", "name": "Facility ID", "data_type": "string", "nullable": false},
      {"code": "well_id", "name": "Well ID", "data_type": "string", "nullable": true},
      {"code": "equipment_id", "name": "Equipment ID", "data_type": "string", "nullable": true},
      {"code": "equipment_type", "name": "Equipment Type", "data_type": "string", "nullable": true},
      {"code": "parameter_id", "name": "Parameter ID", "data_type": "string", "nullable": false},
      {"code": "moment", "name": "Moment", "data_type": "datetime", "nullable": false},
      {"code": "value", "name": "Value", "data_type": "decimal", "nullable": false},
      {"code": "unit", "name": "Unit", "data_type": "string", "nullable": false},
      {"code": "quality", "name": "Quality", "data_type": "string", "nullable": false},
      {"code": "status", "name": "Status", "data_type": "string", "nullable": false}
    ]
  }'
```

## Создать датасет `equipment_states_raw`

```bash
curl -X POST http://localhost:8080/datasets \
  -H "Content-Type: application/json" \
  -d '{
    "code": "equipment_states_raw",
    "name": "Equipment States Raw",
    "mqtt_topic": "oiltech/datasets/equipment_states_raw/rows",
    "columns": [
      {"code": "event_id", "name": "Event ID", "data_type": "string", "nullable": false},
      {"code": "source_system", "name": "Source System", "data_type": "string", "nullable": false},
      {"code": "facility_id", "name": "Facility ID", "data_type": "string", "nullable": false},
      {"code": "equipment_id", "name": "Equipment ID", "data_type": "string", "nullable": false},
      {"code": "equipment_type", "name": "Equipment Type", "data_type": "string", "nullable": false},
      {"code": "state", "name": "State", "data_type": "string", "nullable": false},
      {"code": "mode", "name": "Mode", "data_type": "string", "nullable": false},
      {"code": "scheme_id", "name": "Scheme ID", "data_type": "string", "nullable": true},
      {"code": "moment", "name": "Moment", "data_type": "datetime", "nullable": false},
      {"code": "quality", "name": "Quality", "data_type": "string", "nullable": false}
    ]
  }'
```

## Создать датасет `equipment_work_periods`

```bash
curl -X POST http://localhost:8080/datasets \
  -H "Content-Type: application/json" \
  -d '{
    "code": "equipment_work_periods",
    "name": "Equipment Work Periods",
    "columns": [
      {"code": "facility_id", "name": "Facility ID", "data_type": "string", "nullable": false},
      {"code": "equipment_id", "name": "Equipment ID", "data_type": "string", "nullable": false},
      {"code": "equipment_type", "name": "Equipment Type", "data_type": "string", "nullable": false},
      {"code": "period_start", "name": "Period Start", "data_type": "datetime", "nullable": false},
      {"code": "period_end", "name": "Period End", "data_type": "datetime", "nullable": false},
      {"code": "state", "name": "State", "data_type": "string", "nullable": false},
      {"code": "duration_seconds", "name": "Duration Seconds", "data_type": "int", "nullable": false},
      {"code": "quality", "name": "Quality", "data_type": "string", "nullable": false}
    ]
  }'
```

## Получить все датасеты

```bash
curl http://localhost:8080/datasets
```

## Получить конкретный датасет

```bash
curl http://localhost:8080/datasets/meter_readings_raw
```

```bash
curl http://localhost:8080/datasets/process_measurements_raw
```

```bash
curl http://localhost:8080/datasets/equipment_states_raw
```

```bash
curl http://localhost:8080/datasets/equipment_work_periods
```

## Получить строки датасета

```bash
curl http://localhost:8080/datasets/meter_readings_raw/rows
```

```bash
curl http://localhost:8080/datasets/process_measurements_raw/rows
```

```bash
curl http://localhost:8080/datasets/equipment_states_raw/rows
```

```bash
curl http://localhost:8080/datasets/equipment_work_periods/rows
```

## Получить строки с `limit`

```bash
curl "http://localhost:8080/datasets/meter_readings_raw/rows?limit=10"
```

```bash
curl "http://localhost:8080/datasets/process_measurements_raw/rows?limit=10"
```

```bash
curl "http://localhost:8080/datasets/equipment_states_raw/rows?limit=10"
```

```bash
curl "http://localhost:8080/datasets/equipment_work_periods/rows?limit=10"
```

## Создать алгоритм

```bash
curl -X POST http://localhost:8080/algorithms \
  -H "Content-Type: application/json" \
  -d '{
    "code": "equipment_work_periods_algo",
    "name": "Equipment States to Work Periods",
    "description": "Builds work period intervals from equipment state events.",
    "run_url": "http://equipment-work-periods:9000/run"
  }'
```

## Получить все алгоритмы

```bash
curl http://localhost:8080/algorithms
```

## Получить конкретный алгоритм

```bash
curl http://localhost:8080/algorithms/equipment_work_periods_algo
```

## Создать трансформацию

```bash
curl -X POST http://localhost:8080/transformations \
  -H "Content-Type: application/json" \
  -d '{
    "code": "states_to_work_periods",
    "name": "Equipment States to Work Periods",
    "source_dataset_codes": ["equipment_states_raw"],
    "target_dataset_code": "equipment_work_periods",
    "algorithm_code": "equipment_work_periods_algo",
    "period_seconds": 60,
    "enabled": true
  }'
```

## Получить все трансформации

```bash
curl http://localhost:8080/transformations
```

## Получить конкретную трансформацию

```bash
curl http://localhost:8080/transformations/states_to_work_periods
```

## Получить дерево трансформаций

```bash
curl http://localhost:8080/transformations/tree
```

## Запустить трансформацию вручную

```bash
curl -X POST http://localhost:8080/transformations/states_to_work_periods/run
```

## Получить историю запусков трансформации

```bash
curl http://localhost:8080/transformations/states_to_work_periods/runs
```

## Прямой вызов алгоритма

```bash
curl -X POST http://localhost:9000/run \
  -H "Content-Type: application/json" \
  -d '{
    "run_id": "manual-test-run",
    "transformation_code": "states_to_work_periods"
  }'
```
