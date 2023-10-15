BEGIN;
DO $$ 
BEGIN 
    IF NOT EXISTS(SELECT 1 FROM "Migrations" WHERE "Version" = 1) THEN 
        RAISE NOTICE 'Migration create_todos not applied, skipping';
        RETURN;
    END IF;

    DROP TABLE IF EXISTS "Todos";
    
    DELETE FROM "Migrations" WHERE "Version" = 1;
END $$;
COMMIT;