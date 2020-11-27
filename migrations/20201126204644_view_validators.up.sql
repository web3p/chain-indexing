CREATE TABLE view_validators (
    id BIGSERIAL,
    operator_address VARCHAR NOT NULL,
    consensus_node_address VARCHAR,
    initial_delegator_address VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    jailed BOOL NOT NULL,
    joined_at_block_height BIGINT NOT NULL,
    power VARCHAR NOT NULL,
    unbonding_height BIGINT NULL,
    unbonding_completion_time BIGINT NULL,
    moniker VARCHAR NOT NULL,
    identity VARCHAR NULL,
    website VARCHAR NULL,
    security_contact VARCHAR NULL,
    details VARCHAR NULL,
    commission_rate VARCHAR NOT NULL,
    commission_max_rate VARCHAR NOT NULL,
    commission_max_change_rate VARCHAR NOT NULL,
    min_self_delegation VARCHAR NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (operator_address, consensus_node_address)
);

CREATE INDEX view_validators_operator_address_index ON view_validators(operator_address);
CREATE INDEX view_validators_consensus_node_address_index ON view_validators(consensus_node_address);
