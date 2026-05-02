ALTER TABLE channel_configs DROP CONSTRAINT IF EXISTS fk_channel_configs_template_id;

DROP TABLE IF EXISTS templates;
