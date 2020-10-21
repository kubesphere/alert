# Copyright 2018 The KubeSphere Authors. All rights reserved.
# Use of this source code is governed by a Apache license
# that can be found in the LICENSE file.

FROM dhoer/flyway:5.1.4-mysql-8.0.11-alpine

RUN apk add --no-cache mysql-client

COPY ./schema /flyway/sql
COPY ./ddl /flyway/sql/ddl
COPY ./scripts /flyway/sql/ddl

RUN adduser -D -g alerter -u 1002 alerter && \
    chown -R alerter:alerter /flyway
USER alerter