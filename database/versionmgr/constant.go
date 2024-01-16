package versionmgr

var defaultVersionManagementTab = "app_version_management_tab"

var defaultVersionManagementSQL = `
create table if not exists %s (
	id int unsigned not null primary key,
	version bigint unsigned not null,
	ext_info text
);
`
