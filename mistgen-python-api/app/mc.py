import re

from micloud import MiCloud

from app import config
from app import log

LOG = log.get_logger()

class Devices:
    gw = None
    switch = None
    sensor = None
    def __init__(self):
        self.mc = MiCloud(config.env('MI_USER'), config.env('MI_PASS'))
        try:
            self.mc.login()
        except:
            LOG.error('Cannot login on MI Cloud')
            raise
        LOG.info('MI Cloud success init')
    
    
    def get_devices(self):
        try:
            all_device_list = self.mc.get_devices(country=config.env('MI_COUNTRY'))
        except:
            LOG.error('Cannot get devices list')
            raise
        for device in all_device_list:
            LOG.info('Checking device id="%s", model="%s"..', device.get('did'),device.get('model'))
            if device.get('model') == config.env('GW_MODEL'):
                LOG.info('..It is Gateway model "%s"!', device.get('model') )
                self.gw = device
                LOG.info('Found Gateway')
                LOG.info(self.gw)
            elif device.get('model') == config.env('SWITCH_MODEL'):
                LOG.info('..It is Switch model "%s', device.get('model') )
                LOG.info("Checking device \"%s\" for regex \"%s\"",device.get('name'),config.env('DEVICE_NAME_REGEX'))
                if (re.match(config.env('DEVICE_NAME_REGEX'),device.get('name'))):
                    assert(self.switch==None),('Only one of Switchs should pass regex check!',self.switch,device)
                    LOG.info('....Regex matched. Switch id = %s', device.get('did'))
                    self.switch=device
                    LOG.info('Found Switch')
                    LOG.info(self.switch)
                else:
                    LOG.info('....Regex has not matched.')
            elif device.get('model') == config.env('SENSOR_MODEL'):
                LOG.info('..It is Sensor model "%s', device.get('model') )
                LOG.info("Checking device \"%s\" for regex \"%s\"",device.get('name'),config.env('DEVICE_NAME_REGEX'))
                if (re.match(config.env('DEVICE_NAME_REGEX'),device.get('name'))):
                    assert(self.sensor==None),('Only one of Sensors should pass regex check!',self.sensor,device)
                    LOG.info('....Regex matched. Sensor id = %s', device.get('did'))
                    self.sensor=device
                    LOG.info('Found Sensor')
                    LOG.info(self.sensor)
                else:
                    LOG.info('....Regex has not matched.')
        assert (self.sensor!=None), "Unable to collect Sensor!" 
        assert (self.switch!=None), "Unable to collect Switch!" 
        assert (self.gw!=None), "Unable to collect Gateway!" 

