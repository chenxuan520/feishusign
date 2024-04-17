
create table sign
(
    user_id     varchar(256) not null,
    meeting_id  varchar(256) not null,
    user_name   varchar(256) null,
    status      tinyint      not null,
    create_time bigint       not null,
    primary key (user_id, meeting_id)
);

create index idx_meeting_id
    on sign (meeting_id);


create table meeting
(
    meeting_id    varchar(256) not null
        primary key,
    originator_id varchar(256) null,
    year          int          null,
    month         int          null,
    day           int          null,
    create_time   bigint       null,
    url           varchar(256) null
);

