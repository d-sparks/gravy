# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: daily_prices.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import timestamp_pb2 as google_dot_protobuf_dot_timestamp__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='daily_prices.proto',
  package='dailyprices',
  syntax='proto3',
  serialized_options=b'Z?github.com/d-sparks/gravy/data/dailyprices/proto;dailyprices_pb',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x12\x64\x61ily_prices.proto\x12\x0b\x64\x61ilyprices\x1a\x1fgoogle/protobuf/timestamp.proto\"I\n\x07Request\x12-\n\ttimestamp\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0f\n\x07version\x18\x02 \x01(\x05\"\xce\x05\n\x05Stats\x12\r\n\x05\x61lpha\x18\x01 \x01(\x01\x12\x0c\n\x04\x62\x65ta\x18\x02 \x01(\x01\x12?\n\x0fmoving_averages\x18\x03 \x03(\x0b\x32&.dailyprices.Stats.MovingAveragesEntry\x12L\n\x16moving_average_returns\x18\n \x03(\x0b\x32,.dailyprices.Stats.MovingAverageReturnsEntry\x12?\n\x0fmoving_variance\x18\x05 \x03(\x0b\x32&.dailyprices.Stats.MovingVarianceEntry\x12;\n\rmoving_volume\x18\x08 \x03(\x0b\x32$.dailyprices.Stats.MovingVolumeEntry\x12L\n\x16moving_volume_variance\x18\t \x03(\x0b\x32,.dailyprices.Stats.MovingVolumeVarianceEntry\x12\x10\n\x08\x65xchange\x18\x04 \x01(\t\x12\x0c\n\x04mean\x18\x06 \x01(\x01\x12\x10\n\x08variance\x18\x07 \x01(\x01\x1a\x35\n\x13MovingAveragesEntry\x12\x0b\n\x03key\x18\x01 \x01(\x05\x12\r\n\x05value\x18\x02 \x01(\x01:\x02\x38\x01\x1a;\n\x19MovingAverageReturnsEntry\x12\x0b\n\x03key\x18\x01 \x01(\x05\x12\r\n\x05value\x18\x02 \x01(\x01:\x02\x38\x01\x1a\x35\n\x13MovingVarianceEntry\x12\x0b\n\x03key\x18\x01 \x01(\x05\x12\r\n\x05value\x18\x02 \x01(\x01:\x02\x38\x01\x1a\x33\n\x11MovingVolumeEntry\x12\x0b\n\x03key\x18\x01 \x01(\x05\x12\r\n\x05value\x18\x02 \x01(\x01:\x02\x38\x01\x1a;\n\x19MovingVolumeVarianceEntry\x12\x0b\n\x03key\x18\x01 \x01(\x05\x12\r\n\x05value\x18\x02 \x01(\x01:\x02\x38\x01\"P\n\x06Prices\x12\x0c\n\x04open\x18\x01 \x01(\x01\x12\r\n\x05\x63lose\x18\x02 \x01(\x01\x12\x0b\n\x03low\x18\x04 \x01(\x01\x12\x0c\n\x04high\x18\x05 \x01(\x01\x12\x0e\n\x06volume\x18\x06 \x01(\x01\"S\n\tPairStats\x12\r\n\x05\x66irst\x18\x01 \x01(\t\x12\x0e\n\x06second\x18\x02 \x01(\t\x12\x12\n\ncovariance\x18\x03 \x01(\x01\x12\x13\n\x0b\x63orrelation\x18\x04 \x01(\x01\"\xe3\x02\n\tDailyData\x12-\n\ttimestamp\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12\x0f\n\x07version\x18\x03 \x01(\x05\x12\x32\n\x06prices\x18\x01 \x03(\x0b\x32\".dailyprices.DailyData.PricesEntry\x12\x30\n\x05stats\x18\x04 \x03(\x0b\x32!.dailyprices.DailyData.StatsEntry\x12*\n\npair_stats\x18\x05 \x03(\x0b\x32\x16.dailyprices.PairStats\x1a\x42\n\x0bPricesEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\"\n\x05value\x18\x02 \x01(\x0b\x32\x13.dailyprices.Prices:\x02\x38\x01\x1a@\n\nStatsEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12!\n\x05value\x18\x02 \x01(\x0b\x32\x12.dailyprices.Stats:\x02\x38\x01\"W\n\x05Range\x12&\n\x02lb\x18\x01 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\x12&\n\x02ub\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.Timestamp\">\n\x0cTradingDates\x12.\n\ntimestamps\x18\x01 \x03(\x0b\x32\x1a.google.protobuf.Timestamp\":\n\x11NewSessionRequest\x12%\n\tsim_range\x18\x01 \x01(\x0b\x32\x12.dailyprices.Range\"\x14\n\x12NewSessionResponse2\xd6\x01\n\x04\x44\x61ta\x12\x35\n\x03Get\x12\x14.dailyprices.Request\x1a\x16.dailyprices.DailyData\"\x00\x12O\n\nNewSession\x12\x1e.dailyprices.NewSessionRequest\x1a\x1f.dailyprices.NewSessionResponse\"\x00\x12\x46\n\x13TradingDatesInRange\x12\x12.dailyprices.Range\x1a\x19.dailyprices.TradingDates\"\x00\x42\x41Z?github.com/d-sparks/gravy/data/dailyprices/proto;dailyprices_pbb\x06proto3'
  ,
  dependencies=[google_dot_protobuf_dot_timestamp__pb2.DESCRIPTOR,])




