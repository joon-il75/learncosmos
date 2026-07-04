ALTER TABLE policy_documents
ADD COLUMN IF NOT EXISTS is_draft BOOLEAN NOT NULL DEFAULT false;

CREATE UNIQUE INDEX IF NOT EXISTS idx_policy_documents_draft_type
ON policy_documents(type)
WHERE is_draft = true;
