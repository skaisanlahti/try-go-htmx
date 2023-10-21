BEGIN;
DO $$ 
BEGIN 
    IF NOT EXISTS(SELECT 1 FROM "Migrations" WHERE "Version" = 2) THEN 
        RAISE NOTICE 'Migration create_users not applied, skipping';
        RETURN;
    END IF;

    DROP INDEX IF EXISTS "Index_Users_Name";
    DROP TABLE IF EXISTS "Users";
    
    DELETE FROM "Migrations" WHERE "Version" = 2;
END $$;
COMMIT;