DROP TABLE IF EXISTS executor;
DROP TABLE IF EXISTS nf_address_list;

alter table resource_type change rs_uri_tmpl rs_type_param TEXT;
alter table resource_filter change rs_filter_uri rs_filter_param TEXT;