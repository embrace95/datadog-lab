AI Ops Survival Playbook

This playbook is a collection of practical knowledge, debugging workflows, tips, and command references for surviving and thriving in environments that leverage AI-driven observability and operations automation.

1. 🔍 Observability: Metrics, Logs, and Traces

Datadog / Prometheus / OpenTelemetry

Always check which tags are being used in metrics. Wrong tags = wrong grouping.

Use avg(last_5m) instead of raw values to smooth out noise.

Enable trace retention filters: Drop noisy DB spans but keep parent request context.

Prefer high cardinality tags on traces, not metrics.

Debug Tip:

If you see a spike in 5xx errors:

kubectl logs <pod> --since=10m | grep -i error

Then pivot:

Into trace ID (if sampled)

Into APM screen filtered by @http.status_code:500 or service name

2. 🤖 Incident Automation

Runbooks with Context:

Every monitor should have a doc_link, dashboard_link, and pager_policy reference

Tools

PagerDuty + Datadog + Slack integrations

Use EventBridge to route anomaly detection to Lambda responders

Lambda Responder Example:

import boto3

def lambda_handler(event, context):
    if event['detail-type'] == 'Datadog Alert':
        # auto scale, quarantine, or notify
        print("Received alert: ", event['detail'])

3. 🧠 Anomaly Detection + AI Forecasting

Basic Strategy

Use AI forecasts for capacity prediction, not instant alerting

Use outlier() or anomalies() only when baseline is stable for at least 24–48h

Use Cases

Predictive scaling

Memory leaks (steady ramp over 24h)

Network or disk saturation (noisy patterns)

4. ⚙️ AI + Infra Examples

1. Dynamic AutoScaling via AI

Use SageMaker + CloudWatch metrics

Auto-tune HPA thresholds using AI-predicted CPU/RAM

2. GPT-Triggered Playbooks

Use LLMs to summarize logs and run diagnostic flows:

kubectl logs pod-x | llm-summarize

5. 🛡️ AI Security Monitoring

Flag sudden changes in IAM usage patterns

GPT prompt fuzzing to simulate attack chains on AI bots

aws cloudtrail lookup-events --lookup-attributes AttributeKey=EventName,AttributeValue=AssumeRole

Monitor for bursty AssumeRole calls in non-standard hours.

6. 🧪 AIOps Testing Tips

Test AI workflows in sandbox before integrating into real pager flows

Validate anomaly models on backfilled time windows first

Log all AI model decisions into a separate audit log bucket

7. 🧰 Quick CLI/Infra Recipes

See noisy metrics (CPU):

kubectl top pod --sort-by=cpu

View DNS failures:

kubectl get events --field-selector type=Warning | grep DNS

Force APM trace sample:

DD_TRACE_SAMPLE_RATE=1.0 flask run

8. 🧠 Remember:

AI should enhance SRE judgment — not replace it.

Treat AI as another signal stream. Watch it like you'd watch a junior engineer: capable of brilliance or chaos.

📌 To Add:

LangChain pipelines for log summarization

Integration templates: Datadog ↔ LLM ↔ Incident tools

Embedding use cases: clustering outage themes from logs
