import sqlite3
import os

# Step 1: Define path
base_dir = os.path.dirname(os.path.abspath(__file__))
#db_path = os.path.join(base_dir, "../data/rag_faq.db")
db_path = os.path.join(base_dir, "data", "rag_faq.db")

# Step 2: Connect and create table
conn = sqlite3.connect(db_path)
cursor = conn.cursor()

cursor.execute("""
CREATE TABLE IF NOT EXISTS faq (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    embedding BLOB
);
""")

# Step 3: Seed the DB with some entries
faqs = [
    ("What is Amazon Bedrock?", "Amazon Bedrock is a service that allows you to build and scale generative AI applications using foundation models via API."),
    ("How do I use Claude with Bedrock?", "You call the Bedrock runtime API with Claude’s model ID and pass in a prompt."),
    ("What is Retrieval-Augmented Generation (RAG)?", "It’s a technique where external documents are retrieved and used to improve LLM output.")
]

cursor.executemany("INSERT INTO faq (question, answer) VALUES (?, ?);", faqs)

conn.commit()
conn.close()

print("✅ rag_faq.db created at:", db_path)