_REQUEST = _descriptor.Descriptor(
  name='Request',
  full_name='dailyprices.Request',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='dailyprices.Request.timestamp', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='version', full_name='dailyprices.Request.version', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=68,
  serialized_end=141,
)


_STATS_MOVINGAVERAGESENTRY = _descriptor.Descriptor(
  name='MovingAveragesEntry',
  full_name='dailyprices.Stats.MovingAveragesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.Stats.MovingAveragesEntry.key', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.Stats.MovingAveragesEntry.value', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=579,
  serialized_end=632,
)

_STATS_MOVINGAVERAGERETURNSENTRY = _descriptor.Descriptor(
  name='MovingAverageReturnsEntry',
  full_name='dailyprices.Stats.MovingAverageReturnsEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.Stats.MovingAverageReturnsEntry.key', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.Stats.MovingAverageReturnsEntry.value', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=634,
  serialized_end=693,
)

_STATS_MOVINGVARIANCEENTRY = _descriptor.Descriptor(
  name='MovingVarianceEntry',
  full_name='dailyprices.Stats.MovingVarianceEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.Stats.MovingVarianceEntry.key', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.Stats.MovingVarianceEntry.value', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=695,
  serialized_end=748,
)

_STATS_MOVINGVOLUMEENTRY = _descriptor.Descriptor(
  name='MovingVolumeEntry',
  full_name='dailyprices.Stats.MovingVolumeEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.Stats.MovingVolumeEntry.key', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.Stats.MovingVolumeEntry.value', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=750,
  serialized_end=801,
)

_STATS_MOVINGVOLUMEVARIANCEENTRY = _descriptor.Descriptor(
  name='MovingVolumeVarianceEntry',
  full_name='dailyprices.Stats.MovingVolumeVarianceEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.Stats.MovingVolumeVarianceEntry.key', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.Stats.MovingVolumeVarianceEntry.value', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=803,
  serialized_end=862,
)

