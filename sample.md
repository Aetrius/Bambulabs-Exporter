```
# HELP ams_bed_temp_metric temperature of the ams bed
# TYPE ams_bed_temp_metric gauge
ams_bed_temp_metric{ams_number="0",tray_number="0"} 0
ams_bed_temp_metric{ams_number="0",tray_number="1"} 0
ams_bed_temp_metric{ams_number="0",tray_number="2"} 0
ams_bed_temp_metric{ams_number="0",tray_number="3"} 0
# HELP ams_humidity_metric humidity of the ams
# TYPE ams_humidity_metric gauge
ams_humidity_metric{ams_number="0"} 4
# HELP ams_temp_metric temperature of the ams
# TYPE ams_temp_metric gauge
ams_temp_metric{ams_number="0"} 30.7
# HELP ams_tray_color_metric ID of the ams with color hex values
# TYPE ams_tray_color_metric gauge
ams_tray_color_metric{ams_number="0",tray_color="000000FF",tray_number="1",tray_type="PLA"} 1
ams_tray_color_metric{ams_number="0",tray_color="AF7933FF",tray_number="0",tray_type="PLA"} 1
ams_tray_color_metric{ams_number="0",tray_color="FFFFFFFF",tray_number="2",tray_type="PLA"} 1
ams_tray_color_metric{ams_number="0",tray_color="FFFFFFFF",tray_number="3",tray_type="PLA"} 1
# HELP big_fan1_speed_metric Big Fan 1 Speed
# TYPE big_fan1_speed_metric gauge
big_fan1_speed_metric 0
# HELP big_fan2_speed_metric Big Fan 2 Speed
# TYPE big_fan2_speed_metric gauge
big_fan2_speed_metric 0
# HELP chamber_temper_metric Chamber Temperature of Printer
# TYPE chamber_temper_metric gauge
chamber_temper_metric 30
# HELP cooling_fan_speed_metric Cooling Fan Speed
# TYPE cooling_fan_speed_metric gauge
cooling_fan_speed_metric 0
# HELP fail_reason_metric Print Failure Reason
# TYPE fail_reason_metric gauge
fail_reason_metric 0
# HELP fan_gear_metric Fan Gear
# TYPE fan_gear_metric gauge
fan_gear_metric 0
# HELP layer_number_metric layer number of the print head in gcode
# TYPE layer_number_metric gauge
layer_number_metric 261
# HELP mc_percent_metric Percentage of Progress of print
# TYPE mc_percent_metric gauge
mc_percent_metric 36
# HELP mc_print_error_code_metric Print Progress Error Code
# TYPE mc_print_error_code_metric gauge
mc_print_error_code_metric 0
# HELP mc_print_stage_metric Print Progress Stage
# TYPE mc_print_stage_metric gauge
mc_print_stage_metric 2
# HELP mc_print_sub_stage_metric Print Progress Sub Stage
# TYPE mc_print_sub_stage_metric gauge
mc_print_sub_stage_metric 4
# HELP mc_remaining_time_metric Print Progress Remaining Time in minutes
# TYPE mc_remaining_time_metric gauge
mc_remaining_time_metric 1973
# HELP nozzle_target_temper_metric Nozzle Target Temperature Metric
# TYPE nozzle_target_temper_metric gauge
nozzle_target_temper_metric 0
# HELP nozzle_temper_metric Nozzle Temperature Metric
# TYPE nozzle_temper_metric gauge
nozzle_temper_metric 221
# HELP print_error_metric Print error int
# TYPE print_error_metric gauge
print_error_metric 0
# HELP wifi_signal_metric Wifi signal in dBm
# TYPE wifi_signal_metric gauge
wifi_signal_metric -40