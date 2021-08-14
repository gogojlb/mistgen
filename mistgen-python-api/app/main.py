import falcon

from app import log
from app import mc
from app.api import switch
from app.api import sensor

LOG = log.get_logger()


class App(falcon.API):
    def __init__(self, *args, **kwargs):
        super(App, self).__init__(*args, **kwargs)
        LOG.info("API Server is starting")
        devices = mc.Devices()
        devices.get_devices()
        switch_handler = switch.Handler(devices.switch)
        sensor_handler = sensor.Handler(devices.gw, devices.sensor)
        self.add_route("/switch", switch_handler)
        self.add_route("/sensor", sensor_handler)

        LOG.info("API Server has been started")


application = App()
