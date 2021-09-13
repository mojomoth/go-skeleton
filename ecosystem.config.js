module.exports = {
  apps : [{
    name   : "sos-detection-to-protection",
    script : "./sos-detection-to-protection",
    cron_restart: "* * * * *"
  }]
}
