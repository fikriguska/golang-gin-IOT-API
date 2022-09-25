CREATE DATABASE iot;

\c iot

create table user_person(
    id_user serial primary key,
    email varchar (255) not null,
    username varchar (255) not null,
    password varchar (255) not null,
    status boolean not null,
    token varchar (255) not null,
    is_admin boolean not null
);
create table hardware(
    id_hardware serial primary key,
    name varchar (255) not null,
    type varchar (255) not null,
    description varchar (255) not null
);
create table node(
    id_node serial primary key,
    name varchar (255) not null, 
    location varchar (255) not null, 
    id_hardware integer, 
    id_user integer not null, 
    foreign key (id_hardware) references hardware (id_hardware) on update cascade,
    foreign key (id_user) references user_person (id_user) on update cascade on delete cascade
);
create table sensor(
    id_sensor serial primary key, 
    name varchar (255) not null, 
    unit varchar (255) not null,
    id_hardware integer, 
    id_node integer not null, 
    foreign key (id_hardware) references hardware (id_hardware) on update cascade on delete cascade,   
    foreign key (id_node) references node (id_node) on update cascade on delete cascade
);
create table channel(
    time timestamp, 
    value float not null, 
    id_sensor integer not null, 
    foreign key (id_sensor) references sensor (id_sensor) on update cascade on delete cascade
);
