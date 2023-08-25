-- create database Worganizer;

create table if not exists Poster (
    Id serial primary key,
    KpLink varchar not null,
    Rating real,
    Name varchar,
    Year int,
    CreatedAt timestamp default current_timestamp (now() at time zone 'utc-3') not null
);

create table if not exists PosterGenre (
    Id serial primary key,
    Genre varchar
);

create table if not exists Tab (
    Id serial primary key,
    Name varchar not null
);

create table if not exists TabChildren (
    Id1 int,
    Id2 int,
    foreign key (Id1) references Tab(Id) on delete cascade,
    foreign key (Id2) references Tab(Id) on delete cascade
);

create table if not exists TabQueue ( -- all the free fields are the primary key
    TabId int, -- serial primary key
    PosterId int, -- serial primary key
    foreign key (TabId) references Tab(Id) on delete cascade,
    foreign key (PosterId) references Poster(Id) on delete cascade,
    Position int not null
);

insert into tab values (DEFAULT, 'root')
