BEGIN;
DO $$ 
BEGIN 
    IF EXISTS(SELECT 1 FROM "Migrations" WHERE "Version" = 1) THEN
        RAISE NOTICE 'Migration create_todos already applied, skipping';
        RETURN;
    END IF;

    CREATE TABLE IF NOT EXISTS public."Todos"
    (
        "Id" INTEGER NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        "Task" TEXT NOT NULL,
        "Done" BOOLEAN NOT NULL DEFAULT false
    );
    
    INSERT INTO "Migrations" ("Version", "Name") VALUES (1, 'create_todos');
END $$;
COMMIT;