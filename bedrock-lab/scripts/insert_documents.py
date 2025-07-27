import sqlite3
import os

DB_PATH = os.path.join(os.path.dirname(__file__), "../data/rag_data.db")

SAMPLE_DOCS = [
    ("Kubernetes Basics", "Kubernetes is a container orchestration platform that helps automate deployment, scaling, and management of applications."),
    ("Site Reliability Engineering", "SRE is a discipline that incorporates aspects of software engineering and applies them to infrastructure and operations problems."),
    ("AWS Bedrock", "Amazon Bedrock lets you build and scale generative AI applications using foundation models from AI21, Anthropic, Cohere, Meta, and Stability AI via an API."),
    ("DevOps Philosophy", "DevOps is the combination of cultural philosophies, practices, and tools that increases an organization’s ability to deliver applications and services at high velocity."),
    ("Cloud Cost Optimization", "Cloud cost optimization is the practice of reducing your overall cloud spend by identifying mismanaged resources and eliminating waste.")
]

def init_db():
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    cursor.execute("""
        CREATE TABLE IF NOT EXISTS documents (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            content TEXT NOT NULL
        );
    """)
    conn.commit()
    conn.close()
    print("✅ Database initialized with table 'documents'")

def insert_documents(docs):
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    cursor.executemany("""
        INSERT INTO documents (title, content) VALUES (?, ?);
    """, docs)

    conn.commit()
    conn.close()
    print(f"✅ Inserted {len(docs)} documents into 'documents' table")

if __name__ == "__main__":
    init_db()
    insert_documents(SAMPLE_DOCS)

