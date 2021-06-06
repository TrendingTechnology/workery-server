CREATE TABLE tenants (
    -- base
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    schema_name VARCHAR (63) NULL,
    alternate_name VARCHAR (255) NULL,
    description TEXT NOT NULL DEFAULT '',
    name VARCHAR (255) NULL,
    url VARCHAR (255) NULL,
    state SMALLINT NOT NULL,
    timezone VARCHAR (63) NOT NULL DEFAULT 'utc',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    old_id BIGINT NOT NULL DEFAULT 0,

    -- abstract_postal_address.py
    address_country VARCHAR (127) NOT NULL DEFAULT '',
    address_region VARCHAR (127) NOT NULL DEFAULT '',
    address_locality VARCHAR (127) NOT NULL DEFAULT '',
    post_office_box_number VARCHAR (255) NOT NULL DEFAULT '',
    postal_code VARCHAR (127) NOT NULL DEFAULT '',
    street_address VARCHAR (127) NOT NULL DEFAULT '',
    street_address_extra VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_geo_coorindate.py
    elevation FLOAT NOT NULL DEFAULT 0,
    latitude FLOAT NOT NULL DEFAULT 0,
    longitude FLOAT NOT NULL DEFAULT 0,

    -- abstract_contact_point.py
    area_served VARCHAR (127) NOT NULL DEFAULT '',
    available_language VARCHAR (127) NOT NULL DEFAULT '',
    contact_type VARCHAR (127) NOT NULL DEFAULT '',
    email VARCHAR (255) NOT NULL DEFAULT '',
    fax_number VARCHAR (127) NOT NULL DEFAULT '',
    telephone VARCHAR (127) NOT NULL DEFAULT '',
    telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone VARCHAR (127) NOT NULL DEFAULT '',
    other_telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone_type_of SMALLINT NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX idx_tenant_uuid
ON tenants (uuid);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    first_name VARCHAR (50) NULL,
    last_name VARCHAR (50) NULL,
    email VARCHAR (255) UNIQUE NOT NULL,
    password_algorithm VARCHAR (63) NOT NULL,
    password_hash VARCHAR (511) NOT NULL,
    state SMALLINT NOT NULL DEFAULT 0,
    role SMALLINT NOT NULL DEFAULT 0,
    timezone VARCHAR (63) NOT NULL DEFAULT 'utc',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    joined_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    salt VARCHAR (127) NOT NULL DEFAULT '',
    was_email_activated BOOLEAN NOT NULL DEFAULT FALSE,
    pr_access_code VARCHAR (127) NOT NULL DEFAULT '',
    pr_expiry_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_user_uuid
ON users (uuid);
CREATE UNIQUE INDEX idx_user_email
ON users (email);
CREATE INDEX idx_user_tenant_id
ON users (tenant_id);

CREATE TABLE insurance_requirements (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    text VARCHAR (31) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_insurance_requirement_uuid
ON insurance_requirements (uuid);
CREATE UNIQUE INDEX idx_insurance_requirement_text
ON insurance_requirements (tenant_id, text);
CREATE INDEX idx_insurance_requirement_tenant_id
ON insurance_requirements (tenant_id);

CREATE TABLE how_hear_about_us_items (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    text VARCHAR (127) NOT NULL DEFAULT '',
    sort_number SMALLINT NOT NULL DEFAULT 0,
    is_for_associate BOOLEAN NOT NULL DEFAULT FALSE,
    is_for_customer BOOLEAN NOT NULL DEFAULT FALSE,
    is_for_staff BOOLEAN NOT NULL DEFAULT FALSE,
    is_for_partner BOOLEAN NOT NULL DEFAULT FALSE,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_how_hear_about_us_item_uuid
ON how_hear_about_us_items (uuid);
CREATE INDEX idx_how_hear_about_us_item_tenant_id
ON how_hear_about_us_items (tenant_id);

CREATE TABLE skill_sets (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    category VARCHAR (127) NOT NULL DEFAULT '',
    sub_category VARCHAR (127) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_skill_sets_uuid
ON skill_sets (uuid);
CREATE INDEX idx_skill_sets_tenant_id
ON skill_sets (tenant_id);

CREATE TABLE skill_set_insurance_requirements (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    skill_set_id BIGINT NOT NULL,
    insurance_requirement_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (skill_set_id) REFERENCES skill_sets(id),
    FOREIGN KEY (insurance_requirement_id) REFERENCES insurance_requirements(id)
);
CREATE UNIQUE INDEX idx_skill_set_insurance_requirements_uuid
ON skill_set_insurance_requirements (uuid);
CREATE INDEX idx_skill_set_insurance_requirements_tenant_id
ON skill_set_insurance_requirements (tenant_id);
CREATE INDEX idx_skill_set_insurance_requirements_skill_set_id
ON skill_set_insurance_requirements (skill_set_id);
CREATE INDEX idx_skill_set_insurance_requirements_insurance_requirement_id
ON skill_set_insurance_requirements (insurance_requirement_id);

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    text VARCHAR (127) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_tags_uuid
ON tags (uuid);
CREATE INDEX idx_tags_tenant_id
ON tags (tenant_id);

CREATE TABLE vehicle_types (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    text VARCHAR (127) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_vehicle_types_uuid
ON vehicle_types (uuid);
CREATE INDEX idx_vehicle_types_tenant_id
ON vehicle_types (tenant_id);

CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NULL,
    created_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    last_modified_by_id BIGINT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    text TEXT NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_comment_uuid
ON comments (uuid);
CREATE INDEX idx_comment_tenant_id
ON comments (tenant_id);
CREATE INDEX idx_comment_created_by_id
ON comments (created_by_id);
CREATE INDEX idx_comment_last_modified_by_id
ON comments (last_modified_by_id);

CREATE TABLE work_order_service_fees (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NULL,
    created_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    last_modified_by_id BIGINT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    title VARCHAR (63) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    percentage FLOAT NOT NULL DEFAULT 0,
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_work_order_service_fee_uuid
ON work_order_service_fees (uuid);
CREATE INDEX idx_work_order_service_fee_tenant_id
ON work_order_service_fees (tenant_id);

CREATE TABLE bulletin_board_items (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    created_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_by_id BIGINT NOT NULL,
    last_modified_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    text TEXT NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id)
);
CREATE UNIQUE INDEX idx_bulletin_board_item_uuid
ON bulletin_board_items (uuid);
CREATE INDEX idx_bulletin_board_item_tenant_id
ON bulletin_board_items (tenant_id);

CREATE TABLE customers (
    -- customer.py
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    type_of SMALLINT NOT NULL DEFAULT 0,
    indexed_text VARCHAR (511) NOT NULL DEFAULT '',
    is_ok_to_email BOOLEAN NOT NULL DEFAULT FALSE,
    is_ok_to_text BOOLEAN NOT NULL DEFAULT FALSE,
    is_business BOOLEAN NOT NULL DEFAULT FALSE,
    is_senior BOOLEAN NOT NULL DEFAULT FALSE,
    is_support BOOLEAN NOT NULL DEFAULT FALSE,
    job_info_read VARCHAR (127) NOT NULL DEFAULT '',
    how_hear_old SMALLINT NOT NULL DEFAULT 0,
    how_hear_id BIGINT NOT NULL,
    how_hear_other VARCHAR (2055) NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason_other VARCHAR (2055) NOT NULL DEFAULT '',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    created_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_by_id BIGINT NOT NULL,
    last_modified_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    organization_name VARCHAR (255) NOT NULL DEFAULT '',
    organization_type_of SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,

    -- abstract_postal_address.py
    address_country VARCHAR (127) NOT NULL DEFAULT '',
    address_region VARCHAR (127) NOT NULL DEFAULT '',
    address_locality VARCHAR (127) NOT NULL DEFAULT '',
    post_office_box_number VARCHAR (255) NOT NULL DEFAULT '',
    postal_code VARCHAR (127) NOT NULL DEFAULT '',
    street_address VARCHAR (127) NOT NULL DEFAULT '',
    street_address_extra VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_person.py
    given_name VARCHAR (63) NOT NULL DEFAULT '',
    middle_name VARCHAR (63) NOT NULL DEFAULT '',
    last_name VARCHAR (63) NOT NULL DEFAULT '',
    birthdate TIMESTAMP,
    join_date TIMESTAMP DEFAULT (now() AT TIME ZONE 'utc'),
    nationality VARCHAR (63) NOT NULL DEFAULT '',
    gender VARCHAR (31) NOT NULL DEFAULT '',
    tax_id VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_geo_coorindate.py
    elevation FLOAT NOT NULL DEFAULT 0,
    latitude FLOAT NOT NULL DEFAULT 0,
    longitude FLOAT NOT NULL DEFAULT 0,

    -- abstract_contact_point.py
    area_served VARCHAR (127) NOT NULL DEFAULT '',
    available_language VARCHAR (127) NOT NULL DEFAULT '',
    contact_type VARCHAR (127) NOT NULL DEFAULT '',
    email VARCHAR (255) NOT NULL DEFAULT '',
    fax_number VARCHAR (127) NOT NULL DEFAULT '',
    telephone VARCHAR (127) NOT NULL DEFAULT '',
    telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone VARCHAR (127) NOT NULL DEFAULT '',
    other_telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone_type_of SMALLINT NOT NULL DEFAULT 0,

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (how_hear_id) REFERENCES how_hear_about_us_items(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_customer_uuid
ON customers (uuid);
CREATE INDEX idx_customer_user_id
ON customers (user_id);

CREATE TABLE customer_tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
CREATE UNIQUE INDEX idx_customer_tag_uuid
ON customer_tags (uuid);
-- TODO: INDEXES

CREATE TABLE customer_comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);
CREATE UNIQUE INDEX idx_customer_comment_uuid
ON customer_comments (uuid);
CREATE UNIQUE INDEX idx_customer_comment_customer_id_comment_id
ON customer_comments (tenant_id, customer_id, comment_id);
-- TODO: INDEXES

-- TODO: avatar_image -> PrivateImageUpload

CREATE TABLE associates (
    -- customer.py
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    type_of SMALLINT NOT NULL DEFAULT 0,
    organization_name VARCHAR (255) NOT NULL DEFAULT '',
    organization_type_of SMALLINT NOT NULL DEFAULT 0,
    business VARCHAR (63) NOT NULL DEFAULT '',
    indexed_text VARCHAR (511) NOT NULL DEFAULT '',
    is_ok_to_email BOOLEAN NOT NULL DEFAULT FALSE,
    is_ok_to_text BOOLEAN NOT NULL DEFAULT FALSE,
    hourly_salary_desired SMALLINT NOT NULL DEFAULT 0,
    limit_special VARCHAR (255) NOT NULL DEFAULT '',
    dues_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    commercial_insurance_expiry_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    auto_insurance_expiry_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    wsib_number VARCHAR (127) NOT NULL DEFAULT '',
    wsib_insurance_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    police_check TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    drivers_license_class VARCHAR (31) NOT NULL DEFAULT '',
    how_hear_old SMALLINT NOT NULL DEFAULT 0,
    how_hear_id BIGINT NOT NULL,
    how_hear_other VARCHAR (2055) NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason_other VARCHAR (2055) NOT NULL DEFAULT '',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    created_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_by_id BIGINT NOT NULL,
    last_modified_from_ip VARCHAR (50) NOT NULL DEFAULT '',
    score FLOAT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    -- away_log_id BIGINT NOT NULL,

    -- abstract_postal_address.py
    address_country VARCHAR (127) NOT NULL DEFAULT '',
    address_region VARCHAR (127) NOT NULL DEFAULT '',
    address_locality VARCHAR (127) NOT NULL DEFAULT '',
    post_office_box_number VARCHAR (255) NOT NULL DEFAULT '',
    postal_code VARCHAR (127) NOT NULL DEFAULT '',
    street_address VARCHAR (127) NOT NULL DEFAULT '',
    street_address_extra VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_person.py
    given_name VARCHAR (63) NOT NULL DEFAULT '',
    middle_name VARCHAR (63) NOT NULL DEFAULT '',
    last_name VARCHAR (63) NOT NULL DEFAULT '',
    birthdate TIMESTAMP NULL,
    join_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    nationality VARCHAR (63) NOT NULL DEFAULT '',
    gender VARCHAR (31) NOT NULL DEFAULT '',
    tax_id VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_geo_coorindate.py
    elevation FLOAT NOT NULL DEFAULT 0,
    latitude FLOAT NOT NULL DEFAULT 0,
    longitude FLOAT NOT NULL DEFAULT 0,

    -- abstract_contact_point.py
    area_served VARCHAR (127) NOT NULL DEFAULT '',
    available_language VARCHAR (127) NOT NULL DEFAULT '',
    contact_type VARCHAR (127) NOT NULL DEFAULT '',
    email VARCHAR (255) NOT NULL DEFAULT '',
    fax_number VARCHAR (127) NOT NULL DEFAULT '',
    telephone VARCHAR (127) NOT NULL DEFAULT '',
    telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone VARCHAR (127) NOT NULL DEFAULT '',
    other_telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    balance_owing_amount FLOAT NOT NULL DEFAULT 0,
    -- avatar_image_id BIGINT NOT NULL,
    emergency_contact_name VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_relationship VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_telephone VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_alternative_telephone VARCHAR (127) NOT NULL DEFAULT '',
    service_fee_id BIGINT NOT NULL,

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (how_hear_id) REFERENCES how_hear_about_us_items(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id),
    -- FOREIGN KEY (away_log_id) REFERENCES away_log(id),
    FOREIGN KEY (service_fee_id) REFERENCES work_order_service_fees(id)
);
CREATE UNIQUE INDEX idx_associate_uuid
ON associates (uuid);
-- TODO: INDEXES

CREATE TABLE associate_vehicle_types (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    associate_id BIGINT NOT NULL,
    vehicle_type_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types(id)
);
CREATE UNIQUE INDEX idx_associate_vehicle_type_uuid
ON associate_vehicle_types (uuid);
-- TODO: INDEXES

CREATE TABLE associate_skill_sets (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    associate_id BIGINT NOT NULL,
    skill_set_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (skill_set_id) REFERENCES skill_sets(id)
);
CREATE UNIQUE INDEX idx_associate_skill_set_uuid
ON associate_skill_sets (uuid);
-- TODO: INDEXES

CREATE TABLE associate_tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    associate_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
CREATE UNIQUE INDEX idx_associate_tag_uuid
ON associate_tags (uuid);
-- TODO: INDEXES

CREATE TABLE associate_comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    associate_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
CREATE UNIQUE INDEX idx_associate_comment_uuid
ON associate_comments (uuid);
-- TODO: INDEXES

CREATE TABLE associate_insurance_requirements (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    associate_id BIGINT NOT NULL,
    insurance_requirement_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (insurance_requirement_id) REFERENCES insurance_requirements(id)
);
CREATE UNIQUE INDEX idx_associate_insurance_requirement_uuid
ON associate_insurance_requirements (uuid);
-- TODO: INDEXES

CREATE TABLE ongoing_work_orders (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    associate_id BIGINT NULL,
    state SMALLINT NOT NULL DEFAULT 0,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NULL,
    created_from_ip VARCHAR (50) NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_by_id BIGINT NULL,
    last_modified_from_ip VARCHAR (50) NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_ongoing_work_order_uuid
ON ongoing_work_orders (uuid);
-- TODO: INDEXES

CREATE TABLE work_orders (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    associate_id BIGINT NULL,
    description TEXT NOT NULL DEFAULT '',
    assignment_date TIMESTAMP NULL DEFAULT (now() AT TIME ZONE 'utc'),
    is_ongoing BOOLEAN NOT NULL DEFAULT FALSE,
    is_home_support_service BOOLEAN NOT NULL DEFAULT FALSE,
    start_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    completion_date TIMESTAMP NULL DEFAULT (now() AT TIME ZONE 'utc'),
    hours FLOAT NOT NULL DEFAULT 0,
    type_of SMALLINT NOT NULL DEFAULT 0,
    indexed_text VARCHAR (2055) NOT NULL DEFAULT '',
    closing_reason SMALLINT NOT NULL DEFAULT 0,
    closing_reason_other VARCHAR (1024) NULL,
    closing_reason_comment VARCHAR (1024) NOT NULL DEFAULT '',
    latest_pending_task_id BIGINT NULL,
    ongoing_work_order_id BIGINT NULL,
    state SMALLINT NOT NULL DEFAULT 0,
    currency VARCHAR (3) NOT NULL DEFAULT 'CAD',
    was_survey_conducted BOOLEAN NOT NULL DEFAULT FALSE,
    no_survey_conducted_reason SMALLINT NULL DEFAULT 0,
    no_survey_conducted_reason_other VARCHAR (1024) NULL DEFAULT '',
    was_job_satisfactory BOOLEAN NOT NULL DEFAULT FALSE,
    was_job_finished_on_time_and_on_budget BOOLEAN NOT NULL DEFAULT FALSE,
    was_associate_punctual BOOLEAN NOT NULL DEFAULT FALSE,
    was_associate_professional BOOLEAN NOT NULL DEFAULT FALSE,
    would_customer_refer_our_organization BOOLEAN NOT NULL DEFAULT FALSE,
    score FLOAT NOT NULL DEFAULT 0,
    was_there_financials_inputted BOOLEAN NOT NULL DEFAULT FALSE,
    invoice_paid_to SMALLINT NULL DEFAULT 0,
    invoice_date TIMESTAMP NULL DEFAULT (now() AT TIME ZONE 'utc'),
    invoice_ids VARCHAR (127) NULL DEFAULT '',
    invoice_quote_amount FLOAT NOT NULL DEFAULT 0,
    invoice_labour_amount FLOAT NOT NULL DEFAULT 0,
    invoice_material_amount FLOAT NOT NULL DEFAULT 0,
    invoice_other_costs_amount FLOAT NOT NULL DEFAULT 0,
    invoice_quoted_material_amount FLOAT NOT NULL DEFAULT 0,
    invoice_quoted_labour_amount FLOAT NOT NULL DEFAULT 0,
    invoice_quoted_other_costs_amount FLOAT NOT NULL DEFAULT 0,
    invoice_total_quote_amount FLOAT NOT NULL DEFAULT 0,
    invoice_sub_total_amount FLOAT NOT NULL DEFAULT 0,
    invoice_tax_amount FLOAT NOT NULL DEFAULT 0,
    invoice_total_amount FLOAT NOT NULL DEFAULT 0,
    invoice_deposit_amount FLOAT NOT NULL DEFAULT 0,
    invoice_amount_due FLOAT NOT NULL DEFAULT 0,
    invoice_service_fee_amount FLOAT NOT NULL DEFAULT 0,
    invoice_actual_service_fee_amount_paid FLOAT NOT NULL DEFAULT 0,
    invoice_service_fee_id BIGINT NULL,
    invoice_service_fee_payment_date TIMESTAMP NULL,
    invoice_balance_owing_amount FLOAT NOT NULL DEFAULT 0,
    visits SMALLINT NOT NULL DEFAULT 0,
    cloned_from_id BIGINT NULL,
    ongoing_order_id BIGINT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NULL,
    created_from_ip VARCHAR (50) NULL,
    last_modified_by_id BIGINT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_from_ip VARCHAR (50) NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id),
    FOREIGN KEY (invoice_service_fee_id) REFERENCES work_order_service_fees(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (ongoing_order_id) REFERENCES ongoing_work_orders(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_work_order_uuid
ON work_orders (uuid);
CREATE INDEX idx_work_order_tenant_id
ON work_orders (tenant_id);
CREATE INDEX idx_work_order_customer_id
ON work_orders (customer_id);
CREATE INDEX idx_work_order_associate_id
ON work_orders (associate_id);
CREATE INDEX idx_work_order_invoice_service_fee_id
ON work_orders (invoice_service_fee_id);
CREATE INDEX idx_work_order_created_by_id
ON work_orders (created_by_id);
CREATE INDEX idx_work_order_ongoing_order_id
ON work_orders (ongoing_order_id);
CREATE INDEX idx_work_order_last_modified_by_id
ON work_orders (last_modified_by_id);
-- TODO: INDEXES

-- CREATE TABLE work_order_ongoings (
--     id BIGSERIAL PRIMARY KEY,
--     uuid VARCHAR (36) UNIQUE NOT NULL,
--     tenant_id BIGINT NOT NULL,
--     order_id BIGINT NOT NULL,
--     ongoing_order_id BIGINT NOT NULL,
--     created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
--     created_by_id BIGINT NOT NULL,
--     last_modified_by_id BIGINT NOT NULL,
--     last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
--     old_id BIGINT NOT NULL DEFAULT 0,
--     FOREIGN KEY (tenant_id) REFERENCES tenants(id),
--     FOREIGN KEY (order_id) REFERENCES ongoing_work_orders(id),
--     FOREIGN KEY (ongoing_order_id) REFERENCES associates(id),
--     FOREIGN KEY (created_by_id) REFERENCES users(id),
--     FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
-- );
-- CREATE UNIQUE INDEX idx_work_order_ongoing_uuid
-- ON work_order_ongoings (uuid);

CREATE TABLE activity_sheet_items (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    comment TEXT NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NULL,
    created_from_ip VARCHAR (50) NULL,
    associate_id BIGINT NOT NULL,
    order_id BIGINT NULL,
    state SMALLINT NOT NULL DEFAULT 0,
    ongoing_order_id BIGINT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (ongoing_order_id) REFERENCES ongoing_work_orders(id),
    FOREIGN KEY (associate_id) REFERENCES associates(id)
);
CREATE UNIQUE INDEX idx_activity_sheet_item_uuid
ON activity_sheet_items (uuid);
CREATE INDEX idx_activity_sheet_item_tenant_id
ON activity_sheet_items (tenant_id);
CREATE INDEX idx_activity_sheet_item_order_id
ON activity_sheet_items (order_id);
CREATE INDEX idx_activity_sheet_item_associate_id
ON activity_sheet_items (associate_id);

CREATE TABLE work_order_tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
CREATE UNIQUE INDEX idx_work_order_tag_uuid
ON work_order_tags (uuid);
CREATE INDEX idx_work_order_tag_tenant_id
ON work_order_tags (tenant_id);
CREATE INDEX idx_work_order_tag_order_id
ON work_order_tags (order_id);
CREATE INDEX idx_work_order_tag_tag_id
ON work_order_tags (tag_id);

CREATE TABLE work_order_skill_sets (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    skill_set_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (skill_set_id) REFERENCES skill_sets(id)
);
CREATE UNIQUE INDEX idx_work_order_skill_set_uuid
ON work_order_skill_sets (uuid);
CREATE INDEX idx_work_order_skill_set_tenant_id
ON work_order_skill_sets (tenant_id);
CREATE INDEX idx_work_order_skill_set_order_id
ON work_order_skill_sets (order_id);
CREATE INDEX idx_work_order_skill_set_skill_set_id
ON work_order_skill_sets (skill_set_id);

CREATE TABLE work_order_comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    order_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
CREATE UNIQUE INDEX idx_work_order_comment_uuid
ON work_order_comments (uuid);
CREATE INDEX idx_work_order_comment_tenant_id
ON work_order_comments (tenant_id);
CREATE INDEX idx_work_order_comment_order_id
ON work_order_comments (order_id);
CREATE INDEX idx_work_order_comment_comment_id
ON work_order_comments (comment_id);

CREATE TABLE work_order_invoices (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    revision_version SMALLINT NOT NULL DEFAULT 0,
    invoice_id VARCHAR (127) NOT NULL DEFAULT '',
    invoice_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    associate_name VARCHAR (26) NOT NULL DEFAULT '',
    associate_telephone VARCHAR (127) NOT NULL DEFAULT '',
    client_name VARCHAR (63) NOT NULL DEFAULT '',
    client_address VARCHAR (63) NOT NULL DEFAULT '',
    client_telephone VARCHAR (63) NOT NULL DEFAULT '',
    client_email VARCHAR (63) NULL DEFAULT '',
    line_01_qty FLOAT NOT NULL DEFAULT 0,
    line_01_desc VARCHAR (45) NOT NULL DEFAULT '',
    line_01_price FLOAT NOT NULL DEFAULT 0,
    line_01_amount FLOAT NOT NULL DEFAULT 0,
    line_02_qty FLOAT NULL DEFAULT 0,
    line_02_desc VARCHAR (45) NULL DEFAULT '',
    line_02_price FLOAT NULL DEFAULT 0,
    line_02_amount FLOAT NULL DEFAULT 0,
    line_03_qty FLOAT NULL DEFAULT 0,
    line_03_desc VARCHAR (45) NULL DEFAULT '',
    line_03_price FLOAT NULL DEFAULT 0,
    line_03_amount FLOAT NULL DEFAULT 0,
    line_04_qty FLOAT NULL DEFAULT 0,
    line_04_desc VARCHAR (45) NULL DEFAULT '',
    line_04_price FLOAT NULL DEFAULT 0,
    line_04_amount FLOAT NULL DEFAULT 0,
    line_05_qty FLOAT NULL DEFAULT 0,
    line_05_desc VARCHAR (45) NULL DEFAULT '',
    line_05_price FLOAT NULL DEFAULT 0,
    line_05_amount FLOAT NULL DEFAULT 0,
    line_06_qty FLOAT NULL DEFAULT 0,
    line_06_desc VARCHAR (45) NULL DEFAULT '',
    line_06_price FLOAT NULL DEFAULT 0,
    line_06_amount FLOAT NULL DEFAULT 0,
    line_07_qty FLOAT NULL DEFAULT 0,
    line_07_desc VARCHAR (45) NULL DEFAULT '',
    line_07_price FLOAT NULL DEFAULT 0,
    line_07_amount FLOAT NULL DEFAULT 0,
    line_08_qty FLOAT NULL DEFAULT 0,
    line_08_desc VARCHAR (45) NULL DEFAULT '',
    line_08_price FLOAT NULL DEFAULT 0,
    line_08_amount FLOAT NULL DEFAULT 0,
    line_09_qty FLOAT NULL DEFAULT 0,
    line_09_desc VARCHAR (45) NULL DEFAULT '',
    line_09_price FLOAT NULL DEFAULT 0,
    line_09_amount FLOAT NULL DEFAULT 0,
    line_10_qty FLOAT NULL DEFAULT 0,
    line_10_desc VARCHAR (45) NULL DEFAULT '',
    line_10_price FLOAT NULL DEFAULT 0,
    line_10_amount FLOAT NULL DEFAULT 0,
    line_11_qty FLOAT NULL DEFAULT 0,
    line_11_desc VARCHAR (45) NULL DEFAULT '',
    line_11_price FLOAT NULL DEFAULT 0,
    line_11_amount FLOAT NULL DEFAULT 0,
    line_12_qty FLOAT NULL DEFAULT 0,
    line_12_desc VARCHAR (45) NULL DEFAULT '',
    line_12_price FLOAT NULL DEFAULT 0,
    line_12_amount FLOAT NULL DEFAULT 0,
    line_13_qty FLOAT NULL DEFAULT 0,
    line_13_desc VARCHAR (45) NULL DEFAULT '',
    line_13_price FLOAT NULL DEFAULT 0,
    line_13_amount FLOAT NULL DEFAULT 0,
    line_14_qty FLOAT NULL DEFAULT 0,
    line_14_desc VARCHAR (45) NULL DEFAULT '',
    line_14_price FLOAT NULL DEFAULT 0,
    line_14_amount FLOAT NULL DEFAULT 0,
    line_15_qty FLOAT NULL DEFAULT 0,
    line_15_desc VARCHAR (45) NULL DEFAULT '',
    line_15_price FLOAT NULL DEFAULT 0,
    line_15_amount FLOAT NULL DEFAULT 0,
    invoice_quote_days SMALLINT NOT NULL DEFAULT 0,
    invoice_associate_tax VARCHAR (127) NULL DEFAULT '',
    invoice_quote_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    invoice_customers_approval VARCHAR (20) NOT NULL DEFAULT '',
    line_01_notes VARCHAR (80) NULL DEFAULT '',
    line_02_notes VARCHAR (40) NULL DEFAULT '',
    total_labour FLOAT NOT NULL DEFAULT 0,
    total_materials FLOAT NOT NULL DEFAULT 0,
    other_costs FLOAT NOT NULL DEFAULT 0,
    sub_total FLOAT NOT NULL DEFAULT 0,
    tax FLOAT NOT NULL DEFAULT 0,
    total FLOAT NOT NULL DEFAULT 0,
    deposit FLOAT NOT NULL DEFAULT 0,
    amount_due FLOAT NOT NULL DEFAULT 0,
    payment_amount FLOAT NOT NULL DEFAULT 0,
    payment_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    is_cash BOOLEAN NOT NULL DEFAULT FALSE,
    is_cheque BOOLEAN NOT NULL DEFAULT FALSE,
    is_debit BOOLEAN NOT NULL DEFAULT FALSE,
    is_credit BOOLEAN NOT NULL DEFAULT FALSE,
    is_other BOOLEAN NOT NULL DEFAULT FALSE,
    client_signature VARCHAR (47) NOT NULL DEFAULT '',
    associate_sign_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    associate_signature VARCHAR (29) NOT NULL DEFAULT '',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    last_modified_by_id BIGINT NOT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    old_id BIGINT NOT NULL DEFAULT 0,
    state SMALLINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_work_order_invoice_uuid
ON work_order_invoices (uuid);
CREATE INDEX idx_work_order_invoice_tenant_id
ON work_order_invoices (tenant_id);
CREATE INDEX idx_work_order_invoice_order_id
ON work_order_invoices (order_id);

CREATE TABLE work_order_deposits (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    paid_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    deposit_method SMALLINT NOT NULL DEFAULT 0,
    paid_to SMALLINT NOT NULL DEFAULT 0,
    currency VARCHAR (3) NOT NULL DEFAULT 'CAD',
    amount FLOAT NOT NULL DEFAULT 0,
    paid_for SMALLINT NOT NULL DEFAULT 0,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NULL,
    created_from_ip VARCHAR (50) NULL DEFAULT '',
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    last_modified_by_id BIGINT NULL,
    last_modified_from_ip VARCHAR (50) NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_work_order_deposit_uuid
ON work_order_deposits (uuid);
CREATE INDEX idx_work_order_deposit_order_id
ON work_order_deposits (order_id);
-- CREATE INDEX idx_work_order_deposit_tenant_id
-- ON work_order_deposits (tenant_id);

CREATE TABLE task_items (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    type_of SMALLINT NOT NULL DEFAULT 0,
    title VARCHAR (63) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    due_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    was_postponed BOOLEAN NOT NULL DEFAULT FALSE,
    closing_reason SMALLINT NOT NULL DEFAULT 0,
    closing_reason_other VARCHAR (1024) NOT NULL DEFAULT '',
    order_id BIGINT NOT NULL,
    ongoing_order_id BIGINT NULL,
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_from_ip VARCHAR (50) NULL DEFAULT '',
    created_by_id BIGINT NULL,
    last_modified_by_id BIGINT NULL,
    last_modified_from_ip VARCHAR (50) NULL DEFAULT '',
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    old_id BIGINT NOT NULL DEFAULT 0,
    state SMALLINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (order_id) REFERENCES work_orders(id),
    FOREIGN KEY (ongoing_order_id) REFERENCES ongoing_work_orders(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE INDEX idx_task_item_tenant_id
ON task_items (tenant_id);
CREATE INDEX idx_task_item_order_id
ON task_items (order_id);
CREATE INDEX idx_task_item_ongoing_order_id
ON task_items (ongoing_order_id);
CREATE INDEX idx_task_item_created_by_id
ON task_items (created_by_id);
CREATE INDEX idx_task_item_last_modified_by_id
ON task_items (last_modified_by_id);
CREATE UNIQUE INDEX idx_task_item_uuid
ON task_items (uuid);


-- ######################### --
-- CONTNUE CODING FROM BELOW --
-- ######################### --

-- TODO: work_order_activity_sheets

-------------


CREATE TABLE customers (
    -- customer.py
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    type_of SMALLINT NOT NULL DEFAULT 0,
    indexed_text VARCHAR (511) NOT NULL DEFAULT '',
    is_ok_to_email BOOLEAN NOT NULL DEFAULT FALSE,
    is_ok_to_text BOOLEAN NOT NULL DEFAULT FALSE,
    is_business BOOLEAN NOT NULL DEFAULT FALSE,
    is_senior BOOLEAN NOT NULL DEFAULT FALSE,
    is_support BOOLEAN NOT NULL DEFAULT FALSE,
    job_info_read VARCHAR (127) NOT NULL DEFAULT '',
    how_hear_old SMALLINT NOT NULL DEFAULT 0,
    how_hear_id BIGINT NOT NULL,
    how_hear_other VARCHAR (2055) NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason_other VARCHAR (2055) NOT NULL DEFAULT '',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    last_modified_by_id BIGINT NOT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    organization_name VARCHAR (255) NOT NULL DEFAULT '',
    organization_type_of SMALLINT NOT NULL DEFAULT 0,
    old_id BIGINT NOT NULL DEFAULT 0,

    -- abstract_postal_address.py
    address_country VARCHAR (127) NOT NULL DEFAULT '',
    address_region VARCHAR (127) NOT NULL DEFAULT '',
    address_locality VARCHAR (127) NOT NULL DEFAULT '',
    post_office_box_number VARCHAR (255) NOT NULL DEFAULT '',
    postal_code VARCHAR (127) NOT NULL DEFAULT '',
    street_address VARCHAR (127) NOT NULL DEFAULT '',
    street_address_extra VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_person.py
    given_name VARCHAR (63) NOT NULL DEFAULT '',
    middle_name VARCHAR (63) NOT NULL DEFAULT '',
    last_name VARCHAR (63) NOT NULL DEFAULT '',
    birthdate VARCHAR (31) NOT NULL DEFAULT '',
    join_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    nationality VARCHAR (63) NOT NULL DEFAULT '',
    gender VARCHAR (31) NOT NULL DEFAULT '',
    tax_id VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_geo_coorindate.py
    elevation FLOAT NOT NULL DEFAULT 0,
    latitude FLOAT NOT NULL DEFAULT 0,
    longitude FLOAT NOT NULL DEFAULT 0,

    -- abstract_contact_point.py
    area_served VARCHAR (127) NOT NULL DEFAULT '',
    available_language VARCHAR (127) NOT NULL DEFAULT '',
    contact_type VARCHAR (127) NOT NULL DEFAULT '',
    email VARCHAR (255) NOT NULL DEFAULT '',
    fax_number VARCHAR (127) NOT NULL DEFAULT '',
    telephone VARCHAR (127) NOT NULL DEFAULT '',
    telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone VARCHAR (127) NOT NULL DEFAULT '',
    other_telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone_type_of SMALLINT NOT NULL DEFAULT 0,

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (how_hear_id) REFERENCES how_hear_about_us_items(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_customer_uuid
ON customers (uuid);

CREATE TABLE customer_tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
CREATE UNIQUE INDEX idx_customer_tag_uuid
ON customer_tags (uuid);

CREATE TABLE customer_comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    customer_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id),
    FOREIGN KEY (customer_id) REFERENCES customers(id)
);
CREATE UNIQUE INDEX idx_customer_comment_uuid
ON customer_comments (uuid);

-- TODO: avatar_image -> PrivateImageUpload

CREATE TABLE staff (
    -- customer.py
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    type_of SMALLINT NOT NULL DEFAULT 0,
    indexed_text VARCHAR (511) NOT NULL DEFAULT '',
    is_ok_to_email BOOLEAN NOT NULL DEFAULT FALSE,
    is_ok_to_text BOOLEAN NOT NULL DEFAULT FALSE,
    how_hear_old SMALLINT NOT NULL DEFAULT 0,
    how_hear_id BIGINT NOT NULL,
    how_hear_other VARCHAR (2055) NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason_other VARCHAR (2055) NOT NULL DEFAULT '',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    last_modified_by_id BIGINT NOT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    old_id BIGINT NOT NULL DEFAULT 0,

    -- abstract_postal_address.py
    address_country VARCHAR (127) NOT NULL DEFAULT '',
    address_region VARCHAR (127) NOT NULL DEFAULT '',
    address_locality VARCHAR (127) NOT NULL DEFAULT '',
    post_office_box_number VARCHAR (255) NOT NULL DEFAULT '',
    postal_code VARCHAR (127) NOT NULL DEFAULT '',
    street_address VARCHAR (127) NOT NULL DEFAULT '',
    street_address_extra VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_person.py
    given_name VARCHAR (63) NOT NULL DEFAULT '',
    middle_name VARCHAR (63) NOT NULL DEFAULT '',
    last_name VARCHAR (63) NOT NULL DEFAULT '',
    birthdate VARCHAR (31) NOT NULL DEFAULT '',
    join_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    nationality VARCHAR (63) NOT NULL DEFAULT '',
    gender VARCHAR (31) NOT NULL DEFAULT '',
    tax_id VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_geo_coorindate.py
    elevation FLOAT NOT NULL DEFAULT 0,
    latitude FLOAT NOT NULL DEFAULT 0,
    longitude FLOAT NOT NULL DEFAULT 0,

    -- abstract_contact_point.py
    area_served VARCHAR (127) NOT NULL DEFAULT '',
    available_language VARCHAR (127) NOT NULL DEFAULT '',
    contact_type VARCHAR (127) NOT NULL DEFAULT '',
    email VARCHAR (255) NOT NULL DEFAULT '',
    fax_number VARCHAR (127) NOT NULL DEFAULT '',
    telephone VARCHAR (127) NOT NULL DEFAULT '',
    telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone VARCHAR (127) NOT NULL DEFAULT '',
    other_telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    -- avatar_image_id BIGINT NOT NULL,
    personal_email VARCHAR (255) NOT NULL DEFAULT '',
    emergency_contact_name VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_relationship VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_telephone VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_alternative_telephone VARCHAR (127) NOT NULL DEFAULT '',
    police_check TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (how_hear_id) REFERENCES how_hear_about_us_items(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_staff_uuid
ON staff (uuid);

CREATE TABLE staff_tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    staff_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (staff_id) REFERENCES staff(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
CREATE UNIQUE INDEX idx_staff_tag_uuid
ON staff_tags (uuid);

CREATE TABLE staff_comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    staff_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (staff_id) REFERENCES staff(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
CREATE UNIQUE INDEX idx_staff_comment_uuid
ON staff_comments (uuid);

------------------------------------------------------------

CREATE TABLE partners (
    -- customer.py
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    type_of SMALLINT NOT NULL DEFAULT 0,
    organization_name VARCHAR (255) NOT NULL DEFAULT '',
    organization_type_of SMALLINT NOT NULL DEFAULT 0,
    business VARCHAR (63) NOT NULL DEFAULT '',
    indexed_text VARCHAR (511) NOT NULL DEFAULT '',
    is_ok_to_email BOOLEAN NOT NULL DEFAULT FALSE,
    is_ok_to_text BOOLEAN NOT NULL DEFAULT FALSE,
    hourly_salary_desired SMALLINT NOT NULL DEFAULT 0,
    limit_special VARCHAR (255) NOT NULL DEFAULT '',
    dues_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    commercial_insurance_expiry_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    auto_insurance_expiry_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    wsib_number VARCHAR (127) NOT NULL DEFAULT '',
    wsib_insurance_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    police_check TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    drivers_license_class VARCHAR (31) NOT NULL DEFAULT '',
    how_hear_old SMALLINT NOT NULL DEFAULT 0,
    how_hear_id BIGINT NOT NULL,
    how_hear_other VARCHAR (2055) NOT NULL DEFAULT '',
    is_business BOOLEAN NOT NULL DEFAULT FALSE,
    is_senior BOOLEAN NOT NULL DEFAULT FALSE,
    is_support BOOLEAN NOT NULL DEFAULT FALSE,
    job_info_read VARCHAR (127) NOT NULL DEFAULT '',
    state SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason SMALLINT NOT NULL DEFAULT 0,
    deactivation_reason_other VARCHAR (2055) NOT NULL DEFAULT '',
    created_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    created_by_id BIGINT NOT NULL,
    last_modified_by_id BIGINT NOT NULL,
    last_modified_time TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    score FLOAT NOT NULL DEFAULT 0,
    -- away_log_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,

    -- abstract_postal_address.py
    address_country VARCHAR (127) NOT NULL DEFAULT '',
    address_region VARCHAR (127) NOT NULL DEFAULT '',
    address_locality VARCHAR (127) NOT NULL DEFAULT '',
    post_office_box_number VARCHAR (255) NOT NULL DEFAULT '',
    postal_code VARCHAR (127) NOT NULL DEFAULT '',
    street_address VARCHAR (127) NOT NULL DEFAULT '',
    street_address_extra VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_person.py
    given_name VARCHAR (63) NOT NULL DEFAULT '',
    middle_name VARCHAR (63) NOT NULL DEFAULT '',
    last_name VARCHAR (63) NOT NULL DEFAULT '',
    birthdate VARCHAR (31) NOT NULL DEFAULT '',
    join_date TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'utc'),
    nationality VARCHAR (63) NOT NULL DEFAULT '',
    gender VARCHAR (31) NOT NULL DEFAULT '',
    tax_id VARCHAR (127) NOT NULL DEFAULT '',

    -- abstract_geo_coorindate.py
    elevation FLOAT NOT NULL DEFAULT 0,
    latitude FLOAT NOT NULL DEFAULT 0,
    longitude FLOAT NOT NULL DEFAULT 0,

    -- abstract_contact_point.py
    area_served VARCHAR (127) NOT NULL DEFAULT '',
    available_language VARCHAR (127) NOT NULL DEFAULT '',
    contact_type VARCHAR (127) NOT NULL DEFAULT '',
    email VARCHAR (255) NOT NULL DEFAULT '',
    fax_number VARCHAR (127) NOT NULL DEFAULT '',
    telephone VARCHAR (127) NOT NULL DEFAULT '',
    telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone VARCHAR (127) NOT NULL DEFAULT '',
    other_telephone_extension VARCHAR (31) NOT NULL DEFAULT '',
    other_telephone_type_of SMALLINT NOT NULL DEFAULT 0,
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    balance_owing_amount FLOAT NOT NULL DEFAULT 0,
    -- avatar_image_id BIGINT NOT NULL,
    emergency_contact_name VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_relationship VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_telephone VARCHAR (127) NOT NULL DEFAULT '',
    emergency_contact_alternative_telephone VARCHAR (127) NOT NULL DEFAULT '',
    service_fee_id BIGINT NOT NULL,

    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (how_hear_id) REFERENCES how_hear_about_us_items(id),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (last_modified_by_id) REFERENCES users(id),
    -- FOREIGN KEY (away_log_id) REFERENCES away_log(id),
    FOREIGN KEY (service_fee_id) REFERENCES work_order_service_fees(id)
);
CREATE UNIQUE INDEX idx_partner_uuid
ON partners (uuid);

CREATE TABLE partner_vehicle_types (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    partner_id BIGINT NOT NULL,
    vehicle_type_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id),
    FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types(id)
);
CREATE UNIQUE INDEX idx_partner_vehicle_type_uuid
ON partner_vehicle_types (uuid);

CREATE TABLE partner_skill_sets (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    partner_id BIGINT NOT NULL,
    skill_set_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id),
    FOREIGN KEY (skill_set_id) REFERENCES skill_sets(id)
);
CREATE UNIQUE INDEX idx_partner_skill_set_uuid
ON partner_skill_sets (uuid);

CREATE TABLE partner_tags (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    partner_id BIGINT NOT NULL,
    tag_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id),
    FOREIGN KEY (tag_id) REFERENCES tags(id)
);
CREATE UNIQUE INDEX idx_partner_tag_uuid
ON partner_tags (uuid);

CREATE TABLE partner_comments (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    partner_id BIGINT NOT NULL,
    comment_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id),
    FOREIGN KEY (comment_id) REFERENCES comments(id)
);
CREATE UNIQUE INDEX idx_partner_comment_uuid
ON partner_comments (uuid);

CREATE TABLE partner_insurance_requirements (
    id BIGSERIAL PRIMARY KEY,
    uuid VARCHAR (36) UNIQUE NOT NULL,
    tenant_id BIGINT NOT NULL,
    partner_id BIGINT NOT NULL,
    insurance_requirement_id BIGINT NOT NULL,
    old_id BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (tenant_id) REFERENCES tenants(id),
    FOREIGN KEY (partner_id) REFERENCES partners(id),
    FOREIGN KEY (insurance_requirement_id) REFERENCES insurance_requirements(id)
);
CREATE UNIQUE INDEX idx_partner_insurance_requirement_uuid
ON partner_insurance_requirements (uuid);

-- TODO: Private File Upload
-- TODO: Private Image Upload
-- TODO: Public Image Upload
-- TODO: Unified Search Item
