DROP INDEX IF EXISTS idx_policy_documents_draft_type;

CREATE UNIQUE INDEX IF NOT EXISTS idx_policy_documents_draft_type_locale
ON policy_documents(type, locale)
WHERE is_draft = true;