_STATS = _descriptor.Descriptor(
  name='Stats',
  full_name='dailyprices.Stats',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='alpha', full_name='dailyprices.Stats.alpha', index=0,
      number=1, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='beta', full_name='dailyprices.Stats.beta', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='moving_averages', full_name='dailyprices.Stats.moving_averages', index=2,
      number=3, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='moving_average_returns', full_name='dailyprices.Stats.moving_average_returns', index=3,
      number=10, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='moving_variance', full_name='dailyprices.Stats.moving_variance', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='moving_volume', full_name='dailyprices.Stats.moving_volume', index=5,
      number=8, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='moving_volume_variance', full_name='dailyprices.Stats.moving_volume_variance', index=6,
      number=9, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='exchange', full_name='dailyprices.Stats.exchange', index=7,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='mean', full_name='dailyprices.Stats.mean', index=8,
      number=6, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='variance', full_name='dailyprices.Stats.variance', index=9,
      number=7, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[_STATS_MOVINGAVERAGESENTRY, _STATS_MOVINGAVERAGERETURNSENTRY, _STATS_MOVINGVARIANCEENTRY, _STATS_MOVINGVOLUMEENTRY, _STATS_MOVINGVOLUMEVARIANCEENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=144,
  serialized_end=862,
)


_PRICES = _descriptor.Descriptor(
  name='Prices',
  full_name='dailyprices.Prices',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='open', full_name='dailyprices.Prices.open', index=0,
      number=1, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='close', full_name='dailyprices.Prices.close', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='low', full_name='dailyprices.Prices.low', index=2,
      number=4, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='high', full_name='dailyprices.Prices.high', index=3,
      number=5, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='volume', full_name='dailyprices.Prices.volume', index=4,
      number=6, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=864,
  serialized_end=944,
)


_PAIRSTATS = _descriptor.Descriptor(
  name='PairStats',
  full_name='dailyprices.PairStats',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='first', full_name='dailyprices.PairStats.first', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='second', full_name='dailyprices.PairStats.second', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='covariance', full_name='dailyprices.PairStats.covariance', index=2,
      number=3, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='correlation', full_name='dailyprices.PairStats.correlation', index=3,
      number=4, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=946,
  serialized_end=1029,
)


_DAILYDATA_PRICESENTRY = _descriptor.Descriptor(
  name='PricesEntry',
  full_name='dailyprices.DailyData.PricesEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.DailyData.PricesEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.DailyData.PricesEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1255,
  serialized_end=1321,
)

_DAILYDATA_STATSENTRY = _descriptor.Descriptor(
  name='StatsEntry',
  full_name='dailyprices.DailyData.StatsEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dailyprices.DailyData.StatsEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='value', full_name='dailyprices.DailyData.StatsEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1323,
  serialized_end=1387,
)

_DAILYDATA = _descriptor.Descriptor(
  name='DailyData',
  full_name='dailyprices.DailyData',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='timestamp', full_name='dailyprices.DailyData.timestamp', index=0,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='version', full_name='dailyprices.DailyData.version', index=1,
      number=3, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='prices', full_name='dailyprices.DailyData.prices', index=2,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='stats', full_name='dailyprices.DailyData.stats', index=3,
      number=4, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='pair_stats', full_name='dailyprices.DailyData.pair_stats', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[_DAILYDATA_PRICESENTRY, _DAILYDATA_STATSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1032,
  serialized_end=1387,
)


_RANGE = _descriptor.Descriptor(
  name='Range',
  full_name='dailyprices.Range',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='lb', full_name='dailyprices.Range.lb', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ub', full_name='dailyprices.Range.ub', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1389,
  serialized_end=1476,
)


_TRADINGDATES = _descriptor.Descriptor(
  name='TradingDates',
  full_name='dailyprices.TradingDates',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='timestamps', full_name='dailyprices.TradingDates.timestamps', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1478,
  serialized_end=1540,
)


_NEWSESSIONREQUEST = _descriptor.Descriptor(
  name='NewSessionRequest',
  full_name='dailyprices.NewSessionRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='sim_range', full_name='dailyprices.NewSessionRequest.sim_range', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1542,
  serialized_end=1600,
)


_NEWSESSIONRESPONSE = _descriptor.Descriptor(
  name='NewSessionResponse',
  full_name='dailyprices.NewSessionResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1602,
  serialized_end=1622,
)

_REQUEST.fields_by_name['timestamp'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_STATS_MOVINGAVERAGESENTRY.containing_type = _STATS
_STATS_MOVINGAVERAGERETURNSENTRY.containing_type = _STATS
_STATS_MOVINGVARIANCEENTRY.containing_type = _STATS
_STATS_MOVINGVOLUMEENTRY.containing_type = _STATS
_STATS_MOVINGVOLUMEVARIANCEENTRY.containing_type = _STATS
_STATS.fields_by_name['moving_averages'].message_type = _STATS_MOVINGAVERAGESENTRY
_STATS.fields_by_name['moving_average_returns'].message_type = _STATS_MOVINGAVERAGERETURNSENTRY
_STATS.fields_by_name['moving_variance'].message_type = _STATS_MOVINGVARIANCEENTRY
_STATS.fields_by_name['moving_volume'].message_type = _STATS_MOVINGVOLUMEENTRY
_STATS.fields_by_name['moving_volume_variance'].message_type = _STATS_MOVINGVOLUMEVARIANCEENTRY
_DAILYDATA_PRICESENTRY.fields_by_name['value'].message_type = _PRICES
_DAILYDATA_PRICESENTRY.containing_type = _DAILYDATA
_DAILYDATA_STATSENTRY.fields_by_name['value'].message_type = _STATS
_DAILYDATA_STATSENTRY.containing_type = _DAILYDATA
_DAILYDATA.fields_by_name['timestamp'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_DAILYDATA.fields_by_name['prices'].message_type = _DAILYDATA_PRICESENTRY
_DAILYDATA.fields_by_name['stats'].message_type = _DAILYDATA_STATSENTRY
_DAILYDATA.fields_by_name['pair_stats'].message_type = _PAIRSTATS
_RANGE.fields_by_name['lb'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_RANGE.fields_by_name['ub'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_TRADINGDATES.fields_by_name['timestamps'].message_type = google_dot_protobuf_dot_timestamp__pb2._TIMESTAMP
_NEWSESSIONREQUEST.fields_by_name['sim_range'].message_type = _RANGE
DESCRIPTOR.message_types_by_name['Request'] = _REQUEST
DESCRIPTOR.message_types_by_name['Stats'] = _STATS
DESCRIPTOR.message_types_by_name['Prices'] = _PRICES
DESCRIPTOR.message_types_by_name['PairStats'] = _PAIRSTATS
DESCRIPTOR.message_types_by_name['DailyData'] = _DAILYDATA
DESCRIPTOR.message_types_by_name['Range'] = _RANGE
DESCRIPTOR.message_types_by_name['TradingDates'] = _TRADINGDATES
DESCRIPTOR.message_types_by_name['NewSessionRequest'] = _NEWSESSIONREQUEST
DESCRIPTOR.message_types_by_name['NewSessionResponse'] = _NEWSESSIONRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Request = _reflection.GeneratedProtocolMessageType('Request', (_message.Message,), {
  'DESCRIPTOR' : _REQUEST,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.Request)
  })
_sym_db.RegisterMessage(Request)

Stats = _reflection.GeneratedProtocolMessageType('Stats', (_message.Message,), {

  'MovingAveragesEntry' : _reflection.GeneratedProtocolMessageType('MovingAveragesEntry', (_message.Message,), {
    'DESCRIPTOR' : _STATS_MOVINGAVERAGESENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.Stats.MovingAveragesEntry)
    })
  ,

  'MovingAverageReturnsEntry' : _reflection.GeneratedProtocolMessageType('MovingAverageReturnsEntry', (_message.Message,), {
    'DESCRIPTOR' : _STATS_MOVINGAVERAGERETURNSENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.Stats.MovingAverageReturnsEntry)
    })
  ,

  'MovingVarianceEntry' : _reflection.GeneratedProtocolMessageType('MovingVarianceEntry', (_message.Message,), {
    'DESCRIPTOR' : _STATS_MOVINGVARIANCEENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.Stats.MovingVarianceEntry)
    })
  ,

  'MovingVolumeEntry' : _reflection.GeneratedProtocolMessageType('MovingVolumeEntry', (_message.Message,), {
    'DESCRIPTOR' : _STATS_MOVINGVOLUMEENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.Stats.MovingVolumeEntry)
    })
  ,

  'MovingVolumeVarianceEntry' : _reflection.GeneratedProtocolMessageType('MovingVolumeVarianceEntry', (_message.Message,), {
    'DESCRIPTOR' : _STATS_MOVINGVOLUMEVARIANCEENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.Stats.MovingVolumeVarianceEntry)
    })
  ,
  'DESCRIPTOR' : _STATS,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.Stats)
  })
