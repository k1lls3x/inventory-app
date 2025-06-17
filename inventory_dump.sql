--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: User; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."User" (
    user_id integer NOT NULL,
    username character varying(50) NOT NULL,
    password text NOT NULL,
    full_name character varying(100),
    role character varying(20) NOT NULL,
    CONSTRAINT "User_role_check" CHECK (((role)::text = ANY ((ARRAY['admin'::character varying, 'manager'::character varying, 'worker'::character varying])::text[])))
);


ALTER TABLE public."User" OWNER TO postgres;

--
-- Name: User_user_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."User_user_id_seq"
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."User_user_id_seq" OWNER TO postgres;

--
-- Name: User_user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."User_user_id_seq" OWNED BY public."User".user_id;


--
-- Name: inbound; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.inbound (
    inbound_id integer NOT NULL,
    item_id integer,
    supplier_id integer,
    quantity integer NOT NULL,
    received_at timestamp without time zone DEFAULT now(),
    warehouse_id integer
);


ALTER TABLE public.inbound OWNER TO postgres;

--
-- Name: inbound_inbound_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.inbound_inbound_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.inbound_inbound_id_seq OWNER TO postgres;

--
-- Name: inbound_inbound_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.inbound_inbound_id_seq OWNED BY public.inbound.inbound_id;


--
-- Name: item; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.item (
    item_id integer NOT NULL,
    sku character varying(50) NOT NULL,
    name character varying(200) NOT NULL,
    description text,
    uom character varying(20) NOT NULL,
    reorder_level integer DEFAULT 0,
    reorder_qty integer DEFAULT 0,
    price numeric(12,2),
    cost numeric(12,2),
    created_at timestamp without time zone DEFAULT now()
);


ALTER TABLE public.item OWNER TO postgres;

--
-- Name: item_item_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.item_item_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.item_item_id_seq OWNER TO postgres;

--
-- Name: item_item_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.item_item_id_seq OWNED BY public.item.item_id;


--
-- Name: outbound; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.outbound (
    outbound_id integer NOT NULL,
    item_id integer,
    quantity integer NOT NULL,
    shipped_at timestamp without time zone DEFAULT now(),
    destination text,
    warehouse_id integer
);


ALTER TABLE public.outbound OWNER TO postgres;

--
-- Name: outbound_outbound_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.outbound_outbound_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.outbound_outbound_id_seq OWNER TO postgres;

--
-- Name: outbound_outbound_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.outbound_outbound_id_seq OWNED BY public.outbound.outbound_id;


--
-- Name: stock; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.stock (
    stock_id integer NOT NULL,
    item_id integer,
    quantity integer DEFAULT 0 NOT NULL,
    warehouse_id integer,
    last_updated timestamp without time zone DEFAULT now()
);


ALTER TABLE public.stock OWNER TO postgres;

--
-- Name: stock_stock_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.stock_stock_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.stock_stock_id_seq OWNER TO postgres;

--
-- Name: stock_stock_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.stock_stock_id_seq OWNED BY public.stock.stock_id;


--
-- Name: supplier; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.supplier (
    supplier_id integer NOT NULL,
    name character varying(100) NOT NULL,
    contact_info text,
    inn character varying(20),
    contact_person character varying(100),
    phone character varying(32),
    email character varying(64)
);


ALTER TABLE public.supplier OWNER TO postgres;

--
-- Name: supplier_supplier_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.supplier_supplier_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.supplier_supplier_id_seq OWNER TO postgres;

--
-- Name: supplier_supplier_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.supplier_supplier_id_seq OWNED BY public.supplier.supplier_id;


--
-- Name: warehouse; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.warehouse (
    warehouse_id integer NOT NULL,
    name character varying(100) NOT NULL,
    location text
);


ALTER TABLE public.warehouse OWNER TO postgres;

--
-- Name: warehouse_warehouse_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.warehouse_warehouse_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.warehouse_warehouse_id_seq OWNER TO postgres;

