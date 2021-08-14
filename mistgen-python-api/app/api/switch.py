import falcon

from falcon.media.validators import jsonschema
from miio import ChuangmiPlug
from app import log

LOG = log.get_logger()

fields = {
    "is_on": {
        "type": "boolean"
    },
    "success": {
        "type": "boolean"
    },
    "error": {
        "type": "string"
    }
}


resp_schema = {
    "title": "responce",
    "type": "object",
    "properties": {
        "is_on": fields["is_on"],
        "success": fields["success"],
        "error": fields["error"],
    }
}


req_schema = {
    "title": "request",
    "type": "object",
    "properties": {
        "is_on": fields["is_on"]
    },
    "required": ["is_on"]
}


class Handler:
    switch = None

    def __init__(self, device):
        self.switch = device

    @jsonschema.validate(req_schema, resp_schema)
    def on_put(self, req, resp):
        plug = ChuangmiPlug(self.switch.get('localip'),
                            self.switch.get('token'))
        try:
            if req.media["is_on"]:
                LOG.info("Received request to switch on the plug.")
                plug.on()
                LOG.info("The plug has been switched on.")
            elif not req.media["is_on"]:
                LOG.info("Received request to switch off the plug.")
                plug.off()
                LOG.info("The plug has been switched off.")
            else:
                raise falcon.HTTPMissingParam("is_on")
            status = plug.status()
            resp.media = {
                "is_on": status.is_on,
                "success": True
            }
        except Exception:
            resp.status = 500
            LOG.error("Cannot execute command.")
            try:
                status = plug.status()
                resp.media = {
                    "is_on": status.is_on,
                    "success": False
                }
            except Exception:
                resp.media = {"success": False}
