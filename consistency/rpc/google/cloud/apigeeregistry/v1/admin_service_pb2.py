# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: google/cloud/apigeeregistry/v1/admin_service.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.api import annotations_pb2 as google_dot_api_dot_annotations__pb2
from google.api import client_pb2 as google_dot_api_dot_client__pb2
from google.api import field_behavior_pb2 as google_dot_api_dot_field__behavior__pb2
from google.api import resource_pb2 as google_dot_api_dot_resource__pb2
from google.cloud.apigeeregistry.v1 import admin_models_pb2 as google_dot_cloud_dot_apigeeregistry_dot_v1_dot_admin__models__pb2
from google.longrunning import operations_pb2 as google_dot_longrunning_dot_operations__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2
from google.protobuf import field_mask_pb2 as google_dot_protobuf_dot_field__mask__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n2google/cloud/apigeeregistry/v1/admin_service.proto\x12\x1egoogle.cloud.apigeeregistry.v1\x1a\x1cgoogle/api/annotations.proto\x1a\x17google/api/client.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x19google/api/resource.proto\x1a\x31google/cloud/apigeeregistry/v1/admin_models.proto\x1a#google/longrunning/operations.proto\x1a\x1bgoogle/protobuf/empty.proto\x1a google/protobuf/field_mask.proto\"&\n\x16MigrateDatabaseRequest\x12\x0c\n\x04kind\x18\x01 \x01(\t\"\x19\n\x17MigrateDatabaseMetadata\"*\n\x17MigrateDatabaseResponse\x12\x0f\n\x07message\x18\x01 \x01(\t\"L\n\x13ListProjectsRequest\x12\x11\n\tpage_size\x18\x01 \x01(\x05\x12\x12\n\npage_token\x18\x02 \x01(\t\x12\x0e\n\x06\x66ilter\x18\x03 \x01(\t\"j\n\x14ListProjectsResponse\x12\x39\n\x08projects\x18\x01 \x03(\x0b\x32\'.google.cloud.apigeeregistry.v1.Project\x12\x17\n\x0fnext_page_token\x18\x02 \x01(\t\"P\n\x11GetProjectRequest\x12;\n\x04name\x18\x01 \x01(\tB-\xe0\x41\x02\xfa\x41\'\n%apigeeregistry.googleapis.com/Project\"i\n\x14\x43reateProjectRequest\x12=\n\x07project\x18\x01 \x01(\x0b\x32\'.google.cloud.apigeeregistry.v1.ProjectB\x03\xe0\x41\x02\x12\x12\n\nproject_id\x18\x02 \x01(\t\"\x9d\x01\n\x14UpdateProjectRequest\x12=\n\x07project\x18\x01 \x01(\x0b\x32\'.google.cloud.apigeeregistry.v1.ProjectB\x03\xe0\x41\x02\x12/\n\x0bupdate_mask\x18\x02 \x01(\x0b\x32\x1a.google.protobuf.FieldMask\x12\x15\n\rallow_missing\x18\x03 \x01(\x08\"b\n\x14\x44\x65leteProjectRequest\x12;\n\x04name\x18\x01 \x01(\tB-\xe0\x41\x02\xfa\x41\'\n%apigeeregistry.googleapis.com/Project\x12\r\n\x05\x66orce\x18\x02 \x01(\x08\x32\xb0\t\n\x05\x41\x64min\x12_\n\tGetStatus\x12\x16.google.protobuf.Empty\x1a&.google.cloud.apigeeregistry.v1.Status\"\x12\x82\xd3\xe4\x93\x02\x0c\x12\n/v1/status\x12\x62\n\nGetStorage\x12\x16.google.protobuf.Empty\x1a\'.google.cloud.apigeeregistry.v1.Storage\"\x13\x82\xd3\xe4\x93\x02\r\x12\x0b/v1/storage\x12\xba\x01\n\x0fMigrateDatabase\x12\x36.google.cloud.apigeeregistry.v1.MigrateDatabaseRequest\x1a\x1d.google.longrunning.Operation\"P\x82\xd3\xe4\x93\x02\x15\"\x13/v1/migrateDatabase\xca\x41\x32\n\x17MigrateDatabaseResponse\x12\x17MigrateDatabaseMetadata\x12\x8f\x01\n\x0cListProjects\x12\x33.google.cloud.apigeeregistry.v1.ListProjectsRequest\x1a\x34.google.cloud.apigeeregistry.v1.ListProjectsResponse\"\x14\x82\xd3\xe4\x93\x02\x0e\x12\x0c/v1/projects\x12\x8e\x01\n\nGetProject\x12\x31.google.cloud.apigeeregistry.v1.GetProjectRequest\x1a\'.google.cloud.apigeeregistry.v1.Project\"$\x82\xd3\xe4\x93\x02\x17\x12\x15/v1/{name=projects/*}\xda\x41\x04name\x12\xa2\x01\n\rCreateProject\x12\x34.google.cloud.apigeeregistry.v1.CreateProjectRequest\x1a\'.google.cloud.apigeeregistry.v1.Project\"2\x82\xd3\xe4\x93\x02\x17\"\x0c/v1/projects:\x07project\xda\x41\x12project,project_id\x12\xb4\x01\n\rUpdateProject\x12\x34.google.cloud.apigeeregistry.v1.UpdateProjectRequest\x1a\'.google.cloud.apigeeregistry.v1.Project\"D\x82\xd3\xe4\x93\x02(2\x1d/v1/{project.name=projects/*}:\x07project\xda\x41\x13project,update_mask\x12\x83\x01\n\rDeleteProject\x12\x34.google.cloud.apigeeregistry.v1.DeleteProjectRequest\x1a\x16.google.protobuf.Empty\"$\x82\xd3\xe4\x93\x02\x17*\x15/v1/{name=projects/*}\xda\x41\x04name\x1a \xca\x41\x1d\x61pigeeregistry.googleapis.comB]\n\"com.google.cloud.apigeeregistry.v1B\x11\x41\x64minServiceProtoP\x01Z\"github.com/apigee/registry/rpc;rpcb\x06proto3')



