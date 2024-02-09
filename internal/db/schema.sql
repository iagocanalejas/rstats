CREATE TABLE entity (
    id bigint NOT NULL,
    from_date timestamp with time zone NOT NULL,
    to_date timestamp with time zone,
    is_active boolean NOT NULL,
    name character varying(150) NOT NULL,
    known_names character varying(150)[] NOT NULL,
    type character varying(50) NOT NULL,
    symbol character varying(10),
    metadata jsonb NOT NULL,
    is_partnership boolean NOT NULL,
    normalized_name character varying(150) NOT NULL,
    parent_id bigint
);

CREATE TABLE entity_partnership (
    id bigint NOT NULL,
    is_active boolean NOT NULL,
    part_id bigint NOT NULL,
    target_id bigint NOT NULL
);

CREATE TABLE flag (
    id bigint NOT NULL,
    creation_date timestamp with time zone NOT NULL,
    name character varying(150) NOT NULL,
    tokens character varying(50)[] NOT NULL,
    verified boolean NOT NULL,
    qualifies_for_id bigint
);

CREATE TABLE league (
    id bigint NOT NULL,
    from_date timestamp with time zone NOT NULL,
    to_date timestamp with time zone,
    is_active boolean NOT NULL,
    name character varying(150) NOT NULL,
    symbol character varying(10) NOT NULL,
    gender character varying(10),
    parent_id bigint
);

CREATE TABLE participant (
    id bigint NOT NULL,
    club_name character varying(150),
    distance integer,
    laps time without time zone[] NOT NULL,
    lane smallint,
    series smallint,
    gender character varying(10) NOT NULL,
    category character varying(10) NOT NULL,
    club_id bigint NOT NULL,
    race_id bigint NOT NULL,
    handicap time without time zone,
    CONSTRAINT participant_distance_check CHECK ((distance >= 0)),
    CONSTRAINT participant_lane_check CHECK ((lane >= 0)),
    CONSTRAINT participant_series_check CHECK ((series >= 0))
);

CREATE TABLE penalty (
    id bigint NOT NULL,
    penalty integer NOT NULL,
    disqualification boolean NOT NULL,
    reason character varying(500),
    participant_id bigint NOT NULL,
    CONSTRAINT penalty_penalty_check CHECK ((penalty >= 0))
);

CREATE TABLE race (
    id bigint NOT NULL,
    creation_date timestamp with time zone NOT NULL,
    laps smallint,
    lanes smallint,
    town character varying(100),
    type character varying(50) NOT NULL,
    date date NOT NULL,
    day smallint NOT NULL,
    cancelled boolean NOT NULL,
    cancellation_reasons character varying(200)[] NOT NULL,
    race_name character varying(200),
    sponsor character varying(200),
    trophy_edition smallint,
    flag_edition smallint,
    modality character varying(15) NOT NULL,
    metadata jsonb NOT NULL,
    flag_id bigint,
    league_id bigint,
    organizer_id bigint,
    trophy_id bigint,
    associated_id bigint,
    gender character varying(15) NOT NULL,
    CONSTRAINT race_day_check CHECK ((day >= 0)),
    CONSTRAINT race_flag_edition_check CHECK ((flag_edition >= 0)),
    CONSTRAINT race_lanes_check CHECK ((lanes >= 0)),
    CONSTRAINT race_laps_check CHECK ((laps >= 0)),
    CONSTRAINT race_trophy_edition_check CHECK ((trophy_edition >= 0))
);

CREATE TABLE trophy (
    id bigint NOT NULL,
    creation_date timestamp with time zone NOT NULL,
    name character varying(150) NOT NULL,
    tokens character varying(50)[] NOT NULL,
    verified boolean NOT NULL,
    qualifies_for_id bigint
);
