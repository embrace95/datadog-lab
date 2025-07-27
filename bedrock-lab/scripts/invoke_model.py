import boto3
import json
import os

def load_prompt(key):
    # 🔍 Get absolute path to `config/prompts.json`, no matter where you're calling from
    script_dir = os.path.dirname(os.path.abspath(__file__))  # scripts/
    prompt_path = os.path.join(script_dir, "../config/prompts.json")

    # Debug: print the resolved prompt path
    print("🔍 Looking for prompts at:", os.path.abspath(prompt_path))

    # Load the prompts
    with open(prompt_path, "r") as f:
        prompts = json.load(f)

    # Handle missing key gracefully
    if key not in prompts:
        raise KeyError(f"Prompt key '{key}' not found in prompts.json.")
    
    return prompts[key]

def invoke_bedrock(prompt_key):
    prompt = load_prompt(prompt_key)

    body = json.dumps({
        "prompt": prompt,
        "max_tokens_to_sample": 200,
        "temperature": 0.7,
        "stop_sequences": ["\n\nHuman:"]
    })

    # 🔐 Replace this with any other model you’re using (e.g., Claude 3)
    client = boto3.client("bedrock-runtime", region_name="us-east-1")
    response = client.invoke_model(
        modelId="anthropic.claude-instant-v1",
        body=body,
        contentType="application/json",
        accept="application/json"
    )
    #print(response)
    # Convert streaming body to string and parse JSON
    response_body = json.loads(response["body"].read())
    return response_body["completion"]
