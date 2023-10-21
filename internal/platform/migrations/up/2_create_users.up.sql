BEGIN;
DO $$ 
BEGIN 
    IF EXISTS(SELECT 1 FROM "Migrations" WHERE "Version" = 2) THEN
        RAISE NOTICE 'Migration create_users already applied, skipping';
        RETURN;
    END IF;

    CREATE TABLE IF NOT EXISTS "Users"
    (
        "Id" INTEGER NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
        "Name" TEXT NOT NULL UNIQUE,
        "Password" BYTEA NOT NULL
    );

    CREATE UNIQUE INDEX IF NOT EXISTS "Index_Users_Name" ON "Users"("Name");
    
    INSERT INTO "Migrations" ("Version", "Name") VALUES (2, 'create_users');
END $$;
COMMIT;