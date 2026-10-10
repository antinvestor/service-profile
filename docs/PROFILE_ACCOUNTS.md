# Profile accounts

A person's profile is the only identity Stawi keeps. The profile service also
owns the person's chain account address: it derives the address, stores it on
the profile, and publishes it. This implements the Stawi chain-first rooms
design §4 and the profile accounts design §2.1.

## Derivation

```
identity_salt      = HMAC-SHA256(K_identity, "stawi/identity/v1" ‖ profile_id)
identity_salt_hash = keccak256(identity_salt)
address            = derive.AccountAddress(FamilyEVM, account_version, identity_salt_hash,
                       factory, keccak256(creation_code ‖ identity_salt_hash))
```

* `derive` is `github.com/stawilabs/stawi/pkg/protocol/derive`. It is the same
  code the contracts are vector-tested against. `apps/default/service/accounts`
  calls it and copies nothing. The vector in
  `apps/default/service/accounts/testdata/account_derivation.json` is stawi's
  `packages/contracts/test/vectors/account_derivation.json`. The test asserts
  that profile derives the address the Solidity `AccountFactory` computes.
* `K_identity` is a root key and never rotates, because rotating it would move
  every address. In production it stays inside Vault Transit. Only
  `identity_salt_hash` and the address leave the service. The salt itself is
  never stored or published.
* The init code takes the identity salt hash as its constructor argument, so the
  init code hash is different for every profile. Because of that, the service
  needs the StawiAccount **creation code**; a single init code hash is not
  enough.

## Who gets an account

* **PERSON** profiles: they get a primary account in the same transaction that
  creates the profile. The salt is computed before that transaction opens. If
  Vault is unavailable, the create fails with `UNAVAILABLE`, so a person never
  exists without an account.
* **BOT** and **INSTITUTION** profiles: they get no account.
* **Existing profiles**: the `account-backfill` setup step derives an account
  for every PERSON profile that has no primary account for the configured
  version. It works in batches, is idempotent, and stages one fact per account
  it creates. A legacy migrate job (argv `migrate` / `DO_MIGRATION`) runs only
  the well-known setup steps, so in that mode the backfill runs as part of
  `bootstrap`.
* **Merge**: the surviving profile keeps its primary account. The merged
  profile's accounts move to the survivor as secondary accounts
  (`primary=false`).
  * One database transaction moves and demotes the accounts, saves the
    survivor, deletes the merged profile and stages `profile.accounts_merged`.
    A failure at any step changes nothing.
  * Both profiles must be in the same tenant and partition.
  * An authenticated caller can only merge inside its own tenancy, where
    `profile_merge` was granted. Naming a profile from another tenancy returns
    `NOT_FOUND`.
  * The account rows are selected and updated with the source profile's tenant
    and partition in the WHERE clause.
* **New canonical factory**: when the protocol publishes a new factory, raise
  `STAWI_ACCOUNT_VERSION`. The backfill then adds an account for that version,
  and the existing rows stay.

The accounts are stored in `profile_accounts`:

| Column | Contents |
|---|---|
| `profile_id` | the owning profile |
| `family` | `EVM` |
| `account_version` | the account version |
| `address` | the address; unique |
| `identity_salt_hash` | the keccak256 of the identity salt |
| `factory` | the factory address |
| `is_primary` | whether this is the profile's primary account |

A partial unique index allows one live primary account per profile, family and
version.

## API (`profile.v1`)

* `ProfileObject.accounts` (field 7): a list of `ProfileAccount`, with fields
  `address`, `family`, `version`, `identity_salt_hash` and `primary`. The
  primary account comes first. Addresses are lowercase `0x` hex.
  * The link between an address and a profile is the privacy boundary of the
    design. The field is filled only when the caller is the profile's owner
    (token subject = profile id) or a service principal (holds
    `account_resolve`, which only `ROLE_SERVICE` has).
  * Any other caller, including a user who holds `profile_view`, gets the
    profile without accounts. This applies on every read path: `GetById`,
    `GetByIDAndPartition`, `GetByContact`, `Search`, `Create`, `Update`,
    `Merge`, and the relationship listings.
  * Account reads are scoped to the caller's tenant and partition, so a caller
    in another tenancy sees none.
* `ResolveAccounts(ResolveAccountsRequest{addresses})`: returns
  `ResolveAccountsResponse{data: [AccountOwner{address, profile_id}]}`.
  * It accepts at most 500 addresses, in any case.
  * Addresses that no profile owns are left out of the response. So are
    addresses outside the caller's tenancy.
  * It requires permission `account_resolve`, which is bound to `ROLE_SERVICE`
    only. Owners and admins are refused.

