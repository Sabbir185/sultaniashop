SELECT pid, state, query_start, query
FROM pg_stat_activity
WHERE
    query LIKE '%pg_sleep%'
    AND state != 'idle';


-- pg_sleep(20) --> 20 seconds delay for performance testing
-- http://localhost:8080/list --> 40 seconds delay (FOR 2 ROWS in DB)
