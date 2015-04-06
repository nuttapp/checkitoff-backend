-- The strategy assumed is the Push on change model where an individual
-- change is pushed to all the appropriate users
-- http://highscalability.com/blog/2009/10/13/why-are-facebook-digg-and-twitter-so-hard-to-scale.html
SELECT * FROM system.schema_keyspaces;


-- DROP KEYSPACE demodb;

-- replication factor of 1 for test env
CREATE KEYSPACE demodb WITH replication = {
    'class': 'NetworkTopologyStrategy', 
    'datacenter1': '1'
};

USE demodb;

--
-- list
--
USE demodb;
DROP TABLE list;
CREATE TABLE list (
    list_id uuid,
    category varchar,
    title varchar,
    users set<varchar>,
    is_hidden boolean,
    created_at timeuuid,
    updated_at timeuuid,
    PRIMARY KEY (list_id)
);

INSERT INTO list (list_id, category, title, created_at, updated_at, users, is_hidden) VALUES (62c36092-82a1-3a00-93d1-46196ee77204, 'task', 'Today', now(), now(), {'greg', 'jess'}, false );
INSERT INTO list (list_id, category, title, created_at, updated_at, users, is_hidden) VALUES (62c36092-82a1-3a00-93d1-46196ee77205, 'shopping', 'Trader Joes', now(), now(), {'greg', 'jess'}, false);

UPDATE list SET category = 'task', title = 'updated title' WHERE list_id = 62c36092-82a1-3a00-93d1-46196ee77204;
UPDATE list SET users = users + {'son'} WHERE list_id = 62c36092-82a1-3a00-93d1-46196ee77204;

DELETE FROM List where list_id = 62c36092-82a1-3a00-93d1-46196ee77204;

DROP INDEX list_users_idx;
CREATE INDEX list_users_idx ON demodb.list (users);

select * from list where users contains 'greg';

SELECT list_id, category, title, dateOf(created_at) as created_at, dateOf(updated_at) as updated_at, users, is_hidden FROM list;

--
-- list
--
USE demodb;
DROP TABLE list_item;
CREATE TABLE list_item (
    list_id uuid,
    list_item_id uuid,
    title varchar,
    photo varchar,
    is_hidden boolean,
    created_at timeuuid,
    updated_at timeuuid,
    PRIMARY KEY (list_id, title, list_item_id)
);



-- DROP INDEX list_category_idx;
-- CREATE INDEX list_category_idx ON demodb.list (category);



-------------------------------------------------------- 
-- list_event - Stores create, update, delete history, 
-------------------------------------------------------- 

DROP TABLE list_event;
CREATE TABLE list_event (
    list_id uuid,
    event_id timeuuid,
    user_id uuid,
    event_type varchar,
    data varchar,
    PRIMARY KEY (list_id, event_id)
) WITH CLUSTERING ORDER BY (event_id DESC);

INSERT INTO list_event (list_id, user_id, event_id, event_type, data) VALUES (
    62c36092-82a1-3a00-93d1-46196ee77204, D68FBB1C-D7A9-40A5-A083-92EFE2B18BF4, now(),
    'create-list', '1');
INSERT INTO list_event (list_id, user_id, event_id, event_type, data) VALUES (
    62c36092-82a1-3a00-93d1-46196ee77204, D68FBB1C-D7A9-40A5-A083-92EFE2B18BF4, now(), 
    'update-list', '2');

SELECT list_id, event_id, user_id, event_type, data FROM list_event;

SELECT list_id, event_id, user_id, event_type, data FROM list_event where list_id = 62c36092-82a1-3a00-93d1-46196ee77204;

-------------------------------------------------------- 
-- user_timeline - Stores history of events related to this user
-------------------------------------------------------- 

DROP TABLE user_timeline;
CREATE TABLE user_timeline (
    user_id uuid,
    event_id timeuuid,
    event_type varchar,
    data varchar,
    PRIMARY KEY (user_id, event_id)
) WITH CLUSTERING ORDER BY (event_id DESC);
INSERT INTO user_timeline (user_id, event_id, event_type, data) VALUES (
    16FAFA86-FEB3-4DDD-AFED-A8BAA63F1D73, now(), 'create-list', '{"json": "stuff"}'
);

SELECT user_id, event_id, event_type, data FROM user_timeline;
--
-- list_user: users with access to particular lists
--

DROP TABLE list_user;
CREATE TABLE list_event (
    list_id uuid,
    id timeuuid,
    user_id uuid,
    type varchar,
    data varchar,
    PRIMARY KEY (list_id, id, user_id, type)
) WITH CLUSTERING ORDER BY (id DESC);


-- 
-- list_item
-- 

DROP TABLE list_item;
CREATE TABLE list_item (
    list_id uuid,
    id timeuuid,
    sort_order int,
    listTitle varchar,
    category varchar,
    title varchar,
    quantity int,
    imageUrl varchar,
    PRIMARY KEY (list_id, id)
);

INSERT INTO list_item (list_id, id, listTitle, category, title, quantity, imageUrl) VALUES (
    62c36092-82a1-3a00-93d1-46196ee77205, 58a31570-798d-11e4-b46d-8bd8f4d32221, 
    'Whole Foods', 'meats', 'Steak', 1, 'http://blah.blah');
INSERT INTO list_item (list_id, id, listTitle, category, title, quantity, imageUrl) VALUES (
    62c36092-82a1-3a00-93d1-46196ee77205, 35540780-79ad-11e4-b46d-8bd8f4d32221, 
    'Whole Foods', 'meats', 'Chicken', 1, 'http://blah.blah');
INSERT INTO list_item (list_id, id, listTitle, category, title, quantity, imageUrl) VALUES (
    62c36092-82a1-3a00-93d1-46196ee77205, 3554f1e0-79ad-11e4-b46d-8bd8f4d32221, 
    'Whole Foods', 'meats', 'Pork', 1, 'http://blah.blah');
INSERT INTO list_item (list_id, id, listTitle, category, title, quantity, imageUrl) VALUES (
    62c36092-82a1-3a00-93d1-46196ee77205, now(), 
    'Whole Foods', 'vegetables', 'brussle sprouts', 1, 'http://blah.blah');
select * from list_item;
SELECT list_id, id, listTitle, category, title, quantity, imageUrl FROM list_item
where list_id = 62c36092-82a1-3a00-93d1-46196ee77205 and id = 58a31570-798d-11e4-b46d-8bd8f4d32221;

SELECT list_id, id, listTitle, category, title, quantity, imageUrl FROM list_item
where list_id = 62c36092-82a1-3a00-93d1-46196ee77205;

DROP TABLE list_item_activity;
CREATE TABLE list_activity (
    list_id uuid,
    itemId timeuuid,
    user_id uuid,
    eventId timeuuid,
    eventType varchar,
    eventData varchar,
    PRIMARY KEY ((list_id, itemId, user_id, eventId, eventType)
);