## Facts (profile events queue, transactional outbox)

| Fact | Payload |
|---|---|
| `profile.account_created` | `profile_id`, `address`, `family`, `version`, `identity_salt_hash`, `factory`, `primary` |
| `profile.accounts_merged` | `surviving_profile_id`, `merged_profile_id`, `addresses` |

These facts go only to the profile events queue, which services consume; no
user-facing API exposes them. Each fact also carries the outbox envelope fields `event_id`, `name` and
`occurred_at`. Delivery is at least once, so consumers must deduplicate on
`event_id`.

## Jurisdiction

The `jurisdiction` profile property is an ISO 3166-1 alpha-2 code. The service
trims it, upper-cases it and checks it against the countries table. A value
that is not a string, is not two letters, or names an unknown country is
rejected with `INVALID_ARGUMENT`. The check applies on `Create` and on
`Update`, for both scoped and global properties.

## Configuration

| Variable | Default | Meaning |
|---|---|---|
| `VAULT_ADDR` | | Vault address; turns on the Transit salter |
| `VAULT_K8S_AUTH_ROLE` | | Kubernetes auth role. If empty, `VAULT_TOKEN` is used instead (dev servers) |
| `VAULT_K8S_AUTH_MOUNT` | `kubernetes` | Kubernetes auth mount |
| `VAULT_K8S_TOKEN_PATH` | in-pod default | service-account token path |
| `VAULT_TRANSIT_MOUNT` | `transit` | Transit mount |
| `STAWI_IDENTITY_TRANSIT_KEY` | `stawi-identity` | Transit key |
| `STAWI_IDENTITY_KEY_VERSION` | | **required with Vault**; the pinned key version |
| `STAWI_IDENTITY_STATIC_KEY` | | hex key (32 bytes or more) for tests and local stacks; refused when `ENVIRONMENT` or `SERVICE_ENVIRONMENT` is `prod`/`production` |
| `STAWI_ACCOUNT_FACTORY` | | canonical `AccountFactory` address |
| `STAWI_ACCOUNT_CREATION_CODE` | | StawiAccount creation code as hex (`AccountFactory.accountCreationCode()`) |
| `STAWI_ACCOUNT_CREATION_CODE_FILE` | | the same, read from a file (for example a mounted ConfigMap) |
| `STAWI_ACCOUNT_CREATION_CODE_HASH` | | optional; the manifest's `creation_code_hash`, checked at startup |
| `STAWI_ACCOUNT_VERSION` | `1` | account version |
| `ACCOUNT_BACKFILL_BATCH_SIZE` | `200` | profiles per backfill batch |

How the service behaves with each configuration:

* **No salter configured** (neither `VAULT_ADDR` nor
  `STAWI_IDENTITY_STATIC_KEY`): the service logs a warning and creates profiles
  without accounts. This is the development posture.
* **Account variables set, no salter**: the service refuses to start.
* **Salter configured, account variables missing or malformed**: the service
  refuses to start.

The setup job also needs the Vault and account variables, because it runs the
backfill.

## Vault setup

Create the Transit key once. It must not be exportable, and it must never be
rotated:

```sh
vault secrets enable transit   # if not already enabled
vault write -f transit/keys/stawi-identity type=aes256-gcm96 exportable=false allow_plaintext_backup=false
vault read -field=latest_version transit/keys/stawi-identity   # -> STAWI_IDENTITY_KEY_VERSION (1)
```

Create a policy that allows the HMAC endpoint of this key and nothing else:

```hcl
# stawi-identity-hmac.hcl
path "transit/hmac/stawi-identity" {
  capabilities = ["update"]
}
path "transit/hmac/stawi-identity/sha2-256" {
  capabilities = ["update"]
}
```

Bind the policy to the service with Kubernetes auth:

```sh
vault policy write stawi-identity-hmac stawi-identity-hmac.hcl
vault write auth/kubernetes/role/service-profile \
  bound_service_account_names=service-profile \
  bound_service_account_namespaces=profile \
  token_policies=stawi-identity-hmac \
  token_ttl=1h token_max_ttl=4h
```

Then set `VAULT_ADDR`, `VAULT_K8S_AUTH_ROLE=service-profile` and
`STAWI_IDENTITY_KEY_VERSION=1` on both the service and its setup job.

The salter calls `transit/hmac/stawi-identity/sha2-256` with `key_version`
pinned and `batch_input`, so the backfill makes one call per batch. It logs in
again before the token expires.
