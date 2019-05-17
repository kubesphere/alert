ALTER TABLE metric MODIFY metric_name varchar(100) NOT NULL;
ALTER TABLE history MODIFY resource_name varchar(255);
ALTER TABLE comment MODIFY content varchar(255) NOT NULL;