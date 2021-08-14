from envparse import Env

env = Env(
    MI_USER = str,
    MI_PASS = str,
    MI_COUNTRY = dict(cast=str,default='cn'),
    DEVICE_NAME_REGEX = str,
    LOG_LEVEL = dict(cast=str,default='INFO'),
    GW_MODEL = dict(cast=str, default='lumi.gateway.v3'),
    SWITCH_MODEL = dict(cast=str,default='chuangmi.plug.m3'),
    SENSOR_MODEL = dict(cast=str, default='lumi.sensor_ht.v1')
)
env.read_envfile()

