# 1. Use Client-Side AES-256-GCM Zero-Knowledge Encryption

* Status: Accepted
* Date: 2026-08-16
* Deciders: Platform Engineering Team

## Context and Problem Statement
Environment variables often contain high-risk credentials, connection strings, and API secrets. Transmitting or storing plaintext secrets on central servers creates a high-risk single point of failure.

## Decision Outcome
We decided to adopt **Client-Side Zero-Knowledge Encryption using AES-256-GCM with PBKDF2 Key Derivation**.

### Positive Consequences
- **Zero-Knowledge Architecture:** Central server & database only store encrypted binary payloads. Server operators cannot view secrets.
- **Data Safety:** Nonce-authenticated encryption prevents ciphertext tampering and replay attacks.
- **Compliance:** Satisfies SOC2 and DevSecOps security standards.

### Negative Consequences
- Users must manage master encryption passphrases securely. Loss of passphrase prevents secret decryption.