--
-- Name: warehouse_warehouse_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.warehouse_warehouse_id_seq OWNED BY public.warehouse.warehouse_id;


--
-- Name: User user_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User" ALTER COLUMN user_id SET DEFAULT nextval('public."User_user_id_seq"'::regclass);


--
-- Name: inbound inbound_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbound ALTER COLUMN inbound_id SET DEFAULT nextval('public.inbound_inbound_id_seq'::regclass);


--
-- Name: item item_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item ALTER COLUMN item_id SET DEFAULT nextval('public.item_item_id_seq'::regclass);


--
-- Name: outbound outbound_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.outbound ALTER COLUMN outbound_id SET DEFAULT nextval('public.outbound_outbound_id_seq'::regclass);


--
-- Name: stock stock_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock ALTER COLUMN stock_id SET DEFAULT nextval('public.stock_stock_id_seq'::regclass);


--
-- Name: supplier supplier_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.supplier ALTER COLUMN supplier_id SET DEFAULT nextval('public.supplier_supplier_id_seq'::regclass);


--
-- Name: warehouse warehouse_id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouse ALTER COLUMN warehouse_id SET DEFAULT nextval('public.warehouse_warehouse_id_seq'::regclass);


--
-- Data for Name: User; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public."User" (user_id, username, password, full_name, role) FROM stdin;
13	a	$2a$10$okwmJjU1KBJwpe7FUlOAQ.Zm1nDmYAX9/KYce7iTV3VBzKm7LysJO	Дима Крут	admin
14	1	$2a$10$7yWsfnh6meaNPxH2jMUEgucV/w2Xa7.vzNsk1Mhg3PKO0jyANWvE2	Дима Манага	manager
\.


--
-- Data for Name: inbound; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.inbound (inbound_id, item_id, supplier_id, quantity, received_at, warehouse_id) FROM stdin;
15	7	1	31	2025-05-23 03:00:00	3
16	1	1	25	2025-05-17 03:00:00	3
17	1	2	26	2025-05-24 03:00:00	1
18	7	6	122	2025-05-25 03:00:00	4
\.


--
-- Data for Name: item; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.item (item_id, sku, name, description, uom, reorder_level, reorder_qty, price, cost, created_at) FROM stdin;
1	SKU001	Клавиатура Redragon	Геймерская клавиатура с RGB	шт	10	30	3500.00	2500.00	2025-05-12 15:40:43.417666
2	SKU002	Мышь Logitech	Мышь беспроводная	шт	15	50	2000.00	1400.00	2025-05-12 15:40:43.417666
3	SKU003	Монитор Samsung 24"	IPS, 75 Гц	шт	5	10	12000.00	9500.00	2025-05-12 15:40:43.417666
4	SKU004	Коврик SteelSeries	Игровой коврик для мыши	шт	25	100	700.00	400.00	2025-05-12 15:40:43.417666
5	SKU005	Гарнитура HyperX	Наушники с микрофоном	шт	10	20	4500.00	3000.00	2025-05-12 15:40:43.417666
7	BIBI	Бим барам	Тот самый	шт.	0	0	10000.00	1000.00	2025-05-13 14:55:11.679406
\.


--
-- Data for Name: outbound; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.outbound (outbound_id, item_id, quantity, shipped_at, destination, warehouse_id) FROM stdin;
6	7	12	2025-05-25 03:00:00	12	2
7	7	122	2025-05-24 03:00:00	НоРМ	2
8	5	122	2025-05-14 03:00:00	-	3
5	5	112	2025-05-09 00:00:00	112	3
\.


--
-- Data for Name: stock; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.stock (stock_id, item_id, quantity, warehouse_id, last_updated) FROM stdin;
1	1	40	1	2025-05-12 15:41:00.132548
4	4	100	1	2025-05-12 15:41:00.132548
3	3	53	2	2025-05-14 17:41:03.868774
2	2	120	1	2025-05-14 17:41:07.786347
16	1	222	3	2025-05-24 16:06:21.239764
17	7	111	2	2025-05-24 16:24:30.0504
21	1	111	2	2025-05-24 16:27:22.279714
22	4	111	3	2025-05-24 16:27:24.538717
5	5	222	3	2025-05-24 16:31:11.567574
\.


