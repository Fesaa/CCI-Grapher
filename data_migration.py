import os
import json
import asyncio
import asyncpg
import asqlite

def chunk_list(lst, chunk_size):
    return [lst[i:i + chunk_size] for i in range(0, len(lst), chunk_size)]


async def migrate(conn, file: str, i: int) -> int:
    print(file)
    db = await asqlite.connect(f"./sql/{file}")
    async with db.cursor() as cursor:
        res = await cursor.execute("SELECT user_id,time,roles FROM msgs")
        rows = await res.fetchall()
        if rows == []:
            return i
        for chunk in chunk_list(rows, 10000):
            query = "INSERT INTO messages (channel_id,message_id,user_id,roles,time) VALUES"
            if chunk == []:
                return i
            for row in chunk:
                query += f"({file.removesuffix('.sql')},{i},'{row[0]}', '{row[2]}','{row[1]}'),"
                i += 1
            query = query[:-1] + " ON CONFLICT DO NOTHING;"
            await conn.execute(query)
    return i

async def migrate_usernames(conn):
    db = await asqlite.connect(f"./sql/usernames.sql")
    async with db.cursor() as cursor:
        res = await cursor.execute("SELECT user_id,name FROM usernames;")
        rows = await res.fetchall()
        query = "INSERT INTO usernames (user_id,username) VALUES"
        for row in rows:
            query += f"""({row[0]},'{row[1].replace("'", "''")}'),"""
        query = query[:-1] + "ON CONFLICT DO NOTHING;"
        await conn.execute(query)

async def main():
    c = json.load(open('config.json'))
    conn = await asyncpg.connect(c["psql"])

    i = 0
    for file in os.listdir("./sql"):
        if not file.endswith(".sql"):
            continue
        if file != "usernames.sql":
            i = await migrate(conn, file, i)
        else:
            await migrate_usernames(conn)

if __name__ == '__main__':
    asyncio.run(main())