_sym_db.RegisterMessage(Stats)
_sym_db.RegisterMessage(Stats.MovingAveragesEntry)
_sym_db.RegisterMessage(Stats.MovingAverageReturnsEntry)
_sym_db.RegisterMessage(Stats.MovingVarianceEntry)
_sym_db.RegisterMessage(Stats.MovingVolumeEntry)
_sym_db.RegisterMessage(Stats.MovingVolumeVarianceEntry)

Prices = _reflection.GeneratedProtocolMessageType('Prices', (_message.Message,), {
  'DESCRIPTOR' : _PRICES,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.Prices)
  })
_sym_db.RegisterMessage(Prices)

PairStats = _reflection.GeneratedProtocolMessageType('PairStats', (_message.Message,), {
  'DESCRIPTOR' : _PAIRSTATS,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.PairStats)
  })
_sym_db.RegisterMessage(PairStats)

DailyData = _reflection.GeneratedProtocolMessageType('DailyData', (_message.Message,), {

  'PricesEntry' : _reflection.GeneratedProtocolMessageType('PricesEntry', (_message.Message,), {
    'DESCRIPTOR' : _DAILYDATA_PRICESENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.DailyData.PricesEntry)
    })
  ,

  'StatsEntry' : _reflection.GeneratedProtocolMessageType('StatsEntry', (_message.Message,), {
    'DESCRIPTOR' : _DAILYDATA_STATSENTRY,
    '__module__' : 'daily_prices_pb2'
    # @@protoc_insertion_point(class_scope:dailyprices.DailyData.StatsEntry)
    })
  ,
  'DESCRIPTOR' : _DAILYDATA,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.DailyData)
  })
