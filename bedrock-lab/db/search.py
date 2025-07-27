import sqlite3
import os

def search_documents(keyword):
    # Resolve path to rag_data.db no matter where the script is run from
    base_dir = os.path.dirname(os.path.abspath(__file__))
    db_path = os.path.join(base_dir, "../data/rag_data.db")  # Go up one dir, then into data/

    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    query = """
        SELECT title, content FROM documents
        WHERE content LIKE ?
    """
    cursor.execute(query, (f"%{keyword}%",))

    results = cursor.fetchall()
    conn.close()
    return results