from falcon.media.validators import jsonschema
from miio import Gateway

from app import log

LOG = log.get_logger()

fields = {
    "props": {
        "type": "array",
        "items": {
            "type": "string",
            "enum": ["temperature", "humidity"]}
    },
    "success": {
        "type": "boolean"
    },
    "error": {
        "type": "string"
    },
    "data": {
        "type": "object",
        "properties": {
            "temperature": {"type": "number"},
            "humidity": {"type": "number"}
        }
    }
}


resp_schema = {
    "title": "responce",
    "type": "object",
    "properties": {
        "data": fields["data"],
        "success": fields["success"],
        "error": fields["error"],
    },
    "required": ["success"]
}


req_schema = {
    "title": "request",
    "type": "object",
    "properties": {
        "props": fields["props"]
    },
    "required": ["props"]
}


class Handler:
    switch = None

    def __init__(self, mc_gw, mc_sensor):
        gw = Gateway(mc_gw['localip'], mc_gw['token'])
        gw.discover_devices()
        devices = gw.devices
        sensor_id = mc_sensor['did']
        self.sensor = devices[sensor_id]

    @jsonschema.validate(req_schema, resp_schema)
    def on_get(self, req, resp):

        try:
            LOG.info("Received request to sensor for %s.", req.media["props"])
            props = req.media["props"]
            collected_props = self.sensor.get_property_exp(props)
            data = dict(zip(props, collected_props))
            LOG.info("The data has been collected %s.", data)
            resp.media = {
                "data": data,
                "success": True
            }
        except Exception():
            resp.status = 500
            LOG.error("Cannot collect data.")
            resp.media = {"success": False}
