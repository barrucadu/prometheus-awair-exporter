prometheus-awair-exporter
=========================

A simple [Awair][] exporter for [Prometheus][].

For example:

```
# HELP awair_abs_humid Absolute humidity (g/m^3)
# TYPE awair_abs_humid gauge
awair_abs_humid{sensor="living-room"} 8.86
# HELP awair_co2 CO2 level (ppm)
# TYPE awair_co2 gauge
awair_co2{sensor="living-room"} 970
# HELP awair_dew_point Dew point (C)
# TYPE awair_dew_point gauge
awair_dew_point{sensor="living-room"} 9.7
# HELP awair_humid Relative humidity (%)
# TYPE awair_humid gauge
awair_humid{sensor="living-room"} 49.62
# HELP awair_pm25 Particulate matter (ug/m^3)
# TYPE awair_pm25 gauge
awair_pm25{sensor="living-room"} 5
# HELP awair_score Awair score (%)
# TYPE awair_score gauge
awair_score{sensor="living-room"} 91
# HELP awair_temp Temperature (C)
# TYPE awair_temp gauge
awair_temp{sensor="living-room"} 20.6
# HELP awair_voc Volatile organic compounds (ppb)
# TYPE awair_voc gauge
awair_voc{sensor="living-room"} 54
# HELP awair_co2_est ?
# TYPE awair_co2_est gauge
awair_co2_est{sensor="living-room"} 400
# HELP awair_pm10_est ?
# TYPE awair_pm10_est gauge
awair_pm10_est{sensor="living-room"} 6
# HELP awair_voc_baseline ?
# TYPE awair_voc_baseline gauge
awair_voc_baseline{sensor="living-room"} 2567673398
# HELP awair_voc_ethanol_raw ?
# TYPE awair_voc_ethanol_raw gauge
awair_voc_ethanol_raw{sensor="living-room"} 38
# HELP awair_voc_h2_raw ?
# TYPE awair_voc_h2_raw gauge
awair_voc_h2_raw{sensor="living-room"} 28
```

Enable the Local API for your Awair devices and set the `SENSORS` env
var to a list of the form `name1=ip1,name2=ip2,...`.

[Awair]: https://uk.getawair.com/
[Prometheus]: https://prometheus.io/
