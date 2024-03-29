# JWT Decode

![Build Status](https://github.com/micahhausler/jwtdecode/actions/workflows/test.yml/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/micahhausler/jwtdecode/branch/master/graph/badge.svg)](https://codecov.io/gh/micahhausler/jwtdecode)
[![Go Reference](https://pkg.go.dev/badge/github.com/micahhausler/jwtdecode.svg)](https://pkg.go.dev/github.com/micahhausler/jwtdecode)
Decodes JWT tokens without verification of keys.

## Example

The following example prints a token file to stderr, and jwtdecode reads the
token from stdin. jwtdecode can read from stdin, a single file, or a list of
files. jwtdecode can also read multiple tokens in the same file that are
separated by newlines.

```bash
cat token.jwt | tee -a /dev/stderr | jwtdecode
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
{
    "alg": "HS256",
    "typ": "JWT"
}
{
    "iat": 1516239022,
    "name": "John Doe",
    "sub": "1234567890"
}
```

## Install

```bash
go get github.com/micahhausler/jwtdecode
```

## License

MIT License. See [LICENSE](LICENSE) for full text.
