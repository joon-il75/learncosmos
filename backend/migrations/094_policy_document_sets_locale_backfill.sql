CREATE TABLE IF NOT EXISTS policy_document_sets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_type TEXT NOT NULL CHECK (document_type IN ('terms', 'privacy')),
    version INT NOT NULL,
    required BOOLEAN NOT NULL DEFAULT true,
    active BOOLEAN NOT NULL DEFAULT false,
    effective_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (document_type, version)
);

INSERT INTO policy_document_sets (
    document_type,
    version,
    required,
    active,
    effective_at,
    created_at,
    updated_at
)
SELECT
    pd.type,
    pd.version,
    bool_or(pd.is_required),
    bool_or(pd.is_active),
    min(pd.effective_at),
    min(pd.effective_at),
    max(pd.updated_at)
FROM policy_documents pd
GROUP BY pd.type, pd.version
ON CONFLICT (document_type, version) DO UPDATE
SET required = EXCLUDED.required,
    active = EXCLUDED.active,
    effective_at = EXCLUDED.effective_at,
    updated_at = NOW();

ALTER TABLE policy_documents
    ADD COLUMN IF NOT EXISTS set_id UUID REFERENCES policy_document_sets(id),
    ADD COLUMN IF NOT EXISTS locale TEXT NOT NULL DEFAULT 'ko',
    ADD COLUMN IF NOT EXISTS translation_status TEXT NOT NULL DEFAULT 'source';

UPDATE policy_documents pd
SET set_id = pds.id,
    locale = COALESCE(NULLIF(pd.locale, ''), 'ko'),
    translation_status = COALESCE(NULLIF(pd.translation_status, ''), 'source')
FROM policy_document_sets pds
WHERE pds.document_type = pd.type
  AND pds.version = pd.version
  AND pd.set_id IS NULL;

ALTER TABLE policy_documents
    ALTER COLUMN set_id SET NOT NULL;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'policy_documents_type_version_key'
          AND conrelid = 'public.policy_documents'::regclass
    ) THEN
        ALTER TABLE policy_documents
            DROP CONSTRAINT policy_documents_type_version_key;
    END IF;
END $$;

DROP INDEX IF EXISTS idx_policy_documents_active_type;

CREATE UNIQUE INDEX IF NOT EXISTS idx_policy_document_sets_active_required_type
ON policy_document_sets(document_type)
WHERE active = true AND required = true;

CREATE UNIQUE INDEX IF NOT EXISTS idx_policy_documents_active_set_locale
ON policy_documents(set_id, locale)
WHERE is_active = true;

CREATE INDEX IF NOT EXISTS idx_policy_documents_type_locale_active
ON policy_documents(type, locale, is_active);

ALTER TABLE user_policy_consents
    ADD COLUMN IF NOT EXISTS set_id UUID REFERENCES policy_document_sets(id),
    ADD COLUMN IF NOT EXISTS agreed_locale TEXT NOT NULL DEFAULT 'ko',
    ADD COLUMN IF NOT EXISTS requested_locale TEXT NOT NULL DEFAULT 'ko';

UPDATE user_policy_consents upc
SET set_id = pd.set_id,
    agreed_locale = COALESCE(NULLIF(upc.agreed_locale, ''), pd.locale, 'ko'),
    requested_locale = COALESCE(NULLIF(upc.requested_locale, ''), pd.locale, 'ko')
FROM policy_documents pd
WHERE pd.id = upc.policy_document_id
  AND upc.set_id IS NULL;

ALTER TABLE user_policy_consents
    ALTER COLUMN set_id SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_policy_consents_user_set
ON user_policy_consents(user_id, set_id);

CREATE INDEX IF NOT EXISTS idx_user_policy_consents_set_id
ON user_policy_consents(set_id);