--
-- Data for Name: supplier; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.supplier (supplier_id, name, contact_info, inn, contact_person, phone, email) FROM stdin;
2	ЗАО "СнабЭксперт"	snab@example.com	\N	\N	\N	\N
3	ИП "Доставкин"	deliveryman@sample.org	\N	\N	\N	\N
5	ООО Ромашка	\N	1234567890	Иванов И.И.	+7 900 1112233	contact@romashka.ru
1	ООО "Поставка+", Москва	info@postavka-plus.ru	1488	Нет	+7995472004	takesxq77@gmail.com
6	Овеликий	\N	161	МЯУ	+7788873773	sto@mail.com
\.


--
-- Data for Name: warehouse; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.warehouse (warehouse_id, name, location) FROM stdin;
1	Центральный склад	г. Москва
3	Резервный склад	г. Екатеринбург
2	Запад склад	г. Смоленск
4	Четкий	Четкая
\.


--
-- Name: User_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public."User_user_id_seq"', 15, true);


--
-- Name: inbound_inbound_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.inbound_inbound_id_seq', 18, true);


--
-- Name: item_item_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.item_item_id_seq', 8, true);


--
-- Name: outbound_outbound_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.outbound_outbound_id_seq', 8, true);


--
-- Name: stock_stock_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.stock_stock_id_seq', 24, true);


--
-- Name: supplier_supplier_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.supplier_supplier_id_seq', 6, true);


--
-- Name: warehouse_warehouse_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.warehouse_warehouse_id_seq', 4, true);


--
-- Name: User User_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_pkey" PRIMARY KEY (user_id);


--
-- Name: User User_username_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."User"
    ADD CONSTRAINT "User_username_key" UNIQUE (username);


--
-- Name: inbound inbound_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbound
    ADD CONSTRAINT inbound_pkey PRIMARY KEY (inbound_id);


--
-- Name: item item_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item
    ADD CONSTRAINT item_pkey PRIMARY KEY (item_id);


--
-- Name: item item_sku_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.item
    ADD CONSTRAINT item_sku_key UNIQUE (sku);


--
-- Name: outbound outbound_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.outbound
    ADD CONSTRAINT outbound_pkey PRIMARY KEY (outbound_id);


--
-- Name: stock stock_item_id_warehouse_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_item_id_warehouse_id_key UNIQUE (item_id, warehouse_id);


--
-- Name: stock stock_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_pkey PRIMARY KEY (stock_id);


--
-- Name: supplier supplier_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.supplier
    ADD CONSTRAINT supplier_pkey PRIMARY KEY (supplier_id);


--
-- Name: warehouse warehouse_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.warehouse
    ADD CONSTRAINT warehouse_pkey PRIMARY KEY (warehouse_id);


--
-- Name: inbound inbound_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbound
    ADD CONSTRAINT inbound_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.item(item_id);


--
-- Name: inbound inbound_supplier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbound
    ADD CONSTRAINT inbound_supplier_id_fkey FOREIGN KEY (supplier_id) REFERENCES public.supplier(supplier_id);


--
-- Name: inbound inbound_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.inbound
    ADD CONSTRAINT inbound_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouse(warehouse_id);


--
-- Name: outbound outbound_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.outbound
    ADD CONSTRAINT outbound_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.item(item_id);


--
-- Name: outbound outbound_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.outbound
    ADD CONSTRAINT outbound_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouse(warehouse_id);


--
-- Name: stock stock_item_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_item_id_fkey FOREIGN KEY (item_id) REFERENCES public.item(item_id);


--
-- Name: stock stock_warehouse_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.stock
    ADD CONSTRAINT stock_warehouse_id_fkey FOREIGN KEY (warehouse_id) REFERENCES public.warehouse(warehouse_id);


--
-- PostgreSQL database dump complete
--

