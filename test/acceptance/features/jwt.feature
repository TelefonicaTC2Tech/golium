Feature: JWT

  @jwt
  Scenario: Create and process a signed JWT with HS512
    Given the JWT signature algorithm "HS512"
      And the JWT symmetric key
          """
          [CONF:signSymmetricKey]
          """
    Given the JWT payload with the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | iat           | [NOW]            |
          | exp           | [NOW:+24h:unix]  |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |
      And I generate a signed JWT and store it in context "jwt.jws"
     When I process the signed JWT
          """
          [CTXT:jwt.jws]
          """
     Then the JWT must be valid
      And the JWT payload must have the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |

  @jwt
  Scenario: Create and process a signed JWT with RS256
    Given the JWT signature algorithm "RS256"
      And the JWT payload with the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | iat           | [NOW]            |
          | exp           | [NOW:+24h:unix]  |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |
      And the JWT private key
          """
          [CONF:signPrivateKey]
          """
     When I generate a signed JWT and store it in context "jwt.jws"
    Given the JWT public key
          """
          [CONF:signPublicKey]
          """
     When I process the signed JWT
          """
          [CTXT:jwt.jws]
          """
     Then the JWT must be valid
      And the JWT payload must have the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |

  @jwt
  Scenario: Create and process an expired signed JWT
    Given the JWT signature algorithm "RS256"
      And the JWT payload with the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | iat           | [NOW]            |
          | exp           | [NOW:-1m:unix]   |
      And the JWT private key
          """
          [CONF:signPrivateKey]
          """
     When I generate a signed JWT and store it in context "jwt.jws"
    Given the JWT public key
          """
          [CONF:signPublicKey]
          """
     When I process the signed JWT
          """
          [CTXT:jwt.jws]
          """
     Then the JWT must be invalid by "exp not satisfied"

  @jwt
  Scenario: Create and process a signed JWT with invalid signature
    Given the JWT signature algorithm "RS256"
      And the JWT public key
          """
          [CONF:signPublicKey]
          """
     When I process the signed JWT
          """
          eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibm
          FtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.j4Bf5hc9Vt2F8Ru5xsGOrN6CrufK
          48t3ycha440HWj-LphOGFSgcXUhYi0ZUZCxJt4BuIhkZkSiQAryZKXfFf8w9YiaA_Zw
          rupFr-QJXLRLWIUAymVrZdGf5Z5V6xK7nWNX1yifkHhTCzWm_z5jj3Dr-9VieUeK1qo
          SBQ2H0uFP8n3UeR5CGxuUAlWJvAaRNHUC9lCvl_bFETmAtKZ8bzhsIHEdgPzatNnxyX
          zuNfIOceDlFQkHFRXY5TUW17tssie1u-sY993MQwagnZZY9CHQ41EizBswNziUIALJk
          -gvvViXhysXQX4GdgW9_a8SXijYMmN65GUGti1aPCEP_mA
          """
     Then the JWT must be invalid by "crypto/rsa: verification error"

  @jwt
  Scenario: Create and process a JWE token (encrypted token)
    Given the JWT key encryption algorithm "RSA1_5"
      And the JWT content encryption algorithm "A128CBC-HS256"
      And the JWT payload with the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | iat           | [NOW]            |
          | exp           | [NOW:+24h:unix]  |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |
      And the JWT public key
          """
          [CONF:encryptPublicKey]
          """
     When I generate an encrypted JWT and store it in context "jwt.jwe"
    Given the JWT private key
          """
          [CONF:encryptPrivateKey]
          """
     When I process the encrypted JWT
          """
          [CTXT:jwt.jwe]
          """
      And the JWT payload must have the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |

  @jwt
  Scenario: Create and process a signed encrypted JWT token
    Given the JWT signature algorithm "RS256"
      And the JWT key encryption algorithm "RSA1_5"
      And the JWT content encryption algorithm "A128CBC-HS256"
      And the JWT payload with the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | iat           | [NOW]            |
          | exp           | [NOW:+24h:unix]  |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |
      And the JWT public key
          """
          [CONF:encryptPublicKey]
          """
      And the JWT private key
          """
          [CONF:signPrivateKey]
          """
     When I generate a signed encrypted JWT and store it in context "jwt.jwse"
    Given the JWT public key
          """
          [CONF:signPublicKey]
          """
      And the JWT private key
          """
          [CONF:encryptPrivateKey]
          """
     When I process the signed encrypted JWT
          """
          [CTXT:jwt.jwse]
          """
     Then the JWT must be valid
      And the JWT payload must have the JSON properties
          | sub           | golium           |
          | iss           | issuer           |
          | object.string | test with golium |
          | object.bool   | [TRUE]           |
          | object.number | [NUMBER:123.4]   |
