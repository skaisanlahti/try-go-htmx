BEGIN;
DO $$ 
BEGIN 
    IF EXISTS(SELECT 1 FROM "Migrations" WHERE "Version" = 3) THEN
        RAISE NOTICE 'Migration create_todo_lists already applied, skipping';
        RETURN;
    END IF;

    CREATE TABLE IF NOT EXISTS "TodoLists"
    (
        "Id" INTEGER NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        "Name" TEXT NOT NULL UNIQUE,
        "UserId" INTEGER NOT NULL,
        CONSTRAINT "UserId" FOREIGN KEY ("UserId") REFERENCES "Users"("Id") ON DELETE CASCADE
    );

    ALTER TABLE IF EXISTS "Todos" ADD COLUMN "TodoListId" INTEGER NOT NULL;
    ALTER TABLE IF EXISTS "Todos" ADD CONSTRAINT "TodoListId" FOREIGN KEY ("TodoListId") REFERENCES "TodoLists"("Id") ON DELETE CASCADE;

    INSERT INTO "Migrations" ("Version", "Name") VALUES (3, 'create_todo_lists');
END $$;
COMMIT;