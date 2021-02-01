#!/usr/bin/env python3

import flask
import requests
import os

METRICS = {
    "abs_humid": "Absolute humidity (g/m^3)",
    "co2": "CO2 level (ppm)",
    "dew_point": "Dew point (C)",
    "humid": "Relative humidity (%)",
    "pm25": "Particulate matter (ug/m^3)",
    "score": "Awair score (%)",
    "temp": "Temperature (C)",
    "voc": "Volatile organic compounds (ppb)",
    "co2_est": "?",
    "pm10_est": "?",
    "voc_baseline": "?",
    "voc_ethanol_raw": "?",
    "voc_h2_raw": "?",
}

SENSORS = {}
for sensor in os.environ["SENSORS"].split(","):
    name, ip = sensor.split("=")
    SENSORS[name] = ip

app = flask.Flask(__name__)


def generate_metrics():
    global METRICS, SENSORS

    sensor_data = {}
    for name, ip in SENSORS.items():
        r = requests.get(f"http://{ip}/air-data/latest")
        r.raise_for_status()
        sensor_data[name] = r.json()

    prom_metrics = []

    for metric, description in METRICS.items():
        metric_name = f"awair_{metric}"
        prom_metrics.append(f"# HELP {metric_name} {description}")
        prom_metrics.append(f"# TYPE {metric_name} gauge")
        for sensor, values in sensor_data.items():
            value = values[metric]
            prom_metrics.append(f'{metric_name}{{sensor="{sensor}"}} {value}')

    return "\n".join(prom_metrics) + "\n"


@app.route("/metrics")
def serve_metrics():
    metrics = generate_metrics()
    if metrics is None:
        flask.abort(500)

    response = flask.make_response(metrics, 200)
    response.mimetype = "text/plain"
    return response


app.run(host="0.0.0.0", port=8888)
