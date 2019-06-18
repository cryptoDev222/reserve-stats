--
-- PostgreSQL database dump
--

-- Dumped from database version 9.5.14
-- Dumped by pg_dump version 9.5.14

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: reserve; Type: TABLE; Schema: public; Owner: reserve_stats
--

CREATE TABLE public.reserve (
    id integer NOT NULL,
    address text NOT NULL
);


ALTER TABLE public.reserve OWNER TO reserve_stats;

--
-- Name: reserve_id_seq; Type: SEQUENCE; Schema: public; Owner: reserve_stats
--

CREATE SEQUENCE public.reserve_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.reserve_id_seq OWNER TO reserve_stats;

--
-- Name: reserve_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reserve_stats
--

ALTER SEQUENCE public.reserve_id_seq OWNED BY public.reserve.id;


--
-- Name: token; Type: TABLE; Schema: public; Owner: reserve_stats
--

CREATE TABLE public.token (
    id integer NOT NULL,
    address text NOT NULL
);


ALTER TABLE public.token OWNER TO reserve_stats;

--
-- Name: token_id_seq; Type: SEQUENCE; Schema: public; Owner: reserve_stats
--

CREATE SEQUENCE public.token_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.token_id_seq OWNER TO reserve_stats;

--
-- Name: token_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reserve_stats
--

ALTER SEQUENCE public.token_id_seq OWNED BY public.token.id;


--
-- Name: tradelogs; Type: TABLE; Schema: public; Owner: reserve_stats
--

CREATE TABLE public.tradelogs (
    id integer NOT NULL,
    "timestamp" timestamp without time zone,
    block_number integer,
    tx_hash text,
    eth_amount double precision,
    user_address_id bigint NOT NULL,
    src_address_id bigint NOT NULL,
    dest_address_id bigint NOT NULL,
    src_reserveaddress_id bigint NOT NULL,
    dst_reserveaddress_id bigint NOT NULL,
    src_amount double precision,
    dest_amount double precision,
    wallet_address_id bigint NOT NULL,
    src_burn_amount double precision,
    dst_burn_amount double precision,
    src_wallet_fee_amount double precision,
    dst_wallet_fee_amount double precision,
    integration_app text,
    ip text,
    country text,
    ethusd_rate double precision,
    ethusd_provider text,
    index integer
);


ALTER TABLE public.tradelogs OWNER TO reserve_stats;

--
-- Name: tradelogs_id_seq; Type: SEQUENCE; Schema: public; Owner: reserve_stats
--

CREATE SEQUENCE public.tradelogs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.tradelogs_id_seq OWNER TO reserve_stats;

--
-- Name: tradelogs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reserve_stats
--

ALTER SEQUENCE public.tradelogs_id_seq OWNED BY public.tradelogs.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: reserve_stats
--

CREATE TABLE public.users (
    id integer NOT NULL,
    address text NOT NULL,
    "timestamp" timestamp without time zone
);


ALTER TABLE public.users OWNER TO reserve_stats;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: reserve_stats
--

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO reserve_stats;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reserve_stats
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: wallet; Type: TABLE; Schema: public; Owner: reserve_stats
--

CREATE TABLE public.wallet (
    id integer NOT NULL,
    address text NOT NULL
);


ALTER TABLE public.wallet OWNER TO reserve_stats;

--
-- Name: wallet_id_seq; Type: SEQUENCE; Schema: public; Owner: reserve_stats
--

CREATE SEQUENCE public.wallet_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.wallet_id_seq OWNER TO reserve_stats;

--
-- Name: wallet_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: reserve_stats
--

