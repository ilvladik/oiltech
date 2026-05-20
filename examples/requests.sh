set -euo pipefail

CORE_URL="${http://localhost:8080}"

curl -sS -X POST "$CORE_URL/datasets" \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "meter_readings_raw",
    "name": "Meter Readings Raw",
    "mqtt_topic": "oiltech/datasets/meter_readings_raw",
    "columns": [
      {"code": "source_system",    "name": "Source System",    "data_type": "string",   "nullable": false},
      {"code": "source_tag",       "name": "Source Tag",       "data_type": "string",   "nullable": false},
      {"code": "meter_id",         "name": "Meter ID",         "data_type": "string",   "nullable": false},
      {"code": "facility_id",      "name": "Facility ID",      "data_type": "string",   "nullable": false},
      {"code": "equipment_id",     "name": "Equipment ID",     "data_type": "string",   "nullable": true},
      {"code": "equipment_type",   "name": "Equipment Type",   "data_type": "string",   "nullable": true},
      {"code": "resource_type",    "name": "Resource Type",    "data_type": "string",   "nullable": false},
      {"code": "measurement_kind", "name": "Measurement Kind", "data_type": "string",   "nullable": false},
      {"code": "moment",           "name": "Moment",           "data_type": "datetime", "nullable": true},
      {"code": "period_start",     "name": "Period Start",     "data_type": "datetime", "nullable": true},
      {"code": "period_end",       "name": "Period End",       "data_type": "datetime", "nullable": true},
      {"code": "value",            "name": "Value",            "data_type": "decimal",  "nullable": false},
      {"code": "unit",             "name": "Unit",             "data_type": "string",   "nullable": false},
      {"code": "quality",          "name": "Quality",          "data_type": "string",   "nullable": false},
      {"code": "status",           "name": "Status",           "data_type": "string",   "nullable": false}
    ]
  }'

echo

curl -sS -X POST "$CORE_URL/datasets" \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "process_measurements_raw",
    "name": "Process Measurements Raw",
    "mqtt_topic": "oiltech/datasets/process_measurements_raw",
    "columns": [
      {"code": "source_system",  "name": "Source System",  "data_type": "string",   "nullable": false},
      {"code": "source_tag",     "name": "Source Tag",     "data_type": "string",   "nullable": false},
      {"code": "sensor_id",      "name": "Sensor ID",      "data_type": "string",   "nullable": false},
      {"code": "facility_id",    "name": "Facility ID",    "data_type": "string",   "nullable": false},
      {"code": "well_id",        "name": "Well ID",        "data_type": "string",   "nullable": true},
      {"code": "equipment_id",   "name": "Equipment ID",   "data_type": "string",   "nullable": true},
      {"code": "equipment_type", "name": "Equipment Type", "data_type": "string",   "nullable": true},
      {"code": "parameter_id",   "name": "Parameter ID",   "data_type": "string",   "nullable": false},
      {"code": "moment",         "name": "Moment",         "data_type": "datetime", "nullable": false},
      {"code": "value",          "name": "Value",          "data_type": "decimal",  "nullable": false},
      {"code": "unit",           "name": "Unit",           "data_type": "string",   "nullable": false},
      {"code": "quality",        "name": "Quality",        "data_type": "string",   "nullable": false},
      {"code": "status",         "name": "Status",         "data_type": "string",   "nullable": false}
    ]
  }'

echo

curl -sS -X POST "$CORE_URL/datasets" \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "equipment_states_raw",
    "name": "Equipment States Raw",
    "mqtt_topic": "oiltech/datasets/equipment_states_raw",
    "columns": [
      {"code": "source_system",  "name": "Source System",  "data_type": "string",   "nullable": false},
      {"code": "facility_id",    "name": "Facility ID",    "data_type": "string",   "nullable": false},
      {"code": "equipment_id",   "name": "Equipment ID",   "data_type": "string",   "nullable": false},
      {"code": "equipment_type", "name": "Equipment Type", "data_type": "string",   "nullable": false},
      {"code": "state",          "name": "State",          "data_type": "string",   "nullable": false},
      {"code": "mode",           "name": "Mode",           "data_type": "string",   "nullable": false},
      {"code": "scheme_id",      "name": "Scheme ID",      "data_type": "string",   "nullable": true},
      {"code": "moment",         "name": "Moment",         "data_type": "datetime", "nullable": false},
      {"code": "quality",        "name": "Quality",        "data_type": "string",   "nullable": false}
    ]
  }'

echo

curl -sS -X POST "$CORE_URL/datasets" \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "equipment_work_periods",
    "name": "Equipment Work Periods",
    "columns": [
      {"code": "facility_id",      "name": "Facility ID",       "data_type": "string",   "nullable": false},
      {"code": "equipment_id",     "name": "Equipment ID",      "data_type": "string",   "nullable": false},
      {"code": "equipment_type",   "name": "Equipment Type",    "data_type": "string",   "nullable": false},
      {"code": "period_start",     "name": "Period Start",      "data_type": "datetime", "nullable": false},
      {"code": "period_end",       "name": "Period End",        "data_type": "datetime", "nullable": false},
      {"code": "state",            "name": "State",             "data_type": "string",   "nullable": false},
      {"code": "duration_seconds", "name": "Duration (seconds)","data_type": "int",      "nullable": false},
      {"code": "quality",          "name": "Quality",           "data_type": "string",   "nullable": false}
    ]
  }'

echo

curl -sS -X POST "$CORE_URL/algorithms" \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "equipment_work_periods_algo",
    "name": "Equipment States to Work Periods",
    "description": "Builds work period intervals from consecutive equipment state events.",
    "run_url": "http://equipment-work-periods:9000/run"
  }'

echo

curl -sS -X POST "$CORE_URL/transformations" \
  -H 'Content-Type: application/json' \
  -d '{
    "code": "states_to_work_periods",
    "name": "Equipment States -> Work Periods",
    "source_dataset_codes": ["equipment_states_raw"],
    "target_dataset_code": "equipment_work_periods",
    "algorithm_code": "equipment_work_periods_algo",
    "period_seconds": 30,
    "enabled": true
  }'

echo
