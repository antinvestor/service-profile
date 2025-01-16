

CREATE INDEX ON devices USING hnsw (embedding vector_l2_ops);

--        db.Exec("CREATE INDEX ON items USING ivfflat (embedding vector_l2_ops) WITH (lists = 100)")