# rate-limit

Best for production systems

Each user has:
Bucket size = 5
Refill rate = 5 tokens/min

Flow:

Request comes → consume 1 token
If no token → reject

👉 Pros:

Smooth rate limiting
Handles bursts properly
Efficient