_MIGRATEDATABASEREQUEST = DESCRIPTOR.message_types_by_name['MigrateDatabaseRequest']
_MIGRATEDATABASEMETADATA = DESCRIPTOR.message_types_by_name['MigrateDatabaseMetadata']
_MIGRATEDATABASERESPONSE = DESCRIPTOR.message_types_by_name['MigrateDatabaseResponse']
_LISTPROJECTSREQUEST = DESCRIPTOR.message_types_by_name['ListProjectsRequest']
_LISTPROJECTSRESPONSE = DESCRIPTOR.message_types_by_name['ListProjectsResponse']
_GETPROJECTREQUEST = DESCRIPTOR.message_types_by_name['GetProjectRequest']
_CREATEPROJECTREQUEST = DESCRIPTOR.message_types_by_name['CreateProjectRequest']
_UPDATEPROJECTREQUEST = DESCRIPTOR.message_types_by_name['UpdateProjectRequest']
_DELETEPROJECTREQUEST = DESCRIPTOR.message_types_by_name['DeleteProjectRequest']
MigrateDatabaseRequest = _reflection.GeneratedProtocolMessageType('MigrateDatabaseRequest', (_message.Message,), {
  'DESCRIPTOR' : _MIGRATEDATABASEREQUEST,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.MigrateDatabaseRequest)
  })
_sym_db.RegisterMessage(MigrateDatabaseRequest)

MigrateDatabaseMetadata = _reflection.GeneratedProtocolMessageType('MigrateDatabaseMetadata', (_message.Message,), {
  'DESCRIPTOR' : _MIGRATEDATABASEMETADATA,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.MigrateDatabaseMetadata)
  })
_sym_db.RegisterMessage(MigrateDatabaseMetadata)

MigrateDatabaseResponse = _reflection.GeneratedProtocolMessageType('MigrateDatabaseResponse', (_message.Message,), {
  'DESCRIPTOR' : _MIGRATEDATABASERESPONSE,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.MigrateDatabaseResponse)
  })
_sym_db.RegisterMessage(MigrateDatabaseResponse)

ListProjectsRequest = _reflection.GeneratedProtocolMessageType('ListProjectsRequest', (_message.Message,), {
  'DESCRIPTOR' : _LISTPROJECTSREQUEST,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.ListProjectsRequest)
  })
_sym_db.RegisterMessage(ListProjectsRequest)

ListProjectsResponse = _reflection.GeneratedProtocolMessageType('ListProjectsResponse', (_message.Message,), {
  'DESCRIPTOR' : _LISTPROJECTSRESPONSE,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.ListProjectsResponse)
  })
_sym_db.RegisterMessage(ListProjectsResponse)

GetProjectRequest = _reflection.GeneratedProtocolMessageType('GetProjectRequest', (_message.Message,), {
  'DESCRIPTOR' : _GETPROJECTREQUEST,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.GetProjectRequest)
  })
_sym_db.RegisterMessage(GetProjectRequest)

CreateProjectRequest = _reflection.GeneratedProtocolMessageType('CreateProjectRequest', (_message.Message,), {
  'DESCRIPTOR' : _CREATEPROJECTREQUEST,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.CreateProjectRequest)
  })
_sym_db.RegisterMessage(CreateProjectRequest)

UpdateProjectRequest = _reflection.GeneratedProtocolMessageType('UpdateProjectRequest', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEPROJECTREQUEST,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.UpdateProjectRequest)
  })
_sym_db.RegisterMessage(UpdateProjectRequest)