ALTER SEQUENCE public.wallet_id_seq OWNED BY public.wallet.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.reserve ALTER COLUMN id SET DEFAULT nextval('public.reserve_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.token ALTER COLUMN id SET DEFAULT nextval('public.token_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs ALTER COLUMN id SET DEFAULT nextval('public.tradelogs_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.wallet ALTER COLUMN id SET DEFAULT nextval('public.wallet_id_seq'::regclass);


--
-- Data for Name: reserve; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

COPY public.reserve (id, address) FROM stdin;
1	0x63825c174ab367968EC60f061753D3bbD36A0D8F
2	0x0000000000000000000000000000000000000000
3	0x21433Dec9Cb634A23c6A4BbcCe08c83f5aC2EC18
4	0x56e37b6b79d4E895618B8Bb287748702848Ae8c0
\.


--
-- Name: reserve_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.reserve_id_seq', 5, true);


--
-- Data for Name: token; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

COPY public.token (id, address) FROM stdin;
1	0x595832F8FC6BF59c85C527fEC3740A1b7a361269
2	0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE
3	0xfe5F141Bf94fE84bC28deD0AB966c16B17490657
6	0xdd974D5C2e2928deA5F71b9825b8b646686BD200
7	0xF433089366899D83a9f26A773D59ec7eCF30355e
12	0x89d24A6b4CcB1B6fAA2625fE562bDD9a23260359
13	0x0F5D2fB29fb7d3CFeE444a200298f468908cC942
16	0x4156D3342D5c385a87D264F90653733592000581
18	0x514910771AF9Ca656af840dff83E8264EcF986CA
19	0xC5bBaE50781Be1669306b9e001EFF57a2957b09d
22	0x23Ccc43365D9dD3882eab88F43d515208f832430
\.


--
-- Name: token_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.token_id_seq', 23, true);


--
-- Data for Name: tradelogs; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

COPY public.tradelogs (id, "timestamp", block_number, tx_hash, eth_amount, user_address_id, src_address_id, dest_address_id, src_reserveaddress_id, dst_reserveaddress_id, src_amount, dest_amount, wallet_address_id, src_burn_amount, dst_burn_amount, src_wallet_fee_amount, dst_wallet_fee_amount, integration_app, ip, country, ethusd_rate, ethusd_provider, index) FROM stdin;
1	2018-10-11 08:45:11	6494045	0xce56f715862b458bfe9a2fc7059707efee5827f36f5be01edddc013277aa99fa	1.25605787600702801	1	1	2	1	2	1424.55751099999998	1.25605787600702801	1	1.74278030295975128	0	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	91
2	2018-10-11 08:48:02	6494053	0xce85c30aaffd7d13722a6e3f3b4014c575a4c6b5f9973d66a0547e76ce4aba41	0.594180157001091391	2	3	2	1	2	3170.0215575518755	0.594180157001091391	1	0.8244249678390142	0	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	56
3	2018-10-11 08:48:26	6494055	0xeabbfa4e070cd23a9be7b8bbec47ee8e07f037440d72faa983178642b9d491fc	0.00100000000000000002	3	2	6	2	1	0.00100000000000000002	0.532114987602218537	1	0	0.00138750000000000006	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	18
4	2018-10-11 08:54:03	6494072	0x963fe56aaa2add83bcd8ed1326c0b150c738afb462f4711a25b8ffa83d464d6a	1.00454822773370611	1	7	2	3	2	341.230954139999994	1.00454822773370611	1	0.669029119670648376	0	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	126
5	2018-10-11 08:54:39	6494074	0x190eb1f5432253005e95f5291d9acf8c8c13b752c1c80cd161271d07ba5a7b84	0.00111867653591824802	3	6	2	1	2	0.599999999999999978	0.00111867653591824802	1	0.0015521636935865691	0	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	22
6	2018-10-11 08:56:37	6494088	0x3ef81b263846d9a8674938ad23c346685e34bea48cd5cb3110727498b4f298f8	0.00111747870101999995	3	6	12	3	3	0.599999999999999978	0.220337540440739582	1	0.000744240814879320017	0.000744240814879320017	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	17
7	2018-10-11 08:56:37	6494088	0x00b0343a905007b8a1b89c73b19a85b3a7307286841570c14643c81644ad8b39	1.00519751875436181	1	13	2	1	2	3061.07267668513168	1.00519751875436181	1	1.39471155727167684	0	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	46
8	2018-10-11 08:59:32	6494106	0x2c7c31a01431e6d9abb1769b87b99f8c0d6f660e7fcada13276a2f7cce52abfa	12	8	2	16	2	1	12	4397.85565915000006	2	0	9.99000000000000021	0	6.66000000000000014	ThirdParty	\N	\N	225.662016177461652	coingecko	5
9	2018-10-11 08:59:32	6494106	0xb05789dd7a2f5431fb8ff70e15ba89732b6f9eb8382debb2f00f5fd42f27d960	0.471510913161030443	9	2	18	2	1	0.471510913161030443	295.712956760479472	1	0	0.65422139201092977	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	48
10	2018-10-11 09:04:26	6494129	0x31b7bdc67dcd63a6175086eb34ce37559f7bca7788d6ebf08f193fdddeb92e45	1.00617464264818923	1	19	2	1	2	3125.2175400000001	1.00617464264818923	1	1.39606731667436268	0	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	104
11	2018-10-11 09:04:41	6494131	0xd847a7247826203560a393af2e0a64e732e9c8bdc817bd92758d2ab26804c659	0.0500000000000000028	11	2	22	2	4	0.0500000000000000028	489.501869250098025	1	0	0.0693750000000000061	0	0	KyberSwap	\N	\N	225.662016177461652	coingecko	5
\.


--
-- Name: tradelogs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.tradelogs_id_seq', 11, true);


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

COPY public.users (id, address, "timestamp") FROM stdin;
1	0x64C0372Ebc1F812398Bf3e475117Fd48D5EfB035	2018-10-11 08:45:11
2	0x96159285B88d578BEb3e241E5F8a1bBe31e4C683	2018-10-11 08:48:02
3	0x8fA07F46353A2B17E92645592a94a0Fc1CEb783F	2018-10-11 08:48:26
8	0x17D79F467243c5DB655282Ce6187127c42986413	2018-10-11 08:59:32
9	0xA41983E9baa92bA284A75dB1dB2bCbAfb763B033	2018-10-11 08:59:32
11	0x0826601F28B691CEEa2Be05EC1c922Ea0eC2d82D	2018-10-11 09:04:41
\.


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.users_id_seq', 11, true);


--
-- Data for Name: wallet; Type: TABLE DATA; Schema: public; Owner: reserve_stats
--

COPY public.wallet (id, address) FROM stdin;
1	0x0000000000000000000000000000000000000000
2	0xDECAF9CD2367cdbb726E904cD6397eDFcAe6068D
\.


--
-- Name: wallet_id_seq; Type: SEQUENCE SET; Schema: public; Owner: reserve_stats
--

SELECT pg_catalog.setval('public.wallet_id_seq', 3, true);


--
-- Name: reserve_address_key; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.reserve
    ADD CONSTRAINT reserve_address_key UNIQUE (address);


--
-- Name: reserve_pkey; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.reserve
    ADD CONSTRAINT reserve_pkey PRIMARY KEY (id);


--
-- Name: token_address_key; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.token
    ADD CONSTRAINT token_address_key UNIQUE (address);


--
-- Name: token_pkey; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.token
    ADD CONSTRAINT token_pkey PRIMARY KEY (id);


--
-- Name: tradelogs_pkey; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_pkey PRIMARY KEY (id);


--
-- Name: users_address_key; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_address_key UNIQUE (address);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: wallet_address_key; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.wallet
    ADD CONSTRAINT wallet_address_key UNIQUE (address);


--
-- Name: wallet_pkey; Type: CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.wallet
    ADD CONSTRAINT wallet_pkey PRIMARY KEY (id);


--
-- Name: tradelogs_dest_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_dest_address_id_fkey FOREIGN KEY (dest_address_id) REFERENCES public.token(id);


--
-- Name: tradelogs_dst_reserveaddress_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_dst_reserveaddress_id_fkey FOREIGN KEY (dst_reserveaddress_id) REFERENCES public.reserve(id);


--
-- Name: tradelogs_src_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_src_address_id_fkey FOREIGN KEY (src_address_id) REFERENCES public.token(id);


--
-- Name: tradelogs_src_reserveaddress_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_src_reserveaddress_id_fkey FOREIGN KEY (src_reserveaddress_id) REFERENCES public.reserve(id);


--
-- Name: tradelogs_user_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_user_address_id_fkey FOREIGN KEY (user_address_id) REFERENCES public.users(id);


--
-- Name: tradelogs_wallet_address_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: reserve_stats
--

ALTER TABLE ONLY public.tradelogs
    ADD CONSTRAINT tradelogs_wallet_address_id_fkey FOREIGN KEY (wallet_address_id) REFERENCES public.wallet(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: reserve_stats
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM reserve_stats;
GRANT ALL ON SCHEMA public TO reserve_stats;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

