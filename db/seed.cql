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

DROP TABLE list;
CREATE TABLE list (
    list_id uuid,
    title varchar,
    created_at timeuuid,
    updated_at timeuuid,
    users set<varchar>,
    PRIMARY KEY (list_id)
);
INSERT INTO list (list_id, title, created_at, updated_at, users) VALUES (62c36092-82a1-3a00-93d1-46196ee77204, 'Trader Joes', now(), now(), {'1', '2'} );
INSERT INTO list (list_id, title, created_at, updated_at, users) VALUES (62c36092-82a1-3a00-93d1-46196ee77205, 'Trader Joes', now(), now(), {'10', '11'});

UPDATE list SET users = users + {'4'} WHERE list_id = 62c36092-82a1-3a00-93d1-46196ee77204;

SELECT list_id, title, created_at, updated_at, users FROM list;

-- DELETE FROM list WHERE accountId = 3BF39D34-4622-4FE8-89F1-1E56B875A50E
-- and list_id = 62c36092-82a1-3a00-93d1-46196ee77205 ;

-- DROP INDEX list_category_idx;
-- CREATE INDEX list_category_idx ON demodb.list (category);

DROP INDEX list_users_idx;
CREATE INDEX list_users_idx ON demodb.list (users);
select * from list where users contains '1';

SELECT * FROM list; --WHERE accountId = 3BF39D34-4622-4FE8-89F1-1E56B875A50E;

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