DeleteProjectRequest = _reflection.GeneratedProtocolMessageType('DeleteProjectRequest', (_message.Message,), {
  'DESCRIPTOR' : _DELETEPROJECTREQUEST,
  '__module__' : 'google.cloud.apigeeregistry.v1.admin_service_pb2'
  # @@protoc_insertion_point(class_scope:google.cloud.apigeeregistry.v1.DeleteProjectRequest)
  })
_sym_db.RegisterMessage(DeleteProjectRequest)

_ADMIN = DESCRIPTOR.services_by_name['Admin']
if _descriptor._USE_C_DESCRIPTORS == False:

  DESCRIPTOR._options = None
  DESCRIPTOR._serialized_options = b'\n\"com.google.cloud.apigeeregistry.v1B\021AdminServiceProtoP\001Z\"github.com/apigee/registry/rpc;rpc'
  _GETPROJECTREQUEST.fields_by_name['name']._options = None
  _GETPROJECTREQUEST.fields_by_name['name']._serialized_options = b'\340A\002\372A\'\n%apigeeregistry.googleapis.com/Project'
  _CREATEPROJECTREQUEST.fields_by_name['project']._options = None
  _CREATEPROJECTREQUEST.fields_by_name['project']._serialized_options = b'\340A\002'
  _UPDATEPROJECTREQUEST.fields_by_name['project']._options = None
  _UPDATEPROJECTREQUEST.fields_by_name['project']._serialized_options = b'\340A\002'
  _DELETEPROJECTREQUEST.fields_by_name['name']._options = None
  _DELETEPROJECTREQUEST.fields_by_name['name']._serialized_options = b'\340A\002\372A\'\n%apigeeregistry.googleapis.com/Project'
  _ADMIN._options = None
  _ADMIN._serialized_options = b'\312A\035apigeeregistry.googleapis.com'
  _ADMIN.methods_by_name['GetStatus']._options = None
  _ADMIN.methods_by_name['GetStatus']._serialized_options = b'\202\323\344\223\002\014\022\n/v1/status'
  _ADMIN.methods_by_name['GetStorage']._options = None
  _ADMIN.methods_by_name['GetStorage']._serialized_options = b'\202\323\344\223\002\r\022\013/v1/storage'
  _ADMIN.methods_by_name['MigrateDatabase']._options = None
  _ADMIN.methods_by_name['MigrateDatabase']._serialized_options = b'\202\323\344\223\002\025\"\023/v1/migrateDatabase\312A2\n\027MigrateDatabaseResponse\022\027MigrateDatabaseMetadata'
  _ADMIN.methods_by_name['ListProjects']._options = None
  _ADMIN.methods_by_name['ListProjects']._serialized_options = b'\202\323\344\223\002\016\022\014/v1/projects'
  _ADMIN.methods_by_name['GetProject']._options = None
  _ADMIN.methods_by_name['GetProject']._serialized_options = b'\202\323\344\223\002\027\022\025/v1/{name=projects/*}\332A\004name'
  _ADMIN.methods_by_name['CreateProject']._options = None
  _ADMIN.methods_by_name['CreateProject']._serialized_options = b'\202\323\344\223\002\027\"\014/v1/projects:\007project\332A\022project,project_id'
  _ADMIN.methods_by_name['UpdateProject']._options = None
  _ADMIN.methods_by_name['UpdateProject']._serialized_options = b'\202\323\344\223\002(2\035/v1/{project.name=projects/*}:\007project\332A\023project,update_mask'
  _ADMIN.methods_by_name['DeleteProject']._options = None
  _ADMIN.methods_by_name['DeleteProject']._serialized_options = b'\202\323\344\223\002\027*\025/v1/{name=projects/*}\332A\004name'
  _MIGRATEDATABASEREQUEST._serialized_start=352
  _MIGRATEDATABASEREQUEST._serialized_end=390
  _MIGRATEDATABASEMETADATA._serialized_start=392
  _MIGRATEDATABASEMETADATA._serialized_end=417
  _MIGRATEDATABASERESPONSE._serialized_start=419
  _MIGRATEDATABASERESPONSE._serialized_end=461
  _LISTPROJECTSREQUEST._serialized_start=463
  _LISTPROJECTSREQUEST._serialized_end=539
  _LISTPROJECTSRESPONSE._serialized_start=541
  _LISTPROJECTSRESPONSE._serialized_end=647
  _GETPROJECTREQUEST._serialized_start=649
  _GETPROJECTREQUEST._serialized_end=729
  _CREATEPROJECTREQUEST._serialized_start=731
  _CREATEPROJECTREQUEST._serialized_end=836
  _UPDATEPROJECTREQUEST._serialized_start=839
  _UPDATEPROJECTREQUEST._serialized_end=996
  _DELETEPROJECTREQUEST._serialized_start=998
  _DELETEPROJECTREQUEST._serialized_end=1096
  _ADMIN._serialized_start=1099
  _ADMIN._serialized_end=2299
# @@protoc_insertion_point(module_scope)
