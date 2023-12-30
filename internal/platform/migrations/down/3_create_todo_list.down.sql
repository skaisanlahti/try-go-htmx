BEGIN;
DO $$ 
BEGIN 
    IF NOT EXISTS(SELECT 1 FROM "Migrations" WHERE "Version" = 3) THEN 
        RAISE NOTICE 'Migration create_users not applied, skipping';
        RETURN;
    END IF;

    ALTER TABLE IF EXISTS "Todos" DROP CONSTRAINT "TodoListId";
    ALTER TABLE IF EXISTS "Todos" DROP COLUMN "TodoListId";
    DROP TABLE IF EXISTS "TodoLists";
    
    DELETE FROM "Migrations" WHERE "Version" = 3;
END $$;
COMMIT;