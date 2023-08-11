# watch_organizer

## DB:

Posters: id | kp_link, rating, name, year, created_date
Poster_genres: id | genre
Tabs: id | name
Tab_childrens: id1 | id2
Tab_queues: tab_id, poster_id | position (in a queue) 
Tab_posters: tab_id | poster_id

different tabs can have the same poster_id


## Mechanics:

search (in any tab) poster by name: select all in choosen tab & just sort by attribute
search tab by name
sort current queue (by rating, by name, by year, by created_poster_date)

add new poster:
1. form the poster (data): usually just kp_link. mb search in kp by name via API
2. form the path:
	2.1. let user write down the path like tab1/tab2/tab3/.../last_tab. on each backslash 		show the user possible tabs
	2.2. create all the tabs that don't exist 
3. save created poster and path to the DB.

reorg queue: just like delete & insert element in the array.
move poster to other tab: delete & add row in tab_posters.
move tab: the same with tab_childrens.

used poster checks as watched and goes to "watched" tab
	from "watched" tab posters can be deleted
!hide the watched posters from root tab in UI!


## Problems:
1. Rating updating. Possible sollution: update all posters at once using link.
2. How to reduce queries to DB? Cache?