_sym_db.RegisterMessage(DailyData)
_sym_db.RegisterMessage(DailyData.PricesEntry)
_sym_db.RegisterMessage(DailyData.StatsEntry)

Range = _reflection.GeneratedProtocolMessageType('Range', (_message.Message,), {
  'DESCRIPTOR' : _RANGE,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.Range)
  })
_sym_db.RegisterMessage(Range)

TradingDates = _reflection.GeneratedProtocolMessageType('TradingDates', (_message.Message,), {
  'DESCRIPTOR' : _TRADINGDATES,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.TradingDates)
  })
_sym_db.RegisterMessage(TradingDates)

NewSessionRequest = _reflection.GeneratedProtocolMessageType('NewSessionRequest', (_message.Message,), {
  'DESCRIPTOR' : _NEWSESSIONREQUEST,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.NewSessionRequest)
  })
_sym_db.RegisterMessage(NewSessionRequest)

NewSessionResponse = _reflection.GeneratedProtocolMessageType('NewSessionResponse', (_message.Message,), {
  'DESCRIPTOR' : _NEWSESSIONRESPONSE,
  '__module__' : 'daily_prices_pb2'
  # @@protoc_insertion_point(class_scope:dailyprices.NewSessionResponse)
  })
_sym_db.RegisterMessage(NewSessionResponse)


DESCRIPTOR._options = None
_STATS_MOVINGAVERAGESENTRY._options = None
_STATS_MOVINGAVERAGERETURNSENTRY._options = None
_STATS_MOVINGVARIANCEENTRY._options = None
_STATS_MOVINGVOLUMEENTRY._options = None
_STATS_MOVINGVOLUMEVARIANCEENTRY._options = None
_DAILYDATA_PRICESENTRY._options = None
_DAILYDATA_STATSENTRY._options = None

_DATA = _descriptor.ServiceDescriptor(
  name='Data',
  full_name='dailyprices.Data',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=1625,
  serialized_end=1839,
  methods=[
  _descriptor.MethodDescriptor(
    name='Get',
    full_name='dailyprices.Data.Get',
    index=0,
    containing_service=None,
    input_type=_REQUEST,
    output_type=_DAILYDATA,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='NewSession',
    full_name='dailyprices.Data.NewSession',
    index=1,
    containing_service=None,
    input_type=_NEWSESSIONREQUEST,
    output_type=_NEWSESSIONRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='TradingDatesInRange',
    full_name='dailyprices.Data.TradingDatesInRange',
    index=2,
    containing_service=None,
    input_type=_RANGE,
    output_type=_TRADINGDATES,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_DATA)

DESCRIPTOR.services_by_name['Data'] = _DATA

# @@protoc_insertion_point(module_scope)
