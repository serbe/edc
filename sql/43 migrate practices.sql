ALTER TABLE
    practices
ADD
    CONSTRAINT practices_company_id_kind_id_date_of_practice_key UNIQUE (company_id, kind_id, date_of_practice);